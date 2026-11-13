# Sales-Pilot 前端规范化 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 把 Sales-Pilot 前端改造成与 ALog 同构的分层结构，并让所有列表、容器、状态、筛选统一使用 ReUI 组件。

**Architecture:** 路由页只做取数与拼装，页面级状态移入 `features/`，通用交互进 `components/patterns/`，业务词汇组件进 `components/domain/`，ReUI 产物进 `components/reui/`。数据访问经 `lib/api/{client,endpoints,adapters,types}`，DTO 与视图模型分离。

**Tech Stack:** Next.js 16 App Router、React 19、TypeScript strict、Tailwind v4、shadcn（base-nova）、ReUI（data-grid / frame / badge / alert / filters / icon-tile / icon-stack / number-field / cascader / phone-input / autocomplete）、TanStack Table v9、Vitest。

**Spec:** `docs/superpowers/specs/2026-11-13-frontend-backend-conventions-design.md`

## Global Constraints

- `@/` 别名指向 `frontend/src/`；保留 `src/` 目录。
- 容器只用 ReUI `Frame`；列表只用 ReUI `DataGrid`；状态只用 ReUI `Badge`（`-light` 变体）；提示只用 ReUI `Alert`；筛选只用 ReUI `Filters`。
- 间距用 `gap-*`，不用 `space-y-*`；字号不低于 `text-xs`；动效只用 `transition-colors`；不做 `transition-all`；每个视图最多一个入场动画；表格行不做 stagger。
- 页面与组件禁止写死色值（`bg-amber-100`、`text-emerald-900` 等），一律使用语义令牌。
- `app/` 下页面文件 ≤ 60 行（`auth/callback` 已达标除外）。
- 不引入 react-query、react-hook-form、zod、openapi 代码生成。
- 每个任务结束必须通过：`cd frontend && npx tsc --noEmit && npm run build`。
- HTTP 契约本轮不变；前端只消费现有端点。

## Review Focus

- **空/加载/错误三态**：任何列表在 0 条数据、请求失败、加载中都必须渲染对应状态，不能白屏或空白表格。
- **窄屏与长文本**：产品名、分类、厂商在 375px 宽度下必须截断且可读，表格横向滚动不能撑破页面。
- **危险操作**：删除产品/用户/密钥、恢复备份必须二次确认，且确认后按钮进入禁用态，不能双击重复提交。
- **未登录与无权限**：401 跳转登录，403 渲染无权限提示，而不是空白页或控制台报错。
- **深色模式**：所有语义令牌在 `.dark` 下有足够对比度，不出现浅底浅字。

---

### Task 1: 设计令牌

**Files:**
- Create: `frontend/src/styles/tokens.css`
- Modify: `frontend/src/app/globals.css`

**Interfaces:**
- Produces: CSS 变量 `--brand`、`--success`、`--success-foreground`、`--warning`、`--warning-foreground`、`--info`、`--info-foreground`；Tailwind 工具类 `bg-success`、`text-warning` 等（经 `@theme inline`）。

- [ ] **Step 1: 新建 `src/styles/tokens.css`，把 `globals.css` 中 `:root` 与 `.dark` 的两段变量整体迁入**

保留现有企业蓝 `--primary: #0052cc`，新增语义色（Light / Dark 各一套）：

```css
:root {
  --brand: #0052cc;
  --success: oklch(0.60 0.13 155);
  --success-foreground: oklch(0.98 0 0);
  --warning: oklch(0.72 0.15 75);
  --warning-foreground: oklch(0.20 0.02 75);
  --info: oklch(0.60 0.12 245);
  --info-foreground: oklch(0.98 0 0);
}
.dark {
  --brand: oklch(0.68 0.16 260);
  --primary: oklch(0.68 0.16 260);        /* 修复暗色 primary 仍是灰白 */
  --primary-foreground: oklch(0.14 0 0);
  --success: oklch(0.70 0.14 155);
  --success-foreground: oklch(0.14 0 0);
  --warning: oklch(0.78 0.15 75);
  --warning-foreground: oklch(0.16 0.02 75);
  --info: oklch(0.70 0.12 245);
  --info-foreground: oklch(0.14 0 0);
}
```

- [ ] **Step 2: 在 `globals.css` 顶部加 `@import "../styles/tokens.css";`，并在 `@theme inline` 增加映射**

