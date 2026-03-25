# 基于官方 Nginx，将同域反代配置打入镜像，便于 NAS 等环境仅拉取公网镜像部署（无需挂载 default.conf）
FROM nginx:1.26-bookworm
COPY default.conf /etc/nginx/conf.d/default.conf
