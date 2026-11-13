# Sales-Pilot 后端分层重构 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 把 Go 后端从「handler 直接操作 GORM」重构为 `repository → service → httpapi` 三层，统一错误响应、参数校验与后台列表分页。

**Architecture:** `repository` 只做 GORM 查询，`service` 承载业务规则并返回 `apperr` 领域错误，`httpapi` 负责路由、绑定与响应映射。依赖方向固定为 `httpapi → service → repository → models`；`service` 不 import gin，`repository` 不 import service/httpapi。

**Tech Stack:** Go 1.23、Gin、GORM、PostgreSQL、`httptest` + 手写 fake repository 做单测。

**Spec:** `docs/superpowers/specs/2026-11-13-frontend-backend-conventions-design.md`

## Global Constraints

- **HTTP 契约只增不改**：成功响应形状不变；错误体固定为 `{"error": "...", "code": "..."}`。
- `/api/v1/admin/products` 与 `/api/v1/admin/users/local` 不带 `page` 参数时必须返回裸数组（`/admin/meta/product-maintain` 与外部 `X-API-Key` 自动化依赖它）；带 `page` 才返回 `{items,total,page,page_size}`。
- 500 响应不得回显 `err.Error()`，统一文案「服务器内部错误」，原始错误写日志。
- 路径参数一律经 `httpapi.ParseID`，非法返回 400 + `invalid_id`。
- PATCH 请求体使用指针字段，允许「只传一个字段」。
- 依赖方向：`httpapi → service → repository → models`，禁止反向 import。
- 不修改数据库表结构与 GORM 模型字段。
- 每个任务结束必须通过：`cd backend && go build ./... && go vet ./... && go test ./...`。

## Review Focus

- **分页兼容**：外部调用方不带 `page` 时必须仍拿到数组；带上 `page` 才拿到分页对象。两种形状都要有测试。
- **错误不透传**：数据库错误、SQL 文本、连接串不得出现在响应体里，只能出现在日志。
- **部分更新**：PATCH 只改传入字段，未传字段保持原值；空字符串与 `false` 必须能被显式写入。
- **删除最后一个管理员**：禁止删除或停用最后一名可用本地管理员，返回 409。
- **备份恢复原子性**：任一步失败必须整体回滚，且恢复后不能没有任何可用管理员。

---

### Task B1: 领域错误与统一响应

**Files:**
- Create: `backend/internal/apperr/apperr.go`
- Test: `backend/internal/apperr/apperr_test.go`
- Create: `backend/internal/httpapi/response.go`
- Create: `backend/internal/httpapi/params.go`
- Test: `backend/internal/httpapi/response_test.go`

**Interfaces:**
- Produces:
  - `apperr.Kind`（`Invalid` / `Unauthorized` / `Forbidden` / `NotFound` / `Conflict` / `Internal`）
  - `apperr.Error{Kind, Code, Message string, Err error}` 实现 `Error()` / `Unwrap()`
  - `apperr.New(kind Kind, code, message string) *Error`
  - `apperr.Invalid(code, message string) *Error`、`apperr.NotFound(...)`、`apperr.Conflict(...)`、`apperr.Internal(err error) *Error`
  - `apperr.From(err error) (*apperr.Error, bool)`
  - `apperr.CodeOf(err error) string`（非领域错误返回空字符串）
  - `httpapi.OK(c *gin.Context, v any)`、`httpapi.Created(c *gin.Context, v any)`、`httpapi.NoContent(c *gin.Context)`
  - `httpapi.Fail(c *gin.Context, err error)`
  - `httpapi.ParseID(c *gin.Context, name string) (uint, bool)`

- [ ] **Step 1: 写失败测试**

`apperr_test.go`：

```go
func TestFromUnwraps(t *testing.T) {
    base := apperr.NotFound("product_not_found", "产品不存在")
    wrapped := fmt.Errorf("load: %w", base)
    got, ok := apperr.From(wrapped)
    if !ok || got.Code != "product_not_found" { t.Fatalf("got %+v ok=%v", got, ok) }
}
```

`response_test.go`：

