# 构建阶段
FROM m.daocloud.io/docker.io/library/golang:1.21-alpine AS builder

WORKDIR /app

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/server ./cmd/server

# 运行阶段
FROM m.daocloud.io/docker.io/library/alpine:3.19

WORKDIR /app

# 安装 CA 证书（用于 HTTPS 出站）和时区数据
RUN apk --no-cache add ca-certificates tzdata

# 复制编译好的二进制文件
COPY --from=builder /app/bin/server /app/server

# 上传目录（外部卷挂载到此处）
RUN mkdir -p /app/uploads

EXPOSE 8080

# config.yaml 通过卷挂载到 /app/config/config.yaml（外挂，不烤进镜像）
CMD ["/app/server"]
