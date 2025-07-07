# 基于php:8-apache安装mysql支持。
# 由于php:8开始淘汰了mysqli，使用mysqlnd，支持PDO，但没有mysql驱动。
# 使用命令安装

# 使用官方 php:8-apache 镜像作为基础镜像
FROM php:8-apache

COPY html/* /var/www/html/
COPY init.sql /tmp
COPY start.sh /tmp

# 安装系统依赖
RUN \
    sed -i 's@deb.debian.org@mirrors.tuna.tsinghua.edu.cn@g' /etc/apt/sources.list.d/debian.sources \
    && apt-get update && apt-get install -y mariadb-server && apt-get -y clean \
    && docker-php-ext-install pdo_mysql \
    && chown -R root:root /var/www/html \
    && chmod +x /tmp/start.sh

# 暴露端口
EXPOSE 80 3306

# 创建启动脚本
# COPY start.sh /start.sh
# RUN chmod +x /start.sh


# 设置启动命令
# CMD ["sh", "-c", "service mariadb start && apache2-foreground"]
CMD ["/bin/sh", "/tmp/start.sh"]
