# 备案号运行时配置实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 实现备案号和页脚文字的运行时配置：后端提供 API，前端全局状态管理，三处页面（首页、登录页、管理后台）从 store 读取，管理员可在后台页面修改。

**Architecture:** 后端新增 `SystemSettings` 单例模型，提供公开 GET 和管理 GET/PATCH 接口。前端用 React Context 存储 footer 配置，根 Layout 初始化时获取，三处组件消费该状态。

**Tech Stack:** Go + GORM, Next.js + React Context, Gin

---

## 文件变更概览

| 层级 | 文件 | 操作 |
|------|------|------|
| 后端 | `backend/internal/models/models.go` | 新增 SystemSettings 结构体 |
| 后端 | `backend/internal/handlers/handlers.go` | 新增 3 个路由 |
| 前端 | `frontend/src/lib/footer-config.tsx` | 新建：React Context |
| 前端 | `frontend/src/components/marketing/marketing-auth-layout.tsx` | 获取配置并注入 Context |
| 前端 | `frontend/src/components/marketing/public-shell.tsx` | 从 Context 读取配置 |
| 前端 | `frontend/src/app/(marketing)/login/page.tsx` | 从 Context 读取配置 |
| 前端 | `frontend/src/components/admin/admin-shell.tsx` | 从 Context 读取配置 |
| 前端 | `frontend/src/app/admin/settings/page.tsx` | 新增系统配置 Card |

---

## Task 1: 后端 — 新增 SystemSettings 模型

**Files:**
- Modify: `backend/internal/models/models.go`

- [ ] **Step 1: 在 models.go 末尾追加 SystemSettings 结构体**

在 `AuthSettings` 结构体之后、`}` 之前添加：

```go
// SystemSettings 系统运行时配置（单例 id=1）
type SystemSettings struct {
    ID           uint      `gorm:"primaryKey"`
    FooterText   string    `gorm:"size:255"`
    FilingNumber string    `gorm:"size:64"`
    UpdatedAt    time.Time
}
```

- [ ] **Step 2: Commit**

```bash
git add backend/internal/models/models.go
git commit -m "feat(models): add SystemSettings struct for runtime footer config"
```

---

## Task 2: 后端 — 注册路由

**Files:**
- Modify: `backend/internal/handlers/handlers.go`

- [ ] **Step 1: 在 RegisterRoutes 方法中添加新路由**

在 `api.GET("/health"...` 之后（公开路由区）添加：
```go
api.GET("/config", s.getConfig)
```

在 `admin.GET("/auth-settings"...` 之后（admin 路由区）添加：
```go
admin.GET("/system-settings", s.getSystemSettings)
admin.PATCH("/system-settings", s.patchSystemSettings)
```

- [ ] **Step 2: Commit**

```bash
git add backend/internal/handlers/handlers.go
git commit -m "feat(handlers): register system-settings routes"
```

---

## Task 3: 后端 — 实现 API Handler

**Files:**
- Modify: `backend/internal/handlers/handlers.go`

- [ ] **Step 1: 在 handlers.go 中追加三个 handler 方法**

在文件末尾（任意 handler 方法之后）添加：

```go
// getConfig 公开接口：返回 footer 配置（无需登录）
func (s *Server) getConfig(c *gin.Context) {
    var st models.SystemSettings
    if err := s.DB.FirstOrCreate(&st, models.SystemSettings{ID: 1}).Error; err != nil {
        c.JSON(http.StatusOK, gin.H{"footer_text": "", "filing_number": ""})
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "footer_text":  st.FooterText,
        "filing_number": st.FilingNumber,
    })
}

// getSystemSettings 管理接口：获取完整配置
func (s *Server) getSystemSettings(c *gin.Context) {
    var st models.SystemSettings
    if err := s.DB.FirstOrCreate(&st, models.SystemSettings{ID: 1}).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, st)
}

// patchSystemSettings 管理接口：更新配置
func (s *Server) patchSystemSettings(c *gin.Context) {
    var body struct {
        FooterText   *string `json:"footer_text"`
        FilingNumber *string `json:"filing_number"`
    }
    if err := c.ShouldBindJSON(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
        return
    }
    var st models.SystemSettings
    if err := s.DB.FirstOrCreate(&st, models.SystemSettings{ID: 1}).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    if body.FooterText != nil {
        st.FooterText = *body.FooterText
    }
    if body.FilingNumber != nil {
        st.FilingNumber = *body.FilingNumber
    }
    if err := s.DB.Save(&st).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, st)
}
```

- [ ] **Step 2: Commit**

```bash
git add backend/internal/handlers/handlers.go
git commit -m "feat(api): add getConfig, getSystemSettings, patchSystemSettings handlers"
```

---

## Task 4: 前端 — 创建 FooterConfig Context

**Files:**
- Create: `frontend/src/lib/footer-config.tsx`

- [ ] **Step 1: 创建 footer-config.tsx**

