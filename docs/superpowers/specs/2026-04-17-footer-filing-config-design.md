# 备案号与页脚配置设计

## 背景

当前应用在登录页、首页/方案库、管理后台均无备案号和页脚文字配置，现需要通过环境变量让运维可在部署时统一配置。

## 环境变量

在 `frontend/.env.example` 中新增：

```env
# 页脚文字，可选，如"企业销售赋能"
NEXT_PUBLIC_FOOTER_TEXT=

# 备案号，可选，如"京ICP备xxxx号"
NEXT_PUBLIC_FILING_NUMBER=
```

采用 `NEXT_PUBLIC_` 前缀：两者均为公开信息，不涉及敏感数据，可随环境变量变更而生效（无需重新构建）。

## 配置层

`frontend/src/lib/config.ts` 追加导出：

```ts
export const FOOTER_TEXT = process.env.NEXT_PUBLIC_FOOTER_TEXT;
export const FILING_NUMBER = process.env.NEXT_PUBLIC_FILING_NUMBER;
```

## 消费位置与显示逻辑

### 1. PublicShell（首页/方案库）

文件：`frontend/src/components/marketing/public-shell.tsx`

当前 footer：
```tsx
<footer className="border-t border-border/60 py-8 text-center text-xs text-muted-foreground">
  <p>Sales-Pilot · 企业销售赋能</p>
</footer>
```

改为：仅当 `FOOTER_TEXT` 或 `FILING_NUMBER` 存在时渲染；两者都存在时分两行显示。

### 2. LoginPage（登录页）

文件：`frontend/src/app/(marketing)/login/page.tsx`

在登录卡片底部加一行 footer：
- 仅 `FOOTER_TEXT` 存在 → 显示文字
- 仅 `FILING_NUMBER` 存在 → 显示备案号
- 两者都存在 → 分行显示

### 3. AdminShell（管理后台）

文件：`frontend/src/components/admin/admin-shell.tsx`

在移动端底部导航栏上方加一行 footer（`h-16 md:hidden` 占位块之上），显示逻辑同上。

## 变更文件清单

| 文件 | 操作 |
|------|------|
| `frontend/.env.example` | 新增两个环境变量说明 |
| `frontend/src/lib/config.ts` | 新增 `FOOTER_TEXT` 和 `FILING_NUMBER` 导出 |
| `frontend/src/components/marketing/public-shell.tsx` | Footer 改为条件渲染 |
| `frontend/src/app/(marketing)/login/page.tsx` | 登录卡片底部新增 footer 信息 |
| `frontend/src/components/admin/admin-shell.tsx` | 移动端底部导航上方新增 footer |

## 显示组合

| FOOTER_TEXT | FILING_NUMBER | 效果 |
|-------------|---------------|------|
| 有 | 有 | 两行 |
| 有 | 无 | 仅文字 |
| 无 | 有 | 仅有备案号 |
| 无 | 无 | 不渲染 footer 区域 |
