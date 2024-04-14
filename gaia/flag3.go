package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"html/template"
	"log"
	"math/big"
	mathRand "math/rand"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	SessionName = "PHPSESSID"
)

var (
	gConfig  = Config{}
	gCrypto  = Crypto{}
	gUtility = Utility{}
	gWeb     = Web{}
)

// Crypto //////////////////////////////////////////////////////////////////////////////////////////////////////////////
type Crypto struct{}

func (_ *Crypto) Random(n int) []byte {
	buffer := make([]byte, n)
	if _, err := rand.Read(buffer); nil != err {
		panic(err)
	}
	return buffer
}

func (c *Crypto) Encrypt(key, clear []byte) ([]byte, error) {
	var err error
	var block cipher.Block
	encKey := c.SHA256(key)
	if block, err = aes.NewCipher(encKey); nil != err {
		return nil, err
	}
	var gcm cipher.AEAD
	if gcm, err = cipher.NewGCM(block); nil != err {
		return nil, err
	}
	nonce := c.Random(gcm.NonceSize())
	buffer := new(bytes.Buffer)
	buffer.Write(nonce)
	buffer.Write(gcm.Seal(nil, nonce, clear, key))
	return buffer.Bytes(), nil
}

func (c *Crypto) Decrypt(key, encrypted []byte) ([]byte, error) {
	var err error
	var block cipher.Block
	encKey := c.SHA256(key)
	if block, err = aes.NewCipher(encKey); nil != err {
		return nil, err
	}
	var gcm cipher.AEAD
	if gcm, err = cipher.NewGCM(block); nil != err {
		return nil, err
	}
	nonce := encrypted[:gcm.NonceSize()]
	enc := encrypted[gcm.NonceSize():]
	var buffer []byte
	if buffer, err = gcm.Open(nil, nonce, enc, key); nil != err {
		return nil, err
	}
	return buffer, nil
}

func (_ *Crypto) MD5(in []byte) []byte {
	var hashMaker = md5.New()
	hashMaker.Write(in)
	return hashMaker.Sum(nil)
}

func (_ *Crypto) SHA256(in []byte) []byte {
	var hashMaker = sha256.New()
	hashMaker.Write(in)
	return hashMaker.Sum(nil)
}

func (_ *Crypto) HMAC256(key, in []byte) []byte {
	var hashMaker = hmac.New(sha256.New, key)
	hashMaker.Write(in)
	return hashMaker.Sum(nil)
}

// Utility /////////////////////////////////////////////////////////////////////////////////////////////////////////////
type Utility struct{}

