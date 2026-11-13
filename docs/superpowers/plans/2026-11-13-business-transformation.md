# Sales-Pilot 双视图 / 服务 / 场景 / MCP Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 在不改变首页架构图布局的前提下，加入「我方服务」、销售/客户双视图、独立场景页与 MCP 维护接口。

**Architecture:** 后端在 `repository → service → httpapi` 三层之上增加视图投影层（`view.Mode`），同一接口按 `?view=` 返回销售或客户 DTO；前端同一套页面按视图切换内容，内部字段由服务端裁剪。服务挂在解决方案分类上，场景为两级分类并与产品、服务关联。

**Tech Stack:** Go 1.23 / Gin / GORM / PostgreSQL；Next.js 16 / React 19 / Tailwind v4 / ReUI；Node + TypeScript MCP SDK。

**Spec:** `docs/superpowers/specs/2026-11-13-public-catalog-services-scenarios-design.md`

**前置条件：** `docs/superpowers/plans/2026-11-13-backend-layering.md` 全部完成（repository/service/httpapi 已就位）。本计划不再重复分层搭建。

## Global Constraints

- 首页保持现有三栏架构图布局，不新增区块；服务行插在厂商芯片之上，为空时整行隐藏。
- 两种视图共用同一路由与组件，通过 `?view=sales|customer` 切换，缺省 `sales`。
- 客户视图的裁剪在服务端完成；响应体不得包含联系人、黄金三问、目标画像、触发事件、竞品、ROI、话术。
- 颜色语义固定：品牌蓝 = 客户可见；琥珀 = 仅销售可见。页面、后台、文档一致。
- 场景固定两级，组数上限 6；场景必须归组。
- 新增表沿用 GORM AutoMigrate；不破坏现有表字段。
- 每个任务结束必须通过对应侧检查：
  - 后端：`cd backend && go build ./... && go vet ./... && go test ./...`
  - 前端：`cd frontend && npx tsc --noEmit && npx vitest run && npm run build`

## Review Focus

- **客户视图泄露内部字段**：任何客户视图响应中都不得出现 `sales_contact_*`、`presales_contact_*`、`discovery_questions`、`target_personas`、`trigger_events`、`competitor_analysis`、`roi_metrics`、`scripts`。
- **服务行为空/过长**：没有服务时整行不渲染；服务名过长时芯片换行而不撑破方框。
- **场景页空态与零计数**：`product_count + service_count = 0` 的场景不得出现在前台；整组为空时显示空态而不是空白。
- **可见性配色在深色模式**：蓝/琥珀在 `.dark` 下仍有足够对比度，不出现浅底浅字。
- **备份兼容**：导入 v1 备份不得清空或破坏新表；导入 v2 必须完整恢复新表。

---

### Task N1: 数据模型与迁移

**Files:**
- Create: `backend/internal/models/service.go`
- Create: `backend/internal/models/scenario.go`
- Modify: `backend/internal/models/topology.go`（SolutionCategory 增加客户字段）
- Modify: `backend/internal/models/models.go`（Case 增加 Industry / IsPublic）
- Modify: `backend/internal/database/migrations/migrations.go`

**Interfaces:**
- Produces: `models.Service`、`models.ServiceCategory`、`models.ScenarioGroup`、`models.Scenario`、`models.ScenarioProduct`、`models.ScenarioService`；`models.SolutionCategory.CustomerValue` / `CustomerFeatures`；`models.Case.Industry` / `IsPublic`。

- [ ] **Step 1: 定义六个模型**

字段与 json tag 严格按 spec 5.1、5.2 定义；`TableName()` 使用 `services`、`service_categories`、`scenario_groups`、`scenarios`、`scenario_products`、`scenario_services`。

- [ ] **Step 2: 迁移顺序**

在 `migrations.Migrate` 中按 `services → service_categories → scenario_groups → scenarios → scenario_products → scenario_services` 追加；`SolutionCategory`、`Case` 的 AutoMigrate 保持原位（GORM 只加列）。

- [ ] **Step 3: 验证**

Run: `cd backend && go build ./... && go vet ./...`
Expected: 退出码 0。