```css
--color-brand: var(--brand);
--color-success: var(--success);
--color-success-foreground: var(--success-foreground);
--color-warning: var(--warning);
--color-warning-foreground: var(--warning-foreground);
--color-info: var(--info);
--color-info-foreground: var(--info-foreground);
```

- [ ] **Step 3: 验证构建**

Run: `cd frontend && npx tsc --noEmit && npm run build`
Expected: 两条命令退出码均为 0。

- [ ] **Step 4: 浏览器验证深浅色**

Run: `cd frontend && npm run dev`
打开 `/login` 与 `/admin`，切到 `.dark` 后确认主色是蓝色而非灰白，文字可读。

- [ ] **Step 5: Commit**

```bash
git add frontend/src/styles/tokens.css frontend/src/app/globals.css
git commit -m "style(frontend): centralize design tokens and add semantic colors"
```

---

### Task 2: ReUI registry 与组件安装

**Files:**
- Modify: `frontend/components.json`
- Modify: `frontend/package.json`
- Create: `frontend/src/components/reui/*`（CLI 生成）

**Interfaces:**
- Produces: `@/components/reui/frame`、`badge`、`alert`、`data-grid/data-grid`、`data-grid/data-grid-pagination`、`data-grid/data-grid-table`、`filters/filters`、`icon-tile`、`icon-stack`、`number-field`、`cascader`、`phone-input`、`autocomplete`。

- [ ] **Step 1: 在 `components.json` 加 registry**

```json
"registries": {
  "@reui": "https://reui.io/r/{style}/{name}.json"
}
```

- [ ] **Step 2: 安装 peer 依赖**

Run: `cd frontend && npm install @tanstack/react-table@^9 @tanstack/react-virtual@^3 @dnd-kit/core @dnd-kit/modifiers @dnd-kit/sortable @dnd-kit/utilities cn`

- [ ] **Step 3: 安装 ReUI 组件**

Run: `cd frontend && npx shadcn@latest add @reui/alert @reui/badge @reui/frame @reui/data-grid @reui/filters @reui/icon-tile @reui/icon-stack @reui/number-field @reui/cascader @reui/phone-input @reui/autocomplete --yes`
Expected: `src/components/reui/` 下生成对应文件；若提示覆盖 `components/ui/*`，选择覆盖（新版本兼容）。

- [ ] **Step 4: 读一遍 data-grid 的真实 API**

用 ReUI MCP `get_component("data-grid", sections: ["DataGridPagination"])` 确认分页组件 props；`get_examples("data-grid")` 选一个带分页的 `c-data-grid-*` 作为列定义参考。

- [ ] **Step 5: 验证构建**

Run: `cd frontend && npx tsc --noEmit && npm run build`
Expected: 退出码 0。

- [ ] **Step 6: Commit**

```bash
git add frontend/components.json frontend/package.json frontend/package-lock.json frontend/src/components/reui frontend/src/components/ui
git commit -m "chore(frontend): add ReUI registry and install core components"
```

---

### Task 3: API 层拆分与视图模型

**Files:**
- Create: `frontend/src/lib/api/types.ts`
- Create: `frontend/src/lib/api/problem.ts`
- Create: `frontend/src/lib/api/client.ts`
- Create: `frontend/src/lib/api/adapters.ts`
- Create: `frontend/src/lib/api/endpoints.ts`
- Create: `frontend/src/lib/api/index.ts`
- Create: `frontend/src/lib/api/adapters.test.ts`
- Delete: `frontend/src/lib/api.ts`
- Modify: 19 个 `from "@/lib/api"` 引用点
- Modify: `frontend/package.json`（新增 `test` 脚本与 vitest）

**Interfaces:**
- Produces:
  - `apiFetch(path: string, init?: RequestInit, opts?: { token?: string | null }): Promise<Response>`
  - `apiJson<T>(path: string, init?: RequestInit): Promise<T>`（非 2xx 抛 `ApiError`）
  - `class ApiError extends Error { status: number; code: string }`
  - 视图模型：`Product`（camelCase）、`Topology`、`LocalUser`、`AdminApiKey`、`SalesScript`、`Case`、`AuthSettings`、`SystemSettings`
  - 适配器：`toProduct(dto: ProductDTO): Product`、`toProducts(dto: ProductDTO[]): Product[]`、`toTopology(dto: TopologyDTO): Topology`、`toLocalUser(dto: LocalUserDTO): LocalUser`、`toApiKey(dto: AdminApiKeyDTO): AdminApiKey`
  - 端点模块：`authApi`、`productsApi`、`adminProductsApi`、`topologyApi`、`adminTopologyApi`、`usersApi`、`apiKeysApi`、`settingsApi`、`backupApi`

