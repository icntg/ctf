<?php
//全局session_start
session_start();
//全局居设置时区
date_default_timezone_set('Asia/Shanghai');
//全局设置默认字符
header('Content-type:text/html;charset=utf-8');
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
$html = '';
$uid = 'current_user_credential';
$attach_logout = true;
if (!isset($_COOKIE[$uid])) {
    $html = '还没有登录，请点击<a href="login.php">这里</a>登录';
    $attach_logout = false;
} else if ($_COOKIE[$uid] === 'T86') {
    $flag = getenv('CTF_FLAG');
    $html = "欢迎使用教学管理系统（教师版）！$flag";
} else if (preg_match('/^T\d{2}$/', $_COOKIE[$uid])) {
            // 无flag教师
    $html = '欢迎使用教学管理系统（教师版）！';
} else if (preg_match('/^S\d{2}$/', $_COOKIE[$uid])) {
    // 学生
    $html = '欢迎使用教学管理系统（学生版）！';
} else {
    $html = '凭证信息错误！请点击<a href="logout.php">这里</a>重新登录！';
    $attach_logout = false;
}

$html = '<h2>' . $html . '</h2>';
if ($attach_logout) {
    $html = $html . '<div><a href="logout.php">退出</a></div>';
}
?>

<div class="login-container">
<?php echo $html; ?>
</div>

</body>
</html>