- [ ] **Step 4: Commit**

```bash
git add backend/internal/models backend/internal/database
git commit -m "feat(backend): add service and scenario models"
```

---

### Task N2: repository 层

**Files:**
- Create: `backend/internal/repository/service_repository.go`
- Create: `backend/internal/repository/scenario_repository.go`
- Test: `backend/internal/repository/scenario_repository_test.go`（纯逻辑部分，如 slug 归一化）

**Interfaces:**
- Produces:
  - `NewServiceRepository(db) *ServiceRepository`：`List(ctx, kind string, activeOnly bool) ([]models.Service, error)`、`FindByID`、`FindBySlug`、`Create`、`Save`、`Delete`、`SetCategories(ctx, serviceID uint, categoryIDs []uint) error`、`CategoryIDs(ctx, serviceID uint) ([]uint, error)`、`ListByCategoryID(ctx, categoryID uint) ([]models.Service, error)`、`ListByCategoryIDs(ctx, ids []uint) (map[uint][]models.Service, error)`
  - `NewScenarioRepository(db) *ScenarioRepository`：`Groups`、`Scenarios`、`FindBySlug`、`CountProducts(ctx) (map[uint]int64, error)`、`CountServices(ctx) (map[uint]int64, error)`、`ProductIDs`、`ServiceIDs`、`SetLinks(ctx, scenarioID uint, productIDs, serviceIDs []uint) error`、`Create`、`Save`、`Delete`、`DeleteGroup`、`CountGroups`、`CountScenariosInGroup`

- [ ] **Step 1: 实现两个仓储**

`SetCategories` / `SetLinks` 使用事务：先删后插。`ListByCategoryIDs` 一次查询后按 `category_id` 分组，避免 N+1。

- [ ] **Step 2: 写一个纯逻辑测试**

`TestSetLinksDeduplicates` 用 fake 或内存结构验证重复 ID 只保留一份（若直接依赖 GORM，则改为对参数归一化函数 `dedupeIDs` 的测试）。

- [ ] **Step 3: 验证**

Run: `cd backend && go build ./... && go vet ./... && go test ./internal/repository`
Expected: 退出码 0。

- [ ] **Step 4: Commit**

```bash
git add backend/internal/repository
git commit -m "feat(backend): add service and scenario repositories"
```

---

### Task N3: 视图投影层

**Files:**
- Create: `backend/internal/view/view.go`
- Create: `backend/internal/httpapi/dto/product.go`
- Create: `backend/internal/httpapi/dto/service.go`
- Create: `backend/internal/httpapi/dto/scenario.go`
- Create: `backend/internal/httpapi/dto/case.go`
- Test: `backend/internal/httpapi/dto/product_test.go`
- Test: `backend/internal/httpapi/dto/scenario_test.go`

**Interfaces:**
- Produces:
  - `view.Mode`（`Sales` / `Customer`）、`view.Parse(raw string) Mode`（非法值回退 `Sales`）
  - `dto.ProjectProduct(p models.Product, mode view.Mode) Product`
  - `dto.ProjectService(s models.Service, mode view.Mode) Service`
  - `dto.ProjectCase(c models.Case, mode view.Mode) Case`
  - `dto.ProjectScenario(s models.Scenario, mode view.Mode) Scenario`

- [ ] **Step 1: 写失败测试（客户视图不泄露）**

```go
func TestCustomerProductHasNoInternalFields(t *testing.T) {
    p := models.Product{
        ID: 1, Name: "A", SalesContactName: "张三", DiscoveryQuestions: datatypes.JSON(`["q"]`),
        CompetitorAnalysis: "竞品", ROIMetrics: "35%",
    }
    b, _ := json.Marshal(dto.ProjectProduct(p, view.Customer))
    s := string(b)
    for _, key := range []string{"sales_contact_name", "discovery_questions", "competitor_analysis", "roi_metrics"} {
        if strings.Contains(s, key) { t.Fatalf("customer view leaked %s: %s", key, s) }
    }
}

func TestSalesProductKeepsInternalFields(t *testing.T) {
    p := models.Product{ID: 1, Name: "A", SalesContactName: "张三"}
    b, _ := json.Marshal(dto.ProjectProduct(p, view.Sales))
    if !strings.Contains(string(b), "sales_contact_name") { t.Fatalf("sales view lost field: %s", b) }
}
```

