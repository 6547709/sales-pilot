# Phase 1: 搭建 gin-vue-admin 基础框架

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 搭建基于 gin-vue-admin 的基础框架，配置 PostgreSQL，验证权限体系（JWT + Casbin），交付可运行的基础项目。

**Architecture:** 基于 gin-vue-admin 的 server/web 双目录结构。后端采用分层架构（api/service/model/router/middleware），前端采用 Vue 3 + Element Plus + Pinia。将默认 MySQL 替换为 PostgreSQL，保留 gin-vue-admin 的 JWT + Casbin 权限模型。

**Tech Stack:** Go 1.20+, Gin 1.9+, GORM 1.25+, PostgreSQL 15+, Vue 3.3+, Element Plus 2.3+, Vite 5.x, Redis, Casbin 3.x

---

## 前置准备

### 1.1 克隆 gin-vue-admin 模板

```bash
# 在 sales-pilot 项目根目录下
git clone https://github.com/flipped-aurora/gin-vue-admin.git server_temp
mv server_temp/server .
mv server_temp/web .
rm -rf server_temp
```

### 1.2 目录结构

```
sales-pilot/
├── legacy/                    # 现有代码（仅供参考）
├── server/                    # 新后端
│   ├── api/                  # API 层
│   ├── service/              # 业务逻辑层
│   ├── model/                # 模型层
│   ├── router/               # 路由注册
│   ├── middleware/           # 中间件
│   ├── config/               # 配置
│   ├── initialize/           # 初始化
│   ├── core/                 # 核心文件
│   └── utils/               # 工具函数
├── web/                      # 新前端
│   ├── src/
│   │   ├── api/
│   │   ├── view/
│   │   ├── router/
│   │   ├── pinia/
│   │   └── components/
└── docs/
```

---

## Task 1: 后端 - 数据库配置切换到 PostgreSQL

**Files:**
- Modify: `server/config.yaml` — 数据库连接配置
- Modify: `server/go.mod` — 添加 PostgreSQL 驱动
- Modify: `server/initialize/plugin.go` — 可能的插件调整

- [ ] **Step 1: 修改 config.yaml 数据库配置**

将 `server/config.yaml` 中的数据库配置从 MySQL 改为 PostgreSQL：

```yaml
# 原有 MySQL 配置（删除）
# database:
#   type: mysql
#   qjm: root:123456@tcp(127.0.0.1:3306)/gin-vue-admin?charset=utf8mb4&parseTime=True&loc=Local

# 新增 PostgreSQL 配置
database:
  type: postgres
  postgres:
    dsn: "host=127.0.0.1 user=postgres password=postgres dbname=sales_pilot port=5432 sslmode=disable TimeZone=Asia/Shanghai"
```

- [ ] **Step 2: 修改 go.mod 添加 PostgreSQL 驱动**

```bash
cd server
go get gorm.io/driver/postgres@v1.5.7
go mod tidy
```

- [ ] **Step 3: 修改 initialize/gorm.go 适配 PostgreSQL**

找到数据库初始化代码，确保使用 `postgres.Open` 而非 `mysql.Open`：

```go
// initialize/gorm.go 或类似文件
var db *gorm.DB

func Gorm() *gorm.DB {
    // 读取 config.yaml 中的 database.postgres.dsn
    c := config.GlobalConfig
    dsn := c.Database.Postgres.Dsn
    db = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    return db
}
```

- [ ] **Step 4: 验证数据库连接**

```bash
cd server
go run main.go
# 预期：控制台输出 "Mysql 连接成功" 或类似日志（注意日志可能仍写死为 Mysql）
# 如有错误，检查 DSN 格式和 PostgreSQL 服务状态
```

- [ ] **Step 5: 提交**

```bash
git add server/config.yaml server/go.mod server/go.sum
git commit -m "feat(server): switch database from MySQL to PostgreSQL"
```

---

## Task 2: 后端 - 配置 Redis

