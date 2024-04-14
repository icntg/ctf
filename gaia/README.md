> 2022/11 浙江红队试题

1. 密码暴力破解（6位数字。6位数字以IP为随机数种子生成，保证对每位学员来说，这6位数字是固定不变的）
2. 大整数加法，超时时间5s（超时时间可在配置文件中配置）
3. MD5信息缺失，暴力破解

---

运行方法：

```bash
# 进入docker目录
docker-compose up -d

# 如果修改了flag.txt或config.toml配置文件，需要重新构建镜像
docker-compose up -d --build

```

---

如对`flag3.go`源码进行修改，需要重新编译。

编译方法：

建议使用Debian进行编译。需要安装以下工具：

+ golang 请从官网上下载最新版，apt安装的不够新
+ musl-tools 使用musl-libc，alpine只支持musl
+ upx-ucl 压缩可执行文件

```bash
export CC=musl-gcc && go build -tags netgo \
--ldflags '-s -w -linkmode external -extldflags "-static"' \
flag3.go && unset CC && upx -9f --best flag3 && mv flag3 ./web/start 
```
