# Sales Pilot：AI 智能体（OpenClaw 等）产品维护 API 指南

本文档说明如何通过 **REST JSON API** 维护**方案/产品**（含话术、案例、草稿发布），供 **OpenClaw、MCP 工具链、自动化脚本**等调用。字段语义、认证方式与端点与当前仓库代码一致。

---

## 1. 生产环境与基址

当前部署示例（HTTPS）：

```text
https://sales-pilot.guoqiangli.com
```

所有 API 挂在 **`/api/v1`** 下（经 Nginx 同域反代，无需单独端口）：

| 用途 | URL |
|------|-----|
| **API 根前缀** | `https://sales-pilot.guoqiangli.com/api/v1` |
| 管理后台（人类操作） | `https://sales-pilot.guoqiangli.com/admin` |
| 创建 API 密钥 | `https://sales-pilot.guoqiangli.com/admin/api-keys` |
| 健康检查 | `GET https://sales-pilot.guoqiangli.com/api/v1/health` |

下文用 **`{BASE}`** 表示站点根，例如 `https://sales-pilot.guoqiangli.com`，则接口路径为 **`{BASE}/api/v1/...`**。

---

## 2. 认证方式

### 2.1 JWT（用户名密码登录）

适用于临时操作或人类会话；智能体长期运行更推荐 **API Key**。

```http
POST {BASE}/api/v1/auth/login
Content-Type: application/json

{"username":"admin","password":"你的密码"}
```

响应中的 **`access_token`** 用于：

```http
Authorization: Bearer <access_token>
```

- **有效期**：默认 **24 小时**（`backend/internal/auth/jwt.go` 中 `TokenTTL`）。
- 过期后需重新 `POST /auth/login`（或使用 API Key，无需频繁登录）。

另支持 LDAP / OIDC 等，见 `GET {BASE}/api/v1/auth/login-options`；自动化维护产品以 **本地管理员账号 + API Key** 最常见。

### 2.2 API Key（推荐：OpenClaw / 脚本 / MCP）

1. 使用管理员账号登录 **`{BASE}/admin`**。
2. 打开 **「API 密钥」**（`/admin/api-keys`），新建密钥，**立即复制完整密钥**（仅创建时明文展示一次）。
3. 每次调用 **管理端** ` /api/v1/admin/*` 时在请求头携带：

```http
X-API-Key: sp_<64 位十六进制>
```

- 服务端仅存 **`SHA256(完整密钥)`**，无法找回明文。
- **`expires_at`**：创建时可选 RFC3339；省略则**永不过期**，直至在后台作废。
- 作废：`DELETE {BASE}/api/v1/admin/api-keys/:id`（或后台操作）。

**管理密钥（须已登录且为管理员）：**

| 方法 | 路径 | 请求体 | 说明 |
|------|------|--------|------|
| GET | `/api/v1/admin/api-keys` | — | 列表（不含明文密钥） |
| POST | `/api/v1/admin/api-keys` | `{"name":"openclaw-prod","expires_at":"2030-12-31T23:59:59Z"}` | `expires_at` 可选 |
| DELETE | `/api/v1/admin/api-keys/:id` | — | 软作废（`is_active=false`） |

**创建第一把 API Key** 仍须先通过 **JWT 登录**（或已有管理员会话）；之后自动化可**只使用 `X-API-Key`**。

### 2.3 权限说明

- **`/api/v1/admin/*`**：须 **`role: admin`**。使用 **JWT** 时任意有效用户令牌可访问「需登录」的只读接口，但 **管理写操作** 仍由 `RequireAdmin` 保护；使用 **`X-API-Key`** 时，密钥绑定创建者，且该用户须为**在职管理员**（被降级/禁用后密钥失效）。
- **`X-API-Key` 校验逻辑**见 `backend/internal/middleware/admin_auth.go`：密钥有效且关联用户为 `admin` + `is_active`。

---

## 3. 机器可读元数据（智能体优先拉取）

智能体启动或进入「产品维护」任务前，建议先拉取官方 JSON 说明（字段、路径、推荐流程）：

```http
GET {BASE}/api/v1/admin/meta/product-maintain
X-API-Key: <key>
# 或
Authorization: Bearer <access_token>
```