- [ ] **Step 1: 写失败测试 `adapters.test.ts`**

```ts
import { describe, expect, it } from "vitest";
import { toProduct, toLocalUser } from "./adapters";

describe("toProduct", () => {
  it("maps snake_case DTO to camelCase model and parses json fields", () => {
    const p = toProduct({
      id: 1, name: "A", category: "c", description: "",
      highlights: "[\"h1\",\"h2\"]", target_personas: { x: "y" },
      discovery_questions: "", competitor_analysis: "", roi_metrics: "",
      manufacturer_name: "M", vendor_market: "domestic", is_draft: true,
      updated_at: "2026-01-01T00:00:00Z", created_at: "2026-01-01T00:00:00Z",
    });
    expect(p.manufacturerName).toBe("M");
    expect(p.vendorMarket).toBe("domestic");
    expect(p.isDraft).toBe(true);
    expect(p.highlights).toEqual(["h1", "h2"]);
    expect(p.targetPersonas).toEqual({ x: "y" });
  });
});

describe("toLocalUser", () => {
  it("keeps active flag and iso times", () => {
    const u = toLocalUser({
      id: 2, username: "u", email: "e", role: "admin", is_active: false,
      created_at: "2026-01-01T00:00:00Z", updated_at: "2026-01-01T00:00:00Z",
    });
    expect(u.isActive).toBe(false);
    expect(u.createdAt).toBe("2026-01-01T00:00:00Z");
  });
});
```

- [ ] **Step 2: 安装 vitest、补脚本并运行，确认失败**

Run: `cd frontend && npm install -D vitest && npm pkg set scripts.test="vitest run" scripts.typecheck="tsc --noEmit" && npx vitest run src/lib/api/adapters.test.ts`
Expected: FAIL，提示 `Cannot find module './adapters'`。

- [ ] **Step 3: 实现 `types.ts` 与 `adapters.ts`**

`types.ts` 只声明后端 JSON 形状，字段保持 snake_case，与 `backend/internal/models` 的 `json` tag 一致。`adapters.ts` 用 `parseJsonArray` / `parsePersonas`（从旧 `lib/api.ts` 迁移）做健壮解析，输出 camelCase 视图模型。

- [ ] **Step 4: 运行测试确认通过**

Run: `cd frontend && npx vitest run src/lib/api/adapters.test.ts`
Expected: PASS，2 passed。

- [ ] **Step 5: 实现 `problem.ts`、`client.ts`**

`client.ts` 保留现有 `apiFetch` 行为：自动带 Bearer、字符串 body 自动 `Content-Type: application/json`；新增 `apiJson<T>`，非 2xx 时解析 `{error, code}` 抛 `ApiError`；401 时调用 `clearToken()`。

- [ ] **Step 6: 实现 `endpoints.ts` 与 `index.ts`，删除 `src/lib/api.ts`，更新 19 处 import，并同步改字段名**

机械替换 import 路径（目录形式，`index.ts` 统一导出）。同时把 6 个消费者文件里共 81 处 snake_case 字段改为 camelCase 视图模型字段（`manufacturer_name → manufacturerName`、`is_draft → isDraft`、`updated_at → updatedAt` 等）：

```
src/app/(marketing)/products/page.tsx
src/app/(marketing)/products/[id]/page.tsx
src/app/admin/products/page.tsx
src/app/admin/products/[id]/page.tsx
src/components/marketing/product-spotlight.tsx
src/components/marketing/solution-atlas.tsx
```

本任务只改字段名，**不做文件移动**（移动在 Task 15）。改完必须能过 `tsc`。

- [ ] **Step 7: 验证**

Run: `cd frontend && npx tsc --noEmit && npx vitest run && npm run build`
Expected: 全部退出码 0。

- [ ] **Step 8: Commit**

```bash
git add frontend/package.json frontend/package-lock.json frontend/src/lib frontend/src/app frontend/src/components
git commit -m "refactor(frontend): split api layer into client, types, adapters and endpoints"
```