**Files:**
- Modify: `server/config.yaml` — 添加 Redis 配置
- Modify: `server/initialize/redis.go` — Redis 初始化逻辑

- [ ] **Step 1: 添加 Redis 配置到 config.yaml**

```yaml
redis:
  db: 0
  addr: 127.0.0.1:6379
  password: ""
```

- [ ] **Step 2: 验证 Redis 连接**

 gin-vue-admin 使用 Redis 存储 JWT 令牌和登录会话。确保 Redis 服务运行：

```bash
redis-cli ping
# 预期：PONG
```

- [ ] **Step 3: 提交**

```bash
git add server/config.yaml
git commit -m "feat(server): add Redis configuration for JWT session storage"
```

---

## Task 3: 后端 - 验证权限体系（JWT + Casbin）

**Files:**
- Modify: `server/config.yaml` — JWT 密钥配置
- Test: `server/middleware/` — 验证中间件工作正常

- [ ] **Step 1: 修改 JWT 密钥**

```yaml
jwt:
  secret: "please-change-me-to-a-random-string"
```

- [ ] **Step 2: 启动服务并测试登录**

```bash
cd server
go run main.go
```

使用 curl 测试登录接口：

```bash
# 默认管理员账号: admin / 123456 (或查看初始化数据)
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456"}'
# 预期: 返回 JWT token
```

- [ ] **Step 3: 测试带 Token 的接口访问**

```bash
curl -X GET http://localhost:8080/user \
  -H "Authorization: Bearer <your-token>"
# 预期: 返回用户信息
```

- [ ] **Step 4: 提交**

```bash
git add server/config.yaml
git commit -m "feat(server): verify JWT + Casbin permission system works"
```

---

## Task 4: 前端 - 初始化 Vue 项目

**Files:**
- Modify: `web/package.json` — 项目名称修改
- Modify: `web/vite.config.js` — API 代理配置
- Modify: `web/src/core/config.js` — 后端 API 地址

- [ ] **Step 1: 修改 package.json 中的项目名称**

```json
{
  "name": "sales-pilot-web",
  "version": "0.1.0"
}
```

- [ ] **Step 2: 配置 API 代理**

```js
// web/vite.config.js
server: {
  port: 5173,
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true
    }
  }
}
```

- [ ] **Step 3: 修改后端 API 地址**

```js
// web/src/core/config.js
export const EnvBaseURL = {
  BASE_URL: '/api',  // 开发环境使用代理
}
```

- [ ] **Step 4: 安装依赖并验证启动**

```bash
cd web
npm install
npm run dev
# 预期: http://localhost:5173 可访问，显示 gin-vue-admin 登录页
```

- [ ] **Step 5: 提交**

```bash
git add web/package.json web/vite.config.js web/src/core/config.js
git commit -m "feat(web): initialize Vue project with proper API proxy config"
```

---

## Task 5: 业务模型 - 创建 PostgreSQL 数据表

**Files:**
- Create: `server/model/product.go` — 产品模型
- Create: `server/model/sales_script.go` — 销售话术模型
- Create: `server/model/case.go` — 客户案例模型
- Create: `server/model/topology.go` — 拓扑模型
- Create: `server/model/api_key.go` — API 密钥模型
- Create: `server/model/solution_category.go` — 解决方案分类模型

- [ ] **Step 1: 创建产品模型**

