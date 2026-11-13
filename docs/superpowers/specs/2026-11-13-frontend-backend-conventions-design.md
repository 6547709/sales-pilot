# Sales-Pilot 前后端规范化改造设计（对齐 ALog）

- 日期：2026-11-13
- 状态：待评审
- 参考基准：`/Users/liguoqiang/project/alog`（前端 `CLAUDE.md`「目录约定」「UI 组件优先级」「契约工作流」）

## 1. 背景

Sales-Pilot 与 ALog 同栈（Next.js 16 + React 19 + Tailwind v4 + shadcn + Go/Gin/GORM），但两边的工程规范差异很大：

| 维度 | ALog | Sales-Pilot 现状 |
|------|------|------------------|
| 前端分层 | `app / features / components{ui,patterns,domain,admin,reui} / lib/api / styles` | `src/app + components{ui,admin,marketing} + 单文件 lib/api.ts` |
| 列表页 | 全部 ReUI `data-grid` | 4 个后台页共 6 处手写 `<table>`（products / users / api-keys / topology），空态/错误/加载各写一套 |
| 设计令牌 | `styles/tokens.css` 集中管理，禁止写死色值 | 集中在 `globals.css`，但页面里仍有硬编码状态色（`amber-100`、`emerald-100`） |
| 数据层 | `lib/api/{client,endpoints,adapters,schema}` | `lib/api.ts` 一个文件 248 行，DTO 与视图模型混在一起 |
| 后端分层 | `config / httpapi / store / migrations` | `handlers` 包 1.7k 行直接操作 GORM，无 service/repository |
| 后端响应 | 统一 problem+json 与错误码 | 141 处 `gin.H{"error": ...}`，无错误码；`strconv.Atoi(c.Param)` 25 处无校验 |
| 规范文档 | `CLAUDE.md` 写清目录/组件/契约约定 | `frontend/CLAUDE.md` 只有一行 `@AGENTS.md` |

结果：新页面每家一套写法，改一个状态色要翻十个文件，后端 handler 无法单测。

## 2. 目标

1. 前端建立与 ALog 同构的分层：路由只拼装、`features` 承载页面状态、`patterns` 承载通用交互、`domain` 承载业务词汇。
2. 后台所有列表统一到 ReUI `data-grid`，删除手写表格、手写空态和 `window.confirm`。
3. 全站 UI 以 ReUI 为准：容器 `Frame`、列表 `DataGrid`、筛选 `Filters`、状态 `Badge`、提示 `Alert`，不再出现「同一个界面元素两套写法」。
4. 设计令牌集中到 `styles/tokens.css`，页面不再出现硬编码状态色。
5. 前端 API 层拆成 `client / endpoints / adapters / types`，snake_case DTO 与 camelCase 视图模型分离。
6. 后端引入 `repository → service → httpapi` 三层，统一错误响应与错误码，统一分页与参数校验。
7. 两边各补一份可执行的规范文档，后续新增页面按文档走。

## 3. 非目标

- 不改信息架构与路由（`/admin/*`、`/(marketing)/*` 保持）。
- 不引入 ALog 的 openapi 代码生成、react-query、react-hook-form + zod、i18n 字典（本次不做，留作后续）。
- 不改数据库表结构与 GORM 模型字段（迁移目录只做组织，不换迁移工具）。
- 不重做前台视觉风格，只做令牌化、组件抽取与一致性修正。
- 不改鉴权语义（JWT / X-API-Key / 管理员判定逻辑保持不变）。

## 4. 约束

- **HTTP 契约保持向后兼容**：成功后响应体形状不变（数组仍是数组、对象仍是对象）；错误体在现有 `{"error": "..."}` 基础上**新增** `code` 字段。前端旧代码不会因后端重构而失效。
- `@/` 别名保持指向 `src/`。保留 `src/` 目录（ALog 的根级 `app/` 是个人习惯，不是规范核心；分层才是）。
- 现有 Next 16 与 shadcn `base-nova` 风格保持，ReUI 组件按当前 style 安装；若 registry 不支持该 style，则从 ALog 拷贝已 vendor 的 ReUI 源码并适配 import。
- 每个阶段结束必须 `npm run typecheck && npm run build`、`go build ./... && go vet ./...` 通过。

## 5. 前端改造

### 5.1 目录分层

