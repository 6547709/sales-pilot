# Sales-Pilot 双视图 / 服务 / 场景 / MCP 改造设计（v3）

- 日期：2026-11-13
- 状态：已定稿（2026-11-13 用户确认：首页保持现状、服务行在厂商之上、蓝/琥珀可见性配色）
- 关系：在 `2026-11-13-frontend-backend-conventions-design.md`（分层与规范）之上叠加业务改造；分层重构是前置条件。
- 可视化样张：`docs/superpowers/specs/mockups/scenario-design.html`

## 1. 背景

项目最初面向内部销售。现在需要同一套内容同时服务两种场景：

- **销售视图**：销售自己查资料，需要卖点、黄金三问、竞品、话术。
- **客户视图**：销售把屏幕转向客户，需要解决方案价值、产品/服务特点、行业案例；内部销售资料不能出现。

同时补两类内容：**我公司的服务**（V2V、备份、容灾编排等），以及**场景**（客户真实需求，如核心数据库备份、VMware 替代）。

## 2. 目标

1. 首页保持现有「解决方案架构图 + 产品聚焦区」的形态，两种视图共用，只切换内容。
2. 每个解决方案方框内，从「国内产品 / 国外产品」两组扩展为「国内产品 / 国外产品 / 服务」三组。
3. 客户视图下，产品与服务详情用「解决方案价值 + 特点 + 行业案例」替换内部销售内容。
4. 新增独立的「按场景看」页面：以场景为大类，聚合具体场景、服务、产品与案例。
5. 提供 MCP 服务，让 AI 工具维护产品、服务、场景、拓扑。

## 3. 非目标

- 不把场景做成首页的一个区块；场景是独立页面。
- 不改现有架构图（SolutionAtlas）的交互与视觉结构，只增加「服务」组。
- 不做客户账号体系、商机/CRM、多租户、计费。
- 不做匿名公开页；客户视图需要销售登录后展示。
- 不引入搜索引擎或向量检索。

## 4. 约束

- 两种视图使用同一套路由与页面，通过 `view=sales|customer` 切换；不做两套前端。
- 内部字段由服务端按视图裁剪，前端不承担隐藏责任。
- 新增实体沿用 GORM + AutoMigrate。
- 场景固定两级，组数上限 6。
- 复用现有 `X-API-Key` 鉴权。

## 5. 领域模型

### 5.1 Service 与归属

```go
type Service struct {
    ID           uint           `gorm:"primaryKey" json:"id"`
    Slug         string         `gorm:"uniqueIndex;size:128;not null" json:"slug"`
    Name         string         `gorm:"size:255;not null" json:"name"`
    Kind         string         `gorm:"size:32;not null;index" json:"kind"` // v2v|backup|dr|migration|upgrade|ops
    Summary      string         `gorm:"type:text" json:"summary"`
    Highlights   datatypes.JSON `gorm:"type:jsonb" json:"highlights"`    // []string 能力标签
    Deliverables datatypes.JSON `gorm:"type:jsonb" json:"deliverables"`  // []string 交付物
    ValueProp    string         `gorm:"type:text" json:"value_prop"`     // 客户视图：解决方案价值
    Features     datatypes.JSON `gorm:"type:jsonb" json:"features"`      // 客户视图：特点
    IconKey      string         `gorm:"size:64" json:"icon_key"`
    SortOrder    int            `gorm:"default:0" json:"sort_order"`
    IsActive     bool           `gorm:"default:true" json:"is_active"`
    CreatedAt    time.Time      `json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
}

// ServiceCategory 服务出现在哪些解决方案方框里
type ServiceCategory struct {
    ServiceID  uint `gorm:"primaryKey" json:"service_id"`
    CategoryID uint `gorm:"primaryKey" json:"category_id"`
}
```

**我的决定：** 服务挂在「解决方案分类」上（如虚拟化、数据保护），因此每个架构方框能直接展示三组内容。不建「产品 ↔ 服务」直连：产品详情所需服务由所属分类推导，少一张表和一套维护界面。

**我的决定：** `Service.Kind` 用固定枚举，不建表，避免第二个分类膨胀源。

### 5.2 场景

```go
type ScenarioGroup struct {
    ID          uint   `gorm:"primaryKey" json:"id"`
    Slug        string `gorm:"uniqueIndex;size:128;not null" json:"slug"`
    Name        string `gorm:"size:255;not null" json:"name"`
    Description string `gorm:"size:512" json:"description"`
    SortOrder   int    `gorm:"default:0" json:"sort_order"`
}