```go
// server/model/product.go
package model

import (
    "time"
    "gorm.io/datatypes"
)

type Product struct {
    ID                 uint           `gorm:"primaryKey" json:"id"`
    Name               string         `gorm:"size:255;not null" json:"name"`
    Category           string         `gorm:"size:128" json:"category"`
    SolutionCategoryID *uint          `gorm:"index" json:"solution_category_id,omitempty"`
    ManufacturerName   string         `gorm:"size:255" json:"manufacturer_name"`
    VendorMarket       string         `gorm:"size:16;default:all" json:"vendor_market"` // all/domestic/foreign
    IsDraft            bool           `gorm:"default:false;index" json:"is_draft"`
    Highlights         datatypes.JSON `gorm:"type:jsonb" json:"highlights"`
    TargetPersonas     datatypes.JSON `gorm:"type:jsonb" json:"target_personas"`
    TriggerEvents      string         `gorm:"type:text" json:"trigger_events"`
    DiscoveryQuestions datatypes.JSON `gorm:"type:jsonb" json:"discovery_questions"`
    CompetitorAnalysis string         `gorm:"type:text" json:"competitor_analysis"`
    ROIMetrics         string         `gorm:"type:text" json:"roi_metrics"`
    CreatedAt          time.Time      `json:"created_at"`
    UpdatedAt          time.Time      `json:"updated_at"`
}

func (Product) TableName() string {
    return "products"
}
```

- [ ] **Step 2: 创建销售话术模型**

```go
// server/model/sales_script.go
package model

import "time"

type SalesScript struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    ProductID uint      `gorm:"index;not null" json:"product_id"`
    Scenario  string    `gorm:"size:255" json:"scenario"`
    Content   string    `gorm:"type:text" json:"content"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func (SalesScript) TableName() string {
    return "sales_scripts"
}
```

- [ ] **Step 3: 创建客户案例模型**

```go
// server/model/case.go
package model

import "time"