```go
func TestFailMapsNotFound(t *testing.T) {
    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    httpapi.Fail(c, apperr.NotFound("product_not_found", "产品不存在"))
    if w.Code != 404 { t.Fatalf("status=%d", w.Code) }
    if !strings.Contains(w.Body.String(), `"code":"product_not_found"`) { t.Fatalf("body=%s", w.Body.String()) }
}

func TestFailHidesInternalDetail(t *testing.T) {
    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    httpapi.Fail(c, errors.New("pq: relation \"users\" does not exist"))
    if w.Code != 500 { t.Fatalf("status=%d", w.Code) }
    if strings.Contains(w.Body.String(), "pq:") { t.Fatalf("leaked: %s", w.Body.String()) }
}
```

- [ ] **Step 2: 运行确认失败**

Run: `cd backend && go test ./internal/apperr ./internal/httpapi`
Expected: FAIL，包或符号不存在。

- [ ] **Step 3: 实现 `apperr` 与 `httpapi`**

`Fail` 映射：`Invalid→400`、`Unauthorized→401`、`Forbidden→403`、`NotFound→404`、`Conflict→409`、其余→500（文案「服务器内部错误」，`log.Printf` 原始错误）。响应体：

```json
{ "error": "产品不存在", "code": "product_not_found" }
```

`ParseID` 解析失败时自行写出 400 + `invalid_id` 并返回 `false`。

- [ ] **Step 4: 运行确认通过**

Run: `cd backend && go test ./internal/apperr ./internal/httpapi`
Expected: PASS。

- [ ] **Step 5: Commit**

```bash
git add backend/internal/apperr backend/internal/httpapi
git commit -m "feat(backend): add domain errors and unified http responses"
```

---

### Task B2: 分页类型

**Files:**
- Create: `backend/internal/paging/paging.go`
- Test: `backend/internal/paging/paging_test.go`

**Interfaces:**
- Produces:
  - `paging.Request{Page, PageSize int, Sort, Order string}` 及 `Offset() int`、`Limit() int`
  - `paging.New(page, pageSize int) Request`（`page<1→1`，`pageSize<1→20`，`pageSize>100→100`）
  - `paging.Page[T any]{Items []T, Total int64, Page, PageSize int}`，JSON tag 为 `items` / `total` / `page` / `page_size`

- [ ] **Step 1: 写失败测试**

```go
func TestNewClamps(t *testing.T) {
    r := paging.New(0, 0)
    if r.Page != 1 || r.PageSize != 20 { t.Fatalf("%+v", r) }
    r = paging.New(3, 500)
    if r.PageSize != 100 || r.Offset() != 200 { t.Fatalf("%+v", r) }
}
```

- [ ] **Step 2: 运行确认失败，然后实现**

Run: `cd backend && go test ./internal/paging`
Expected: 先 FAIL，实现后 PASS。

- [ ] **Step 3: Commit**

```bash
git add backend/internal/paging
git commit -m "feat(backend): add pagination request and page types"
```

---

### Task B3: 迁移目录整理

**Files:**
- Create: `backend/internal/database/migrations/migrations.go`
- Modify: `backend/internal/database/db.go`

**Interfaces:**
- Produces: `migrations.Migrate(db *gorm.DB) error`，内部按 `users → topology → products → settings` 顺序调用 `AutoMigrate`。

- [ ] **Step 1: 迁移模型清单**

把 `db.go` 中的 `AutoMigrate(...)` 列表整体搬到 `migrations.Migrate`，按依赖顺序分组并加注释。`Connect` 改为调用 `migrations.Migrate(db)`。

- [ ] **Step 2: 验证**

Run: `cd backend && go build ./... && go vet ./...`
Expected: 退出码 0。

- [ ] **Step 3: Commit**

```bash
git add backend/internal/database
git commit -m "refactor(backend): centralize automigrate order in migrations package"
```

---

### Task B4: repository 层

**Files:**
- Create: `backend/internal/repository/user_repository.go`
- Create: `backend/internal/repository/apikey_repository.go`
- Create: `backend/internal/repository/product_repository.go`
- Create: `backend/internal/repository/topology_repository.go`
- Create: `backend/internal/repository/settings_repository.go`
- Create: `backend/internal/repository/backup_repository.go`
- Create: `backend/internal/models/backup.go`

