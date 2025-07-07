<?php
$uid = 'current_user_credential';
setcookie($uid, '', time() - 360000);
$html = '<script>window.location.href="/"</script>';
echo $html;
?>