```
frontend/src/
├── app/                         # 只做路由、取数、拼装
│   ├── (marketing)/...
│   ├── admin/...
│   ├── globals.css              # 只 import 令牌与 base 层
│   ├── error.tsx                # 全局错误边界（新增）
│   └── not-found.tsx            # 全局 404（新增）
├── features/                    # 页面级组合与状态（新增）
│   ├── products/
│   │   ├── admin-products-list.tsx
│   │   ├── product-editor.tsx
│   │   ├── scripts-editor.tsx
│   │   └── cases-editor.tsx
│   ├── topology/topology-editor.tsx
│   ├── users/users-data-grid.tsx
│   ├── api-keys/api-keys-data-grid.tsx
│   ├── backup/backup-view.tsx
│   ├── settings/auth-settings-form.tsx
│   ├── settings/system-settings-form.tsx
│   └── marketing/...
├── components/
│   ├── ui/                      # shadcn 原语，不含产品知识
│   ├── patterns/                # 通用交互模式（新增）
│   ├── domain/                  # 产品/厂商/发布状态等领域组件（新增）
│   ├── admin/                   # admin-shell、admin-panel、admin-badge、admin-alert
│   ├── marketing/               # 现有营销组件
│   └── reui/                    # ReUI 组件（新增）
├── lib/
│   ├── api/                     # client / endpoints / adapters / types / problem
│   ├── auth.ts  config.ts  utils.ts  footer-config.tsx  ...
└── styles/
    └── tokens.css               # 品牌与语义令牌（新增）
```

分层规则（写进规范文档）：

- `app/` 页面文件原则上 ≤ 60 行，只做参数解析、数据获取、拼装 `features` 组件；目前只有 `auth/callback`（40 行）达标，其余超长路由都在 5.7 清单里拆分。
- `features/` 可以知道产品或拓扑，可以持有状态。
- `components/patterns/` 不知道产品词汇，只认 `title/description/action` 这类通用 props。
- `components/domain/` 允许知道 `Product`、`TopologyCategory` 等领域类型。
- `components/ui/` 只依赖 `lib/utils`，不 import 领域类型。

### 5.2 设计令牌

新增 `src/styles/tokens.css`，`globals.css` 首行 `@import "../styles/tokens.css"`。内容：

- 品牌：`--brand: #0052cc`（保留现有企业蓝），`--brand-foreground`。
- 语义状态：`--success / --warning / --info / --destructive` 及各自 `-foreground`、`-subtle`。
- 中性色阶、圆角、阴影、侧边栏令牌。
- `@theme inline` 暴露 `--color-success` 等，使 `bg-success/10 text-success` 可用。
- 修复 `.dark` 下 `--primary` 仍是灰白的问题：暗色沿用品牌蓝的提亮版本。

约束：页面与组件不得出现 `bg-amber-100`、`text-emerald-900` 这类硬编码色；状态一律走语义令牌或 ReUI `Badge` 的 `light` 变体。

### 5.3 UI 组件：ReUI 优先，一致性硬约束

**规则一：能用 ReUI 就用 ReUI。** 顺序永远是：

1. ReUI 组件 + 免费 `c-*` 示例（先 `get_component` 读 API，再 `get_examples` 抄真实组合，最后 `validate_usage`）
2. shadcn/ui 原语
3. 自己组合（仅在 ReUI 与 shadcn 都没有时）

**规则二：一个屏幕只用一种容器族。** 全站统一 `Frame`，不出现 Card 套 Frame 或 Frame 套 Card。

**规则三：状态一律 `Badge` 的 `light` 变体**，不手写 `rounded-full bg-amber-100`。

**规则四：语义令牌，不写死色值**；间距用 `gap-*` 不用 `space-y-*`；字号不低于 `text-xs`；动效只用 `transition-colors`，每个视图最多一个入场动画，表格行不做 stagger。

#### ReUI 组件落位

| 界面 | ReUI 组件 | 说明 |
|------|-----------|------|
| 页面区块 / 表单分组 / 营销卡片 | `frame` | 全站唯一容器族 |
| 后台 4 个页面的 6 处表格（products / users / api-keys / topology） | `data-grid` | 排序、列可见性、分页、虚拟滚动内置 |
| 列表筛选（状态/市场/分类） | `filters` | 产品列表、用户列表筛选栏 |
| 状态标签（草稿/已发布/启停/国内国外） | `badge` | `light` 变体 + 语义 token |
| 成功/失败/警告提示 | `alert` | 替换手写提示条 |
| 空数据视觉 | `icon-stack` + 文案 | 配合 `patterns/empty-state` |
| 列表/卡片图标承载 | `icon-tile` | 统一图标底色与圆角 |
| 排序号等数字输入 | `number-field` | 拓扑层级/分类/厂商排序 |
| 产品分类选择（层 → 分类） | `cascader` | 替换原生多级 `<select>` |
| 销售/售前电话 | `phone-input` | 产品编辑器联系人区块 |
| 产品/厂商搜索建议 | `autocomplete` | 产品编辑器厂商与分类输入 |
| 备份/操作记录 | — | 本轮不做：后端没有审计表，非目标禁止改表结构，留待后续 |