type Case struct {
    ID             uint      `gorm:"primaryKey" json:"id"`
    ProductID      uint      `gorm:"index;not null" json:"product_id"`
    ClientName     string    `gorm:"size:255" json:"client_name"`
    PainPoints     string    `gorm:"type:text" json:"pain_points"`
    Solution       string    `gorm:"type:text" json:"solution"`
    ValueDelivered string    `gorm:"type:text" json:"value_delivered"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
}

func (Case) TableName() string {
    return "cases"
}
```

- [ ] **Step 4: 创建拓扑模型**

```go
// server/model/topology.go
package model

import "time"

type TopologyLayer struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Name      string    `gorm:"size:128;not null" json:"name"`
    SortOrder int       `gorm:"default:0" json:"sort_order"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func (TopologyLayer) TableName() string {
    return "topology_layers"
}

type TopologyCategory struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    LayerID   uint      `gorm:"index;not null" json:"layer_id"`
    Name      string    `gorm:"size:128;not null" json:"name"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func (TopologyCategory) TableName() string {
    return "topology_categories"
}

type TopologyVendor struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    CategoryID  uint      `gorm:"index;not null" json:"category_id"`
    Name        string    `gorm:"size:255;not null" json:"name"`
    Website     string    `gorm:"size:512" json:"website"`
    Description string    `gorm:"type:text" json:"description"`
    VendorMarket string   `gorm:"size:16;default:all" json:"vendor_market"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

func (TopologyVendor) TableName() string {
    return "topology_vendors"
}
```

- [ ] **Step 5: 创建 API 密钥模型**

```go
// server/model/api_key.go
package model

import "time"

type APIKey struct {
    ID          uint       `gorm:"primaryKey" json:"id"`
    Name        string     `gorm:"size:128;not null" json:"name"`
    KeyHash     string     `gorm:"size:255;not null" json:"-"`
    IsActive    bool       `gorm:"default:true" json:"is_active"`
    LastUsedAt  *time.Time `json:"last_used_at,omitempty"`
    ExpiresAt   *time.Time `json:"expires_at,omitempty"`
    CreatedBy   uint       `gorm:"index" json:"created_by"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
}

func (APIKey) TableName() string {
    return "api_keys"
}
```

- [ ] **Step 6: 创建解决方案分类模型**

```go
// server/model/solution_category.go
package model

import "time"

type SolutionCategory struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Label     string    `gorm:"size:128;not null" json:"label"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func (SolutionCategory) TableName() string {
    return "solution_categories"
}
```

- [ ] **Step 7: 自动迁移数据库表**

在 `server/initialize/db.go` 或 `main.go` 中添加自动迁移：

```go
// 在数据库连接后调用
func AutoMigrate(db *gorm.DB) {
    db.AutoMigrate(
        // gin-vue-admin 系统表...
        &model.Product{},
        &model.SalesScript{},
        &model.Case{},
        &model.TopologyLayer{},
        &model.TopologyCategory{},
        &model.TopologyVendor{},
        &model.APIKey{},
        &model.SolutionCategory{},
    )
}
```

- [ ] **Step 8: 验证表创建**

```bash
cd server
go run main.go
# 检查 PostgreSQL 中的表是否创建成功
psql -U postgres -d sales_pilot -c "\dt"
# 预期: 看到 products, sales_scripts, cases, topology_layers 等表
```

- [ ] **Step 9: 提交**

```bash
git add server/model/
git commit -m "feat(server): add business models with PostgreSQL auto-migration
- Product (产品)
- SalesScript (销售话术)
- Case (客户案例)
- TopologyLayer/Category/Vendor (拓扑管理)
- APIKey (API密钥)
- SolutionCategory (解决方案分类)"
```

---

## Task 6: Docker Compose 配置

**Files:**
- Create: `docker-compose.yml` — 开发环境 Docker 配置

- [ ] **Step 1: 创建 docker-compose.yml**

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: sales_pilot
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 5s
      retries: 5

  server:
    build:
      context: ./server
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    environment:
      - CONFIG_FILE=config.yaml
    volumes:
      - ./server:/app

  web:
    build:
      context: ./web
      dockerfile: Dockerfile
    ports:
      - "5173:5173"
    depends_on:
      - server
    volumes:
      - ./web:/app
      - /app/node_modules

volumes:
  postgres_data:
  redis_data:
```

- [ ] **Step 2: 创建 server/Dockerfile**

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/config.yaml .
EXPOSE 8080
CMD ["./server"]
```

- [ ] **Step 3: 创建 web/Dockerfile**

```dockerfile
FROM node:20-alpine AS builder

WORKDIR /app
COPY package*.json ./
RUN npm install

COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

- [ ] **Step 4: 创建 web/nginx.conf**

```nginx
server {
    listen 80;
    server_name localhost;
    root /usr/share/nginx/html;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api/ {
        proxy_pass http://server:8080/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

- [ ] **Step 5: 测试 Docker Compose**

```bash
docker compose up -d postgres redis
docker compose up -d server
# 验证服务启动
docker compose logs server | tail -20
```

- [ ] **Step 6: 提交**

```bash
git add docker-compose.yml server/Dockerfile web/Dockerfile web/nginx.conf
git commit -m "feat: add Docker Compose configuration for local development"
```

---

## Task 7: Phase 1 收尾 - 验证清单

**验证以下功能正常工作：**

- [ ] PostgreSQL 连接正常
- [ ] Redis 连接正常
- [ ] JWT 登录/登出正常
- [ ] Casbin 权限拦截正常（访问管理接口无权限报错）
- [ ] 数据库自动迁移成功（所有业务表已创建）
- [ ] 前端 Vue 项目可启动
- [ ] API 代理正常工作

---

## 交付物

Phase 1 完成后，应有以下内容：

1. **server/** 目录 — gin-vue-admin 后端框架，配置为 PostgreSQL
2. **web/** 目录 — gin-vue-admin 前端框架，可正常启动
3. **docker-compose.yml** — 一键启动开发环境
4. **数据库表** — products, sales_scripts, cases, topology_layers, topology_categories, topology_vendors, api_keys, solution_categories
5. **权限体系** — JWT + Casbin 正常工作

---

## 备注

- gin-vue-admin 默认有一些示例表（如 sys_user, sys_role, sys_menu 等），这些表保留用于权限管理
- 业务表与系统表分开，业务表使用 `products` 等命名（英文复数），系统表使用 `sys_` 前缀
- Phase 1 不实现具体业务逻辑，只确保框架跑通
