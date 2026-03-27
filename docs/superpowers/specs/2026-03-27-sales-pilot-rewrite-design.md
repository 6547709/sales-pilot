# Sales-Pilot 重构设计文档

**日期**：2026-03-27
**项目**：基于 gin-vue-admin 重构 Sales-Pilot
**状态**：已批准

---

## 1. 项目目标

将现有 Sales-Pilot（裸代码开发）重构为基于 gin-vue-admin 架构的企业级应用，保留现有功能，提升代码可扩展性和可维护性。

### 1.1 重构原则

- **技术栈重构**：后端 Go + Gin → 沿用，前端 Next.js → Vue 3 + Element Plus
- **数据库**：PostgreSQL（保持不变）
- **认证模块**：直接采用 gin-vue-admin 成熟方案，不参考原有代码
- **保留现有代码**：移动到 `legacy/` 目录，仅供参考

---

## 2. 目录结构

```
sales-pilot/
├── legacy/                    # 现有代码（仅供参考）
│   ├── backend/              # 原 Go 后端（Handler 集中式）
│   └── frontend/             # 原 Next.js 前端
│
├── server/                    # 新后端（gin-vue-admin 分层架构）
│   ├── api/                  # API 层（路由 + 参数校验）
│   ├── service/              # 业务逻辑层
│   ├── model/                # 模型层（GORM + PostgreSQL）
│   ├── router/               # 路由注册
│   ├── middleware/           # 中间件（JWT、Casbin）
│   ├── config/               # 配置（Viper + YAML）
│   └── utils/                # 工具函数
│
├── web/                       # 新前端（Vue 3 + Element Plus）
│   ├── src/
│   │   ├── api/             # API 调用层
│   │   ├── view/            # 页面视图
│   │   ├── router/         # 路由（动态菜单）
│   │   ├── pinia/          # 状态管理
│   │   └── components/     # 组件
│   └── vite.config.js
│
├── docs/                    # 文档
├── docs/superpowers/specs/  # 设计文档
└── docker-compose.yml       # 开发/部署
```

---

## 3. 数据库设计（PostgreSQL）

### 3.1 系统表（沿用 gin-vue-admin）

| 表名 | 说明 |
|------|------|
| sys_user | 用户表 |
| sys_role | 角色表 |
| sys_menu | 菜单权限表 |
| sys_api | API 权限表 |
| sys_dictionary | 字典表 |
| casbin_rule | Casbin 策略表 |
| sys_file | 文件附件表 |

### 3.2 业务表

```sql
-- 解决方案分类
CREATE TABLE solution_categories (
    id SERIAL PRIMARY KEY,
    label VARCHAR(128) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 产品（核心业务表）
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(128),                    -- 与 solution_category.label 同步
    solution_category_id INTEGER REFERENCES solution_categories(id),
    manufacturer_name VARCHAR(255),
    vendor_market VARCHAR(16) DEFAULT 'all',  -- all/domestic/foreign
    is_draft BOOLEAN DEFAULT FALSE,
    highlights JSONB,                        -- []string，3大亮点
    target_personas JSONB,                    -- []string
    trigger_events TEXT,
    discovery_questions JSONB,                -- 黄金三问
    competitor_analysis TEXT,
    roi_metrics TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 销售话术
CREATE TABLE sales_scripts (
    id SERIAL PRIMARY KEY,
    product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
    scenario VARCHAR(255),
    content TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 客户案例
CREATE TABLE cases (
    id SERIAL PRIMARY KEY,
    product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
    client_name VARCHAR(255),
    pain_points TEXT,
    solution TEXT,
    value_delivered TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 拓扑层级
CREATE TABLE topology_layers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 拓扑分类
CREATE TABLE topology_categories (
    id SERIAL PRIMARY KEY,
    layer_id INTEGER REFERENCES topology_layers(id) ON DELETE CASCADE,
    name VARCHAR(128) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 拓扑厂商
CREATE TABLE topology_vendors (
    id SERIAL PRIMARY KEY,
    category_id INTEGER REFERENCES topology_categories(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    website VARCHAR(512),
    description TEXT,
    vendor_market VARCHAR(16) DEFAULT 'all',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- API 密钥（用于 AI/MCP 机器访问）
CREATE TABLE api_keys (
    id SERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    key_hash VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    last_used_at TIMESTAMP,
    expires_at TIMESTAMP,
    created_by INTEGER REFERENCES sys_user(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 联系人信息（嵌入到 products 中，也可独立）
-- （已在 products 表中以 sales_contact_*, presales_contact_* 字段存储）
```

---

## 4. 权限模型

采用 gin-vue-admin 的 **JWT + Casbin** 方案。