---

### Task 4: 通用 patterns 组件

**Files:**
- Create: `frontend/src/components/patterns/page-header.tsx`
- Create: `frontend/src/components/patterns/empty-state.tsx`
- Create: `frontend/src/components/patterns/error-state.tsx`
- Create: `frontend/src/components/patterns/loading-state.tsx`
- Create: `frontend/src/components/patterns/confirm-dialog.tsx`
- Create: `frontend/src/components/patterns/stat-card.tsx`
- Create: `frontend/src/components/patterns/form-field.tsx`

**Interfaces:**
- Produces:
  - `PageHeader({ title: string; description?: string; actions?: React.ReactNode })`
  - `EmptyState({ icon?: LucideIcon; title: string; description?: string; action?: React.ReactNode })`
  - `ErrorState({ message: string; onRetry?: () => void })`
  - `LoadingState({ rows?: number })`
  - `ConfirmDialog({ open, onOpenChange, title, description, confirmLabel, onConfirm, destructive? })`
  - `StatCard({ label: string; value: React.ReactNode; hint?: string; icon: LucideIcon })`
  - `FormField({ label: string; htmlFor?: string; description?: string; error?: string; children: React.ReactNode })`

- [ ] **Step 1: 实现 7 个组件**

`EmptyState` 用 `IconStack` + `IconTile`；`ErrorState` 用 ReUI `Alert variant="destructive"`；`PageHeader` 用 `flex items-end justify-between gap-4`；`ConfirmDialog` 用 shadcn `Dialog` + ReUI `Alert`，`onConfirm` 为 async 时内部维护 `pending` 禁用确认按钮。

- [ ] **Step 2: 验证**

Run: `cd frontend && npx tsc --noEmit && npm run build`
Expected: 退出码 0。

- [ ] **Step 3: Commit**

```bash
git add frontend/src/components/patterns
git commit -m "feat(frontend): add shared page patterns built on ReUI"
```

---

### Task 5: domain 组件

**Files:**
- Create: `frontend/src/components/domain/publish-status-badge.tsx`
- Create: `frontend/src/components/domain/market-badge.tsx`
- Create: `frontend/src/components/domain/vendor-tag.tsx`
- Create: `frontend/src/components/domain/contact-card.tsx`
- Create: `frontend/src/components/domain/product-card.tsx`

**Interfaces:**
- Produces:
  - `PublishStatusBadge({ isDraft: boolean })`
  - `MarketBadge({ market?: "all" | "domestic" | "foreign" })`
  - `VendorTag({ name?: string; logo?: string })`
  - `ContactCard({ kind: "sales" | "presales"; name?: string; phone?: string; email?: string })`
  - `ProductCard({ product: Product; href: string })`

- [ ] **Step 1: 实现组件**

`PublishStatusBadge`：`<Badge variant="warning-light">草稿</Badge>` / `<Badge variant="success-light">已发布</Badge>`。
`MarketBadge`：国内 `info-light`、国外 `secondary`、不区分 `outline`。
`ProductCard` 外层必须是 ReUI `Frame`（不得用 Card）。

- [ ] **Step 2: 验证并提交**

Run: `cd frontend && npx tsc --noEmit && npm run build`

```bash
git add frontend/src/components/domain
git commit -m "feat(frontend): add domain components for products and statuses"
```

---

### Task 6: 路由级状态页与 AdminShell

**Files:**
- Create: `frontend/src/app/admin/loading.tsx`
- Create: `frontend/src/app/admin/error.tsx`
- Create: `frontend/src/app/admin/not-found.tsx`
- Create: `frontend/src/app/(marketing)/error.tsx`
- Create: `frontend/src/app/(marketing)/not-found.tsx`
- Create: `frontend/src/app/not-found.tsx`
- Modify: `frontend/src/components/admin/admin-shell.tsx`

**Interfaces:**
- Consumes: `patterns/loading-state`、`patterns/error-state`。
- Produces: 路由级错误/404 兜底；AdminShell 使用语义令牌与 ReUI `Badge` 展示角色。

- [ ] **Step 1: 实现 error/loading/not-found**