func (_ *Utility) CurrentAbsolutePathOfExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("get current path failed: %v", err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

func (_ *Utility) ReadFile(fn string) []byte {
	var err error
	var stream []byte
	if stream, err = os.ReadFile(fn); nil != err {
		_ = err
		return nil
	}
	return stream
}

func (_ *Utility) Serialize(o interface{}) []byte {
	var err error
	var buf = new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	if err = encoder.Encode(o); nil != err {
		panic(err)
	}
	return buf.Bytes()
}

func (u *Utility) SerializeEncrypt(key []byte, o interface{}) string {
	var err error
	var enc []byte
	stream := u.Serialize(o)
	if enc, err = gCrypto.Encrypt(key, stream); nil != err {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(enc)
}

func (_ *Utility) Deserialize(stream []byte, out interface{}) error {
	var err error
	decoder := gob.NewDecoder(bytes.NewReader(stream))
	if err = decoder.Decode(out); nil != err {
		return err
	}
	return nil
}

func (u *Utility) DecryptDeserialize(key []byte, b64stream string, out interface{}) error {
	var err error
	var dec, stream []byte
	if stream, err = base64.URLEncoding.DecodeString(b64stream); nil != err {
		return err
	}
	if dec, err = gCrypto.Decrypt(key, stream); nil != err {
		return err
	}
	return u.Deserialize(dec, out)
}

// Config //////////////////////////////////////////////////////////////////////////////////////////////////////////////
type Config struct {
	Host           string
	Port           uint16
	Timeout        int
	PasswordLength int
	Debug          bool
	SecretHex      string
	secret         []byte
}

func (c *Config) ReadFrom(TomlFilename string) {
	ext := path.Ext(TomlFilename)
	fn := TomlFilename[:len(TomlFilename)-len(ext)]
	viper.SetConfigName(fn) //设置文件名时不要带后缀
	viper.SetConfigType(ext[1:])
	viper.AddConfigPath(gUtility.CurrentAbsolutePathOfExecutable())                   //搜索路径可以设置多个，viper 会根据设置顺序依次查找
	viper.AddConfigPath(path.Join(gUtility.CurrentAbsolutePathOfExecutable(), "web")) //搜索路径可以设置多个，viper 会根据设置顺序依次查找
	viper.AddConfigPath(".")                                                          //搜索路径可以设置多个，viper 会根据设置顺序依次查找
	//viper.AddConfigPath(path.Join(".", "web"))                                        //搜索路径可以设置多个，viper 会根据设置顺序依次查找
	if err := viper.ReadInConfig(); nil != err {
		log.Fatalf("read config failed: %v", err)
	}
	if err := viper.Unmarshal(c); nil != err {
		log.Fatalf("unmarshal config failed: %v", err)
	}
}

func (c *Config) Secret() []byte {
	var err error
	if nil == c.secret {
		if c.secret, err = hex.DecodeString(c.SecretHex); nil != err {
			log.Fatalf("cannot decode SecretHex: %v", err)
			return nil
		}
		c.secret = gCrypto.SHA256(c.secret)
	}
	return c.secret
}

// Web /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type Web struct {
	Config
	Password       string
	Big0           uint64
	Big1           uint64
	Hash10         string // 真正的结果
	Hash10Cover    string // 打码
	Hash10Value    string // MD5打码
	StartTimestamp int64
	Step           int
}

func (w *Web) init(remoteIp string) {
	w.Config = gConfig
	{ // 根据远端IP地址生成固定的6位数字。
		ds := strings.Split(remoteIp, ".")
		var ip uint32
		ip = 0
		for i := 0; i < 4; i++ {
			a, _ := strconv.Atoi(ds[i])
			ip = (ip << 8) | uint32(a&0xff)
		}
		mathRand.Seed(int64(ip))
		w.Password = strconv.Itoa(mathRand.Int())
		if len(w.Password) > gConfig.PasswordLength {
			w.Password = w.Password[:gConfig.PasswordLength]
		} else if len(w.Password) < gConfig.PasswordLength {
			sb := strings.Builder{}
			for i := 0; i < gConfig.PasswordLength-len(w.Password); i++ {
				sb.WriteString("0")
			}
			sb.WriteString(w.Password)
			w.Password = sb.String()
		}
	}
	{ // 生成两个较长的整数
		const N = 16
		var buffer = gCrypto.Random(N)
		_ = binary.Read(bytes.NewBuffer(buffer), binary.LittleEndian, &w.Big0)
		_ = binary.Read(bytes.NewBuffer(buffer[8:]), binary.LittleEndian, &w.Big1)
		w.Big0 = w.Big0 >> 1 // 除以2防止两者相加结果大于uint64范围
		w.Big1 = w.Big1 >> 1 // 除以2防止两者相加结果大于uint64范围
	}
	{ // 生成MD5值与掩码
		const N = 10
		var buffer = gCrypto.Random(N)
		for i := 0; i < N; i++ {
			buffer[i] = (buffer[i] % 10) + '0'
		}

		w.Hash10 = string(buffer)
		var indexes = gCrypto.Random(2)
		var bufCov = make([]byte, N)
		copy(bufCov, buffer)
		indexes[0] = indexes[0] % N
		indexes[1] = indexes[1] % (N - 1)
		bufCov[indexes[0]] = '*'
		if indexes[1] < indexes[0] {
			bufCov[indexes[1]] = '*'
		} else {
			bufCov[indexes[1]+1] = '*'
		}
		w.Hash10Cover = string(bufCov)

		val := hex.EncodeToString(gCrypto.MD5(buffer))
		tmp := []byte(val)
		tmp[indexes[0]] = '*'
		w.Hash10Value = string(tmp)
	}
	{
		w.Step = 0
		w.StartTimestamp = 0
	}
}

func (_ *Web) serviceCookieSet(c *gin.Context, value string) {
	cookieName := &http.Cookie{
		Name:     SessionName,
		Value:    value,
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
		SameSite: 3,
	}
	http.SetCookie(c.Writer, cookieName)
}

func (_ *Web) serviceCookieGet(c *gin.Context) string {
	if cookie, err := c.Request.Cookie(SessionName); err == nil {
		return cookie.Value
	} else {
		return ""
	}
}
func (_ *Web) serviceCookieClear(c *gin.Context) {
	clearCookieName := &http.Cookie{
		Name:     SessionName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1, //<0 意味着删除cookie
		Secure:   false,
		HttpOnly: true,
	}
	http.SetCookie(c.Writer, clearCookieName)
}

func (w *Web) IndexGet(c *gin.Context) {
	w.init(c.RemoteIP())
	w.serviceCookieSet(c, gUtility.SerializeEncrypt(gConfig.Secret(), w))
	log.Printf("[DEBUG] remoteIp = %s, session = %v\n", c.RemoteIP(), *w)

	//html := w.render(path.Join(gUtility.CurrentAbsolutePathOfExecutable(), "web", "page0.html"), w)
	html := w.render(path.Join(gUtility.CurrentAbsolutePathOfExecutable(), "page0.html"), w)
	c.Header("Content-Type", "text/html; charset=UTF-8")
	c.String(http.StatusOK, "%v", html)
}

func (_ *Web) render(pagePath string, parameters any) string {
	var err error
	var stream []byte
	var tpl *template.Template
	var buffer bytes.Buffer
	if stream, err = os.ReadFile(pagePath); nil != err {
		return fmt.Sprintf("cannot read file [%v]: %v", pagePath, err)
	}
	if tpl, err = template.New(pagePath).Parse(string(stream)); nil != err {
		return fmt.Sprintf("cannot parse template: %v", err)
	}
	if err = tpl.Execute(&buffer, parameters); nil != err {
		return fmt.Sprintf("cannot execute template: %v", err)
	}
	return buffer.String()
}

func (w *Web) postPage0(c *gin.Context) string {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "admin" && password == w.Password {
		w.Step = 1
		w.StartTimestamp = time.Now().Unix()
		//return w.render(path.Join(gUtility.CurrentAbsolutePathOfExecutable(), "web", "page1.html"), w)
		return w.render(path.Join(gUtility.CurrentAbsolutePathOfExecutable(), "page1.html"), w)
	} else {
		return "账号或密码错误，请重试！"
	}
}

func (w *Web) postPage1(c *gin.Context) string {
	var now = time.Now().Unix()
	if now-w.StartTimestamp > int64(w.Timeout) {
		return "超时了！请重试！"
	}
	var bigAns *big.Int
	if bigAns, _ = new(big.Int).SetString(c.PostForm("answer"), 10); nil == bigAns {
		log.Printf("[DEBUG] SetString bigAns = %v", bigAns)
		return "计算结果错误！请重试！"
	}

	big0 := big.NewInt(int64(w.Big0))
	big1 := big.NewInt(int64(w.Big1))
	expAns := new(big.Int).Add(big0, big1)
	log.Printf("[DEBUG] %v + %v bigAns = %v , expAns = %v", w.Big0, w.Big1, bigAns, expAns)
	if bigAns.Cmp(expAns) == 0 {
		w.Step = 2
		w.StartTimestamp = 0
		//return w.render(path.Join(gUtility.CurrentAbsolutePathOfExecutable(), "web", "page2.html"), w)
		return w.render(path.Join(gUtility.CurrentAbsolutePathOfExecutable(), "page2.html"), w)
	} else {
		return "计算结果错误！请重试！"
	}
}

func (w *Web) postPage2(c *gin.Context) string {
	answer := c.PostForm("answer")
	if answer == w.Hash10 {
		//return w.render(path.Join(gUtility.CurrentAbsolutePathOfExecutable(), "web", "flag.txt"), w)
		return w.render(path.Join(gUtility.CurrentAbsolutePathOfExecutable(), "flag.txt"), w)
	}
	return "数据错误！请重试！"
}

func (w *Web) IndexPost(c *gin.Context) {
	var err error
	var html string

	s := w.serviceCookieGet(c)
	if err = gUtility.DecryptDeserialize(gConfig.Secret(), s, &w); nil != err {
		log.Printf("[ERROR] cookie decode error %v\n", err)
		c.Redirect(200, "/index.php")
		return
	}
	switch w.Step {
	case 0:
		html = w.postPage0(c)
		break
	case 1:
		html = w.postPage1(c)
		break
	case 2:
		html = w.postPage2(c)
		break
	}
	w.serviceCookieSet(c, gUtility.SerializeEncrypt(gConfig.Secret(), w))
	c.Header("Content-Type", "text/html; charset=UTF-8")
	c.String(http.StatusOK, "%v", html)
}

func main() {
	var err error
	gConfig.ReadFrom("config.toml")

	if gConfig.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	//store := memstore.NewStore()
	//store.Options(sessions.Options{
	//	Secure:   false,
	//	HttpOnly: true,
	//	SameSite: 3,
	//})
	//r.Use(sessions.Sessions(SessionName, store))
	//r.StaticFile("/favicon.ico", path.Join(gUtility.CurrentAbsolutePathOfExecutable(), "web", "favicon.ico"))
	r.StaticFile("/favicon.ico", path.Join(gUtility.CurrentAbsolutePathOfExecutable(), "favicon.ico"))
	//r.StaticFile("/css/pure.css", path.Join(gUtility.CurrentAbsolutePathOfExecutable(), "web", "css", "pure.css"))
	r.StaticFile("/css/pure.css", path.Join(gUtility.CurrentAbsolutePathOfExecutable(), "css", "pure.css"))
	r.GET("/", gWeb.IndexGet)
	r.POST("/", gWeb.IndexPost)
	r.GET("/index.php", gWeb.IndexGet)
	r.POST("/index.php", gWeb.IndexPost)

	if err = r.Run(fmt.Sprintf("%s:%d", gConfig.Host, gConfig.Port)); nil != err {
		log.Fatalf("start web server failed: %v", err)
	}
}
