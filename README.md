# Sales Pilot

面向销售与方案场景的 Web 应用：公开站点（方案/产品展示、登录）与管理后台（产品、用户、拓扑、备份、API 密钥等），后端为 Go + PostgreSQL，前端为 Next.js。

## 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | Next.js 16、React 19、Tailwind CSS |
| 后端 | Go 1.23、Gin、GORM、PostgreSQL |
| 编排 | Docker Compose（PostgreSQL + API + Next + Nginx 同域反代） |

## 仓库结构

```
sales-pilot/
├── backend/          # Go API（cmd/server）
├── frontend/         # Next.js 应用
├── nginx/            # 同域反代：/api → 后端，/ → 前端
├── docs/             # 补充文档
└── docker-compose.yml
```

## 本地一键启动（推荐）

需要已安装 [Docker](https://docs.docker.com/get-docker/) 与 Docker Compose。

在项目根目录执行：

```bash
docker compose up --build -d
```

浏览器访问 **http://localhost**（默认映射宿主 **80** 端口）。

常用命令：

```bash
docker compose logs -f    # 查看日志
docker compose down       # 停止（数据卷保留）
docker compose down -v    # 停止并删除数据卷（清空数据库）
```

### Compose 环境变量（可选）

| 变量 | 说明 | 默认 |
|------|------|------|
| `HTTP_PORT` | Nginx 对外端口 | `80` |
| `JWT_SECRET` | JWT 签名密钥，**生产环境务必修改** | `please-change-me` |
| `APP_BASE_URL` | 对外基址（OIDC 回调等） | `http://localhost` |
| `FRONTEND_ORIGIN` | CORS 允许的前端来源 | `http://localhost` |

示例：

```bash
JWT_SECRET="$(openssl rand -hex 32)" HTTP_PORT=8080 docker compose up --build -d
```

前端在镜像构建时使用相对路径请求 API（`NEXT_PUBLIC_API_URL` 为空），与 Nginx 同域部署一致。

## 本地开发（前后端分离）

1. 准备 PostgreSQL，连接串与后端一致（见 `backend/internal/config/config.go` 中 `DATABASE_URL` 默认值），或自行导出环境变量。
2. **注意**：默认 `docker-compose.yml` 中的 `db` 服务未将 `5432` 映射到宿主机；若只用 Compose 起数据库、在本机跑 `go run`，需为 `db` 增加 `ports: ["5432:5432"]`，否则本机进程无法通过 `localhost:5432` 连接容器内数据库。
3. 后端：`cd backend && go run ./cmd/server`（默认 `:8080`）。
4. 前端：`cd frontend && npm install && npm run dev`（默认 [http://localhost:3000](http://localhost:3000)）。未设置 `NEXT_PUBLIC_API_URL` 时，前端会请求 `http://localhost:8080`。

敏感配置请使用根目录 `.env` / `frontend/.env.local` 等（已列入 `.gitignore`），勿提交密钥。

## 相关文档

- [AI 产品维护 API 说明](docs/api-ai-product-maintenance.md)
- [Nginx 外部部署说明](docs/nginx-external.md)

## 上游仓库

<https://github.com/6547709/sales-pilot>
