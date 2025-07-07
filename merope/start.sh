#!/bin/sh
chown root:root /etc/vsftpd.conf
useradd -m $CTF_USER 
echo "$CTF_USER:$CTF_PASSWORD"|chpasswd
rm -rf /home/$CTF_USER/*
rm -rf /home/$CTF_USER/.*
echo $FLAG > /home/$CTF_USER/flag.txt
/etc/init.d/vsftpd start
rm -rf /tmp/start.sh
tail -f /dev/null