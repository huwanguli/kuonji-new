# 构建阶段
FROM m.daocloud.io/docker.io/library/node:20-alpine AS builder

WORKDIR /app

# 复制 package.json 和锁文件
COPY package.json package-lock.json ./
RUN npm ci

# 复制源代码并构建
COPY . .
RUN npm run build

# 运行阶段：Caddy 托管 SPA 并反代 API / 上传
FROM m.daocloud.io/docker.io/library/caddy:2-alpine

WORKDIR /app

# 复制构建产物
COPY --from=builder /app/dist /usr/share/caddy

# Caddyfile：kuonji.xyz 域名、SPA fallback、/api 与 /uploads 反代 backend:8080
# 该文件位于 frontend/ 构建上下文中（docker/frontend/Caddyfile 无法被 COPY 到）
COPY Caddyfile /etc/caddy/Caddyfile

EXPOSE 80 443

CMD ["caddy", "run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"]