- [ ] **Step 2: 运行确认失败**

Run: `cd backend && go test ./internal/httpapi/dto`
Expected: FAIL（包不存在）。

- [ ] **Step 3: 实现投影**

客户 DTO 白名单：`id, name, category, vendor_market, manufacturer_name, manufacturer_logo, description, highlights, solution_category_id, customer_value, customer_features, services, scenarios, cases`。销售 DTO 在此基础上追加内部字段。Case 客户投影只含 `id, client_name, industry, solution, value_delivered`。

- [ ] **Step 4: 运行确认通过**

Run: `cd backend && go test ./internal/httpapi/dto`
Expected: PASS。

- [ ] **Step 5: Commit**

```bash
git add backend/internal/view backend/internal/httpapi/dto
git commit -m "feat(backend): add sales/customer view projections"
```

---

### Task N4: service 层

**Files:**
- Create: `backend/internal/service/service_service.go`
- Create: `backend/internal/service/scenario_service.go`
- Test: `backend/internal/service/service_service_test.go`
- Test: `backend/internal/service/scenario_service_test.go`
- Modify: `backend/internal/service/topology_service.go`
- Modify: `backend/internal/service/product_service.go`

**Interfaces:**
- Produces:
  - `NewServiceService(repo ServiceRepo, topology TopologyRepo) *ServiceService`：`List(ctx, kind string) ([]models.Service, error)`、`GetBySlug`、`CategoriesOf(ctx, serviceID uint) ([]models.SolutionCategory, error)`、`Create`、`Update`、`Delete`、`SetCategories`
  - `NewScenarioService(repo ScenarioRepo, products ProductRepo, services ServiceRepo) *ScenarioService`：
    - `type GroupView struct { Group models.ScenarioGroup; Scenarios []ScenarioSummary }`
    - `type ScenarioSummary struct { Scenario models.Scenario; ProductCount, ServiceCount int64 }`
    - `type ScenarioDetail struct { Scenario models.Scenario; Products []models.Product; Services []models.Service; Cases []models.Case }`
    - `Groups(ctx) ([]GroupView, error)`、`GetBySlug(ctx, slug) (*ScenarioDetail, error)`、CRUD、`SetLinks`
  - `TopologyService.Tree` 返回的分类节点增加 `Services []models.Service` 字段。

- [ ] **Step 1: 写失败测试**

```go
func TestScenarioGroupLimit(t *testing.T) {
    svc := service.NewScenarioService(&fakeScenarioRepo{groups: 6}, nil, nil)
    _, err := svc.CreateGroup(context.Background(), "新组", "new-group")
    if apperr.CodeOf(err) != "scenario_group_limit" { t.Fatalf("code=%s", apperr.CodeOf(err)) }
}

func TestServiceListFiltersInactive(t *testing.T) {
    svc := service.NewServiceService(&fakeServiceRepo{items: []models.Service{{IsActive: false}, {IsActive: true}}}, nil)
    got, _ := svc.List(context.Background(), "")
    if len(got) != 1 { t.Fatalf("len=%d", len(got)) }
}
```

- [ ] **Step 2: 运行确认失败，然后实现**

Run: `cd backend && go test ./internal/service`
Expected: 先 FAIL，实现后 PASS。

- [ ] **Step 3: 业务规则**

- 创建第 7 个场景组 → `apperr.Conflict("scenario_group_limit", "场景组最多 6 个")`。
- 删除非空场景组 → `apperr.Conflict("scenario_group_not_empty", ...)`。
- 删除服务：先删 `ServiceCategory` 与 `ScenarioService` 关联，再删服务本体（事务）。
- `Groups` 返回全部组；计数为 0 的过滤由前端做（需要保留后台全量）。

- [ ] **Step 4: 运行确认通过并提交**

