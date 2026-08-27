# 构建阶段
FROM node:20-alpine AS builder

WORKDIR /app

# 复制 package.json 和锁文件
COPY package.json ./
RUN npm install

# 复制源代码并构建
COPY . .
RUN npm run build

# 运行阶段
FROM node:20-alpine

WORKDIR /app

# 复制构建产物
COPY --from=builder /app/.output /app/.output
COPY --from=builder /app/node_modules /app/node_modules
COPY --from=builder /app/package.json /app/package.json

EXPOSE 3000

CMD ["node", ".output/server/index.mjs"]
