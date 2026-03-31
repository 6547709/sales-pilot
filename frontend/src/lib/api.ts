import { API_BASE } from "./config";
import { clearToken, getToken } from "./auth";

export type User = {
  id: number;
  username: string;
  email: string;
  role: string;
  auth_provider: string;
};

/** 管理端：本地用户列表项 */
export type LocalUserRow = {
  id: number;
  username: string;
  email: string;
  role: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
};

/** 全量恢复备份时须在请求中携带的确认令牌（与后端一致） */
export const BACKUP_RESTORE_CONFIRM = "RESTORE_FULL_V1";

export type SolutionCategoryRef = {
  id: number;
  slug: string;
  label: string;
  column_type?: string;
};

export type Product = {
  id: number;
  name: string;
  category: string;
  /** 与首页国内/国外筛选一致：all | domestic | foreign */
  vendor_market?: string;
  is_draft?: boolean;
  solution_category_id?: number | null;
  solution_category?: SolutionCategoryRef | null;
  manufacturer_name?: string;
  manufacturer_logo?: string;
  sales_contact_name?: string;
  sales_contact_phone?: string;
  sales_contact_email?: string;
  presales_contact_name?: string;
  presales_contact_phone?: string;
  presales_contact_email?: string;
  description: string;
  highlights: string | string[];
  target_personas: Record<string, string> | string;
  trigger_events: string;
  discovery_questions: string | string[];
  competitor_analysis: string;
  roi_metrics: string;
  updated_at: string;
};

export type TopologyVendor = {
  id: number;
  name: string;
  sort_order: number;
};

export type TopologyCategoryNode = {
  id: number;
  slug: string;
  label: string;
  icon_key: string;
  keywords: string[];
  hint: string;
  vendors: { domestic: TopologyVendor[]; foreign: TopologyVendor[] };
};

export type TopologyLayerBlock = {
  layer: {
    id: number;
    level: number;
    title: string;
    subtitle: string;
    sort_order: number;
  };
  categories: TopologyCategoryNode[];
};

export type TopologyResponse = {
  security: TopologyCategoryNode[];
  ops: TopologyCategoryNode[];
  central_layers: TopologyLayerBlock[];
};

// 内存缓存，避免重复请求
let topologyCache: TopologyResponse | null = null;
let productsCache: Product[] | null = null;

export async function fetchTopology(): Promise<TopologyResponse> {
  if (topologyCache) return topologyCache;
  const r = await apiFetch("/api/v1/topology");
  if (!r.ok) throw new Error("加载拓扑失败");
  topologyCache = await r.json() as TopologyResponse;
  return topologyCache;
}

export async function fetchProducts(): Promise<Product[]> {
  if (productsCache) return productsCache;
  const r = await apiFetch("/api/v1/products");
  if (!r.ok) throw new Error("加载产品失败");
  productsCache = await r.json() as Product[];
  return productsCache;
}

export function clearProductsCache() {
  productsCache = null;
}

export async function fetchSolutionCategories(
  activeOnly?: boolean,
): Promise<
  {
    id: number;
    slug: string;
    label: string;
    column_type: string;
    layer_id?: number | null;
    is_active: boolean;
  }[]
> {
  const q = activeOnly ? "?active_only=1" : "";
  const r = await apiFetch(`/api/v1/solution-categories${q}`);
  if (!r.ok) return [];
  return r.json();
}

export type SalesScript = {
  id: number;
  product_id: number;
  scenario: string;
  content: string;
};

export type Case = {
  id: number;
  product_id: number;
  client_name: string;
  pain_points: string;
  solution: string;
  value_delivered: string;
};

export async function apiFetch(
  path: string,
  init: RequestInit = {},
  opts?: { token?: string | null },
): Promise<Response> {
  const headers = new Headers(init.headers);
  const token = opts?.token !== undefined ? opts.token : getToken();
  if (token) headers.set("Authorization", `Bearer ${token}`);
  if (
    init.body &&
    typeof init.body === "string" &&
    !headers.has("Content-Type")
  ) {
    headers.set("Content-Type", "application/json");
  }
  return fetch(`${API_BASE}${path}`, { ...init, headers });
}

export async function loginLocal(username: string, password: string) {
  const res = await apiFetch(
    "/api/v1/auth/login",
    {
      method: "POST",
      body: JSON.stringify({ username, password }),
    },
    { token: null },
  );
  if (!res.ok) throw new Error((await res.json()).error || "登录失败");
  const data = await res.json();
  return data as { access_token: string; user: User };
}

export type LoginOptions = {
  local_enabled: boolean;
  ldap_enabled: boolean;
  oidc_enabled: boolean;
};

/** 公开接口：登录页用于隐藏已关闭的登录方式 */
export async function fetchLoginOptions(): Promise<LoginOptions> {
  const r = await apiFetch("/api/v1/auth/login-options", {}, { token: null });
  if (!r.ok) {
    return { local_enabled: true, ldap_enabled: false, oidc_enabled: false };
  }
  return r.json() as Promise<LoginOptions>;
}

export async function loginLdap(username: string, password: string) {
  const res = await apiFetch(
    "/api/v1/auth/login/ldap",
    {
      method: "POST",
      body: JSON.stringify({ username, password }),
    },
    { token: null },
  );
  if (!res.ok) throw new Error((await res.json()).error || "LDAP 登录失败");
  const data = await res.json();
  return data as { access_token: string; user: User };
}

export async function fetchMe(): Promise<User | null> {
  const res = await apiFetch("/api/v1/me");
  if (res.status === 401) {
    clearToken();
    return null;
  }
  if (!res.ok) return null;
  return res.json();
}

export function parseJsonArray(v: unknown): string[] {
  if (Array.isArray(v)) return v.map(String);
  if (typeof v === "string") {
    try {
      const p = JSON.parse(v);
      return Array.isArray(p) ? p.map(String) : [];
    } catch {
      return [];
    }
  }
  return [];
}

export function parsePersonas(v: unknown): Record<string, string> {
  if (v && typeof v === "object" && !Array.isArray(v))
    return v as Record<string, string>;
  if (typeof v === "string") {
    try {
      const p = JSON.parse(v);
      if (p && typeof p === "object" && !Array.isArray(p))
        return p as Record<string, string>;
    } catch {
      /* ignore */
    }
  }
  return {};
}