```bash
git add backend/internal/service
git commit -m "feat(backend): add service and scenario business rules"
```

---

### Task N5: HTTP 接口

**Files:**
- Create: `backend/internal/httpapi/handlers/service_handler.go`
- Create: `backend/internal/httpapi/handlers/scenario_handler.go`
- Modify: `backend/internal/httpapi/handlers/product_handler.go`
- Modify: `backend/internal/httpapi/handlers/topology_handler.go`
- Modify: `backend/internal/httpapi/handlers/admin_product_handler.go`（case publish）
- Modify: `backend/internal/httpapi/router.go`
- Test: `backend/internal/httpapi/view_test.go`

**Interfaces:**
- Consumes: N3 投影、N4 service。
- Produces 路由：

```
GET  /api/v1/services                ?view=&kind=
GET  /api/v1/services/:slug          ?view=
GET  /api/v1/scenarios               ?view=
GET  /api/v1/scenarios/:slug         ?view=
GET  /api/v1/admin/services
POST /api/v1/admin/services
PUT  /api/v1/admin/services/:id
DELETE /api/v1/admin/services/:id
PUT  /api/v1/admin/services/:id/categories
GET  /api/v1/admin/scenario-groups
POST /api/v1/admin/scenario-groups
PUT  /api/v1/admin/scenario-groups/:id
DELETE /api/v1/admin/scenario-groups/:id
GET  /api/v1/admin/scenarios
POST /api/v1/admin/scenarios
PUT  /api/v1/admin/scenarios/:id
DELETE /api/v1/admin/scenarios/:id
PUT  /api/v1/admin/scenarios/:id/links
PATCH /api/v1/admin/cases/:id/publish
PUT  /api/v1/admin/solution-categories/:id/customer
```

- [ ] **Step 1: 写路由测试**

覆盖：客户视图 `/api/v1/products/1?view=customer` 响应不含 `sales_contact_name`；`/api/v1/scenarios` 返回 `groups`；第 7 个组创建返回 409。

- [ ] **Step 2: 注册路由并实现 handler**

handler 只做：解析 `view`、绑定请求、调用 service、投影、`OK/Created`。所有写接口走 `httpapi.Fail`。

- [ ] **Step 3: 运行确认通过**

Run: `cd backend && go build ./... && go vet ./... && go test ./...`
Expected: 退出码 0。

- [ ] **Step 4: Commit**

```bash
git add backend/internal/httpapi
git commit -m "feat(backend): expose service, scenario and view-aware apis"
```

---

### Task N6: Seed 数据

**Files:**
- Create: `backend/internal/seed/services.go`
- Create: `backend/internal/seed/scenarios.go`
- Modify: `backend/internal/seed/seed.go`

**Interfaces:**
- Produces: `seed.EnsureServices(db)`、`seed.EnsureScenarios(db)`，按 slug 幂等 upsert。

- [ ] **Step 1: 实现 seed**

6 个服务、4 个场景组、10 个场景（按 spec 9）。`SEED_DEMO=1` 时按 `solution_category_id` 建立 `ScenarioProduct` 示例关联；生产环境不建立。

- [ ] **Step 2: 验证幂等**

连续启动两次 `seeder`，`services` 与 `scenarios` 行数不变。

Run: `cd backend && go build ./... && go test ./...`
Expected: 退出码 0。

- [ ] **Step 3: Commit**

```bash
git add backend/internal/seed
git commit -m "feat(backend): seed default services and scenarios"
```

---

### Task N7: 备份 v2

**Files:**
- Modify: `backend/internal/models/backup.go`
- Modify: `backend/internal/repository/backup_repository.go`
- Modify: `backend/internal/service/backup_service.go`
- Test: `backend/internal/service/backup_service_test.go`

**Interfaces:**
- Produces: `BackupPayload.Version = 2` 时包含 `services`、`service_categories`、`scenario_groups`、`scenarios`、`scenario_products`、`scenario_services`。

- [ ] **Step 1: 写失败测试**