**Interfaces:**
- Produces（构造函数均接收 `*gorm.DB`）：
  - `repository.NewUserRepository(db) *UserRepository`：`ListLocal(ctx) ([]models.User, error)`、`FindByUsername(ctx, string) (*models.User, error)`、`FindByID(ctx, uint) (*models.User, error)`、`Create(ctx, *models.User) error`、`Save(ctx, *models.User) error`、`Delete(ctx, uint) error`、`CountActiveLocalAdmins(ctx) (int64, error)`
  - `repository.NewAPIKeyRepository(db)`：`List(ctx) ([]models.AdminAPIKey, error)`、`FindByHash(ctx, string) (*models.AdminAPIKey, error)`、`Create(ctx, *models.AdminAPIKey) error`、`Revoke(ctx, uint) error`、`TouchLastUsed(ctx, uint) error`
  - `repository.NewProductRepository(db)`：`ListPublic(ctx, solutionID *uint, keyword, market string) ([]models.Product, error)`、`Search(ctx, q string, limit int) ([]models.Product, error)`、`FindByID(ctx, uint) (*models.Product, error)`、`ListAdmin(ctx, f repository.AdminProductFilter, p paging.Request) ([]models.Product, int64, error)`、`ListAdminAll(ctx, f repository.AdminProductFilter) ([]models.Product, error)`、`Create/Save/Delete`、`ListScripts/CreateScript/SaveScript/DeleteScript`、`ListCases/CreateCase/SaveCase/DeleteCase`、`FindCategoryName(ctx, id uint) (string, error)`
  - `repository.NewTopologyRepository(db)`：`Layers/FindLayer/CreateLayer/SaveLayer/DeleteLayer`、`Categories(activeOnly bool)/FindCategory/CreateCategory/SaveCategory/DeleteCategory`、`Vendors/FindVendor/CreateVendor/SaveVendor/DeleteVendor`、`TreeRows(ctx) (layers, categories, vendors, error)`
  - `repository.NewSettingsRepository(db)`：`Auth(ctx) (*models.AuthSettings, error)`、`SaveAuth(ctx, *models.AuthSettings) error`、`System(ctx) (*models.SystemSettings, error)`、`SaveSystem(ctx, *models.SystemSettings) error`
  - `repository.NewBackupRepository(db)`：`Dump(ctx) (*models.BackupPayload, error)`、`Restore(ctx, *models.BackupPayload) error`（失败返回的错误必须能区分「哪一步失败」）

- [ ] **Step 0: 新增备份载荷类型**

`backend/internal/models/backup.go` 定义 `BackupPayload` 与 `BackupUserRow`，字段与 json tag 与现有 `handlers/backup.go` 的 `backupPayloadV1` / `backupUserRow` **逐字一致**（`version`、`exported_at`、`topology_layers`、`solution_categories`、`solution_vendors`、`users`、`products`、`sales_scripts`、`cases`、`auth_settings`）。

- [ ] **Step 1: 从现有 handler 提取查询逻辑**

逐文件对照现有 SQL：`handlers.go`（products/scripts/cases/settings）、`users_admin.go`、`apikeys.go`、`topology.go`、`backup.go`。**查询条件、排序、预加载（`Preload("SolutionCategory")`）必须逐字保留**，只搬迁不改写。

- [ ] **Step 2: 分页查询实现**

`ListAdmin`：`Count` 后按 `paging.Request` 追加 `Order` / `Limit` / `Offset`。`sort` 只接受白名单 `id|name|updated_at|created_at`，其他回退 `id desc`。`AdminProductFilter{Q, Status, Market}`：`Status` 为 `draft|published`（映射 `is_draft`），`Market` 复用现有 `applyVendorMarketFilter` 语义，`Q` 对 `name` 做 `ILIKE`。

- [ ] **Step 3: 验证**

Run: `cd backend && go build ./... && go vet ./...`
Expected: 退出码 0（此时 service 尚未使用这些仓储）。

- [ ] **Step 4: Commit**

```bash
git add backend/internal/repository
git commit -m "feat(backend): add repository layer extracted from handlers"
```

---

### Task B5: service 层与端口接口

**Files:**
- Create: `backend/internal/service/ports.go`
- Create: `backend/internal/service/auth_service.go` + `auth_service_test.go`
- Create: `backend/internal/service/product_service.go` + `product_service_test.go`
- Create: `backend/internal/service/apikey_service.go` + `apikey_service_test.go`
- Create: `backend/internal/service/user_service.go` + `user_service_test.go`
- Create: `backend/internal/service/topology_service.go`
- Create: `backend/internal/service/settings_service.go`
- Create: `backend/internal/service/backup_service.go`

