# 构建阶段
# 本地默认使用 DaoCloud 镜像源（阿里云加速器失效场景），CI 通过 build-args 覆盖为官方源
ARG BASE_NODE=m.daocloud.io/docker.io/library/node:20-alpine
FROM ${BASE_NODE} AS builder

WORKDIR /app

# 复制 package.json 和锁文件
COPY package.json package-lock.json ./
RUN npm ci

# 复制源代码并构建
COPY . .
RUN npm run build

# 运行阶段：Caddy 托管 SPA 并反代 API / 上传
ARG BASE_CADDY=m.daocloud.io/docker.io/library/caddy:2-alpine
FROM ${BASE_CADDY}

WORKDIR /app

# 复制构建产物
COPY --from=builder /app/dist /usr/share/caddy

# Caddyfile：kuonji.xyz 域名、SPA fallback、/api 与 /uploads 反代 backend:8080
# 该文件位于 frontend/ 构建上下文中（docker/frontend/Caddyfile 无法被 COPY 到）
COPY Caddyfile /etc/caddy/Caddyfile

EXPOSE 80 443

CMD ["caddy", "run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"]
