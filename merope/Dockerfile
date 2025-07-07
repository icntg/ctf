# 基础镜像使用Debian
FROM debian:12

# 维护者信息
LABEL maintainer="yourname@example.com"


COPY ./start.sh /tmp/start.sh

# 安装vsftpd和必要依赖
RUN \
    sed -i 's@deb.debian.org@mirrors.tuna.tsinghua.edu.cn@g' /etc/apt/sources.list.d/debian.sources \
    && apt-get update && apt-get install -y vsftpd net-tools && apt-get -y clean && chmod +x /tmp/start.sh

# 配置vsftpd
COPY ./vsftpd.conf /etc/vsftpd.conf


# 暴露FTP端口
# EXPOSE 20 21 990 40000-40100

# 启动vsftpd
# CMD ["vsftpd", "/etc/vsftpd.conf"]
CMD ["/tmp/start.sh"]