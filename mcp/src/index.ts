import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import { z } from "zod";
import { api, text } from "./client.js";

const server = new McpServer({ name: "sales-pilot", version: "0.1.0" });

function qs(params: Record<string, unknown>): string {
  const search = new URLSearchParams();
  for (const [k, v] of Object.entries(params)) {
    if (v !== undefined && v !== null && v !== "") search.set(k, String(v));
  }
  const s = search.toString();
  return s ? `?${s}` : "";
}

// ---------- 产品 ----------

server.registerTool(
  "products_list",
  {
    description: "列出产品（含草稿），支持分页与筛选",
    inputSchema: {
      page: z.number().optional(),
      page_size: z.number().optional(),
      q: z.string().optional(),
      status: z.enum(["draft", "published"]).optional(),
      market: z.enum(["domestic", "foreign"]).optional(),
    },
  },
  async (args) =>
    text(await api(`/api/v1/admin/products${qs({ page: args.page ?? 1, page_size: args.page_size ?? 20, q: args.q, status: args.status, market: args.market })}`)),
);

server.registerTool(
  "products_get",
  { description: "读取单个产品的完整字段", inputSchema: { id: z.number() } },
  async (args) => text(await api(`/api/v1/admin/products/${args.id}`)),
);

server.registerTool(
  "products_create",
  {
    description: "创建产品（默认草稿）",
    inputSchema: {
      name: z.string(),
      category: z.string().optional(),
      description: z.string().optional(),
      solution_category_id: z.number().nullable().optional(),
      highlights: z.array(z.string()).optional(),
      vendor_market: z.enum(["all", "domestic", "foreign"]).optional(),
    },
  },
  async (args) =>
    text(
      await api("/api/v1/admin/products", {
        method: "POST",
        body: JSON.stringify({ is_draft: true, ...args }),
      }),
    ),
);

server.registerTool(
  "products_update",
  {
    description: "按字段更新产品；只需传要修改的字段",
    inputSchema: { id: z.number(), fields: z.record(z.unknown()) },
  },
  async (args) => {
    const current = (await api(`/api/v1/admin/products/${args.id}`)) as Record<string, unknown>;
    const merged = { ...current, ...args.fields };
    delete merged.solution_category;
    return text(
      await api(`/api/v1/admin/products/${args.id}`, {
        method: "PUT",
        body: JSON.stringify(merged),
      }),
    );
  },
);

server.registerTool(
  "products_set_publish",
  {
    description: "切换产品草稿/发布状态",
    inputSchema: { id: z.number(), is_draft: z.boolean() },
  },
  async (args) =>
    text(
      await api(`/api/v1/admin/products/${args.id}/draft`, {
        method: "PATCH",
        body: JSON.stringify({ is_draft: args.is_draft }),
      }),
    ),
);

// ---------- 服务 ----------

server.registerTool(
  "services_list",
  { description: "列出全部服务（含停用）", inputSchema: {} },
  async () => text(await api("/api/v1/admin/services")),
);

const serviceShape = {
  id: z.number().optional(),
  slug: z.string(),
  name: z.string(),
  kind: z.enum(["v2v", "backup", "dr", "migration", "upgrade", "ops"]),
  summary: z.string().optional(),
  highlights: z.array(z.string()).optional(),
  deliverables: z.array(z.string()).optional(),
  value_prop: z.string().optional(),
  features: z.array(z.string()).optional(),
  icon_key: z.string().optional(),
  sort_order: z.number().optional(),
  is_active: z.boolean().optional(),
  category_ids: z.array(z.number()).optional(),
};

server.registerTool(
  "services_upsert",
  { description: "新建或更新服务", inputSchema: serviceShape },
  async (args) => {
    const { id, ...body } = args;
    const payload = JSON.stringify({
      ...body,
      highlights: body.highlights ?? [],
      deliverables: body.deliverables ?? [],
      features: body.features ?? [],
      sort_order: body.sort_order ?? 0,
      is_active: body.is_active ?? true,
    });
    if (id) {
      return text(
        await api(`/api/v1/admin/services/${id}`, { method: "PUT", body: payload }),
      );
    }
    return text(await api("/api/v1/admin/services", { method: "POST", body: payload }));
  },
);

server.registerTool(
  "services_link_categories",
  {
    description: "设置服务所属的解决方案分类（决定出现在哪些架构方框）",
    inputSchema: { id: z.number(), category_ids: z.array(z.number()) },
  },
  async (args) =>
    text(
      await api(`/api/v1/admin/services/${args.id}/categories`, {
        method: "PUT",
        body: JSON.stringify({ category_ids: args.category_ids }),
      }),
    ),
);

// ---------- 场景 ----------

server.registerTool(
  "scenarios_list",
  { description: "列出场景组与场景（含关联计数）", inputSchema: {} },
  async () => text(await api("/api/v1/admin/scenarios")),
);

server.registerTool(
  "scenarios_groups",
  { description: "列出场景组及场景数量", inputSchema: {} },
  async () => text(await api("/api/v1/admin/scenario-groups")),
);

