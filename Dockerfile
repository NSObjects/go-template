############################
# STEP 1 构建可执行文件
############################

# 指定 GO 版本号，和 go.mod 保持一致。
ARG GO_VERSION=1.26.0

# 指定构建环境
FROM golang:$GO_VERSION-alpine AS builder

#go mod 代理服务器。
#ENV GOPROXY=https://athens.azurefd.net

# ca-certificates is required to call HTTPS endpoints.
# tzdata is required for time zone info.
RUN apk add --no-cache ca-certificates tzdata && update-ca-certificates
#RUN apk update && apk upgrade && apk add --no-cache ca-certificates git
#使用私有仓库时设置GIT
#RUN git config --global url."git@gitlab.com:".insteadOf "http://gitlab.com/"

# 复制源码并指定工作目录
WORKDIR /src/app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# 为 go build 设置环境变量:
# * CGO_ENABLED=0 表示构建一个静态链接的可执行程序
# * GOOS=linux GOARCH=amd64 表示指定linux 64位的运行环境
# * GOPROXY=https://goproxy.io 可按部署环境指定代理地址
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# 构建可执行文件
RUN go build -trimpath -ldflags="-w -s" -o /out/app ./main.go
############################
# STEP 2 构建镜像
############################

# 指定最小镜像源
FROM scratch AS final

# 设置系统语言
ENV LANG=en_US.UTF-8

# 从 builder 中导入时区信息
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# 从 builder 中导入证书
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# 将构建的可执行文件复制到新镜像中
COPY --from=builder /src/app/configs/config.toml /configs/config.toml
COPY --from=builder /out/app /app

# 端口申明
EXPOSE 9322

# 运行
ENTRYPOINT [ "/app" ]