**Interfaces:**
- Consumes: Task B3/B4 的仓储类型；`apperr`、`paging`。
- Produces:
  - `service.Ports` 中定义各仓储接口（`UserRepo`、`APIKeyRepo`、`ProductRepo`、`TopologyRepo`、`SettingsRepo`、`BackupRepo`），`*repository.XxxRepository` 天然满足。
  - `service.NewAuthService(cfg config.Config, users UserRepo, settings SettingsRepo) *AuthService`
    - `LoginLocal(ctx, username, password string) (*models.User, string, error)`（返回用户与 JWT）
    - `LoginLDAP(ctx, username, password string) (*models.User, string, error)`
    - `LoginOptions(ctx) (LoginOptions, error)`
    - `Me(ctx, userID uint) (*models.User, error)`
  - `service.NewProductService(repo ProductRepo) *ProductService`：`ListPublic`、`Search`、`Get`、`ListAdmin`、`ListAdminAll`、`Create`、`Update`、`SetDraft`、`Delete`、`ListScripts/CreateScript/UpdateScript/DeleteScript`、`ListCases/CreateCase/UpdateCase/DeleteCase`
  - `service.NewAPIKeyService(repo APIKeyRepo, users UserRepo) *APIKeyService`：`List`、`Create(ctx, userID uint, name string, expiresAt *time.Time) (string, *models.AdminAPIKey, error)`、`Revoke`、`Authenticate(ctx, raw string) (*models.User, error)`
  - `service.NewUserService(repo UserRepo) *UserService`：`ListLocal`、`CreateLocal(ctx, CreateLocalUserInput) (*models.User, error)`、`UpdateLocal(ctx, id uint, in UpdateLocalUserInput) (*models.User, error)`、`ResetPassword(ctx, id uint, password string) error`、`DeleteLocal(ctx, id uint) error`
    - `service.UpdateLocalUserInput{Email, Role *string; IsActive *bool}` —— 指针字段保证「只传一个字段」的 PATCH 语义
  - `service.NewTopologyService(repo TopologyRepo) *TopologyService`：`Tree`、`ListCategories`、层/分类/厂商 CRUD
  - `service.NewSettingsService(repo SettingsRepo) *SettingsService`：`Auth/PatchAuth/System/PatchSystem`
  - `service.NewBackupService(repo BackupRepo, cache *cache.SearchCache) *BackupService`：`Export(ctx)`、`Import(ctx, payload *models.BackupPayload, confirm string) error`

- [ ] **Step 1: 写失败测试（三个关键业务规则）**

`apikey_service_test.go`：

```go
func TestAuthenticateRejectsExpired(t *testing.T) {
    past := time.Now().Add(-time.Hour)
    svc := service.NewAPIKeyService(&fakeKeyRepo{key: &models.AdminAPIKey{IsActive: true, ExpiresAt: &past, UserID: 1}}, &fakeUserRepo{})
    _, err := svc.Authenticate(context.Background(), "sp_x")
    if code := apperr.CodeOf(err); code != "api_key_expired" { t.Fatalf("code=%s", code) }
}
```

`user_service_test.go`：

```go
func TestDeleteLastAdminConflicts(t *testing.T) {
    svc := service.NewUserService(&fakeUserRepo{admins: 1})
    err := svc.DeleteLocal(context.Background(), 7)
    if code := apperr.CodeOf(err); code != "last_admin" { t.Fatalf("code=%s", code) }
}
```

`product_service_test.go`：

```go
func TestSetDraftNormalizesMarketOnCreate(t *testing.T) {
    repo := &fakeProductRepo{}
    svc := service.NewProductService(repo)
    _, err := svc.Create(context.Background(), &models.Product{Name: "A", VendorMarket: "DOMESTIC "})
    if err != nil { t.Fatal(err) }
    if repo.saved.VendorMarket != "domestic" { t.Fatalf("got %q", repo.saved.VendorMarket) }
}
```

`user_service_test.go` 追加部分更新与备份原子性：