安装：`npx shadcn@latest add @reui/{alert,badge,frame,data-grid,filters,icon-stack,icon-tile,number-field,cascader,phone-input,autocomplete} --yes`，产物落在 `components/reui/`。

- `components.json` 增加 `@reui` registry（`https://reui.io/r/{style}/{name}.json`）。
- 若 registry 不支持当前 `base-nova` style，改为从 ALog 拷贝已 vendor 的 ReUI 源码并适配 import，**不降级成手写组件**。
- ReUI 不覆盖的交互（危险操作二次确认、骨架屏组合）放 `components/patterns/`，但内部仍用 ReUI `Alert`/`Frame` 与 shadcn `Dialog` 组装。

新增依赖（与 ALog 的 ReUI data-grid 实际 import 一致，共 7 个）：`@tanstack/react-table`、`@tanstack/react-virtual`、`@dnd-kit/core`、`@dnd-kit/modifiers`、`@dnd-kit/sortable`、`@dnd-kit/utilities`、`cn`。不引入 react-query。

### 5.4 通用 patterns 组件

这些只放 ReUI 与 shadcn 都没有的组合型模式，内部一律由 ReUI 组件拼装：

| 组件 | 职责 | 内部使用 |
|------|------|----------|
| `page-header.tsx` | 标题 + 描述 + 右侧操作区 | shadcn Button |
| `empty-state.tsx` | 空数据：图标 + 标题 + 说明 + 主操作 | ReUI `IconStack` / `IconTile` |
| `error-state.tsx` | 请求失败：消息 + 重试按钮 | ReUI `Alert` |
| `loading-state.tsx` | 表格/卡片骨架屏 | shadcn Skeleton |
| `confirm-dialog.tsx` | 危险操作二次确认 | shadcn Dialog + ReUI `Alert` |
| `stat-card.tsx` | 后台首页指标卡 | ReUI `Frame` + `IconTile` |
| `form-field.tsx` | label + 说明 + 错误文案的统一表单行 | shadcn Label + Input |

### 5.5 domain 组件

| 组件 | 职责 |
|------|------|
| `product-card.tsx` | 前台产品卡片（从 `product-spotlight` 抽出） |
| `market-badge.tsx` | 国内/国外/不区分 |
| `publish-status-badge.tsx` | 草稿/已发布 |
| `vendor-tag.tsx` | 厂商标签 |
| `contact-card.tsx` | 销售/售前联系人卡片 |

### 5.6 前端 API 层

```
lib/api/
├── client.ts      # apiFetch、鉴权头、Problem 解析、401 处理
├── problem.ts     # Problem 类型与 isProblem()
├── types.ts       # 后端 DTO（snake_case，与 JSON 一一对应）
├── adapters.ts    # DTO → 视图模型（camelCase）
├── endpoints.ts   # auth / products / topology / users / apiKeys / settings / backup
└── index.ts       # 统一出口
```

- 页面只 import 视图模型，不直接读 DTO 字段。
- `lib/api.ts` 删除，调用点全部改为 `@/lib/api`。
- 保留现有内存缓存（topology/products），后续需要时再换。

### 5.7 页面改造清单

| 路由 | 现状 | 改造 |
|------|------|------|
| `/admin` | 两张静态卡片 | `stat-card` 指标区 + 快捷入口 |
| `/admin/products` | 手写 9 列表格 | `features/products/admin-products-list.tsx` + data-grid + filters + 状态 Badge + 确认对话框 |
| `/admin/products/[id]` | 996 行单文件 | 拆 `product-editor` / `scripts-editor` / `cases-editor`，表单用 Frame 分组，分类用 Cascader，电话用 PhoneInput |
| `/admin/topology` | 862 行单文件 | 拆 `features/topology/topology-editor.tsx` 与三个表单，排序用 NumberField，层/分类用 Cascader 或 Tree |
| `/admin/users` | 手写表格 | data-grid + filters + 角色 Badge + 启停操作确认 |
| `/admin/api-keys` | 手写表格 | data-grid + 过期状态 Badge |
| `/admin/backup` | 手写卡片 | Frame + 操作历史 Timeline，危险操作走 confirm-dialog |
| `/admin/settings` | 单文件 215 行 | 拆 auth-settings / system-settings 两个 feature 组件，Frame 分组 |
| `/(marketing)` 首页 | `solution-atlas` 411 行 + `product-spotlight` 263 行 | 卡片统一 ReUI `Frame`，筛选栏用 `filters`，保持 IA |
| `/(marketing)/products` | 299 行，含搜索/筛选 | 拆 `features/marketing/product-catalog.tsx`，Frame 卡片 + `filters` |
| `/(marketing)/products/[id]` | 619 行单文件 | 抽 domain 组件，Frame 容器，页面只拼装 |
| `/(marketing)/login` | 230 行 | 拆 `features/auth/login-form.tsx`，只调整 UI，登录方式与鉴权流程不变 |
| `/(marketing)/auth/callback` | 40 行 | 已符合规范，不动 |

