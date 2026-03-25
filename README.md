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
├── .github/workflows/ # CI：推镜像到 GHCR（推荐，无需本地构建）
├── backend/           # Go API（cmd/server）
├── frontend/          # Next.js 应用
├── nginx/             # 同域反代配置 + Dockerfile（CI 构建 sales-pilot-nginx）
├── docs/              # 补充文档
├── build-and-push.sh  # 可选：本地多架构推 Docker Hub
├── docker-compose.yml     # 本地开发：含 build
├── docker-compose.nas.yml # NAS/生产：仅拉取 GHCR 公网镜像
└── .env.example           # 环境变量模板（复制为 .env）
```

### GitHub Actions → GHCR（推荐）

推送到 **`main`** 时自动构建 **linux/amd64、linux/arm64** 并推送到 **GitHub Container Registry**（`ghcr.io`），**无需在本地执行 `docker build`**。

| 镜像 | 示例地址（将 `OWNER` 换为你的 GitHub 用户名） |
|------|-----------------------------------------------|
| 后端 | `ghcr.io/OWNER/sales-pilot-backend` |
| 前端 | `ghcr.io/OWNER/sales-pilot-frontend` |
| Nginx（内嵌 `default.conf`） | `ghcr.io/OWNER/sales-pilot-nginx` |

说明：根目录 `docker-compose.yml` 里的 Nginx 曾使用**官方** `nginx:1.26-bookworm` 并**挂载**本地 `nginx/default.conf`，因此仓库里原先没有「应用 Nginx」镜像。要在 NAS 上**只拉公网镜像、不挂载配置文件**，需使用 CI 构建的 **`sales-pilot-nginx`**（把配置打进镜像）。数据库仍使用 **Docker Hub 官方** `postgres:17-bookworm`，无需在 GHCR 再构建一份。

**版本标签（便于追踪）**：

| 标签 | 含义 |
|------|------|
| `0.1.<run_number>` | 递增序号：本工作流「第 N 次成功运行」+1（同一推送内多个 commit 只触发一次构建） |
| `sha-<短提交>` | 对应本次构建的 Git 提交 |
| `latest` | 仅 **`main`** 分支构建会更新 |

拉取示例（需 [登录 ghcr.io](https://docs.github.com/packages/working-with-a-github-packages-registry/working-with-the-container-registry)）：

```bash
docker pull ghcr.io/<你的GitHub用户名>/sales-pilot-backend:0.1.1
docker pull ghcr.io/<你的GitHub用户名>/sales-pilot-frontend:0.1.1
docker pull ghcr.io/<你的GitHub用户名>/sales-pilot-nginx:0.1.1
```

工作流文件：`.github/workflows/docker-publish.yml`。若仓库为私有，需在 GitHub 上为 Package 设置可见性或使用带 `read:packages` 的 token 拉取。

### NAS 部署（`docker-compose.nas.yml`）

在 NAS 上准备目录，放入仓库中的 **`docker-compose.nas.yml`**，将 **`.env.example`** 复制为 **`.env`** 并按注释填写（完整说明见 `.env.example`）。

拉取私有 GHCR 包前需登录（[创建 classic PAT](https://github.com/settings/tokens)，勾选 `read:packages`）：

```bash
echo <PAT> | docker login ghcr.io -u <GitHub用户名> --password-stdin
docker compose -f docker-compose.nas.yml pull
docker compose -f docker-compose.nas.yml up -d
```

更新镜像时修改 `IMAGE_TAG` 或保持 `latest` 后再次 `pull` + `up -d`。

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

### 为何本地看不到 `docker.io/liguoqiang/...` 镜像？

- `docker compose up --build` 会给镜像打上**项目名 + 服务名**的本地标签（如 `sales-pilot-backend:latest`），**不会**自动加上 Docker Hub 的 `用户名/` 前缀。
- 推送到 Hub 需单独执行 `build-and-push.sh`（或自行 `docker buildx build -t docker.io/用户名/... --push`）。
- 多架构 **`buildx build --push`** 会把 manifest 推到远端，本地 `docker images` 往往**仍不会出现** `liguoqiang/...`（多架构清单不 `--load` 到本机）；属正常现象，可在 Hub 网页查看，或执行 `docker pull docker.io/liguoqiang/sales-pilot-backend:latest` 再查看本地镜像。

### Docker Hub：本地脚本（可选）

若未使用 GitHub Actions，可自行本地构建并推 **Docker Hub**。需已登录：`docker login docker.io`

```bash
chmod +x build-and-push.sh
./build-and-push.sh
```

默认推送 `docker.io/liguoqiang/sales-pilot-backend` 与 `sales-pilot-frontend`，标签为 `0.1.0` 与 `latest`。可通过环境变量调整：

| 变量 | 说明 | 默认 |
|------|------|------|
| `VERSION` | 版本号标签 | `0.1.0` |
| `REGISTRY` | Docker Hub **用户名**（非完整域名） | `liguoqiang` |
| `PLATFORMS` | 目标平台 | `linux/amd64,linux/arm64` |
| `NEXT_PUBLIC_API_URL` | 前端构建时 API 基址；同域反代可留空 | 空 |

示例：`VERSION=0.2.0 REGISTRY=你的用户名 ./build-and-push.sh`

首次运行会创建名为 `multiarch-builder` 的 buildx 实例；在非 arm64 机器上构建 arm64 时，需 Docker Desktop 或已安装 QEMU/binfmt（一般已就绪）。

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