```go
func TestUpdateLocalKeepsUnsetFields(t *testing.T) {
    role := "admin"
    repo := &fakeUserRepo{user: &models.User{ID: 7, Email: "old@x", Role: "user", IsActive: true}}
    svc := service.NewUserService(repo)
    if _, err := svc.UpdateLocal(context.Background(), 7, service.UpdateLocalUserInput{Role: &role}); err != nil { t.Fatal(err) }
    if repo.user.Email != "old@x" { t.Fatalf("email=%q", repo.user.Email) }
    if repo.user.Role != "admin" { t.Fatalf("role=%q", repo.user.Role) }
}

func TestBackupImportRollsBackOnFailure(t *testing.T) {
    repo := &fakeBackupRepo{failAt: "products"}
    svc := service.NewBackupService(repo, cache.NewSearchCache(time.Second))
    err := svc.Import(context.Background(), &models.BackupPayload{Version: 1}, "RESTORE_FULL_V1")
    if err == nil || !repo.rolledBack { t.Fatalf("err=%v rolledBack=%v", err, repo.rolledBack) }
}
```

- [ ] **Step 2: 运行确认失败**

Run: `cd backend && go test ./internal/service`
Expected: FAIL（包不存在）。

- [ ] **Step 3: 实现 service 层**

把现有 handler 中的业务规则搬进 service，**只搬规则不搬 gin**：

- `AuthService.LoginLocal`：本地登录开关（`AuthSettings.LocalLoginEnabled`）→ 停用用户拒绝 → 密码校验 → 签发 JWT。
- `UserService.DeleteLocal` / `UpdateLocal`（`IsActive=false`）：若目标是本地管理员且 `CountActiveLocalAdmins <= 1`，返回 `apperr.Conflict("last_admin", "至少保留一名可用管理员")`。
- `APIKeyService.Authenticate`：hash 查找 → `IsActive` → 过期 → 用户存在且为 admin；错误码分别为 `api_key_invalid` / `api_key_expired` / `api_key_forbidden`。
- `ProductService.Create/Update`：`vendorMarket` 归一化为 `all|domestic|foreign`；`SolutionCategoryID` 存在时同步 `Category` 名称（复用现有 `syncProductCategoryFromSolution`）。
- `TopologyService.Tree`：返回 `{security, ops, central_layers}` 现有形状。
- `BackupService.Import`：确认令牌 → 版本校验 → 事务内 TRUNCATE/写入 → 至少一名可用本地管理员 → 成功后 `cache.Clear()`；失败 `apperr.Invalid` / `apperr.Conflict`。

- [ ] **Step 4: 运行确认通过**

Run: `cd backend && go test ./internal/service`
Expected: PASS。

- [ ] **Step 5: Commit**

```bash
git add backend/internal/service
git commit -m "feat(backend): add service layer with business rules and unit tests"
```

---

### Task B6: httpapi 路由与 handler

**Files:**
- Create: `backend/internal/httpapi/router.go`
- Create: `backend/internal/httpapi/handlers/auth_handler.go`
- Create: `backend/internal/httpapi/handlers/product_handler.go`
- Create: `backend/internal/httpapi/handlers/admin_product_handler.go`
- Create: `backend/internal/httpapi/handlers/topology_handler.go`
- Create: `backend/internal/httpapi/handlers/admin_user_handler.go`
- Create: `backend/internal/httpapi/handlers/apikey_handler.go`
- Create: `backend/internal/httpapi/handlers/settings_handler.go`
- Create: `backend/internal/httpapi/handlers/backup_handler.go`
- Create: `backend/internal/httpapi/handlers/meta_handler.go`
- Create: `backend/internal/httpapi/middleware/jwt.go`、`admin_auth.go`（从 `internal/middleware` 迁入，改为依赖 `*service.APIKeyService`）
- Modify: `backend/cmd/server/main.go`
- Delete: `backend/internal/handlers/*.go`、`backend/internal/middleware/*.go`

**Interfaces:**
- Consumes: Task B5 全部 service 方法；`httpapi.OK/Created/NoContent/Fail/ParseID`。
- Produces:
  - `httpapi.Deps{Cfg config.Config; Auth *service.AuthService; Products *service.ProductService; Topology *service.TopologyService; Users *service.UserService; APIKeys *service.APIKeyService; Settings *service.SettingsService; Backup *service.BackupService; OidcStates *auth.OidcStateStore}`
  - `httpapi.NewRouter(d Deps) *gin.Engine`

- [ ] **Step 1: 注册路由，路径与方法逐条对照现有 `RegisterRoutes`**

