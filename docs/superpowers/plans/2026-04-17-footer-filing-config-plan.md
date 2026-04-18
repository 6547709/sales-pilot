# Footer 与备案号配置实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 通过环境变量 `NEXT_PUBLIC_FOOTER_TEXT` 和 `NEXT_PUBLIC_FILING_NUMBER` 配置页脚文字与备案号，在登录页、首页/方案库、管理后台三处显示。

**Architecture:** 在 `config.ts` 统一封装两个 env 变量为导出常量，三处组件各自条件渲染。均为 `NEXT_PUBLIC_` 前缀，因其为公开信息无需重构建即可生效。

**Tech Stack:** Next.js (NextPublic env), TypeScript, Tailwind

---

## 文件变更概览

| 文件 | 操作 |
|------|------|
| `frontend/.env.example` | 修改：新增两个环境变量说明 |
| `frontend/src/lib/config.ts` | 修改：新增 `FOOTER_TEXT` 和 `FILING_NUMBER` 导出 |
| `frontend/src/components/marketing/public-shell.tsx` | 修改：footer 改为条件渲染 |
| `frontend/src/app/(marketing)/login/page.tsx` | 修改：登录卡片底部新增 footer 信息 |
| `frontend/src/components/admin/admin-shell.tsx` | 修改：移动端底部导航上方新增 footer |

---

## Task 1: 更新 .env.example

**Files:**
- Modify: `frontend/.env.example`

- [ ] **Step 1: 在 .env.example 末尾追加环境变量说明**

在 `# NEXT_PUBLIC_API_URL=http://localhost:8080` 之后追加：

```env
# -----------------------------------------------------------------------------
# 页脚与备案号配置（可选，部署时按需设置）
# -----------------------------------------------------------------------------
# NEXT_PUBLIC_FOOTER_TEXT：页脚文字，如"企业销售赋能"。不设置则不显示。
# NEXT_PUBLIC_FILING_NUMBER：备案号，如"京ICP备xxxx号"。不设置则不显示。
```

- [ ] **Step 2: Commit**

```bash
git add frontend/.env.example
git commit -m "docs: add NEXT_PUBLIC_FOOTER_TEXT and FILING_NUMBER to .env.example"
```

---

## Task 2: 更新 config.ts 追加配置导出

**Files:**
- Modify: `frontend/src/lib/config.ts`

- [ ] **Step 1: 在 config.ts 末尾追加导出**

在 `export const API_BASE = resolveApiBase();` 后追加：

```ts
/**
 * 页脚文字，可选。如 "企业销售赋能"。
 */
export const FOOTER_TEXT = process.env.NEXT_PUBLIC_FOOTER_TEXT;

/**
 * 备案号，可选。如 "京ICP备xxxx号"。
 */
export const FILING_NUMBER = process.env.NEXT_PUBLIC_FILING_NUMBER;
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/lib/config.ts
git commit -m "feat(config): add FOOTER_TEXT and FILING_NUMBER exports"
```

---

## Task 3: 更新 PublicShell footer 条件渲染

**Files:**
- Modify: `frontend/src/components/marketing/public-shell.tsx`

- [ ] **Step 1: 在文件顶部 import config**

确认 import 区已有：
```ts
import { FOOTER_TEXT, FILING_NUMBER } from "@/lib/config";
```
（如果尚未导入则添加）

- [ ] **Step 2: 替换 footer 区域**

将当前的：
```tsx
<footer className="border-t border-border/60 py-8 text-center text-xs text-muted-foreground">
  <p>Sales-Pilot · 企业销售赋能</p>
</footer>
```

替换为：
```tsx
{(FOOTER_TEXT || FILING_NUMBER) && (
  <footer className="border-t border-border/60 py-8 text-center text-xs text-muted-foreground">
    {FOOTER_TEXT && <p>{FOOTER_TEXT}</p>}
    {FILING_NUMBER && <p>{FILING_NUMBER}</p>}
  </footer>
)}
```

- [ ] **Step 3: Commit**

```bash
git add frontend/src/components/marketing/public-shell.tsx
git commit -m "feat(public-shell): make footer conditional on FOOTER_TEXT/FILING_NUMBER env"
```

---

## Task 4: 更新 LoginPage 添加 footer 信息

**Files:**
- Modify: `frontend/src/app/(marketing)/login/page.tsx`

- [ ] **Step 1: 在文件顶部 import config**

确认 import 区已有：
```ts
import { FOOTER_TEXT, FILING_NUMBER } from "@/lib/config";
```
（如果尚未导入则添加）

- [ ] **Step 2: 在登录卡片底部 TabsContent 之后、CardContent 结束标签之前添加 footer**

在 `</CardContent>` 之前找到合适位置（紧跟在 `</Tabs>` 或登录方式提示文字之后、`</CardContent>` 之前），插入：

```tsx
{(FOOTER_TEXT || FILING_NUMBER) && (
  <div className="mt-4 border-t pt-4 text-center text-xs text-muted-foreground">
    {FOOTER_TEXT && <p>{FOOTER_TEXT}</p>}
    {FILING_NUMBER && <p>{FILING_NUMBER}</p>}
  </div>
)}
```

具体位置在：
```tsx
          <p className="text-center text-xs text-muted-foreground">
            登录后可访问首页与方案库等企业内容。
          </p>
        </CardContent>
```
改为在 `</CardContent>` 之前插入上述 footer div。

- [ ] **Step 3: Commit**

```bash
git add frontend/src/app/\(marketing\)/login/page.tsx
git commit -m "feat(login): add conditional footer info based on env vars"
```

---

## Task 5: 更新 AdminShell 移动端 footer

**Files:**
- Modify: `frontend/src/components/admin/admin-shell.tsx`

- [ ] **Step 1: 在文件顶部 import config**

在 import 区添加：
```ts
import { FOOTER_TEXT, FILING_NUMBER } from "@/lib/config";
```

- [ ] **Step 2: 在移动端底部导航上方添加 footer**

找到：
```tsx
      {/* 移动端底部导航 */}
      <nav className="fixed bottom-0 left-0 right-0 z-40 flex border-t bg-card/95 backdrop-blur md:hidden">
```

在其上方插入 footer：
```tsx
      {(FOOTER_TEXT || FILING_NUMBER) && (
        <div className="fixed bottom-16 left-0 right-0 z-40 border-t bg-muted/50 py-2 text-center text-xs text-muted-foreground md:hidden">
          {FOOTER_TEXT && <p>{FOOTER_TEXT}</p>}
          {FILING_NUMBER && <p>{FILING_NUMBER}</p>}
        </div>
      )}
      {/* 移动端底部导航 */}
      <nav className="fixed bottom-0 left-0 right-0 z-40 flex border-t bg-card/95 backdrop-blur md:hidden">
```

- [ ] **Step 3: Commit**

```bash
git add frontend/src/components/admin/admin-shell.tsx
git commit -m "feat(admin): add conditional footer to mobile admin view"
```

---

## 验证检查

实现完成后，在 `.env.local` 中设置：
```env
NEXT_PUBLIC_FOOTER_TEXT=企业销售赋能
NEXT_PUBLIC_FILING_NUMBER=京ICP备12345678号
```

运行 `npm run dev` 验证：
- [ ] 首页 footer 显示两行文字
- [ ] 登录页卡片底部显示 footer
- [ ] 管理后台移动端（resize 到手机宽度）底部导航上方显示 footer
- [ ] 注释掉任一变量后对应行不显示
- [ ] 两个都不设置时 footer 区域完全消失