```go
func TestImportV1LeavesNewTablesUntouched(t *testing.T) {
    repo := &fakeBackupRepo{}
    svc := service.NewBackupService(repo, cache.NewSearchCache(time.Second))
    _ = svc.Import(context.Background(), &models.BackupPayload{Version: 1}, "RESTORE_FULL_V1")
    if repo.truncatedNewTables { t.Fatal("v1 import must not truncate new tables") }
}
```

- [ ] **Step 2: 实现并确认通过**

导出始终 v2；导入 v1 只清空原有表；导入 v2 清空并恢复全部表。恢复后仍校验「至少一名可用本地管理员」。

Run: `cd backend && go test ./internal/service`
Expected: PASS。

- [ ] **Step 3: Commit**

```bash
git add backend/internal/models/backup.go backend/internal/repository/backup_repository.go backend/internal/service/backup_service.go
git commit -m "feat(backend): extend backup payload to v2"
```

---

### Task N8: 前端 API 与视图模式

**Files:**
- Modify: `frontend/src/lib/api/types.ts`
- Modify: `frontend/src/lib/api/adapters.ts`
- Modify: `frontend/src/lib/api/endpoints.ts`
- Modify: `frontend/src/lib/api/adapters.test.ts`
- Create: `frontend/src/lib/view-mode.tsx`
- Test: `frontend/src/lib/view-mode.test.ts`

**Interfaces:**
- Produces:
  - `type ViewMode = "sales" | "customer"`
  - `withView(path: string, view: ViewMode): string`（已有 query 用 `&` 追加）
  - `useViewMode(): { mode: ViewMode; setMode(m: ViewMode): void }`（读 URL `view`，其次 localStorage，默认 `sales`）
  - 视图模型：`Service`、`ServiceDetail`、`ScenarioGroup`、`ScenarioSummary`、`ScenarioDetail`
  - `servicesApi.list/get`、`scenariosApi.groups/get`、`adminServicesApi.*`、`adminScenariosApi.*`
  - `fetchTopology(view?: ViewMode)`、`fetchProducts(view?: ViewMode)`、`productsApi.get(id, view?)`、`productsApi.search(params, view?)`：所有目录请求接受 `view` 参数并走 `withView`。

- [ ] **Step 1: 写失败测试**

```ts
it("appends view to a path without query", () => {
  expect(withView("/api/v1/products", "customer")).toBe("/api/v1/products?view=customer");
});
it("keeps existing query params", () => {
  expect(withView("/api/v1/products?q=a", "customer")).toBe("/api/v1/products?q=a&view=customer");
});
```

- [ ] **Step 2: 运行确认失败，实现后通过**

Run: `cd frontend && npx vitest run src/lib/view-mode.test.ts`

- [ ] **Step 3: 增加场景/服务适配器测试**

```ts
it("maps scenario summaries with counts", () => {
  const g = toScenarioGroups([{ id:1, slug:"dp", name:"数据保护",
    scenarios:[{ id:11, slug:"db", name:"核心数据库备份", summary:"…", product_count:2, service_count:3 }] }]);
  expect(g[0].scenarios[0].productCount).toBe(2);
});
```

- [ ] **Step 4: 验证并提交**

Run: `cd frontend && npx tsc --noEmit && npx vitest run && npm run build`

```bash
git add frontend/src/lib
git commit -m "feat(frontend): add view mode and service/scenario api"
```

---

### Task N9: 可见性组件与配色

**Files:**
- Create: `frontend/src/components/patterns/visibility-badge.tsx`
- Create: `frontend/src/components/patterns/visibility-section.tsx`
- Create: `frontend/src/components/patterns/visibility-legend.tsx`
- Modify: `frontend/src/styles/tokens.css`（如需暗色下琥珀对比度）

**Interfaces:**
- Produces:
  - `VisibilityBadge({ kind: "customer" | "internal" })`
  - `VisibilitySection({ kind, title, children })`：左侧 3px 蓝/琥珀竖条 + 徽标
  - `VisibilityLegend()`：一行图例「蓝色 = 客户可见；琥珀 = 仅销售可见」

- [ ] **Step 1: 实现**