server.registerTool(
  "scenarios_upsert",
  {
    description: "新建或更新场景",
    inputSchema: {
      id: z.number().optional(),
      group_id: z.number(),
      slug: z.string(),
      name: z.string(),
      summary: z.string().optional(),
      pain_point: z.string().optional(),
      value_prop: z.string().optional(),
      icon_key: z.string().optional(),
      sort_order: z.number().optional(),
      is_active: z.boolean().optional(),
    },
  },
  async (args) => {
    const { id, ...body } = args;
    const payload = JSON.stringify({
      ...body,
      sort_order: body.sort_order ?? 0,
      is_active: body.is_active ?? true,
    });
    if (id) {
      return text(await api(`/api/v1/admin/scenarios/${id}`, { method: "PUT", body: payload }));
    }
    return text(await api("/api/v1/admin/scenarios", { method: "POST", body: payload }));
  },
);

server.registerTool(
  "scenarios_link",
  {
    description: "设置场景关联的产品与服务",
    inputSchema: { id: z.number(), product_ids: z.array(z.number()), service_ids: z.array(z.number()) },
  },
  async (args) =>
    text(
      await api(`/api/v1/admin/scenarios/${args.id}/links`, {
        method: "PUT",
        body: JSON.stringify({ product_ids: args.product_ids, service_ids: args.service_ids }),
      }),
    ),
);

// ---------- 拓扑 ----------

server.registerTool(
  "topology_list",
  { description: "读取拓扑分层、分类与厂商", inputSchema: {} },
  async () => text(await api("/api/v1/topology")),
);

server.registerTool(
  "topology_layer_upsert",
  {
    description: "新建或更新拓扑分层",
    inputSchema: {
      id: z.number().optional(),
      level: z.number(),
      title: z.string(),
      subtitle: z.string().optional(),
      sort_order: z.number().optional(),
    },
  },
  async (args) => {
    const { id, ...body } = args;
    const payload = JSON.stringify({
      ...body,
      subtitle: body.subtitle ?? "",
      sort_order: body.sort_order ?? 0,
    });
    if (id) return text(await api(`/api/v1/admin/topology/layers/${id}`, { method: "PUT", body: payload }));
    return text(await api("/api/v1/admin/topology/layers", { method: "POST", body: payload }));
  },
);

server.registerTool(
  "topology_category_upsert",
  {
    description: "新建或更新解决方案分类（含客户价值/特点）",
    inputSchema: {
      id: z.number().optional(),
      slug: z.string(),
      label: z.string(),
      column_type: z.enum(["security", "ops", "central"]),
      icon_key: z.string().optional(),
      layer_id: z.number().nullable().optional(),
      sort_order: z.number().optional(),
      hint: z.string().optional(),
      keywords: z.array(z.string()).optional(),
      is_active: z.boolean().optional(),
      customer_value: z.string().optional(),
      customer_features: z.array(z.string()).optional(),
    },
  },
  async (args) => {
    const { id, ...body } = args;
    const payload = JSON.stringify({
      ...body,
      icon_key: body.icon_key ?? "layers",
      sort_order: body.sort_order ?? 0,
      hint: body.hint ?? "",
      keywords: body.keywords ?? [],
      is_active: body.is_active ?? true,
      customer_value: body.customer_value ?? "",
      customer_features: body.customer_features ?? [],
    });
    if (id) return text(await api(`/api/v1/admin/topology/categories/${id}`, { method: "PUT", body: payload }));
    return text(await api("/api/v1/admin/topology/categories", { method: "POST", body: payload }));
  },
);

server.registerTool(
  "topology_vendor_upsert",
  {
    description: "新建或更新分类下的厂商",
    inputSchema: {
      id: z.number().optional(),
      category_id: z.number(),
      market: z.enum(["domestic", "foreign"]),
      name: z.string(),
      sort_order: z.number().optional(),
    },
  },
  async (args) => {
    const { id, ...body } = args;
    const payload = JSON.stringify({ ...body, sort_order: body.sort_order ?? 0 });
    if (id) return text(await api(`/api/v1/admin/topology/vendors/${id}`, { method: "PUT", body: payload }));
    return text(await api("/api/v1/admin/topology/vendors", { method: "POST", body: payload }));
  },
);

server.registerTool(
  "topology_delete",
  {
    description: "删除拓扑分层/分类/厂商",
    inputSchema: { kind: z.enum(["layer", "category", "vendor"]), id: z.number() },
  },
  async (args) => {
    const path =
      args.kind === "layer"
        ? `/api/v1/admin/topology/layers/${args.id}`
        : args.kind === "category"
          ? `/api/v1/admin/topology/categories/${args.id}`
          : `/api/v1/admin/topology/vendors/${args.id}`;
    await api(path, { method: "DELETE" });
    return text({ ok: true });
  },
);

// ---------- 备份 ----------

server.registerTool(
  "backup_export",
  { description: "导出全量备份并返回摘要", inputSchema: {} },
  async () => {
    const payload = (await api("/api/v1/admin/backup/export")) as Record<string, unknown>;
    const count = (key: string) => (Array.isArray(payload[key]) ? (payload[key] as unknown[]).length : 0);
    return text({
      version: payload.version,
      exported_at: payload.exported_at,
      products: count("products"),
      users: count("users"),
      services: count("services"),
      scenarios: count("scenarios"),
    });
  },
);

const transport = new StdioServerTransport();
await server.connect(transport);
