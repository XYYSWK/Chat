 # 使用 golang:alpine 作为基础镜像，并且指定一个别名为 builder，这个阶段主要用于构建应用。
FROM golang:alpine AS builder

# 设置工作目录为：
WORKDIR /app

# 为我们的镜像设置必要的环境变量
ENV Go111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

# 将当前目录下的文件复制到容器中
COPY . .

# -ldflags="-s -w" 用于减少可执行文件的大小，-o 用于指定可执行文件的生成地址（包括名字）
RUN go build -ldflags="-s -w" -o bin/chat main.go
# 使用 sed 命令替换 /etc/yum.repos.d/ 中的 apline 源地址，使用阿里云镜像加速
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
# 安装curl
RUN apk add curl
# 使用 curl 下载 golang-migrate 工具
RUN curl -O -L https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz
# 解压缩下载的 golang-migrate 工具
RUN tar xvzf migrate.linux-amd64.tar.gz


# 从 alpine 开始一个新的阶段，用于构建最小化的镜像
FROM alpine

# 设置工作目录为 /app
WORKDIR /app

# 从之前的 builder 阶段复制可执行文件 chat 到当前阶段的 /app 目录
COPY --from=builder /app/bin/chat .
# 从 builder 阶段复制 golang-migrate 工具到当前阶段的 /app/migrate 目录
COPY --from=builder /app/migrate ./migrate
# 从 builder 阶段复制数据库迁移文件到当前阶段的 /app/migration 目录
COPY --from=builder /app/dao/postgresql/migration ./migration
COPY --from=builder /app/start.sh .
COPY --from=builder /app/wait-for.sh .
# 从 builder 阶段复制配置文件到当前阶段的 /app/config 目录
COPY --from=builder /app/config/app ./config/app

RUN chmod +x wait-for.sh
RUN chmod +x start.sh

# 声明容器将会监听的端口号
EXPOSE 8888

# 将 start.sh 设置为容器启动时的入口点
# ENTRYPOINT ["/app/start.sh"]
# 指定容器启动时默认执行的命令
# CMD ["/app/chat","-path=/app/config"]