补齐路由级 `loading.tsx`、`error.tsx`、`not-found.tsx`（`app/` 与 `admin/` 各一份）。

### 5.8 package.json

新增脚本：`"typecheck": "tsc --noEmit"`。保留 `lint`。

## 6. 后端改造

### 6.1 目标分层

```
backend/
├── cmd/server/main.go                 # 只做装配与启动
└── internal/
    ├── config/                        # 不变
    ├── database/
    │   ├── db.go                      # 连接
    │   └── migrations/migrations.go   # Migrate(db)：模型清单集中在此
    ├── models/                        # 不变
    ├── repository/                    # 新增：只做 GORM 查询
    │   ├── user_repository.go
    │   ├── apikey_repository.go
    │   ├── product_repository.go      # products / scripts / cases
    │   ├── topology_repository.go     # layers / categories / vendors
    │   ├── settings_repository.go     # auth_settings / system_settings
    │   └── backup_repository.go       # 导出/导入用全表读写
    ├── service/                       # 新增：业务规则与事务边界
    │   ├── auth_service.go
    │   ├── product_service.go
    │   ├── topology_service.go
    │   ├── user_service.go
    │   ├── apikey_service.go
    │   ├── settings_service.go
    │   └── backup_service.go
    ├── httpapi/                       # 新增：传输层
    │   ├── router.go                  # 路由注册（按资源分组）
    │   ├── response.go                # OK / Created / Fail
    │   ├── errors.go                  # 错误码 → HTTP 状态映射
    │   ├── middleware/                # 由 internal/middleware 迁入
    │   └── handlers/
    │       ├── auth_handler.go
    │       ├── product_handler.go
    │       ├── admin_product_handler.go
    │       ├── topology_handler.go
    │       ├── admin_user_handler.go
    │       ├── apikey_handler.go
    │       ├── settings_handler.go
    │       ├── backup_handler.go
    │       └── meta_handler.go
    ├── apperr/                        # 新增：领域错误
    ├── auth/                          # JWT/LDAP/OIDC 保留
    ├── cache/                         # 保留
    └── seed/                          # 保留
```

依赖方向固定：`httpapi → service → repository → models`；`repository` 不 import `service`/`httpapi`；`service` 不 import `gin`。

### 6.2 统一响应

成功保持原样。错误统一为：

```json
{ "error": "产品不存在", "code": "product_not_found" }
```

- `httpapi.Fail(c, err)`：根据 `apperr` 的 `Code` 映射状态码（`not_found` 404、`conflict` 409、`invalid` 400、`unauthorized` 401、`forbidden` 403，其余 500）。
- 500 只返回通用文案，具体错误写 `log`，不把 `err.Error()` 透给客户端（现状有泄露 SQL 错误的可能）。
- 前端 `lib/api/problem.ts` 读取 `error` 与 `code`，行为与现在一致。

### 6.3 参数与分页

**列表端点契约**（本轮只改 `/api/v1/admin/products` 与 `/api/v1/admin/users/local`）：

- **不带 `page` 参数时保持现状，返回裸数组**。`/api/v1/admin/meta/product-maintain` 与外部 `X-API-Key` 自动化依赖这个形状，不能被破坏。
- 带 `page` 参数时返回 `{ "items": [], "total": 0, "page": 1, "page_size": 20 }`。
- `page_size` 默认 20、上限 100。
- 排序：`sort` 白名单（`id` / `name` / `updated_at` / `created_at`），`order` 取 `asc` | `desc`；白名单外的 `sort` 回退默认排序，非法 `order` 回退 `asc`，均不报错。
- 筛选：`q`（name / username / email 模糊匹配）、`status`（`draft` | `published`）、`market`（`all` | `domestic` | `foreign`）、`role`（`admin` | `user`）、`active`（`true` | `false`）。未识别的值按未传处理。
- 前端 data-grid 对这两个列表开启 **manual pagination / manual sorting / manual filtering**，筛选或排序变化时回到第 1 页；`filters` 组件只暴露上述白名单维度，不把布尔查询树直接映射成请求。
- 其余列表端点本轮不动，标注为「分页待迁移」。

