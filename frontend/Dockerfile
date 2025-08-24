# 构建阶段
FROM node:18-alpine AS builder

WORKDIR /app

# 复制package文件
COPY package.json package-lock.json* ./

# 安装依赖
RUN npm ci

# 复制源码
COPY . .

# 设置环境变量（构建时）
ARG VITE_API_URL=http://localhost:58080
ENV VITE_API_URL=$VITE_API_URL

# 构建应用（跳过TypeScript检查）
RUN npm run build:prod

# 生产阶段
FROM nginx:alpine

# 复制构建的文件
COPY --from=builder /app/dist /usr/share/nginx/html

# 复制nginx配置
COPY nginx.conf /etc/nginx/conf.d/default.conf

# 暴露端口
EXPOSE 80

# 启动nginx
CMD ["nginx", "-g", "daemon off;"]