保留全部现有路径与 HTTP 方法。`session` 组继续用 `AdminAuth`（Bearer 或 X-API-Key），`admin` 组再叠加 `RequireAdmin()`。CORS 中间件随 `main.go` 迁入 `router.go`，行为不变。

- [ ] **Step 2: handler 只做绑定与映射**

每个 handler 结构：绑定请求 → 调用 service → `OK/Created/NoContent`，错误一律 `Fail`。禁止在 handler 里写 `s.DB`。`ParseID` 替换全部 `strconv.Atoi(c.Param(...))`。

- [ ] **Step 3: 后台列表分页形状**

`GET /api/v1/admin/products`、`GET /api/v1/admin/users/local`：

```go
if c.Query("page") == "" {
    httpapi.OK(c, list)                       // 裸数组，兼容外部
    return
}
p := paging.New(queryInt(c, "page", 1), queryInt(c, "page_size", 20))
httpapi.OK(c, paging.Page[models.Product]{Items: items, Total: total, Page: p.Page, PageSize: p.PageSize})
```

- [ ] **Step 4: 删除旧包并接线 `main.go`**

`main.go` 只保留：加载配置 → 连库 → 迁移 → seed → 构造 repository → 构造 service → `httpapi.NewRouter` → `Run`。

- [ ] **Step 5: 验证**

Run: `cd backend && go build ./... && go vet ./... && go test ./...`
Expected: 退出码 0；旧 `internal/handlers`、`internal/middleware` 已不存在。

- [ ] **Step 6: Commit**

```bash
git add backend
git commit -m "refactor(backend): move http layer to httpapi and wire services"
```

---

### Task B7: 路由冒烟测试

**Files:**
- Create: `backend/internal/httpapi/router_test.go`

**Interfaces:**
- Consumes: `httpapi.NewRouter` + fake repository（复用 Task B5 的 fake，放到 `internal/service/testdata_test.go` 供两个包共用，或在新包内重建）。

- [ ] **Step 1: 写测试（四条契约）**

```go
// testAdminJWT：auth.IssueToken(cfg.JWTSecret, 1, "admin", "admin") 的结果
func TestAdminProductsWithoutPageReturnsArray(t *testing.T) {
    r := newTestRouter(t)
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/products", nil)
    req.Header.Set("Authorization", "Bearer "+testAdminJWT(t))
    r.ServeHTTP(w, req)
    if w.Code != 200 { t.Fatalf("status=%d", w.Code) }
    if !strings.HasPrefix(strings.TrimSpace(w.Body.String()), "[") { t.Fatalf("body=%s", w.Body.String()) }
}

func TestAdminProductsWithPageReturnsEnvelope(t *testing.T) { /* 断言 items/total/page/page_size */ }

func TestUnauthorizedReturnsCode(t *testing.T) { /* 无 token → 401，body 含 "code" */ }

func TestInvalidIDReturns400(t *testing.T) { /* /api/v1/admin/products/abc → 400 + invalid_id */ }
```

- [ ] **Step 2: 运行确认通过**

Run: `cd backend && go test ./internal/httpapi`
Expected: PASS。

- [ ] **Step 3: Commit**

```bash
git add backend/internal/httpapi/router_test.go
git commit -m "test(backend): cover pagination compatibility and error contracts"
```

---

### Task B8: 后端规范文档与全量验证

**Files:**
- Create: `backend/CLAUDE.md`
- Modify: `README.md`

- [ ] **Step 1: 写 `backend/CLAUDE.md`**

包含：目录与依赖方向、错误码表（`product_not_found`、`last_admin`、`api_key_expired`、`invalid_id` 等）、统一响应示例、分页约定、新增端点步骤（repository → service → handler → router）、测试要求。

- [ ] **Step 2: README 增加「开发规范」索引**

指向 `frontend/CLAUDE.md` 与 `backend/CLAUDE.md`。

- [ ] **Step 3: 全量验证**

Run: `cd backend && go build ./... && go vet ./... && go test ./...`
Run: `cd frontend && npx tsc --noEmit && npm run build`
Run: `docker compose up --build -d`，按主链路手工回归（登录、产品列表/编辑、拓扑、用户、密钥、备份）。

- [ ] **Step 4: Commit**

```bash
git add backend/CLAUDE.md README.md
git commit -m "docs: document backend layering and error contracts"
```
