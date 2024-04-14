#!/bin/sh
/etc/init.d/ssh start
/etc/init.d/mariadb start
/etc/init.d/apache2 start

# set flags
sed -i "s@flag0000@${flag0000}@g" /var/www/html/robots.txt
sed -i "s@flag0001@${flag0001}@g" /var/www/html/admin/index.php
sed -i "s@flag0002@${flag0002}@g" /var/www/html/admin/welcome.php
mysql -u root -e "use admin;drop table if exists flag;create table flag(flag VARCHAR(100));insert into flag VALUES('flag{${flag0003}}');"
echo "<?php /* flag{${flag0004}} */ ?>" > /var/www/flag.php
echo "flag{${flag0005}}" > /var/www/laravel/flag.txt
echo "flag{${flag0006}}" > /home/user/user.txt
chown user:user /home/user/user.txt
chmod 400 /home/user/user.txt
echo "flag{${flag0007}}" > /root/root.txt
chmod 400 /root/root.txt

unset flag0000
unset flag0001
unset flag0002
unset flag0003
unset flag0004
unset flag0005
unset flag0006
unset flag0007

rm -rf /init.sh
tail -f /var/log/apache2/access.log