`customer` 用 `Badge variant="info-light">客户可见`；`internal` 用 `Badge variant="warning-light">仅销售`，并加锁图标。

- [ ] **Step 2: 深色模式检查**

在 `.dark` 下确认文字与底色对比度足够（浏览器 DevTools 取色，文本对比 ≥ 4.5:1）。

- [ ] **Step 3: 验证并提交**

Run: `cd frontend && npx tsc --noEmit && npm run build`

```bash
git add frontend/src/components/patterns frontend/src/styles
git commit -m "feat(frontend): add visibility badges and legend"
```

---

### Task N10: 首页服务行

**Files:**
- Modify: `frontend/src/features/marketing/solution-atlas.tsx`（若尚未迁入则先迁入）
- Modify: `frontend/src/components/marketing/arch-card.tsx`
- Modify: `frontend/src/features/marketing/product-spotlight.tsx`
- Modify: `frontend/src/components/marketing/product-spotlight.tsx`（删除旧文件）

**Interfaces:**
- Consumes: `Service`、`ArchCardProductChips`。
- Produces: `ArchCard` 新增 `services: { label: string; href: string }[]`；`ProductSpotlight` 同时渲染产品卡与服务卡。

- [ ] **Step 1: 扩展架构卡**

在厂商分组**之上**渲染服务行：`<div className="... border-t border-dashed">` 内放蓝色服务芯片。`services.length === 0` 时整行不渲染。芯片 `max-w-full truncate`，容器 `flex-wrap`。

- [ ] **Step 2: 数据来源**

`fetchTopology(view)` 返回的分类节点读取 `services`；`buildBucketChips` 不变。

- [ ] **Step 3: 验证**

Run: `cd frontend && npx tsc --noEmit && npm run build`
浏览器：无服务分类不显示服务行；长服务名换行不撑破卡片；375px 不横向溢出。

- [ ] **Step 4: Commit**

```bash
git add frontend/src/features/marketing frontend/src/components/marketing
git commit -m "feat(frontend): show services above vendors in solution cards"
```

---

### Task N11: 产品 / 服务详情双视图

**Files:**
- Create: `frontend/src/features/marketing/product-detail.tsx`
- Create: `frontend/src/features/marketing/service-detail.tsx`
- Create: `frontend/src/app/(marketing)/services/[slug]/page.tsx`
- Modify: `frontend/src/app/(marketing)/products/[id]/page.tsx`

**Interfaces:**
- Consumes: `productsApi.get(id, view)`、`servicesApi.get(slug, view)`、`VisibilitySection`。
- Produces: `ProductDetail({ product })`、`ServiceDetail({ service })`。

- [ ] **Step 1: 客户视图内容**

产品：解决方案价值 → 产品特点 → 行业案例（按行业分组）→ 相关服务 → 相关场景。
服务：服务价值 → 服务特点 → 交付物 → 行业案例 → 适用场景 → 涉及产品。

- [ ] **Step 2: 销售视图内容**

在客户内容之后追加琥珀色区块：联系人、黄金三问、目标画像、触发事件、竞品分析、ROI、话术，全部包在 `VisibilitySection kind="internal"`。

- [ ] **Step 3: 验证**

Run: `cd frontend && npx tsc --noEmit && npm run build`
浏览器：客户视图无内部字段；切到销售视图出现琥珀区块。

- [ ] **Step 4: Commit**

```bash
git add frontend/src/features/marketing frontend/src/app/\(marketing\)
git commit -m "feat(frontend): add customer/sales product and service detail views"
```

---

### Task N12: 场景页

**Files:**
- Create: `frontend/src/app/(marketing)/scenarios/page.tsx`
- Create: `frontend/src/app/(marketing)/scenarios/[slug]/page.tsx`
- Create: `frontend/src/features/scenarios/scenario-overview.tsx`
- Create: `frontend/src/features/scenarios/scenario-detail.tsx`
- Create: `frontend/src/components/domain/scenario-card.tsx`

**Interfaces:**
- Consumes: `scenariosApi.groups/get`、`ScenarioCard`、`PageHeader`、`EmptyState`。
- Produces: `ScenarioOverview()`、`ScenarioDetailView({ detail })`。