### 4.1 角色定义

| 角色 | 说明 |
|------|------|
| superadmin | 超管，可管理所有配置 |
| admin | 普通管理员，可管理业务数据 |
| user | 普通用户，仅查看和访问授权资源 |

### 4.2 权限层级

| 层级 | 实现 |
|------|------|
| **菜单权限** | sys_menu + 动态菜单渲染 |
| **API 权限** | Casbin 策略，精确到接口路径 + 方法 |
| **数据权限** | 初期不实现，后续可按部门/租户隔离 |

### 4.3 认证模块

直接使用 gin-vue-admin 自带的：
- JWT 令牌签发与验证
- Casbin 策略执行
- 登录会话管理（Redis）

**不参考原有认证代码**。

---

## 5. 业务功能模块

| 模块 | 功能 | 页面路径 |
|------|------|---------|
| **产品管理** | 列表、创建、编辑、删除、草稿控制 | /admin/products |
| **销售话术** | 场景+内容 CRUD | /admin/products/:id（Tab） |
| **客户案例** | 客户+痛点+方案+价值 CRUD | /admin/products/:id（Tab） |
| **拓扑管理** | 层级、分类、厂商管理 | /admin/topology |
| **用户管理** | 用户 CRUD、状态管理 | /admin/users |
| **角色管理** | 角色+菜单+API 权限配置 | /admin/roles |
| **认证配置** | LDAP/OIDC 开关与配置 | /admin/settings |
| **API 密钥** | 创建、撤销、有效期管理 | /admin/api-keys |
| **备份恢复** | 导入/导出 JSON 备份 | /admin/backup |
| **系统字典** | 配置项管理 | /admin/dictionary |

---

## 6. OpenClaw 集成

在 Service 层预留集成点，不在 Phase 1 实现。

```go
// server/service/product/openclaw.go
type OpenClawService struct {
    DB *gorm.DB
}

func (s *OpenClawService) SyncProduct(ctx context.Context, productID uint) error
func (s *OpenClawService) FetchUpdates(ctx context.Context) ([]ProductUpdate, error)
```

### 6.1 集成端点（预留）

| 端点 | 说明 |
|------|------|
| GET /admin/meta/product-maintain | OpenClaw 获取产品维护字段 |
| POST /admin/openclaw/sync | 触发产品同步 |

---

## 7. 开发阶段划分

| 阶段 | 内容 | 交付物 |
|------|------|--------|
| **Phase 1** | 搭建 gin-vue-admin 框架，配置 PostgreSQL，验证权限体系 | 可运行的基础项目 |
| **Phase 2** | 迁移业务模型（Product、Script、Case、Topology），完成 CRUD | 后端 API 完整 |
| **Phase 3** | 前端 Vue 管理页面重写，Element Plus 组件 | 前端页面完整 |
| **Phase 4** | API 密钥功能、OpenClaw 集成点 | 机器访问能力 |
| **Phase 5** | 测试、修复、文档 | 可发布版本 |

**总工期预估：6-8 周**

---

## 8. 技术选型

| 层级 | 技术 | 版本 |
|------|------|------|
| 后端框架 | Gin | 1.9+ |
| ORM | GORM | 1.25+ |
| 数据库 | PostgreSQL | 15+ |
| 权限 | Casbin | 3.x |
| 前端框架 | Vue 3 | 3.3+ |
| UI 库 | Element Plus | 2.3+ |
| 构建工具 | Vite | 5.x |
| 状态管理 | Pinia | 2.x |
| 配置 | Viper | 1.x |

---

## 9. 风险与注意事项

| 风险 | 缓解措施 |
|------|---------|
| gin-vue-admin 默认 MySQL | GORM 支持 PostgreSQL，需调整部分 SQL |
| Vue vs React 学习曲线 | Element Plus 生态成熟，社区资源丰富 |
| 完全重写期间无法迭代功能 | 分阶段提交代码，保持项目可用性 |
| OpenClaw API 变动 | 抽象 Service 层，便于适配 |
| 原有代码迁移遗漏功能 | 逐模块对照检查，legacy 代码随时可查 |

---

## 10. 迁移清单

从 legacy 代码迁移到新系统时，需逐项核对：

- [ ] 产品 CRUD + 草稿
- [ ] 销售话术 CRUD
- [ ] 客户案例 CRUD
- [ ] 拓扑层级/分类/厂商 CRUD
- [ ] 用户管理（本地/LDAP/OIDC）
- [ ] 认证配置页面
- [ ] API 密钥管理
- [ ] 备份导入/导出
- [ ] 首页营销内容（方案展示）
- [ ] 产品详情页