**参数校验**：

- 路径参数统一用 `httpapi.ParseID(c, "id")`，非法时返回 400 + `invalid_id`（消除 25 处 `strconv.Atoi` 忽略错误）。
- 请求体统一 `ShouldBindJSON` + `binding` 标签，失败返回 400 + `invalid_body`。部分更新（PATCH）继续使用指针字段，不得因为加 `binding` 把「只传一个字段」判成非法。

### 6.4 迁移目录

`internal/database/migrations/migrations.go` 集中 `Migrate(db)`，内部按 `users → topology → products → settings` 顺序调用 `AutoMigrate`。保留 GORM AutoMigrate（不引入 golang-migrate），但把模型清单与顺序从 `db.go` 抽出，后续加表只改这一处。

### 6.5 handler 拆分映射

| 现有文件 | 行数 | 去向 |
|----------|------|------|
| `handlers/handlers.go` | 770 | `router.go` + `handlers/auth_handler.go` + `handlers/product_handler.go` + `handlers/settings_handler.go` |
| `handlers/users_admin.go` | 219 | `handlers/admin_user_handler.go` + `service/user_service.go` |
| `handlers/topology.go` | 309 | `handlers/topology_handler.go` + `service/topology_service.go` |
| `handlers/apikeys.go` | 89 | `handlers/apikey_handler.go` + `service/apikey_service.go` |
| `handlers/backup.go` | 244 | `handlers/backup_handler.go` + `service/backup_service.go` |
| `handlers/meta_ai.go` | 117 | `handlers/meta_handler.go` |

### 6.6 测试与验收

- `service` 层用接口注入 repository，补单元测试覆盖：登录开关、产品草稿切换、厂商市场归一化、API Key 过期判定、备份确认令牌校验。
- `httpapi` 层用 `httptest` + service 桩补路由冒烟：401/403/404/400 各一条。
- 数据库层不做容器化集成测试（本次范围外），靠 `docker compose up` 手工回归主链路。
- 验收命令：

```bash
cd backend && go build ./... && go vet ./... && go test ./...
cd frontend && npm run typecheck && npm run lint && npm run build
```

主链路手工回归：登录 → 首页方案筛选 → 产品详情 → 后台产品列表/编辑/发布 → 拓扑维护 → 用户管理 → API Key → 备份导出。

## 7. 规范文档

- `frontend/CLAUDE.md`：重写为真实规范（目录约定、分层规则、组件优先级、令牌约束、API 层工作流、验收命令）。
- `backend/CLAUDE.md`（新增）：分层与依赖方向、错误码表、分页约定、新增端点步骤、测试要求。
- 根 `README.md` 增补「开发规范」索引，指向上述两份。

## 8. 风险与取舍

| 风险 | 处理 |
|------|------|
| 后端大范围重构引入回归 | HTTP 契约只增不改；按资源逐个迁移，每迁一个跑一次 `go build/test`；先迁 settings/apikey 这类小模块验证模式 |
| ReUI 在 `base-nova` style 下的可用性 | 已实测 `https://reui.io/r/base-nova/{data-grid,frame,badge}.json` 均返回 200；安装失败时再退回拷贝 ALog `components/reui/` 并适配 import |
| data-grid 增加依赖体积 | 只安装 data-grid 实际 import 的 7 个包；不做全局引入 |
| 页面拆分期间新旧代码并存 | 每个页面一次性切换，不在同一路由保留两套实现 |
| 前端无 e2e 测试 | 以 `typecheck + build` 为门禁，主链路人工验收；后续再补 Playwright |

## 9. 实施顺序

1. **前端地基**：tokens、ReUI 组件安装与 registry、patterns、domain、`lib/api` 拆分、typecheck 脚本。
2. **后台页面**：products → users → api-keys → settings → backup → topology → 产品编辑器。
3. **前台页面**：营销首页与产品详情令牌化与组件抽取。
4. **后端分层**：apperr/response → repository → service → httpapi 逐资源迁移。
5. **规范文档与全量验证**。
