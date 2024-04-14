FROM debian:11

# user & db password: KPHZL4KFT2QME4ZXX3KNZLWJ

RUN \
    sed -i 's@deb.debian.org@mirrors.tuna.tsinghua.edu.cn@g' /etc/apt/sources.list &&\
    sed -i 's@security.debian.org@mirrors.tuna.tsinghua.edu.cn@g' /etc/apt/sources.list &&\
    apt-get update &&\
    apt-get full-upgrade -y --fix-missing &&\
    apt-get install -y php php-mysql mariadb-server openssh-server python3 net-tools vim sudo &&\ 
    apt-get autoremove -y
    

ADD ./root/var/www/html /var/www/html
ADD ./root/etc/ssh/sshd_config /etc/ssh/sshd_config
ADD ./root/etc/sudoers /etc/sudoers
ADD ./root/init.sh /init.sh
ADD ./root/init.sql /tmp/init.sql

RUN \
    useradd -m -s /bin/bash -G shadow user &&\ 
    echo user:KPHZL4KFT2QME4ZXX3KNZLWJ | chpasswd &&\
    chmod 440 /etc/shadow &&\
    chmod 700 /home/user &&\
    chown root:root /etc/ssh/sshd_config &&\ 
    chmod 644 /etc/ssh/sshd_config &&\
    chown root:root /etc/sudoers &&\ 
    chmod 440 /etc/sudoers &&\
    chown -R root:root /var/www/html &&\
    mkdir -p /var/www/html/upload && chmod 777 /var/www/html/upload && chown www-data:www-data /var/www/html/upload &&\
    chmod 700 /init.sh &&\
    /etc/init.d/mariadb start &&\
    mysql -u root < /tmp/init.sql && rm -rf /tmp/init.sql &&\
    chmod 4755 /usr/bin/find &&\
    chmod u+w /etc/shadow && sed -i 's@root:.*?:@root:$1$A5aJmH.V$Q3UfFj.GMCfN19qdotwLc/:@g' /etc/shadow && chmod u-w /etc/shadow

EXPOSE 80/tcp

ENTRYPOINT ["/init.sh"]