`error.tsx` 必须是 `"use client"`，接收 `{ error, reset }`，用 `ErrorState` 渲染并提供重试。`ApiError.status === 401` 时在 `client.ts` 已清 token，页面错误边界负责跳转 `/login`；`status === 403` 时渲染「没有访问权限」文案而不是通用错误。

- [ ] **Step 2: AdminShell 对齐令牌**

活动项改用 `bg-primary/10 text-primary`（保留），退出按钮用 `Button variant="outline"`；用户名旁加 `<Badge variant="secondary">{user.role === "admin" ? "管理员" : "用户"}</Badge>`。

- [ ] **Step 3: 验证并提交**

Run: `cd frontend && npx tsc --noEmit && npm run build`

```bash
git add frontend/src/app frontend/src/components/admin
git commit -m "feat(frontend): add route-level loading, error and not-found states"
```

---

### Task 7: 后台首页指标卡

**Files:**
- Modify: `frontend/src/app/admin/page.tsx`
- Create: `frontend/src/features/dashboard/admin-overview.tsx`

**Interfaces:**
- Consumes: `adminProductsApi.list`、`usersApi.list`、`apiKeysApi.list`、`StatCard`。
- Produces: `AdminOverview()` 客户端组件。

- [ ] **Step 1: 实现 `AdminOverview`**

并行取三份数据（`Promise.allSettled`），渲染三张 `StatCard`（产品总数 / 本地用户 / 有效密钥），下方保留快捷入口（Frame 列表）。任一请求失败时该卡显示 `—` 并附 `text-xs text-muted-foreground` 的「加载失败」。

- [ ] **Step 2: 页面只拼装**

`page.tsx` ≤ 15 行：`<AdminOverview />`。

- [ ] **Step 3: 验证并提交**

Run: `cd frontend && npx tsc --noEmit && npm run build`

```bash
git add frontend/src/app/admin/page.tsx frontend/src/features/dashboard
git commit -m "feat(admin): replace static dashboard cards with live stat cards"
```

---

### Task 8: 后台产品列表迁移到 DataGrid

**Files:**
- Create: `frontend/src/features/products/admin-products-list.tsx`
- Modify: `frontend/src/app/admin/products/page.tsx`
- Modify: `frontend/src/lib/api/endpoints.ts`（`adminProductsApi.list` 支持分页参数）

**Interfaces:**
- Consumes: `adminProductsApi.list(params?) => Promise<PageResult<Product>>`；`PublishStatusBadge`、`MarketBadge`、`PageHeader`、`ConfirmDialog`、`DataGrid`。
- Produces: `AdminProductsList()`。

- [ ] **Step 1: 扩展端点**

```ts
export type PageResult<T> = { items: T[]; total: number; page: number; pageSize: number; };

// 适配器统一两种后端形状：带 page 返回 {items,total,page,page_size}；不带 page 返回裸数组
export function toPageResult<T, M>(data: unknown, page: number, pageSize: number, map: (dto: T) => M): PageResult<M>;
```

`adminProductsApi.list(params)` 始终返回 `PageResult<Product>`：后端返回裸数组时包成 `{ items, total: items.length, page: 1, pageSize: items.length }`。同名适配器供 `usersApi.list` 复用。

- [ ] **Step 2: 定义列**

`ColumnDef<DataGridFeatures, Product>[]`：ID（宽 64）、名称（`max-w-[260px] truncate`）、状态（`PublishStatusBadge`）、分类、厂商、市场（`MarketBadge`）、更新时间、操作（编辑 / 发布切换 / 删除）。

- [ ] **Step 3: 接入 manual 分页与筛选**

`useTable({ features: dataGridFeatures, data, columns, rowCount: total, manualPagination: true, manualSorting: true, manualFiltering: true, state: { pagination, sorting }, onPaginationChange, onSortingChange })`。筛选栏用 ReUI `Filters`，只暴露 `q / status / market` 三个维度，变更时重置 `pageIndex = 0`。具体 props 以 `get_component("filters")` 为准。

- [ ] **Step 4: 删除手写表格与 `window.confirm`**

删除按钮打开 `ConfirmDialog`；确认后调用 `adminProductsApi.remove(id)`，成功后刷新当前页。

- [ ] **Step 5: 页面只拼装**

`page.tsx` 只渲染 `PageHeader` + `AdminProductsList`。

- [ ] **Step 6: 验证**