type Scenario struct {
    ID          uint   `gorm:"primaryKey" json:"id"`
    GroupID     uint   `gorm:"index;not null" json:"group_id"`
    Slug        string `gorm:"uniqueIndex;size:128;not null" json:"slug"`
    Name        string `gorm:"size:255;not null" json:"name"`
    Summary     string `gorm:"size:512" json:"summary"`   // 一句话：客户在什么情况下需要
    PainPoint   string `gorm:"type:text" json:"pain_point"`
    ValueProp   string `gorm:"type:text" json:"value_prop"`
    IconKey     string `gorm:"size:64" json:"icon_key"`
    SortOrder   int    `gorm:"default:0" json:"sort_order"`
    IsActive    bool   `gorm:"default:true" json:"is_active"`
}

type ScenarioProduct struct {
    ScenarioID uint `gorm:"primaryKey" json:"scenario_id"`
    ProductID  uint `gorm:"primaryKey" json:"product_id"`
}
type ScenarioService struct {
    ScenarioID uint `gorm:"primaryKey" json:"scenario_id"`
    ServiceID  uint `gorm:"primaryKey" json:"service_id"`
}
```

规则：场景必须归属于一个场景组；组数上限 6；不引入第三级。

### 5.3 客户视图新增内容字段

```go
// SolutionCategory 增加客户视图文案
CustomerValue    string         `gorm:"type:text" json:"customer_value"`
CustomerFeatures datatypes.JSON `gorm:"type:jsonb" json:"customer_features"` // []string

