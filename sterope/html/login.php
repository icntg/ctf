<?php
//全局session_start
session_start();
//全局居设置时区
date_default_timezone_set('Asia/Shanghai');
//全局设置默认字符
header('Content-type:text/html;charset=utf-8');

$html = '<h2>用户名或密码错误！请重试！</h2>';
{
    $uid = 'current_user_credential';
    
    if ($_SERVER['REQUEST_METHOD'] === 'POST') {
        if (isset($_POST['username']) && isset($_POST['password'])) {
            $username = $_POST['username'];
            $password = $_POST['password'];
            if ($username === 'S42' && $password === 'Tsinghua@012357') {
                $html = '<h2>登录成功！欢迎你，韩梅梅同学！点击<a href="index.php">这里</a>继续</h2>';
                setcookie($uid, $username, time() + (86400 * 30));
            } 
        }
    }
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
if ($_SERVER['REQUEST_METHOD'] === 'GET') {
    // connect();
/* 登录界面 */
?>

<div class="login-container">
    <h2>教学管理系统登录</h2>
    <form method="post">
        <div class="form-group">
            <label for="username">工号/学号：</label>
            <input type="text" id="username" name="username" required>
        </div>
        <div class="form-group">
            <label for="password">密码：</label>
            <input type="password" id="password" name="password" required>
        </div>
        <div class="form-group">
            <button type="submit">登录</button>
        </div>
        <p>教师工号T**（*表示数字），默认密码为SJTU#1234567890</p>
        <p>学生学号S**（*表示数字），默认密码为Tsinghua@012357</p>
    </form>
</div>


<?php } else if ($_SERVER['REQUEST_METHOD'] === 'POST') { 
/* 登录验证 */ 
?>

<div class="login-container">
    <?php echo $html; ?>
</div>

<?php } ?>
</body>
</html>