- [ ] **Step 1: 总览页**

组 Tab（全部 + 各组）→ 分组标题 → 场景卡网格。场景卡展示名称、摘要、`服务 N · 产品 M`。过滤规则：`productCount + serviceCount === 0` 的场景不渲染；组内全为空时整组不渲染；全部为空时 `EmptyState`。搜索框匹配名称与摘要。

- [ ] **Step 2: 详情页**

面包屑 → 标题 + 摘要 → 客户痛点 / 我方价值 → 推荐服务 → 推荐产品（国内/国外分组）→ 行业案例 → 相关解决方案。

- [ ] **Step 3: 导航入口**

在 `PublicShell` 增加「按场景看」链接（`/scenarios`）。

- [ ] **Step 4: 验证**

Run: `cd frontend && npx tsc --noEmit && npm run build`
浏览器：搜索、组筛选、空态、详情跳转均正常；375px 单列。

- [ ] **Step 5: Commit**

```bash
git add frontend/src/app/\(marketing\) frontend/src/features/scenarios frontend/src/components/domain frontend/src/components/marketing/public-shell.tsx
git commit -m "feat(frontend): add scenario overview and detail pages"
```

---

### Task N13: 后台服务维护

**Files:**
- Create: `frontend/src/app/admin/services/page.tsx`
- Create: `frontend/src/features/services/services-data-grid.tsx`
- Create: `frontend/src/features/services/service-editor.tsx`
- Modify: `frontend/src/components/admin/admin-shell.tsx`（导航加「服务维护」）

**Interfaces:**
- Consumes: `adminServicesApi.*`、`DataGrid`、`Sheet`、`VisibilitySection`。
- Produces: `ServicesDataGrid()`、`ServiceEditor({ service, onSaved })`。

- [ ] **Step 1: 实现列表与编辑**

列：名称、类型、所属分类、排序、可见性（客户可见）、启用、操作。编辑 Sheet：名称、slug、类型、摘要、亮点、交付物、客户价值、特点、图标、所属分类多选、排序、启用。

- [ ] **Step 2: 验证并提交**

Run: `cd frontend && npx tsc --noEmit && npm run build`

```bash
git add frontend/src/app/admin/services frontend/src/features/services frontend/src/components/admin
git commit -m "feat(admin): add service maintenance"
```

---

### Task N14: 后台场景维护

**Files:**
- Create: `frontend/src/app/admin/scenarios/page.tsx`
- Create: `frontend/src/features/scenarios/scenario-groups-panel.tsx`
- Create: `frontend/src/features/scenarios/scenarios-data-grid.tsx`
- Create: `frontend/src/features/scenarios/scenario-editor.tsx`
- Create: `frontend/src/features/scenarios/scenario-links-dialog.tsx`
- Modify: `frontend/src/components/admin/admin-shell.tsx`（导航加「场景维护」）

**Interfaces:**
- Consumes: `adminScenariosApi.*`、`DataGrid`、`Dialog`、`Sheet`、`ConfirmDialog`。
- Produces: 两栏页面：`ScenarioGroupsPanel` + `ScenariosDataGrid`。

- [ ] **Step 1: 左栏组管理**

组列表显示名称与场景数；达到 6 个时「新建场景组」禁用并提示上限；删除非空组显示后端 409 文案。

- [ ] **Step 2: 右栏场景管理**

列：场景名、摘要、产品数、服务数、状态（启用 `success-light` / 停用 `outline` / 未关联 `warning-light`）、操作。编辑用 `Sheet`；关联用 `Dialog` 内两个 `DataGrid` 多选。

- [ ] **Step 3: 验证并提交**

Run: `cd frontend && npx tsc --noEmit && npm run build`

```bash
git add frontend/src/app/admin/scenarios frontend/src/features/scenarios frontend/src/components/admin
git commit -m "feat(admin): add scenario maintenance"
```

---

### Task N15: 后台可见性标注

**Files:**
- Modify: `frontend/src/features/products/product-editor.tsx`
- Modify: `frontend/src/features/products/cases-editor.tsx`
- Modify: `frontend/src/features/topology/category-form.tsx`

