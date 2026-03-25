# 使用外部 Nginx 终止 SSL 的说明

若 SSL 证书统一由宿主机或其它反向代理管理，可不使用本仓库 `docker-compose.yml` 中的 `nginx` 服务，改为：

1. 将 `frontend`、`backend` 端口映射到内网（或通过 Docker network 仅内部访问）。
2. 外部 Nginx 将 `location /api/` 代理到后端（例如 `http://127.0.0.1:8080/api/`），将 `/` 代理到 Next（例如 `http://127.0.0.1:3000`）。
3. 环境变量：
   - `APP_BASE_URL` 设为浏览器可访问的 **HTTPS 根地址**（不含路径），例如 `https://sales.example.com`。
   - `FRONTEND_ORIGIN` 与上述一致。
4. OIDC 在 IdP 中登记的回调地址应为：  
   `https://sales.example.com/api/v1/auth/oidc/callback`  
   并在管理后台「认证配置」中将 `redirect_url` 设为同一地址。

示例片段（仅供参考，按实际证书路径修改）：

```nginx
server {
    listen 443 ssl http2;
    server_name sales.example.com;

    ssl_certificate     /etc/ssl/fullchain.pem;
    ssl_certificate_key /etc/ssl/privkey.pem;

    location /api/ {
        proxy_pass http://127.0.0.1:8080/api/;
        proxy_set_header Host $host;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    location / {
        proxy_pass http://127.0.0.1:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Compose 中去掉 `nginx` 服务后，请为 `backend` 增加 `ports: "8080:8080"`、`frontend` 增加 `ports: "3000:3000"`（或按需绑定 127.0.0.1）。
