# 构建：Node Bookworm（非 Alpine）
FROM node:20-bookworm AS build
WORKDIR /app
COPY package.json package-lock.json* ./
RUN npm ci
COPY . .
# 与 Nginx 同域时使用空字符串，走相对路径 /api
ARG NEXT_PUBLIC_API_URL=
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL
RUN npm run build

# 运行：UBI9 Node（非 Alpine）
FROM registry.access.redhat.com/ubi9/nodejs-20:latest
WORKDIR /app
ENV NODE_ENV=production
ENV PORT=3000
ENV HOSTNAME=0.0.0.0
COPY --from=build /app/.next/standalone ./
COPY --from=build /app/.next/static ./.next/static
COPY --from=build /app/public ./public
USER 1001
EXPOSE 3000
CMD ["node", "server.js"]