返回体包含 `product_resource.fields`、`admin_endpoints_products`、`suggested_agent_workflow`、`related_sales_script`、`related_case` 等，与本文档互补；**以线上实际 JSON 为准**时可发现新增字段。

实现位置：`backend/internal/handlers/meta_ai.go`。

---

## 4. 产品（Product）字段摘要

| JSON 字段 | 类型 | 说明 |
|-----------|------|------|
| `id` | number | 只读，创建后由系统分配 |
| `name` | string | **创建必填**，前台标题 |
| `solution_category_id` | number \| null | 关联解决方案；可先 `GET /api/v1/solution-categories` 取 `id` |
| `category` | string | 展示分类；若设置了 `solution_category_id`，保存时常由后端同步为对应方案名称 |
| `manufacturer_name` | string | 厂商/品牌 |
| `sales_contact_*` / `presales_contact_*` | string | 姓名、电话、邮箱；前台详情页电话为 `tel:`、邮箱为 `mailto:` |
| `description` | string | 概览，**Markdown（建议 GFM）** |
| `highlights` | string[] | 三大亮点 |
| `discovery_questions` | string[] | 黄金三问 |
| `target_personas` | object | 角色名 → 一句话关注点 |
| `trigger_events` | string | 触发事件，纯文本多行 |
| `competitor_analysis` | string | 竞品分析，Markdown |
| `roi_metrics` | string | ROI / 指标叙述 |
| `vendor_market` | string | **`all`** \| **`domestic`** \| **`foreign`**，与首页国内/外筛选一致 |
| `is_draft` | boolean | `true` = 仅管理端可见；`false` = 对已登录用户可见于方案库（公开列表仍不含草稿，见下节） |

**注意**：`PUT` 更新时不要发送嵌套对象 **`solution_category`**，只发送 **`solution_category_id`**（或 `null`）。

---

## 5. 管理端：产品 / 话术 / 案例 API

以下均需 **`X-API-Key`（管理员）** 或 **管理员 JWT**。

### 5.1 产品 CRUD

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/products` | 全部产品（**含草稿**），预加载 `solution_category` |
| GET | `/api/v1/admin/products/:id` | 单条（含草稿） |
| POST | `/api/v1/admin/products` | 创建，Body 为 Product JSON |
| PUT | `/api/v1/admin/products/:id` | 全量更新常用字段；Body 与创建类似，`id` 须一致 |
| PATCH | `/api/v1/admin/products/:id/draft` | 仅改草稿状态，Body：`{"is_draft": false}` |
| DELETE | `/api/v1/admin/products/:id` | 删除产品及下属话术、案例 |

### 5.2 话术（SalesScript）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/products/:id/scripts` | 该产品下话术列表 |
| POST | `/api/v1/admin/products/:id/scripts` | 创建；Body：`{"scenario":"…","content":"…"}` |
| PUT | `/api/v1/admin/scripts/:id` | 更新 `scenario`、`content` |
| DELETE | `/api/v1/admin/scripts/:id` | 删除 |

### 5.3 案例（Case）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/products/:id/cases` | 该产品下案例列表 |
| POST | `/api/v1/admin/products/:id/cases` | 创建；Body：`client_name`、`pain_points`、`solution`、`value_delivered` |
| PUT | `/api/v1/admin/cases/:id` | 更新上述字段 |
| DELETE | `/api/v1/admin/cases/:id` | 删除 |

---

## 6. 已登录只读接口（校验草稿）

以下路由在 **`/api/v1`** 下，须 **`Authorization: Bearer`** 或 **`X-API-Key`（管理员密钥）**（通过同一套 `AdminAuth` 中间件）。

- **`GET /api/v1/solution-categories`**  
  - Query：`active_only=1` 可选，仅启用项。  
  - 用于写入产品的 `solution_category_id`。

- **`GET /api/v1/products`**  
  - **仅非草稿**（`is_draft=false`）。  
  - Query：`vendor_market`（`domestic` \| `foreign` \| 空）、`solution_category_id`、`manufacturer`（精确匹配厂商名）。

- **`GET /api/v1/products/search`**  
  - Query：`q`、`solution_category_id`、`vendor_market`、`manufacturer`；无 `q` 且无解分类筛选时行为同列表。

- **`GET /api/v1/products/:id`**、**`.../scripts`**、**`.../cases`**  
  - 对**草稿**产品返回 **404**（与前台一致）。  
  - 维护草稿内容请始终使用 **`/api/v1/admin/products/...`**。

