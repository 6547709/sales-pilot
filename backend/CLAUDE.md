# Sales-Pilot 后端开发规范

## 分层与依赖方向

```
cmd/server → httpapi → service → repository → models
```

- `repository/` 只做 GORM 查询，不写业务规则，不 import `service` / `gin`。
- `service/` 承载业务规则与事务边界，不 import `gin`；错误一律返回 `apperr`。
- `httpapi/` 只做参数绑定、调用 service、投影与响应；handler 不直接访问 `*gorm.DB`。
- `models/` 只有数据结构；迁移清单集中在 `internal/database/migrations`。

## 统一响应

成功：保持资源原样（对象或数组）。

错误：固定形状，前端只依赖 `error` 字段。

```json
{ "error": "产品不存在", "code": "product_not_found" }
```

- `httpapi.Fail(c, err)` 负责把 `apperr.Kind` 映射为 HTTP 状态码。
- 500 一律返回「服务器内部错误」，原始错误只写日志。
- 路径参数用 `httpapi.ParseID`，非法返回 400 + `invalid_id`。

常用错误码：`invalid_body`、`invalid_id`、`unauthorized`、`forbidden`、`product_not_found`、`last_admin`、`api_key_expired`、`scenario_group_limit`。

## 分页

- 列表端点默认保持裸数组；只有带 `page` 参数时返回分页对象：

```json
{ "items": [], "total": 0, "page": 1, "page_size": 20 }
```

- `page_size` 默认 20、上限 100；排序字段走白名单（`id` / `name` / `updated_at` / `created_at`）。

## 新增端点步骤

1. 在 `models` 加结构 → `migrations.Migrate` 注册。
2. 在 `repository` 加查询；接口方法放在 `service/ports.go`。
3. 在 `service` 加业务方法，返回 `apperr`。
4. 在 `httpapi` 加 handler，并在 `router.go` 注册路由。
5. 写测试：业务规则进 `service/*_test.go`，路由契约进 `httpapi/router_test.go`。

## 测试要求

- service 层测试用 fake repository（嵌入接口 + 覆写用到的方法）。
- 路由测试覆盖：401、400 `invalid_id`、分页两种响应形状。
- 交付前必须通过：

```bash
go build ./... && go vet ./... && go test ./...
```