// Case 增加行业与展示开关
Industry string `gorm:"size:64" json:"industry"` // 金融/制造/医疗/政府/互联网…
IsPublic bool   `gorm:"default:false" json:"is_public"`
```

**行业案例的归属：** `Case.ProductID` 保持必填。产品详情展示自身 `IsPublic` 案例；服务详情展示「该服务关联分类下产品的公开案例」去重；场景详情展示「该场景关联产品的公开案例」。不为 Case 增加 service_id，避免破坏现有非空约束。

### 5.4 视图边界

| 内容 | 销售视图 | 客户视图 |
|------|----------|----------|
| 产品名称、厂商、分类、亮点 | 显示 | 显示 |
| 解决方案价值 / 特点 | 显示 | 显示 |
| 行业案例（IsPublic） | 显示 | 显示 |
| 联系人（销售/售前） | 显示 | **隐藏** |
| 黄金三问（discovery_questions） | 显示 | **隐藏** |
| 目标画像、触发事件 | 显示 | **隐藏** |
| 竞品分析、ROI 指标 | 显示 | **隐藏** |
| 销售话术 SalesScript | 显示 | **隐藏** |

销售视图 = 现状 + 新增服务/场景/客户字段（便于销售预览客户会看到什么）。客户视图 = 上述白名单。

## 6. API 设计

### 6.1 视图参数

所有目录接口支持 `?view=sales|customer`，缺省为 `sales`：

```
GET /api/v1/products?view=customer
GET /api/v1/products/:id?view=customer
GET /api/v1/services?view=customer
GET /api/v1/services/:slug?view=customer
GET /api/v1/scenarios?view=customer
GET /api/v1/scenarios/:slug?view=customer
GET /api/v1/topology?view=customer        # 分类节点内附带 services
```

同一个 handler 根据 `view` 选择 DTO 投影；`customer` 投影不包含任何内部字段。

### 6.2 首页数据

`GET /api/v1/topology?view=...` 的每个分类节点增加 `services` 数组：

```json
{
  "id": 3, "slug": "virtualization", "label": "虚拟化",
  "vendors": { "domestic": [], "foreign": [] },
  "services": [{ "id": 1, "slug": "v2v", "name": "V2V 迁移服务", "summary": "…" }]
}
```

架构方框按「国内产品 / 国外产品 / 服务」三组渲染；服务为空时该组不显示。

### 6.3 场景接口

```
GET /api/v1/scenarios                     # 场景组 + 场景（两级，含关联计数）
GET /api/v1/scenarios/:slug               # 场景详情：services + products(customer/sales) + cases + categories
```

`GET /api/v1/scenarios` 响应：

```json
{
  "groups": [
    { "id": 1, "slug": "data-protection", "name": "数据保护",
      "scenarios": [
        { "id": 11, "slug": "core-db-backup", "name": "核心数据库备份",
          "summary": "核心库备份窗口不足、恢复不可验证",
          "service_count": 3, "product_count": 2 }
      ] }
  ]
}
```

前台只展示 `product_count > 0 或 service_count > 0` 的场景。

### 6.4 管理接口

```
GET/POST        /api/v1/admin/services
PUT/DELETE      /api/v1/admin/services/:id
PUT             /api/v1/admin/services/:id/categories   # { category_ids: [] }
GET/POST        /api/v1/admin/scenario-groups
PUT/DELETE      /api/v1/admin/scenario-groups/:id
GET/POST        /api/v1/admin/scenarios
PUT/DELETE      /api/v1/admin/scenarios/:id
PUT             /api/v1/admin/scenarios/:id/links       # { product_ids: [], service_ids: [] }
PATCH           /api/v1/admin/cases/:id/publish         # { is_public, industry }
PUT             /api/v1/admin/solution-categories/:id/customer   # { value, features }
```

校验：

- 第 7 个场景组 → `409 + scenario_group_limit`
- 场景组非空时删除 → `409 + scenario_group_not_empty`
- 删除服务时连带删除 `ServiceCategory`、`ScenarioService`
- 分类不存在 → `400 + invalid_category`

## 7. 前端信息架构

### 7.1 双视图切换

- 顶栏右侧放分段控件：`销售视图 | 客户视图`，状态写入 URL `?view=customer` 并记忆到 localStorage。
- 切换只改变数据请求参数与少量布局（例如客户视图隐藏「管理入口」），页面结构不变。
- 任何页面都可通过带 `view=customer` 的链接直接以客户视图打开；销售可把链接发给客户（仍需登录）。

### 7.2 首页：保持现有三栏架构图

首页完全沿用现有布局，不新增任何区块：

```
┌ 安全体系 ───┐ ┌ 云平台层 · 应用与服务 ───────────────────┐ ┌ 运维体系 ──┐
│ 主机安全     │ │ Web/中间件  数据库  AI中间件  AI/Agent/MCP│ │ AIOps      │
│ 边界安全     │ ├ 云平台层 · 基础设施 ────────────────────┤ │ CMDB       │
│ 机密管理     │ │ 虚拟化  容器云  AI云  多云管理           │ │ 拨测       │
│ 身份认证     │ ├ 网络与连接层 ──────────────────────────┤ │ 全链路监控  │
│ 零信任       │ │ 核心网络  SDN/SD-WAN  负载均衡  DNS/IPAM  │ │ 日志分析    │
│ 数据加密     │ ├ 基础硬件层 ────────────────────────────┤ │ 自动化运维  │
│ 态势感知     │ │ 计算  SAN存储  文件/对象  GPU            │ │ 灾备       │
└─────────────┘ └────────────────────────────────────────┘ └────────────┘
```

每个分类方框内部结构：在现有厂商芯片**之上**插入一行服务。

```
[图标]
标题
服务   [V2V 迁移服务] [虚拟化升级服务]   ← 新增；为空时整行隐藏
国内   [SmartX]                        ← 现有
国外   [Redhat] [VMware]               ← 现有
```

- 不改栅格、间距、圆角、层标题、分组与现有交互。
- 服务行数据来自 `ServiceCategory`，一个分类可挂 0-N 个服务。
- 服务芯片点击进入 `/services/[slug]`；厂商芯片维持现有深链（`/products?solution_category_id=…&manufacturer=…`）。
- 首页不新增服务带、场景入口或统计块。

### 7.3 颜色与可见性（维护便捷性的核心）

全站只使用两类可见性颜色，页面、后台、文档保持一致：

| 类型 | 颜色 | 用途 |
|------|------|------|
| 客户可见 | 品牌蓝 `--primary` | 服务芯片、客户价值/特点、公开案例，以及所有客户视图字段 |
| 仅销售可见 | 琥珀 `--warning` + 锁图标 | 联系人、黄金三问、目标画像、竞品、ROI、销售话术 |

落地规则：

- **首页分类方框**：服务芯片用蓝色；产品厂商芯片保持中性灰。销售视图下若分类有内部备注，用琥珀色小标签。
- **销售视图**：内部内容统一用琥珀色标签并带「仅内部」字样；客户视图下不渲染，且服务端不下发。
- **后台表单**：字段按可见性分组。客户可见字段左侧 3px 蓝色竖条 + `Badge variant="info-light">客户可见`；内部字段左侧 3px 琥珀竖条 + `Badge variant="warning-light">仅销售</Badge>`。
- **后台列表**：服务、产品、案例列表增加「可见性」列，用蓝/琥珀徽标，一屏判断哪些内容会出现在客户面前。
- **编辑页顶部**：一个固定说明条——「蓝色 = 客户能看到；琥珀 = 仅内部」。新维护人员不需要读文档就能上手。

### 7.4 产品 / 服务详情

**客户视图（产品）**

```
产品名 · 厂商 · 国内/国外 · 所属方案
────────────────────────────────────
解决方案价值      ← SolutionCategory.CustomerValue
产品特点          ← Product.highlights + Category.CustomerFeatures
行业案例          ← 该产品 IsPublic 的 Case，按 Industry 分组
相关服务          ← 所属分类下的 Service
相关场景          ← 包含该产品的 Scenario
```

**客户视图（服务）**

```
服务名 · 服务类型
────────────────────────────────────
服务价值          ← Service.ValueProp
服务特点          ← Service.Features
交付物            ← Service.Deliverables
行业案例          ← 关联分类下产品的公开案例
适用场景          ← 关联的 Scenario
涉及产品          ← 关联分类下的产品
```

**销售视图**：在客户视图内容之上，恢复联系人、黄金三问、目标画像、竞品、ROI、销售话术，并标注「仅内部可见」。

### 7.5 「按场景看」新页面（重点）

**设计目标：** 销售和客户都能在 10 秒内判断「这里有没有我要的场景」。

**`/scenarios` 总览页**

```
┌──────────────────────────────────────────────────────────────┐
│ 按场景找方案                                                   │
│ 从客户真实需要出发，找到对应的服务、产品与案例                    │
│                                       [搜索场景或关键字______]  │
├──────────────────────────────────────────────────────────────┤
│ [全部] [数据保护] [基础架构升级] [虚拟化替代] [容灾与迁移]        │
├──────────────────────────────────────────────────────────────┤
│ 数据保护                                                       │
│ ┌────────────────┐ ┌────────────────┐ ┌────────────────┐      │
│ │ 核心数据库备份   │ │ 虚拟机备份      │ │ 文件与对象归档   │      │
│ │ 备份窗口不足…    │ │ 多平台统一备份… │ │ 长期留存与合规…  │      │
│ │ 服务 3 · 产品 2 │ │ 服务 2 · 产品 4 │ │ 服务 1 · 产品 1 │      │
│ └────────────────┘ └────────────────┘ └────────────────┘      │
│                                                                │
│ 基础架构升级                                                   │
│ ┌────────────────┐ ┌────────────────┐                          │
│ │ VMware 升级     │ │ 存储扩容        │                          │
│ └────────────────┘ └────────────────┘                          │
└──────────────────────────────────────────────────────────────┘
```

- 组标题始终显示，即使当前筛选只命中部分组——保留「有哪些大类」的全貌。
- 每组最多 6 张场景卡，超出用「查看全部」。
- 搜索命中场景名、摘要、关联产品/服务名。
- 计数让销售确认「点进去有东西」；计数为 0 的场景不展示。

**`/scenarios/[slug]` 详情页**

```
面包屑：按场景 / 数据保护 / 核心数据库备份
H1 核心数据库备份
一句话：核心库备份窗口不足，恢复结果不可验证