- **`GET /api/v1/topology`**、**`GET /api/v1/me`** 等：与架构展示、当前用户相关，智能体一般少用。

---

## 7. 推荐智能体工作流（OpenClaw 类）

1. **配置环境变量**（勿写入版本库）：  
   - `SALES_PILOT_BASE=https://sales-pilot.guoqiangli.com`  
   - `SALES_PILOT_API_KEY=sp_...`
2. **首轮请求**：`GET ${SALES_PILOT_BASE}/api/v1/admin/meta/product-maintain`，缓存字段说明与端点列表。
3. **发现分类**（可选）：`GET .../api/v1/solution-categories`，选定 `solution_category_id`。
4. **创建草稿**：`POST .../api/v1/admin/products`，`is_draft: true`，先填 `name`、`description` 等骨架。
5. **迭代**：`PUT .../api/v1/admin/products/:id` 补全 `highlights`、`discovery_questions`、`target_personas`、`vendor_market` 等。
6. **子资源**：`POST .../scripts`、`POST .../cases`；必要时 `PUT/DELETE` 单条话术或案例。
7. **发布**：`PATCH .../api/v1/admin/products/:id/draft`，Body：`{"is_draft": false}`。
8. **校验**：`GET .../api/v1/products` 或打开前台方案库确认可见。

**安全建议**：默认先 **`is_draft: true`**，由人类或二次检查后再发布；删除产品会级联删除话术与案例，调用前确认 `id`。

---

## 8. curl 示例（生产域名 + API Key）

```bash
export BASE="https://sales-pilot.guoqiangli.com"
export KEY='sp_你的完整密钥'

curl -sS "${BASE}/api/v1/health"

curl -sS "${BASE}/api/v1/admin/meta/product-maintain" -H "X-API-Key: ${KEY}"

curl -sS "${BASE}/api/v1/admin/products" -H "X-API-Key: ${KEY}"

curl -sS -X POST "${BASE}/api/v1/admin/products" \
  -H "X-API-Key: ${KEY}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "示例方案",
    "description": "## 概述\n适用于……",
    "highlights": ["亮点一", "亮点二", "亮点三"],
    "discovery_questions": ["问题一？", "问题二？", "问题三？"],
    "target_personas": {"CEO": "关注投入产出", "CIO": "关注架构与安全"},
    "trigger_events": "",
    "competitor_analysis": "",
    "roi_metrics": "",
    "vendor_market": "all",
    "is_draft": true
  }'
```

---

## 9. 错误码与健康检查

| HTTP | 含义 |
|------|------|
| 401 | 未携带或无效的 JWT / API Key；密钥过期、已作废 |
| 403 | 非管理员访问管理接口；或密钥关联用户非 admin |
| 400 | JSON 非法或校验失败 |
| 404 | 资源不存在；或访问了草稿的「公开」产品接口 |

**健康检查**：`GET {BASE}/api/v1/health` → `{"status":"ok"}`。

---

## 10. 与本仓库代码的对应关系

| 内容 | 路径 |
|------|------|
| 路由注册 | `backend/internal/handlers/handlers.go` |
| 管理员认证（JWT + X-API-Key） | `backend/internal/middleware/admin_auth.go` |
| API 密钥 CRUD | `backend/internal/handlers/apikeys.go` |
| 产品维护元数据 JSON | `backend/internal/handlers/meta_ai.go` |
| JWT 有效期 | `backend/internal/auth/jwt.go` → `TokenTTL` |
| Product / Script / Case 模型 | `backend/internal/models/models.go` |

---

## 11. MCP / OpenClaw 集成说明

本仓库**不包含**内置 MCP Server。请让智能体通过 **HTTPS 客户端** 调用上述 JSON 接口：

- 将 **`{BASE}`** 与 **`X-API-Key`** 配置为工具或 Agent 的 **Secret / 环境变量**。
- 工具描述中可引用本文档路径：`docs/api-ai-product-maintenance.md`（仓库内）或团队内同步的副本。
- 若部署域名变更，只需更新 **`SALES_PILOT_BASE`**，路径 **`/api/v1`** 不变。

---

*文档路径保持稳定，便于自动化与 RAG 引用：`docs/api-ai-product-maintenance.md`。*