```tsx
"use client";

import { createContext, useContext, useState, useEffect, type ReactNode } from "react";
import { apiFetch } from "@/lib/api";

type FooterConfig = {
  footerText: string;
  filingNumber: string;
};

type FooterConfigContextType = FooterConfig & {
  loading: boolean;
};

const FooterConfigContext = createContext<FooterConfigContextType>({
  footerText: "",
  filingNumber: "",
  loading: true,
});

export function FooterConfigProvider({ children }: { children: ReactNode }) {
  const [config, setConfig] = useState<FooterConfig>({ footerText: "", filingNumber: "" });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    apiFetch("/api/v1/config")
      .then((r) => r.json())
      .then((data: FooterConfig) => setConfig(data))
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  return (
    <FooterConfigContext.Provider value={{ ...config, loading }}>
      {children}
    </FooterConfigContext.Provider>
  );
}

export function useFooterConfig() {
  return useContext(FooterConfigContext);
}
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/lib/footer-config.tsx
git commit -m "feat(frontend): add FooterConfig React context"
```

---

## Task 5: 前端 — 注入 FooterConfigProvider 到根 Layout

**Files:**
- Modify: `frontend/src/components/marketing/marketing-auth-layout.tsx`

- [ ] **Step 1: 在文件顶部 import FooterConfigProvider**

确认 import 区有：
```tsx
import { FooterConfigProvider } from "@/lib/footer-config";
```

- [ ] **Step 2: 用 FooterConfigProvider 包裹 children**

找到：
```tsx
  return <PublicShell>{children}</PublicShell>;
```

改为：
```tsx
  return (
    <FooterConfigProvider>
      <PublicShell>{children}</PublicShell>
    </FooterConfigProvider>
  );
```

注意：`<PublicShell>` 本身已经包含了 `<div className="flex min-h-screen flex-col">` 所以 Provider 放在它外面包裹即可。

- [ ] **Step 3: Commit**

```bash
git add frontend/src/components/marketing/marketing-auth-layout.tsx
git commit -m "feat(frontend): wrap PublicShell with FooterConfigProvider"
```

---

## Task 6: 前端 — 更新 PublicShell 使用 store

**Files:**
- Modify: `frontend/src/components/marketing/public-shell.tsx`

- [ ] **Step 1: 替换 import**

将：
```tsx
import { FOOTER_TEXT, FILING_NUMBER } from "@/lib/config";
```

改为：
```tsx
import { useFooterConfig } from "@/lib/footer-config";
```

- [ ] **Step 2: 在 PublicShell 组件内获取配置**

在 `const pathname = usePathname();` 后添加：
```tsx
  const { footerText, filingNumber } = useFooterConfig();
```

- [ ] **Step 3: 更新 footer 渲染逻辑**

将：
```tsx
{(FOOTER_TEXT || FILING_NUMBER) && (
  <footer className="border-t border-border/60 py-8 text-center text-xs text-muted-foreground">
    {FOOTER_TEXT}
    {FOOTER_TEXT && FILING_NUMBER && <span className="mx-2">|</span>}
    {FILING_NUMBER}
  </footer>
)}
```

改为：
```tsx
{(footerText || filingNumber) && (
  <footer className="border-t border-border/60 py-8 text-center text-xs text-muted-foreground">
    {footerText}
    {footerText && filingNumber && <span className="mx-2">|</span>}
    {filingNumber}
  </footer>
)}
```

- [ ] **Step 4: Commit**

```bash
git add frontend/src/components/marketing/public-shell.tsx
git commit -m "feat(public-shell): use FooterConfig context instead of env vars"
```

---

## Task 7: 前端 — 更新 LoginPage 使用 store

**Files:**
- Modify: `frontend/src/app/(marketing)/login/page.tsx`

- [ ] **Step 1: 替换 import**

将：
```tsx
import { API_BASE, FOOTER_TEXT, FILING_NUMBER } from "@/lib/config";
```

改为：
```tsx
import { API_BASE } from "@/lib/config";
import { useFooterConfig } from "@/lib/footer-config";
```

- [ ] **Step 2: 在 LoginPage 组件内获取配置**

在 `const router = useRouter();` 后添加：
```tsx
  const { footerText, filingNumber } = useFooterConfig();
```

- [ ] **Step 3: 更新 footer 渲染逻辑**

将：
```tsx
          {(FOOTER_TEXT || FILING_NUMBER) && (
            <div className="mt-4 border-t pt-4 text-center text-xs text-muted-foreground">
              {FOOTER_TEXT}
              {FOOTER_TEXT && FILING_NUMBER && <span className="mx-2">|</span>}
              {FILING_NUMBER}
            </div>
          )}
```

改为：
```tsx
          {(footerText || filingNumber) && (
            <div className="mt-4 border-t pt-4 text-center text-xs text-muted-foreground">
              {footerText}
              {footerText && filingNumber && <span className="mx-2">|</span>}
              {filingNumber}
            </div>
          )}
```