┌ 客户痛点 ──────────────┐ ┌ 我方价值 ──────────────┐
│ …                      │ │ …                      │
└────────────────────────┘ └────────────────────────┘

推荐服务（3）
[服务卡] [服务卡] [服务卡]

推荐产品（2）    国内 / 国外分组
[产品卡] [产品卡]

行业案例
[金融] 某银行核心库备份改造 …
[制造] 某制造企业备份窗口优化 …

相关解决方案 → 数据保护 / 虚拟化
```

**为什么这样组织：**

1. 以「场景大类」作为第一层，符合客户的语言（「我要备份数据库」），而不是厂商或产品语言。
2. 组标题 + 计数 + 搜索三个入口，分别服务「浏览」「确认」与「已知关键词」三种找法。
3. 详情页先给痛点和价值，再给服务、产品、案例，符合销售讲解顺序。
4. 与首页架构图互为补充：架构图回答「市场上有哪些方案」，场景页回答「我这个需求怎么落地」。

顶部导航顺序：`首页 / 方案库 / 按场景看`。

### 7.6 后台维护

新增导航：`服务维护`、`场景维护`。所有编辑页顶部固定一条说明：**蓝色 = 客户能看到，琥珀 = 仅内部**。

**可见性标注（产品 / 服务 / 案例编辑页通用）**

- 客户可见字段：左侧 3px 蓝色竖条 + `Badge variant="info-light">客户可见`。
- 仅销售字段：左侧 3px 琥珀竖条 + `Badge variant="warning-light">仅销售</Badge>`。
- 列表增加「可见性」列，用同一套蓝/琥珀徽标。
- 产品编辑页字段分组：客户可见 = 名称、分类、厂商、市场、描述、亮点、客户价值/特点、公开案例；仅销售 = 联系人、黄金三问、目标画像、触发事件、竞品、ROI、话术。

**服务维护 `/admin/services`**：`DataGrid` + 编辑 `Sheet`。字段：名称、slug、类型、摘要、亮点、交付物、客户价值、特点、图标、所属分类（多选）、排序、启用。服务当前全部为客户可见；后续若加内部备注，按琥珀分区。

**场景维护 `/admin/scenarios`**：

- 左栏（240px）：场景组列表，显示组名与场景数，支持拖拽排序；达到 6 个后「新建场景组」禁用。
- 右栏：当前组场景 `DataGrid`，列：场景名、摘要、关联产品数、关联服务数、状态（启用 / 停用 / 未关联）、排序、操作。
- 行操作：编辑（`Sheet`）、关联（`Dialog` 内两个 `DataGrid` 多选产品与服务）、删除（`ConfirmDialog`）。
- 场景编辑字段：名称、slug、所属组、摘要、客户痛点、我方价值、图标、排序、启用。场景内容默认客户可见。

**产品编辑**：内部字段区顶部用 `Alert variant="warning"` 标注「仅内部可见」；新增「行业案例公开」开关与行业选择；新增「关联场景」多选。

**解决方案分类编辑**：新增「客户价值」「客户特点」两个字段，均按客户可见样式标注。

## 8. MCP 服务

独立 `mcp/` 进程（Node + TypeScript + `@modelcontextprotocol/sdk`，stdio），复用 `SALES_PILOT_BASE_URL` 与 `SALES_PILOT_API_KEY`。

| 工具 | 作用 |
|------|------|
| `products_list` / `products_get` | 列出/读取产品（销售视图字段） |
| `products_create` / `products_update` | 新建/更新产品 |
| `products_set_publish` | 草稿 ↔ 发布 |
| `services_list` / `services_upsert` | 维护服务 |
| `services_link_categories` | 设置服务所属解决方案 |
| `scenarios_list` / `scenarios_upsert` | 维护场景 |
| `scenarios_link` | 设置场景关联的产品与服务 |
| `topology_list` | 读取层/分类/厂商 |
| `topology_layer_upsert` | 新建/更新分层 |
| `topology_category_upsert` | 新建/更新分类（含客户价值/特点） |
| `topology_vendor_upsert` | 新建/更新厂商 |
| `topology_delete` | 删除层/分类/厂商 |
| `backup_export` | 触发导出 |

密钥等价于管理员权限，只在可信环境使用。

## 9. 数据初始化

- AutoMigrate 新增 6 张表：`services`、`service_categories`、`scenario_groups`、`scenarios`、`scenario_products`、`scenario_services`。
- 现有表加列：`SolutionCategory.customer_value`、`customer_features`；`Case.industry`、`is_public`。
- Seed 场景组 4 个：数据保护、基础架构升级、虚拟化替代、容灾与迁移。
- Seed 场景 10 个：核心数据库备份、虚拟机备份、文件与对象归档、VMware 升级、存储扩容、操作系统升级、VMware → 国产虚拟化、VMware → 超融合、V2V 迁移、双活容灾、上云迁移（示例，可在后台调整）。
- Seed 服务 6 条：V2V 迁移、备份服务、容灾编排、虚拟化升级、上云迁移、运维托管。
- 生产环境不自动猜测场景与产品关联；`SEED_DEMO=1` 时按解决方案分类建立示例关联。
- 备份载荷 v2 覆盖 6 张新表与新增列；导入 v1 时新表保持原样。

## 10. 实施顺序

1. 后端分层重构（`2026-11-13-backend-layering.md`）。
2. 数据模型 + 视图参数：新表、新列、`view` 投影、公开/内部 DTO。
3. 服务：模型、后台维护、分类架构图第三组。
4. 场景：模型、后台维护、`/scenarios` 两页。
5. 客户视图：详情页内容替换、行业案例开关。
6. MCP 服务。
7. 验证。

原前端计划 Task 15（营销页组件化）在第 3-5 步中一并完成，不再单独按旧 IA 重写。

## 11. 风险

| 风险 | 处理 |
|------|------|
| 客户视图泄露内部字段 | 服务端投影 + 路由测试断言响应体不含 `sales_contact` / `scripts` / `competitor` |
| 架构方框信息过载 | 每组最多 4 个 chip，超出显示「+N」并链接到方案聚焦区 |
| 场景页内容空 | 计数为 0 的场景不展示；后台用「未关联」状态提示 |
| 场景分类膨胀 | 组数上限 6、场景必须归组 |
| MCP 密钥权限过大 | 文档标注可信环境；后续加只读 scope |
| 双视图实现分叉 | 同一路由同一组件，只按 `view` 切数据与少量布局 |

## 12. 验收标准

- 首页架构方框内出现「服务」组，为空时自动隐藏。
- 顶栏可在销售视图/客户视图之间切换，URL 带 `view` 参数。
- 客户视图产品详情不含联系人、黄金三问、竞品、ROI、话术；包含解决方案价值、特点、行业案例。
- `/scenarios` 能按大类浏览并以关键词搜索；`/scenarios/[slug]` 展示服务、产品、案例。
- 后台可维护服务、场景、组、关联关系与案例公开开关。
- MCP 至少 14 个工具可被调用。

## 13. 验证记录（2026-11-13）

- 后端：`go build ./... && go vet ./... && go test ./...` → 全部通过（apperr / paging / repository / service / httpapi 均有测试）。
- 前端：`npx tsc --noEmit` → 0；`npx vitest run` → 8 passed；`npm run lint` → 0 errors（22 warnings）；`npm run build` → 成功。
- MCP：`npx tsc --noEmit` → 0。
- 关键测试：客户视图投影不泄露内部字段（DTO 单测 + 路由测试）、场景组上限、备份 v1/v2 表清单、服务行挂载。
- **未执行**：`docker compose` 运行时手工回归——本机没有 Docker CLI；需在有 Docker/PostgreSQL 的环境补跑主链路。