Run: `cd frontend && npx tsc --noEmit && npm run build`
浏览器：空数据渲染 `EmptyState`；接口失败渲染 `ErrorState`；翻页后筛选仍生效；删除弹确认。

- [ ] **Step 7: Commit**

```bash
git add frontend/src/features/products frontend/src/app/admin/products/page.tsx frontend/src/lib/api
git commit -m "feat(admin): migrate products list to ReUI data-grid"
```

---

### Task 9: 产品编辑器拆分

**Files:**
- Create: `frontend/src/features/products/product-editor.tsx`
- Create: `frontend/src/features/products/scripts-editor.tsx`
- Create: `frontend/src/features/products/cases-editor.tsx`
- Modify: `frontend/src/app/admin/products/[id]/page.tsx`

**Interfaces:**
- Consumes: `adminProductsApi.get/update/listScripts/createScript/updateScript/deleteScript/listCases/createCase/updateCase/deleteCase`；`Frame`、`Cascader`、`PhoneInput`、`NumberField`、`FormField`。
- Produces: `ProductEditor({ productId: number })`，内部渲染三个编辑区。

- [ ] **Step 1: 拆出三个组件**

`ProductEditor` 负责产品字段与保存；`ScriptsEditor`、`CasesEditor` 各自维护增删改状态。每个区块用 `<Frame><FramePanel>…</FramePanel></Frame>`，标题用 `FrameHeader`。分类字段用 `Cascader`（层 → 分类），电话用 `PhoneInput`。

- [ ] **Step 2: 页面只拼装（≤ 20 行）**

```tsx
export default function Page({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  return <ProductEditor productId={Number(id)} />;
}
```
（Next 16 的 `params` 为 Promise，使用 `use()` 解包。）

- [ ] **Step 3: 验证并提交**

Run: `cd frontend && npx tsc --noEmit && npm run build`
浏览器：新建草稿 → 编辑 → 保存 → 刷新仍在；话术/案例增删改正常。

```bash
git add frontend/src/features/products frontend/src/app/admin/products/\[id\]/page.tsx
git commit -m "refactor(admin): split product editor into focused feature components"
```

---

### Task 10: 本地用户列表迁移到 DataGrid

**Files:**
- Create: `frontend/src/features/users/users-data-grid.tsx`
- Modify: `frontend/src/app/admin/users/page.tsx`

**Interfaces:**
- Consumes: `usersApi.list({ page, q, role, active }) => PageResult<LocalUser> | LocalUser[]`；`DataGrid`、`Filters`、`Badge`、`ConfirmDialog`。
- Produces: `UsersDataGrid()`。

- [ ] **Step 1: 实现列与操作**

列：用户名、邮箱、角色（`Badge variant={role === "admin" ? "info-light" : "secondary"}`）、状态（启用 `success-light` / 停用 `outline`）、创建时间、操作（编辑邮箱/角色、重置密码、启停、删除）。

- [ ] **Step 2: 危险操作全部走 ConfirmDialog**

删除、停用、重置密码必须二次确认；确认期间按钮 `disabled`。

- [ ] **Step 3: 验证并提交**

Run: `cd frontend && npx tsc --noEmit && npm run build`
浏览器：创建用户 → 修改角色 → 停用 → 删除；空列表显示 `EmptyState`。

```bash
git add frontend/src/features/users frontend/src/app/admin/users/page.tsx
git commit -m "feat(admin): migrate local users to ReUI data-grid"
```

---

### Task 11: API 密钥列表迁移到 DataGrid

**Files:**
- Create: `frontend/src/features/api-keys/api-keys-data-grid.tsx`
- Modify: `frontend/src/app/admin/api-keys/page.tsx`

**Interfaces:**
- Consumes: `apiKeysApi.list/create/revoke`；`DataGrid`、`Badge`、`Alert`、`ConfirmDialog`。
- Produces: `ApiKeysDataGrid()`。

- [ ] **Step 1: 实现**

列：名称、前缀（`font-mono`）、状态（有效 `success-light` / 已作废 `outline` / 已过期 `warning-light`）、过期时间、最后使用、操作（作废）。新建后完整密钥只展示一次，用 ReUI `Alert variant="warning"` 包裹并提示「仅显示一次」。

- [ ] **Step 2: 验证并提交**

Run: `cd frontend && npx tsc --noEmit && npm run build`