- [ ] **Step 4: Commit**

```bash
git add "frontend/src/app/(marketing)/login/page.tsx"
git commit -m "feat(login): use FooterConfig context instead of env vars"
```

---

## Task 8: 前端 — 更新 AdminShell 使用 store

**Files:**
- Modify: `frontend/src/components/admin/admin-shell.tsx`

- [ ] **Step 1: 替换 import**

将：
```tsx
import { FOOTER_TEXT, FILING_NUMBER } from "@/lib/config";
```

改为：
```tsx
import { useFooterConfig } from "@/lib/footer-config";
```

- [ ] **Step 2: 在 AdminShell 组件内获取配置**

在 `const pathname = usePathname();` 后添加：
```tsx
  const { footerText, filingNumber } = useFooterConfig();
```

- [ ] **Step 3: 更新 footer 渲染逻辑**

将：
```tsx
      {(FOOTER_TEXT || FILING_NUMBER) && (
        <div className="fixed bottom-16 left-0 right-0 z-40 border-t bg-muted/50 py-2 text-center text-xs text-muted-foreground md:hidden">
          {FOOTER_TEXT}
          {FOOTER_TEXT && FILING_NUMBER && <span className="mx-2">|</span>}
          {FILING_NUMBER}
        </div>
      )}
```

改为：
```tsx
      {(footerText || filingNumber) && (
        <div className="fixed bottom-16 left-0 right-0 z-40 border-t bg-muted/50 py-2 text-center text-xs text-muted-foreground md:hidden">
          {footerText}
          {footerText && filingNumber && <span className="mx-2">|</span>}
          {filingNumber}
        </div>
      )}
```

- [ ] **Step 4: Commit**

```bash
git add frontend/src/components/admin/admin-shell.tsx
git commit -m "feat(admin): use FooterConfig context instead of env vars"
```

---

## Task 9: 前端 — 管理页面新增系统配置 Card

**Files:**
- Modify: `frontend/src/app/admin/settings/page.tsx`

- [ ] **Step 1: 在文件顶部 import useFooterConfig**

```tsx
import { useFooterConfig } from "@/lib/footer-config";
```

- [ ] **Step 2: 在组件内获取 store 刷新函数**

添加：
```tsx
  const { footerText, filingNumber } = useFooterConfig();
```

（这个页面本身不需要消费配置，只需要在保存后刷新全局状态。但 React Context 不提供强制刷新方法，所以这里用 apiFetch 直接读取当前值来初始化表单。）

- [ ] **Step 3: 添加 state 和初始 fetch**

在 `const [msg, setMsg] = useState("");` 后添加：
```tsx
  const [footerText, setFooterText] = useState("");
  const [filingNumber, setFilingNumber] = useState("");
```

在 `useEffect` 末尾添加获取配置的调用：
```tsx
    apiFetch("/api/v1/config").then(async (r) => {
      if (!r.ok) return;
      const data = await r.json();
      setFooterText(data.footer_text || "");
      setFilingNumber(data.filing_number || "");
    });
```

- [ ] **Step 4: 添加保存函数**

在 `save` 函数内，在 `const res = await apiFetch("/api/v1/admin/auth-settings"...` 之前添加：
```tsx
    // 先保存系统配置
    await apiFetch("/api/v1/admin/system-settings", {
      method: "PATCH",
      body: JSON.stringify({
        footer_text: footerText,
        filing_number: filingNumber,
      }),
    });
```

- [ ] **Step 5: 在 return JSX 中新增 Card**

在最后一个 `<Card>`（OIDC Card）之后、`{msg ? ...}` 之前添加：

```tsx
      <Card>
        <CardHeader>
          <CardTitle>系统配置</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="ft">页脚文字</Label>
            <Input
              id="ft"
              value={footerText}
              onChange={(e) => setFooterText(e.target.value)}
              placeholder="如：企业销售赋能"
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="fn">备案号</Label>
            <Input
              id="fn"
              value={filingNumber}
              onChange={(e) => setFilingNumber(e.target.value)}
              placeholder="如：京ICP备xxxx号"
            />
          </div>
          <p className="text-xs text-muted-foreground">
            将在首页、登录页、管理后台底部显示。
          </p>
        </CardContent>
      </Card>
```

- [ ] **Step 6: Commit**

```bash
git add frontend/src/app/admin/settings/page.tsx
git commit -m "feat(admin): add system settings card to settings page"
```

---

## 验证检查

1. 启动后端和前端
2. 访问 `/admin/settings`，在"系统配置"卡片中填写页脚文字和备案号，保存
3. 访问首页，确认 footer 显示正确
4. 访问 `/login`，确认登录页底部显示正确
5. 调整浏览器宽度到手机尺寸，访问管理后台，确认移动端 footer 显示正确
6. 清空两个字段，保存后确认 footer 区域消失
