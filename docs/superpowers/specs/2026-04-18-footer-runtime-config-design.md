# 备案号与页脚运行时配置设计

## 背景

当前 `NEXT_PUBLIC_` 环境变量在 Docker 构建时被 baked into 镜像，无法在运行时动态修改。本设计改为运行时从后端 API 读取配置，管理员可在后台页面修改，无需重新构建镜像。

## 数据模型

### SystemSettings（单例，id=1）

```go
type SystemSettings struct {
    ID           uint      `gorm:"primaryKey"`
    FooterText   string    `gorm:"size:255"`   // 页脚文字
    FilingNumber string    `gorm:"size:64"`    // 备案号
    UpdatedAt    time.Time
}
```

## API 设计

### 公开接口（无需认证）

| 方法 | 路径 | 响应 |
|------|------|------|
| GET | `/api/v1/config` | `{"footer_text":"", "filing_number":""}` |

### 管理接口（需 admin 认证）

| 方法 | 路径 | 响应 |
|------|------|------|
| GET | `/api/v1/admin/system-settings` | 完整 SystemSettings |
| PATCH | `/api/v1/admin/system-settings` | 更新后完整 SystemSettings |

#### PATCH 请求体

```json
{
  "footer_text": "企业销售赋能",
  "filing_number": "冀ICP备2022027603号-1"
}
```

## 前端架构

### 状态管理

新建 `frontend/src/lib/footer-config.ts`（基于 React Context / zustand）：

```ts
type FooterConfig = {
  footerText: string;
  filingNumber: string;
}
```

根 Layout（`MarketingAuthLayout`）初始化时调用 `GET /api/v1/config`，将结果存入全局状态。三处消费组件（PublicShell、LoginPage、AdminShell）从 store 读取。

### 管理页面

在 `/admin/settings` 底部新增 Card：

```tsx
<Card>
  <CardHeader>
    <CardTitle>系统配置</CardTitle>
  </CardHeader>
  <CardContent className="space-y-4">
    <div className="space-y-2">
      <Label htmlFor="ft">页脚文字</Label>
      <Input id="ft" value={footerText} onChange={...} />
    </div>
    <div className="space-y-2">
      <Label htmlFor="fn">备案号</Label>
      <Input id="fn" value={filingNumber} onChange={...} />
    </div>
    <Button onClick={save}>保存</Button>
  </CardContent>
</Card>
```

## 变更文件清单

| 层级 | 文件 | 操作 |
|------|------|------|
| 后端 | `backend/internal/models/models.go` | 新增 `SystemSettings` 结构体 |
| 后端 | `backend/internal/handlers/handlers.go` | 新增 3 个路由和方法 |
| 前端 | `frontend/src/lib/footer-config.ts` | 新建：footer 配置 context |
| 前端 | `frontend/src/components/marketing/marketing-auth-layout.tsx` | 调用 API 初始化配置 |
| 前端 | `frontend/src/components/marketing/public-shell.tsx` | 改用 store 中的配置 |
| 前端 | `frontend/src/app/(marketing)/login/page.tsx` | 改用 store 中的配置 |
| 前端 | `frontend/src/components/admin/admin-shell.tsx` | 改用 store 中的配置 |
| 前端 | `frontend/src/app/admin/settings/page.tsx` | 新增系统配置 Card |

## 显示组合

| footerText | filingNumber | 效果 |
|-----------|-------------|------|
| 有 | 有 | `企业销售赋能 \| 冀ICP备2022027603号-1` |
| 有 | 无 | `企业销售赋能` |
| 无 | 有 | `冀ICP备2022027603号-1` |
| 无 | 无 | 不渲染 footer |

## 数据库迁移

GORM AutoMigrate 自动创建 `system_settings` 表（单例 id=1 由 FirstOrCreate 保证）。

## 后续清理

原 `docker-publish.yml` 中添加的 `NEXT_PUBLIC_FOOTER_TEXT` 和 `NEXT_PUBLIC_FILING_NUMBER` build-args 应移除（等本功能上线后）。
