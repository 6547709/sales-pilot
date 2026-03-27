# Sales Pilot

面向销售与方案团队的 Web 应用：公开站点（全景图→解决方案→产品）与管理后台（产品、销售话术、客户案例、拓扑、备份等）。

## 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | Vue 3 + Element Plus + Vite（基于 gin-vue-admin） |
| 后端 | Go + Gin + GORM（基于 gin-vue-admin） |
| 数据库 | PostgreSQL 15 |
| 缓存 | Redis 7 |
| 编排 | Docker Compose |

## 仓库结构

```
sales-pilot/
├── .github/workflows/      # CI：推镜像到 GHCR
├── server/                # Go 后端（API、Service、Model、Router）
│   ├── api/               # API 路由处理层
│   ├── service/           # 业务逻辑层
│   ├── model/             # 数据模型
│   ├── router/            # 路由注册
│   ├── source/            # 数据库初始化数据
│   └── config.yaml        # 后端配置
├── web/                   # Vue 前端
│   ├── src/view/          # 页面组件
│   ├── src/api/           # API 调用
│   └── Dockerfile         # 前端镜像构建
├── nginx/                 # 生产 Nginx 反向代理配置
├── docker-compose.yml     # 本地开发（含 build）
├── docker-compose.nas.yml # NAS/生产部署（仅拉取 GHCR 镜像）
└── .env.example           # 环境变量模板
```

## 快速部署

### 方式一：本地开发（Docker 一键启动）

```bash
docker compose up --build -d
```

访问 **http://localhost:5173**

- 后端 API：http://localhost:8888
-Swagger 文档：http://localhost:8888/swagger/index.html

### 方式二：NAS / 生产部署（推荐）

1. 在 NAS 上准备目录，放入 `docker-compose.nas.yml` 和 `.env`（复制自 `.env.example`）
2. 登录 GHCR（拉取私有镜像需要）：
   ```bash
   echo <PAT> | docker login ghcr.io -u <GitHub用户名> --password-stdin
   ```
3. 拉取并启动：
   ```bash
   docker compose -f docker-compose.nas.yml pull
   docker compose -f docker-compose.nas.yml up -d
   ```

## 数据库初始化

后端启动时**自动执行**以下初始化（无需手动操作）：

1. **AutoMigrate** — 根据 Model 定义自动创建/更新表结构
2. **初始化数据** — 预置拓扑层级、分类等基础数据（仅在表为空时插入）

初始化内容：

| 数据 | 说明 |
|------|------|
| 拓扑层级 | 5 层（L5→L1：云平台层→物理环境层） |
| 拓扑分类 | 34 个（7 安全 + 7 运维 + 20 中心层级） |
| 解决方案分类 | 6 类（安全合规、云计算与基础设施等） |

如需重新初始化，删除对应表数据后重启后端即可自动重插。

## 组件说明

| 组件 | 镜像 | 说明 |
|------|------|------|
| PostgreSQL | postgres:15-alpine | 数据库，默认用户 `postgres`，密码 `postgres`，库名 `sales_pilot` |
| Redis | redis:7-alpine | Session 缓存 |
| Server | 本地 build 或 GHCR | Go API 服务，监听 8888 |
| Web | 本地 build 或 GHCR | Vue 单页应用，Nginx 80 端口 |

## 环境变量

| 变量 | 说明 | 默认 |
|------|------|------|
| `JWT_SECRET` | JWT 签名密钥，**生产必须修改** | `please-change-me` |
| `APP_BASE_URL` | 外部访问地址（OIDC 回调等） | `http://localhost` |
| `FRONTEND_ORIGIN` | CORS 允许来源 | `http://localhost` |
| `HTTP_PORT` | 宿主机 HTTP 端口 | `80` |

## 公开功能（无需登录）

| 路径 | 功能 |
|------|------|
| `/home` | 首页全景图，按架构层级浏览解决方案 |
| `/home/product/:id` | 产品详情页（含销售话术、客户案例） |

## 管理后台（需登录）

| 路径 | 功能 |
|------|------|
| `/dashboard` | 业务数据统计 |
| `/product` | 产品管理（增删改查、销售话术、客户案例） |
| `/topology` | 拓扑层级与分类管理 |
| `/solutionCategory` | 解决方案分类管理 |
| `/statistics` | 数据统计 |
| `/backup` | 数据备份导入/导出 |

**默认管理员账号：** `admin` / `123456`（建议首次登录后修改密码）

## GitHub Actions CI/CD

推送到 **`main`** 分支自动构建 **linux/amd64** 镜像推送到 GHCR。

镜像地址：
```
ghcr.io/6547709/sales-pilot-server:latest
ghcr.io/6547709/sales-pilot-web:latest
```

## 常见问题

**Q: 访问 http://localhost:5173 无响应？**
检查容器状态：`docker compose ps`，确认所有容器 running。

**Q: 后端启动失败，提示数据库连接错误？**
确认 PostgreSQL 容器已就绪（`docker compose up -d postgres` 先启动数据库），等待 `service_healthy` 后再启动 server。

**Q: 如何查看后端日志？**
```bash
docker compose logs -f server
```

**Q: 如何清空数据库重新开始？**
```bash
docker compose down -v   # 删除数据卷
docker compose up -d      # 重新创建，数据会被自动初始化
```

## 相关文档

- [AI 产品维护 API 说明](docs/superpowers/api-ai-product-maintenance.md)
