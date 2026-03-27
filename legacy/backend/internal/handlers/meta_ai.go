package handlers

import (
	"fmt"
	"net/http"

	"github.com/6547709/sales-pilot/backend/internal/auth"
	"github.com/gin-gonic/gin"
)

// adminProductMaintainMeta 供 AI / MCP / OpenClaw 等拉取「产品维护」相关字段与端点说明（机器可读）
func (s *Server) adminProductMaintainMeta(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"service":     "sales-pilot",
		"api_version": "v1",
		"base_url":    "将 BASE 替换为部署根地址，如 https://your-domain 或 http://localhost:8080",
		"prefix":      "/api/v1",
		"authentication": gin.H{
			"jwt": gin.H{
				"login_path": "POST /api/v1/auth/login",
				"ttl":        fmt.Sprintf("默认 %v（auth.TokenTTL）", auth.TokenTTL),
				"header":     "Authorization: Bearer <access_token>",
			},
			"api_key": gin.H{
				"description": "长期密钥，适合 MCP/脚本；由管理员在后台「API 密钥」创建",
				"header":      "X-API-Key: <完整密钥，仅创建时明文展示一次>",
				"storage":     "服务端仅存 SHA256(密钥)，无法找回明文",
				"expires":     "创建时可填 expires_at（RFC3339）；省略则永不过期，也可在后台作废",
			},
			"manage_keys": gin.H{
				"list":   "GET /api/v1/admin/api-keys",
				"create": "POST /api/v1/admin/api-keys  body: {name, expires_at?}",
				"revoke": "DELETE /api/v1/admin/api-keys/:id",
			},
		},
		"discovery_before_create": []gin.H{
			{
				"purpose": "获取「所属解决方案」可选 ID，写入产品的 solution_category_id",
				"method":  "GET",
				"path":    "/api/v1/solution-categories",
				"query":   "可选 active_only=1 仅启用项；管理端建议不传以含停用项",
				"response_each": gin.H{
					"id":    "number，写入 solution_category_id",
					"label": "展示名称；保存产品时若带 solution_category_id，后台会把 category 同步为该名称",
					"slug":  "英文唯一标识",
				},
			},
		},
		"product_resource": gin.H{
			"description": "销售赋能「方案/产品」一条记录，前台方案库与详情页的数据来源；is_draft=true 时仅管理端可见。",
			"fields": []gin.H{
				{"json": "id", "type": "number", "read_only": true, "zh": "主键，创建后由系统分配"},
				{"json": "name", "type": "string", "required_on_create": true, "zh": "方案/产品标题，前台列表与详情主标题", "hint": "简洁、可搜索，避免内部项目代号"},
				{"json": "category", "type": "string", "zh": "分类展示文案", "hint": "若填写 solution_category_id，保存时后端会用对应解决方案名称覆盖"},
				{"json": "solution_category_id", "type": "number|null", "zh": "关联首页架构「解决方案」节点 ID", "hint": "先 GET /solution-categories 选 id；与前台拓扑筛选一致"},
				{"json": "manufacturer_name", "type": "string", "zh": "厂商/品牌名", "hint": "可与拓扑卡片厂商区分：此处为方案级签约主体"},
				{"json": "sales_contact_name", "type": "string", "zh": "销售联系人姓名"},
				{"json": "sales_contact_phone", "type": "string", "zh": "销售电话"},
				{"json": "sales_contact_email", "type": "string", "zh": "销售邮箱"},
				{"json": "presales_contact_name", "type": "string", "zh": "售前联系人姓名"},
				{"json": "presales_contact_phone", "type": "string", "zh": "售前电话"},
				{"json": "presales_contact_email", "type": "string", "zh": "售前邮箱"},
				{"json": "description", "type": "string", "zh": "概览正文", "format": "Markdown（推荐 GFM）", "hint": "价值主张、适用场景、与架构关系；支持表格、列表"},
				{"json": "highlights", "type": "array[string]", "zh": "三大亮点", "hint": "通常 3 条短句，每条一个核心利益点；JSON 数组"},
				{"json": "discovery_questions", "type": "array[string]", "zh": "黄金三问（需求挖掘）", "hint": "3 条左右开放式问题，引导客户暴露痛点"},
				{"json": "target_personas", "type": "object", "zh": "目标画像", "hint": "JSON 对象：键为角色名（如 CEO、CIO），值为一句话关注点；与后台预设角色兼容"},
				{"json": "trigger_events", "type": "string", "zh": "触发事件/采购契机", "hint": "多行纯文本，如涨价、信创、容灾演练等"},
				{"json": "competitor_analysis", "type": "string", "zh": "竞品分析", "format": "Markdown", "hint": "对比维度、差异化、埋雷话术；可用表格"},
				{"json": "roi_metrics", "type": "string", "zh": "ROI 与效果指标", "hint": "可量化表述，如 TCO、交付周期区间"},
				{"json": "vendor_market", "type": "string", "zh": "厂商分区", "hint": "all|domestic|foreign；与首页国内/国外筛选一致"},
				{"json": "is_draft", "type": "boolean", "zh": "草稿", "hint": "true=前台不展示；完稿后改 false 或 PATCH .../draft"},
			},
		},
		"related_sales_script": gin.H{
			"description": "按场景拆分的销售话术片段，多条挂在一个 product_id 下",
			"fields": []gin.H{
				{"json": "scenario", "type": "string", "zh": "场景标题，如「开场破冰」「异议：价格」"},
				{"json": "content", "type": "string", "zh": "话术正文", "hint": "口语化、可照读或改编"},
			},
			"create": "POST /api/v1/admin/products/:productId/scripts",
		},
		"related_case": gin.H{
			"description": "客户案例，用于建立信任",
			"fields": []gin.H{
				{"json": "client_name", "type": "string", "zh": "客户名称或脱敏称谓"},
				{"json": "pain_points", "type": "string", "zh": "痛点"},
				{"json": "solution", "type": "string", "zh": "我方方案要点"},
				{"json": "value_delivered", "type": "string", "zh": "交付价值/结果"},
			},
			"create": "POST /api/v1/admin/products/:productId/cases",
		},
		"admin_endpoints_products": []gin.H{
			{"method": "GET", "path": "/api/v1/admin/products", "zh": "全部产品（含草稿），含 solution_category 预加载"},
			{"method": "GET", "path": "/api/v1/admin/products/:id", "zh": "单条详情（含草稿）"},
			{"method": "POST", "path": "/api/v1/admin/products", "zh": "创建；body 为 Product JSON，勿嵌套 solution_category 对象"},
			{"method": "PUT", "path": "/api/v1/admin/products/:id", "zh": "全量更新常用字段；请求体同创建，需带 id 一致"},
			{"method": "PATCH", "path": "/api/v1/admin/products/:id/draft", "zh": "仅切换草稿", "body": gin.H{"is_draft": true}},
			{"method": "DELETE", "path": "/api/v1/admin/products/:id", "zh": "删除产品及下属话术、案例"},
			{"method": "GET", "path": "/api/v1/admin/products/:id/cases", "zh": "案例列表（草稿产品也可用）"},
			{"method": "GET", "path": "/api/v1/admin/products/:id/scripts", "zh": "话术列表"},
		},
		"suggested_agent_workflow": []string{
			"1. 使用 X-API-Key（推荐自动化）或 POST /auth/login 获取 JWT",
			"2. GET /solution-categories 选定 solution_category_id（可选）",
			"3. POST /admin/products 创建 is_draft=true 的骨架（name、description 可先填）",
			"4. PUT /admin/products/:id 迭代补全 highlights、discovery_questions、target_personas 等",
			"5. POST /admin/products/:id/cases 与 scripts 补充案例与话术",
			"6. PATCH /admin/products/:id/draft 设置 is_draft=false 发布",
		},
		"notes": []string{
			"所有 /admin/* 需管理员身份：JWT Bearer 或有效 X-API-Key",
			"PUT 更新时不要发送嵌套 solution_category，只发送 solution_category_id",
			"GET /products、/topology 等需已登录：Authorization Bearer JWT 或 X-API-Key；不含草稿产品",
			"详细人类可读说明见仓库 docs/api-ai-product-maintenance.md",
		},
	})
}
