CC=musl-gcc go build -tags netgo \
--ldflags '-s -w -linkmode external -extldflags "-static"' \
flag3.go


1. 密码暴力破解（6位数字，6位数字需要根据IP生成。）
2. 大整数加法，超时时间3s
3. MD5信息缺失，暴力计算
