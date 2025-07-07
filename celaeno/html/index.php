<?php
//全局session_start
session_start();
//全局居设置时区
date_default_timezone_set('Asia/Shanghai');
//全局设置默认字符
header('Content-type:text/html;charset=utf-8');
//定义数据库连接参数
define('DBHOST', '127.0.0.1');//将localhost或者127.0.0.1修改为数据库服务器的地址
define('DBUSER', 'ctf_user');//将root修改为连接mysql的用户名
define('DBPASS', 'ctf_password');//将root修改为连接mysql的密码，如果改了还是连接不上，请先手动连接下你的数据库，确保数据库服务没问题在说！
define('DBNAME', 'ctf');//自定义，建议不修改
define('DBPORT', '3306');//将3306修改为mysql的连接端口，默认tcp3306

//db connect
function connect($host=DBHOST, $username=DBUSER, $password=DBPASS, $databasename=DBNAME, $port=DBPORT) {
    try {
        $link = new PDO("mysql:host=$host;port=$port;dbname=$databasename;charset=utf8mb4", $username, $password);
    } catch (PDOException $e) {
        exit("数据库连接失败: " . $e->getMessage());
    }
    $link->query("SET NAMES 'UTF8MB4'");
	return $link;
}

?>

<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>LOGIN</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f4f4f4;
            margin: 0;
            padding: 0;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
        }
        .login-container {
            background-color: #fff;
            padding: 20px;
            border-radius: 5px;
            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
        }
        .login-container h2 {
            text-align: center;
            margin-bottom: 20px;
        }
        .form-group {
            margin-bottom: 15px;
        }
        .form-group label {
            display: block;
            margin-bottom: 5px;
        }
        .form-group input {
            width: 100%;
            padding: 10px;
            border: 1px solid #ddd;
            border-radius: 4px;
        }
        .form-group button {
            width: 100%;
            padding: 10px;
            border: none;
            background-color: #5cb85c;
            color: white;
            border-radius: 4px;
            cursor: pointer;
        }
        .form-group button:hover {
            background-color: #4cae4c;
        }
    </style>
</head>
<body>

<?php
if ($_SERVER['REQUEST_METHOD'] === 'GET'):
    // connect();
/* 登录界面 */
?>

<div class="login-container">
    <h2>登录</h2>
    <form method="post">
        <div class="form-group">
            <label for="username">用户名：</label>
            <input type="text" id="username" name="username" required>
        </div>
        <div class="form-group">
            <label for="password">密码：</label>
            <input type="password" id="password" name="password" required>
        </div>
        <div class="form-group">
            <button type="submit">Login</button>
        </div>
    </form>
</div>


<?php else: 
/* 结果界面 */ 

$html = '';

if (isset($_POST['username']) && isset($_POST['password'])) {
    $username = $_POST['username'];
    $password = $_POST['password'];
    $link = connect();
    $sql = "select username from users where username='$username' and password='$password'";
    $stmt = $link->prepare($sql);
    $stmt->execute();

    // 4. 获取条目数量（方法1：使用rowCount）
    // $rowCount =$stmt->rowCount();
    // $result = $link->query($query);
    if($stmt->rowCount() === 1) {
        $html = "<h3>登录成功！欢迎你，{$username}！</h3>";
    } else {
        $html = '<h3>用户名或密码错误！请重试！</h3>';
    }
}
?>

<div class="login-container">
    <h2>登录</h2>
        <div class="form-group">
            <div><?php echo $html; ?></div>
        </div>
</div>

<?php endif; ?>
</body>
</html>

