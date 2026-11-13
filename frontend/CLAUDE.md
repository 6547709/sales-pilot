@AGENTS.md

# Sales-Pilot 前端开发规范

## 目录分层

```
src/
├── app/                 # 只做路由、取数、拼装（页面 ≤ 60 行）
├── features/            # 页面级组合与状态（products / users / api-keys / settings / backup / topology ...）
├── components/
│   ├── ui/              # shadcn 原语，不含产品知识
│   ├── patterns/        # 通用交互：page-header、empty-state、error-state、loading-state、confirm-dialog、stat-card、form-field
│   ├── domain/          # 业务词汇组件：product-card、market-badge、publish-status-badge、vendor-tag、contact-card
│   ├── admin/           # admin-shell
│   ├── marketing/       # 公开站点组件
│   └── reui/            # ReUI 组件（data-grid / frame / badge / alert / filters ...）
├── lib/api/             # client / types / adapters / endpoints / problem
└── styles/tokens.css    # 设计令牌
```

依赖方向：`app → features → components/domain → components/ui`。`components/ui` 不知道产品，`components/domain` 可以知道。

## UI 组件优先级

1. ReUI（`@/components/reui/*`）
2. shadcn 原语（`@/components/ui/*`）
3. 自行组合（仅前两者都没有时）

硬约束：

- 容器统一 `Frame`，不混用 Card。
- 列表统一 `data-grid`；状态统一 `Badge` 的 `-light` 变体；提示统一 `Alert`。
- 间距用 `gap-*`，不用 `space-y-*`；字号不低于 `text-xs`；动效只用 `transition-colors`；表格行不做 stagger。
- 页面与组件禁止写死色值，只用语义令牌（`primary` / `success` / `warning` / `info` / `muted` 等）。

## 颜色语义

| 颜色 | 含义 |
|------|------|
| 品牌蓝 | 客户可见内容（服务、公开案例、客户价值） |
| 琥珀 | 仅销售可见（联系人、黄金三问、竞品、ROI、话术） |

后台编辑页按可见性分区，列表带可见性列，页顶固定图例。

## API 层

- 页面只 import 视图模型，不直接读 DTO 字段。
- DTO 在 `lib/api/types.ts`（snake_case），视图模型在 `lib/api/adapters.ts`（camelCase）。
- 写回后端用 `toProductDTO` 之类的反向适配器。
- 新增端点：`types.ts` → `adapters.ts` → `endpoints.ts`。

## 验收命令

```bash
npx tsc --noEmit
npx vitest run
npm run lint
npm run build
```

涉及页面的改动还要在浏览器确认：空/加载/错误三态、375px 窄屏、深色模式。

## 双视图与场景

- 视图状态由 `ViewModeProvider` 提供：`useViewMode()` 返回 `mode` 与 `setMode`，URL 使用 `?view=sales|customer`。
- 目录请求统一走 `withView(path, mode)`；不要在组件里手拼 query。
- 客户视图隐藏内部区块：给内部区块加 `data-internal="true"`，根容器加 `view-customer` 类（`globals.css` 已定义）。
- 场景页：`/scenarios` 总览、`/scenarios/[slug]` 详情；零计数场景不展示。