```bash
git add frontend/src/features/api-keys frontend/src/app/admin/api-keys/page.tsx
git commit -m "feat(admin): migrate api keys to ReUI data-grid"
```

---

### Task 12: 认证配置页拆分

**Files:**
- Create: `frontend/src/features/settings/auth-settings-form.tsx`
- Create: `frontend/src/features/settings/system-settings-form.tsx`
- Modify: `frontend/src/app/admin/settings/page.tsx`

**Interfaces:**
- Consumes: `settingsApi.auth/patchAuth/system/patchSystem`；`Frame`、`FormField`、`Alert`、`ConfirmDialog`。
- Produces: `AuthSettingsForm()`、`SystemSettingsForm()`。

- [ ] **Step 1: 两个表单各自独立**

每个表单用 `<Frame class="max-w-3xl"><FramePanel>…</FramePanel></Frame>`，保存成功用 `Alert variant="success"`，失败用 `destructive`。切换登录方式（本地/LDAP/OIDC）属于敏感操作，保存前弹 `ConfirmDialog`。

- [ ] **Step 2: 页面只拼装**

`page.tsx` 渲染两个表单。

- [ ] **Step 3: 验证并提交**

Run: `cd frontend && npx tsc --noEmit && npm run build`

```bash
git add frontend/src/features/settings frontend/src/app/admin/settings/page.tsx
git commit -m "refactor(admin): split auth and system settings into feature forms"
```

---

### Task 13: 备份恢复页

**Files:**
- Create: `frontend/src/features/backup/backup-view.tsx`
- Modify: `frontend/src/app/admin/backup/page.tsx`

**Interfaces:**
- Consumes: `backupApi.exportUrl`、`backupApi.import(payload, confirm)`；`Frame`、`Alert`、`ConfirmDialog`、`FormField`。
- Produces: `BackupView()`。

- [ ] **Step 1: 实现**

导出用 `<a href={backupApi.exportUrl()}>`（保持现有鉴权方式）。导入必须选择文件、展示解析后的版本与条数摘要，再通过 `ConfirmDialog` 输入确认令牌 `RESTORE_FULL_V1` 才提交。全程用 `Alert` 展示结果。

- [ ] **Step 2: 验证并提交**

Run: `cd frontend && npx tsc --noEmit && npm run build`
浏览器：导出触发下载；导入非法 JSON 显示 destructive Alert；确认令牌错误被拒。

```bash
git add frontend/src/features/backup frontend/src/app/admin/backup/page.tsx
git commit -m "refactor(admin): rebuild backup page with Frame and confirm dialog"
```

---

### Task 14: 拓扑维护页拆分

**Files:**
- Create: `frontend/src/features/topology/topology-editor.tsx`
- Create: `frontend/src/features/topology/layer-form.tsx`
- Create: `frontend/src/features/topology/category-form.tsx`
- Create: `frontend/src/features/topology/vendor-form.tsx`
- Modify: `frontend/src/app/admin/topology/page.tsx`

**Interfaces:**
- Consumes: `adminTopologyApi.*`、`topologyApi.get`；`DataGrid`、`NumberField`、`Cascader`、`ConfirmDialog`、`Frame`。
- Produces: `TopologyEditor()`。

- [ ] **Step 1: 三个表格改用 DataGrid**

分层、分类、厂商三处手写 `<table>` 全部替换；排序列用 `NumberField` 直接编辑，保存调用 back-end 的 PUT。

- [ ] **Step 2: 表单抽到独立文件**

层/分类/厂商表单各自独立，使用 `FormField` + `Frame`；删除走 `ConfirmDialog`（分类与厂商删除会级联，文案必须写明）。

- [ ] **Step 3: 验证并提交**

Run: `cd frontend && npx tsc --noEmit && npm run build`
浏览器：新增层 → 新增分类 → 新增厂商 → 调整排序 → 删除；前台首页拓扑同步变化。

```bash
git add frontend/src/features/topology frontend/src/app/admin/topology/page.tsx
git commit -m "refactor(admin): split topology editor and migrate tables to data-grid"
```

---

### Task 15: 营销页面组件化

