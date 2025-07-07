#!/bin/sh

service mariadb start
cp /tmp/init.sql /tmp/flag.sql
sed -i "s@__flag__@${CTF_FLAG}@g" /tmp/flag.sql
mysql -u root < /tmp/flag.sql

# rm -rf /tmp/init.sql
rm -rf /tmp/flag.sql
# rm -rf /tmp/start.sh

apache2-foreground