**Interfaces:**
- Consumes: `VisibilitySection`、`VisibilityLegend`、`VisibilityBadge`。

- [ ] **Step 1: 产品编辑分区**

客户可见：名称、分类、厂商、市场、描述、亮点、客户价值/特点。仅销售：联系人、黄金三问、目标画像、触发事件、竞品、ROI、话术。顶部放 `VisibilityLegend`。

- [ ] **Step 2: 案例与分类**

案例列表增加「行业」与「客户视图公开」两列/开关；分类表单增加「客户价值」「客户特点」。

- [ ] **Step 3: 验证并提交**

Run: `cd frontend && npx tsc --noEmit && npm run build`

```bash
git add frontend/src/features/products frontend/src/features/topology
git commit -m "feat(admin): mark customer and internal fields"
```

---

### Task N16: MCP 服务

**Files:**
- Create: `mcp/package.json`
- Create: `mcp/tsconfig.json`
- Create: `mcp/src/index.ts`
- Create: `mcp/src/client.ts`
- Create: `mcp/src/tools/products.ts`
- Create: `mcp/src/tools/services.ts`
- Create: `mcp/src/tools/scenarios.ts`
- Create: `mcp/src/tools/topology.ts`
- Create: `mcp/README.md`

**Interfaces:**
- Produces: MCP server（stdio），工具：`products_list`、`products_get`、`products_create`、`products_update`、`products_set_publish`、`services_list`、`services_upsert`、`services_link_categories`、`scenarios_list`、`scenarios_upsert`、`scenarios_link`、`topology_list`、`topology_layer_upsert`、`topology_category_upsert`、`topology_vendor_upsert`、`topology_delete`、`backup_export`。

- [ ] **Step 1: 搭骨架**

`client.ts` 读取 `SALES_PILOT_BASE_URL` / `SALES_PILOT_API_KEY`，统一 `X-API-Key` 与 `Content-Type`，非 2xx 抛出后端 `error` 文案。

- [ ] **Step 2: 注册工具**

每个工具定义 `inputSchema`（zod 或 JSON schema）与 handler，返回 `content: [{ type: "text", text: JSON.stringify(result) }]`。

- [ ] **Step 3: 验证**

Run: `cd mcp && npm install && npx tsc --noEmit`
Expected: 退出码 0。

- [ ] **Step 4: Commit**

```bash
git add mcp
git commit -m "feat(mcp): add maintenance mcp server"
```

---

### Task N17: 文档与全量验证

**Files:**
- Modify: `frontend/CLAUDE.md`
- Modify: `backend/CLAUDE.md`
- Modify: `README.md`
- Modify: `docs/superpowers/specs/2026-11-13-public-catalog-services-scenarios-design.md`（补「验证记录」）

- [ ] **Step 1: 补文档**

前端：双视图约定、可见性配色、场景页结构。后端：`view` 参数、投影规则、服务/场景模型、备份 v2。README：MCP 启动方式。

- [ ] **Step 2: 全量验证**

```bash
cd backend && go build ./... && go vet ./... && go test ./...
cd frontend && npx tsc --noEmit && npx vitest run && npm run lint && npm run build
cd mcp && npx tsc --noEmit
docker compose up --build -d
```

手工回归：
1. 首页三栏布局与改动前一致；虚拟化、数据保护等有服务的方框出现服务行。
2. 切到客户视图：琥珀内容消失，服务行仍可见。
3. 产品详情客户视图无联系人/黄金三问/竞品/ROI/话术。
4. `/scenarios` 分组、搜索、计数、空态正常；详情页内容完整。
5. 后台新增服务、场景、关联，前台立即可见。
6. 备份导出 v2 → 导入 → 新表数据恢复。
7. MCP 客户端至少调用 6 个工具成功。

- [ ] **Step 3: 记录结果并提交**

```bash
git add frontend/CLAUDE.md backend/CLAUDE.md README.md docs
git commit -m "docs: document dual view, services, scenarios and mcp"
```