**Files:**
- Create: `frontend/src/features/marketing/product-catalog.tsx`
- Create: `frontend/src/features/marketing/solution-atlas.tsx`（从 components 迁入）
- Create: `frontend/src/features/marketing/product-spotlight.tsx`（从 components 迁入）
- Create: `frontend/src/features/auth/login-form.tsx`
- Modify: `frontend/src/app/(marketing)/page.tsx`
- Modify: `frontend/src/app/(marketing)/products/page.tsx`
- Modify: `frontend/src/app/(marketing)/products/[id]/page.tsx`
- Modify: `frontend/src/app/(marketing)/login/page.tsx`
- Delete: `frontend/src/components/marketing/solution-atlas.tsx`、`product-spotlight.tsx`

**Interfaces:**
- Consumes: `productsApi.list/search/get`、`topologyApi.get`、`authApi.*`；`Frame`、`Filters`、`ProductCard`、`EmptyState`、`ErrorState`。
- Produces: `ProductCatalog()`、`SolutionAtlas()`、`ProductSpotlight()`、`LoginForm()`。

- [ ] **Step 1: 迁入 features 并改用视图模型**

三个营销组件从 `components/marketing` 移到 `features/marketing`，字段全部改为 camelCase（`product.manufacturerName` 等）。卡片外层改 ReUI `Frame`，删除本地 `Card` 用法。

- [ ] **Step 2: 产品列表页用 Filters + ProductCard**

`/(marketing)/products` 只保留取数与拼装，筛选（关键词/市场/分类）由 ReUI `Filters` 驱动，列表渲染 `ProductCard`；空结果显示 `EmptyState`。

- [ ] **Step 3: 产品详情页拆分**

`/[id]/page.tsx` ≤ 60 行，内容拆为 `features/marketing/product-detail.tsx`（亮点、画像、话术、案例分区，全部用 `Frame`）。

- [ ] **Step 4: 登录页拆出 LoginForm**

`/login/page.tsx` 只渲染 `LoginForm`；登录方式（本地/LDAP/OIDC）与回调逻辑保持不变，仅替换为 ReUI `Frame` + `FormField` + `Alert`。

- [ ] **Step 5: 验证**

Run: `cd frontend && npx tsc --noEmit && npm run build`
浏览器：首页筛选与跳转、产品列表搜索、产品详情、三种登录方式、错误提示。

- [ ] **Step 6: Commit**

```bash
git add frontend/src/features/marketing frontend/src/features/auth frontend/src/app/\(marketing\) frontend/src/components/marketing
git commit -m "refactor(marketing): move page logic into features and unify on ReUI"
```

---

### Task 16: 前端规范文档

**Files:**
- Modify: `frontend/CLAUDE.md`
- Modify: `README.md`

**Interfaces:**
- Produces: 可执行的目录、组件、令牌、API 层与验收约定。

- [ ] **Step 1: 重写 `frontend/CLAUDE.md`**

保留顶部 `@AGENTS.md` 的 Next.js 提示，补充：目录分层与依赖方向、ReUI 优先顺序、容器/状态/筛选/提示的唯一写法、令牌约束、`lib/api` 工作流、验收命令（`npx tsc --noEmit && npx vitest run && npm run build`）。

- [ ] **Step 2: README 加规范索引**

在「仓库结构」后增加「开发规范」小节，指向 `frontend/CLAUDE.md` 与 `backend/CLAUDE.md`。

- [ ] **Step 3: Commit**

```bash
git add frontend/CLAUDE.md README.md
git commit -m "docs: document frontend conventions"
```

---

### Task 17: 全量验证

**Files:** 无新增。

- [ ] **Step 1: 静态检查**

Run: `cd frontend && npx tsc --noEmit && npx vitest run && npm run lint && npm run build`
Expected: 全部退出码 0；vitest 全绿。

- [ ] **Step 2: 启动前后端做手工回归**

Run: `docker compose up --build -d`（仓库根目录），等待健康后打开 `http://localhost`。
依次验证：登录 → 首页方案筛选 → 产品列表搜索 → 产品详情 → 后台概览 → 产品列表翻页/筛选/删除确认 → 产品编辑保存 → 拓扑增删改 → 用户增删改 → API Key 创建与作废 → 备份导出 → 设置保存。
每页同时检查 375px 窄屏与 `.dark` 深色模式。

- [ ] **Step 3: 记录结果**

把验证结果写进 `docs/superpowers/plans/2026-11-13-frontend-conventions.md` 末尾的「验证记录」小节；失败项不得勾选完成。
