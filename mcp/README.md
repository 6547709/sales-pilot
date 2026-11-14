# Sales-Pilot MCP 使用教程

把 Sales-Pilot 的管理接口包装成 MCP 工具，让 Claude Code / Claude Desktop 等 AI 客户端直接维护产品、服务、场景与拓扑。

## 1. 准备工作

### 1.1 创建 API 密钥

1. 登录后台：`http://<你的域名或 localhost>/admin/api-keys`
2. 填写名称（如 `mcp-server`），可设置过期时间，点击「创建密钥」
3. **完整密钥只显示一次**，形如 `sp_xxxxxxxx...`，立即保存

密钥等价于管理员权限；服务端只存 SHA-256 摘要，丢失只能作废重建。

### 1.2 环境变量

| 变量 | 必填 | 说明 |
|------|------|------|
| `SALES_PILOT_BASE_URL` | 否 | 后端地址，默认 `http://127.0.0.1:8080`；本地开发用 `http://127.0.0.1:18080` |
| `SALES_PILOT_API_KEY` | 是 | 上一步创建的完整密钥 |

## 2. 安装与构建

```bash
cd mcp
npm install
npm run build
```

产物在 `mcp/dist/index.js`，客户端通过 `node dist/index.js` 启动。

## 3. 接入客户端

### Claude Code

```bash
claude mcp add sales-pilot \
  --env SALES_PILOT_BASE_URL=http://127.0.0.1:8080 \
  --env SALES_PILOT_API_KEY=sp_xxx \
  -- node /absolute/path/to/sales-pilot/mcp/dist/index.js
```

### Claude Desktop

编辑配置文件（macOS：`~/Library/Application Support/Claude/claude_desktop_config.json`）：

```json
{
  "mcpServers": {
    "sales-pilot": {
      "command": "node",
      "args": ["/absolute/path/to/sales-pilot/mcp/dist/index.js"],
      "env": {
        "SALES_PILOT_BASE_URL": "http://127.0.0.1:8080",
        "SALES_PILOT_API_KEY": "sp_xxx"
      }
    }
  }
}
```

保存后重启客户端。

## 4. 验证

用官方 Inspector 肉眼确认工具已注册并能调用：

```bash
cd mcp
SALES_PILOT_BASE_URL=http://127.0.0.1:18080 \
SALES_PILOT_API_KEY=sp_xxx \
npx @modelcontextprotocol/inspector node dist/index.js
```

在打开的页面里：

1. 左侧应列出 19 个工具；
2. 选中 `products_list`，参数留空，点击运行；
3. 返回 JSON 数组即表示鉴权与接口都正常。

也可以在 AI 客户端里直接说：「列出 sales-pilot 的所有产品」。

## 5. 工具清单

### 产品

| 工具 | 说明 | 关键参数 |
|------|------|----------|
| `products_list` | 列出产品（含草稿） | `page`、`page_size`、`q`、`status`、`market` |
| `products_get` | 读取单个产品 | `id` |
| `products_create` | 创建产品（默认草稿） | `name`、`category`、`description` 等 |
| `products_update` | 按字段更新 | `id`、`fields`（只传要改的字段） |
| `products_link_scenarios` | 设置产品关联场景 | `id`、`scenario_ids` |
| `products_set_publish` | 草稿 / 发布切换 | `id`、`is_draft` |

### 服务

| 工具 | 说明 | 关键参数 |
|------|------|----------|
| `services_list` | 列出全部服务 | — |
| `services_upsert` | 新建或更新服务 | `id`（更新时）、`slug`、`name`、`kind`、`category_ids` |
| `services_link_categories` | 设置服务所属解决方案 | `id`、`category_ids` |

### 场景

| 工具 | 说明 | 关键参数 |
|------|------|----------|
| `scenarios_list` | 两级场景与关联计数 | — |
| `scenarios_groups` | 场景组及数量 | — |
| `scenarios_upsert` | 新建或更新场景 | `id`、`group_id`、`slug`、`name` |
| `scenarios_link` | 设置场景关联的产品与服务 | `id`、`product_ids`、`service_ids` |

### 拓扑

| 工具 | 说明 |
|------|------|
| `topology_list` | 读取分层、分类与厂商 |
| `topology_layer_upsert` | 新建 / 更新分层 |
| `topology_category_upsert` | 新建 / 更新分类（含客户价值与特点） |
| `topology_vendor_upsert` | 新建 / 更新厂商 |
| `topology_delete` | 删除分层 / 分类 / 厂商（`kind` + `id`） |

### 备份

| 工具 | 说明 |
|------|------|
| `backup_export` | 导出全量备份并返回摘要（版本、产品/用户/服务/场景数量） |

分区 ZIP 备份目前只提供 HTTP 接口（`/admin/backup/export.zip`、`/admin/backup/import-zip`），未包装成 MCP 工具。

## 6. 常见问题

| 现象 | 原因与处理 |
|------|------------|
| `SALES_PILOT_API_KEY 未配置` | 客户端配置里少了 `env`，补上后重启客户端 |
| `API Key 无效或已停用` | 密钥填错或已作废；重新创建 |
| `API Key 已过期` | 密钥超过 `expires_at`；重新创建 |
| 连接超时 / `fetch failed` | `SALES_PILOT_BASE_URL` 不可达；本地开发确认后端端口（18080） |
| 工具返回 403 | 密钥关联用户不是管理员，或账号被停用 |

## 7. 安全

- API 密钥等价于管理员权限，只应在可信环境使用。
- 不要把密钥提交到仓库或写进脚本参数历史。
- 不再使用时在后台「API 密钥」作废。
