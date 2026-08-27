# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/server ./cmd/server

# 运行阶段
FROM alpine:3.19

WORKDIR /app

# 安装 CA 证书（用于 HTTPS）
RUN apk --no-cache add ca-certificates tzdata

# 复制编译好的二进制文件
COPY --from=builder /app/bin/server /app/server
COPY config/config.yaml /app/config/config.yaml

EXPOSE 8080

CMD ["/app/server"]
