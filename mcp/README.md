# Sales-Pilot MCP 服务

把 Sales-Pilot 的管理接口包装成 MCP 工具，供 AI 客户端维护产品、服务、场景与拓扑。

## 环境变量

| 变量 | 说明 |
|------|------|
| `SALES_PILOT_BASE_URL` | 后端地址，默认 `http://127.0.0.1:8080` |
| `SALES_PILOT_API_KEY` | 管理员 API 密钥（后台「API 密钥」创建） |

## 安装与运行

```bash
cd mcp
npm install
npm run build
SALES_PILOT_BASE_URL=http://127.0.0.1:8080 \
SALES_PILOT_API_KEY=sp_xxx \
npm start
```

开发模式：`npm run dev`。

## 客户端配置

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

## 工具

- 产品：`products_list`、`products_get`、`products_create`、`products_update`、`products_set_publish`
- 服务：`services_list`、`services_upsert`、`services_link_categories`
- 场景：`scenarios_list`、`scenarios_groups`、`scenarios_upsert`、`scenarios_link`
- 拓扑：`topology_list`、`topology_layer_upsert`、`topology_category_upsert`、`topology_vendor_upsert`、`topology_delete`
- 备份：`backup_export`

## 安全

API 密钥等价于管理员权限，只应在可信环境使用。服务端仅存密钥摘要，丢失后需作废重建。
