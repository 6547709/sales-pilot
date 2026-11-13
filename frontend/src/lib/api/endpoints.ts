import { apiFetch, apiJson } from "./client";
import {
  toApiKeyCreated,
  toApiKeys,
  toCases,
  toLocalUser,
  toLocalUsers,
  toLoginOptions,
  toPageResult,
  toProduct,
  toProductDTO,
  toProducts,
  toSalesScripts,
  toUser,
  type ApiKeyCreateResult,
  type ApiKeyRow,
  type Case,
  type LocalUserRow,
  type LoginOptions,
  type PageResult,
  type Product,
  type SalesScript,
  type SolutionCategory,
  type TopologyResponse,
  type User,
} from "./adapters";
import type {
  ApiKeyCreateDTO,
  ApiKeyDTO,
  AuthSettingsDTO,
  CaseDTO,
  LocalUserDTO,
  LoginOptionsDTO,
  ProductDTO,
  SalesScriptDTO,
  SolutionCategoryDTO,
  SystemSettingsDTO,
  TopologyDTO,
  UserDTO,
} from "./types";

/** 全量恢复备份时须在请求中携带的确认令牌（与后端一致） */
export const BACKUP_RESTORE_CONFIRM = "RESTORE_FULL_V1";

// 内存缓存，避免重复请求
let topologyCache: TopologyResponse | null = null;
let productsCache: Product[] | null = null;

export async function fetchTopology(): Promise<TopologyResponse> {
  if (topologyCache) return topologyCache;
  const data = await apiJson<TopologyDTO>("/api/v1/topology");
  topologyCache = data;
  return topologyCache;
}

export async function fetchProducts(): Promise<Product[]> {
  if (productsCache) return productsCache;
  const data = await apiJson<ProductDTO[]>("/api/v1/products");
  productsCache = toProducts(data);
  return productsCache;
}

export function clearProductsCache() {
  productsCache = null;
}

export async function fetchSolutionCategories(
  activeOnly?: boolean,
): Promise<SolutionCategory[]> {
  const q = activeOnly ? "?active_only=1" : "";
  const res = await apiFetch(`/api/v1/solution-categories${q}`);
  if (!res.ok) return [];
  return (await res.json()) as SolutionCategoryDTO[];
}

export async function loginLocal(username: string, password: string) {
  const data = await apiJson<{ access_token: string; user: UserDTO }>(
    "/api/v1/auth/login",
    { method: "POST", body: JSON.stringify({ username, password }) },
    { token: null },
  );
  return { accessToken: data.access_token, user: toUser(data.user) };
}

export async function loginLdap(username: string, password: string) {
  const data = await apiJson<{ access_token: string; user: UserDTO }>(
    "/api/v1/auth/login/ldap",
    { method: "POST", body: JSON.stringify({ username, password }) },
    { token: null },
  );
  return { accessToken: data.access_token, user: toUser(data.user) };
}

export async function fetchLoginOptions(): Promise<LoginOptions> {
  const res = await apiFetch("/api/v1/auth/login-options", {}, { token: null });
  if (!res.ok) {
    return { localEnabled: true, ldapEnabled: false, oidcEnabled: false };
  }
  return toLoginOptions((await res.json()) as LoginOptionsDTO);
}

export async function fetchMe(): Promise<User | null> {
  const res = await apiFetch("/api/v1/me");
  if (res.status === 401) {
    const { clearToken } = await import("@/lib/auth");
    clearToken();
    return null;
  }
  if (!res.ok) return null;
  return toUser((await res.json()) as UserDTO);
}

export const productsApi = {
  list: fetchProducts,
  async search(params: {
    q?: string;
    solutionCategoryId?: string | number;
    vendorMarket?: string;
    manufacturer?: string;
  }): Promise<Product[]> {
    const qs = new URLSearchParams();
    if (params.q) qs.set("q", params.q);
    if (params.solutionCategoryId)
      qs.set("solution_category_id", String(params.solutionCategoryId));
    if (params.vendorMarket) qs.set("vendor_market", params.vendorMarket);
    if (params.manufacturer) qs.set("manufacturer", params.manufacturer);
    const path = params.q
      ? `/api/v1/products/search?${qs.toString()}`
      : `/api/v1/products${qs.toString() ? `?${qs.toString()}` : ""}`;
    return toProducts(await apiJson<ProductDTO[]>(path));
  },
  async get(id: number | string): Promise<Product> {
    return toProduct(await apiJson<ProductDTO>(`/api/v1/products/${id}`));
  },
  async scripts(id: number | string): Promise<SalesScript[]> {
    return toSalesScripts(
      await apiJson<SalesScriptDTO[]>(`/api/v1/products/${id}/scripts`),
    );
  },
  async cases(id: number | string): Promise<Case[]> {
    return toCases(await apiJson<CaseDTO[]>(`/api/v1/products/${id}/cases`));
  },
};

export type AdminProductQuery = {
  page?: number;
  pageSize?: number;
  sort?: string;
  order?: "asc" | "desc";
  q?: string;
  status?: "draft" | "published";
  market?: string;
};

export const adminProductsApi = {
  async list(params: AdminProductQuery = {}): Promise<PageResult<Product>> {
    const qs = new URLSearchParams();
    if (params.page) qs.set("page", String(params.page));
    if (params.pageSize) qs.set("page_size", String(params.pageSize));
    if (params.sort) qs.set("sort", params.sort);
    if (params.order) qs.set("order", params.order);
    if (params.q) qs.set("q", params.q);
    if (params.status) qs.set("status", params.status);
    if (params.market) qs.set("market", params.market);
    const path = `/api/v1/admin/products${qs.toString() ? `?${qs.toString()}` : ""}`;
    const data = await apiJson<unknown>(path);
    return toPageResult<ProductDTO, Product>(
      data,
      params.page ?? 1,
      params.pageSize ?? 20,
      toProduct,
    );
  },
  async get(id: number | string): Promise<Product> {
    return toProduct(
      await apiJson<ProductDTO>(`/api/v1/admin/products/${id}`),
    );
  },
  async create(input: Partial<ProductDTO>): Promise<Product> {
    return toProduct(
      await apiJson<ProductDTO>("/api/v1/admin/products", {
        method: "POST",
        body: JSON.stringify(input),
      }),
    );
  },
  async update(id: number | string, product: Product): Promise<Product> {
    return toProduct(
      await apiJson<ProductDTO>(`/api/v1/admin/products/${id}`, {
        method: "PUT",
        body: JSON.stringify(toProductDTO(product)),
      }),
    );
  },
  async setDraft(id: number | string, isDraft: boolean): Promise<void> {
    await apiJson<unknown>(`/api/v1/admin/products/${id}/draft`, {
      method: "PATCH",
      body: JSON.stringify({ is_draft: isDraft }),
    });
  },
  async remove(id: number | string): Promise<void> {
    const res = await apiFetch(`/api/v1/admin/products/${id}`, {
      method: "DELETE",
    });
    if (!res.ok) throw new Error("删除失败");
  },
  async listScripts(id: number | string): Promise<SalesScript[]> {
    return toSalesScripts(
      await apiJson<SalesScriptDTO[]>(`/api/v1/admin/products/${id}/scripts`),
    );
  },
  async createScript(
    id: number | string,
    input: { scenario: string; content: string },
  ): Promise<void> {
    await apiJson<unknown>(`/api/v1/admin/products/${id}/scripts`, {
      method: "POST",
      body: JSON.stringify(input),
    });
  },
  async updateScript(
    scriptId: number | string,
    input: { scenario: string; content: string },
  ): Promise<void> {
    await apiJson<unknown>(`/api/v1/admin/scripts/${scriptId}`, {
      method: "PUT",
      body: JSON.stringify(input),
    });
  },
  async deleteScript(scriptId: number | string): Promise<void> {
    const res = await apiFetch(`/api/v1/admin/scripts/${scriptId}`, {
      method: "DELETE",
    });
    if (!res.ok) throw new Error("删除话术失败");
  },
  async listCases(id: number | string): Promise<Case[]> {
    return toCases(
      await apiJson<CaseDTO[]>(`/api/v1/admin/products/${id}/cases`),
    );
  },
  async createCase(
    id: number | string,
    input: Omit<CaseDTO, "id" | "product_id">,
  ): Promise<void> {
    await apiJson<unknown>(`/api/v1/admin/products/${id}/cases`, {
      method: "POST",
      body: JSON.stringify(input),
    });
  },
  async updateCase(
    caseId: number | string,
    input: Omit<CaseDTO, "id" | "product_id">,
  ): Promise<void> {
    await apiJson<unknown>(`/api/v1/admin/cases/${caseId}`, {
      method: "PUT",
      body: JSON.stringify(input),
    });
  },
  async deleteCase(caseId: number | string): Promise<void> {
    const res = await apiFetch(`/api/v1/admin/cases/${caseId}`, {
      method: "DELETE",
    });
    if (!res.ok) throw new Error("删除案例失败");
  },
};

export const topologyApi = {
  get: fetchTopology,
  categories: fetchSolutionCategories,
};

export type LayerInput = {
  level: number;
  title: string;
  subtitle: string;
  sort_order: number;
};

export type CategoryInput = {
  label: string;
  icon_key: string;
  column_type: string;
  layer_id?: number | null;
  sort_order: number;
  hint?: string;
  keywords?: string[];
  is_active: boolean;
};

export type VendorInput = { name: string; sort_order: number };

async function writeJSON(
  path: string,
  method: "POST" | "PUT" | "PATCH" | "DELETE",
  body?: unknown,
  token?: string | null,
): Promise<void> {
  const res = await apiFetch(path, {
    method,
    body: body === undefined ? undefined : JSON.stringify(body),
  }, token === undefined ? undefined : { token });
  if (!res.ok) {
    const data: unknown = await res.json().catch(() => null);
    const msg =
      data && typeof data === "object" && "error" in data
        ? String((data as { error: unknown }).error)
        : `请求失败（${res.status}）`;
    throw new Error(msg);
  }
}

export const adminTopologyApi = {
  createLayer: (input: LayerInput) =>
    writeJSON("/api/v1/admin/topology/layers", "POST", input),
  updateLayer: (id: number | string, input: LayerInput) =>
    writeJSON(`/api/v1/admin/topology/layers/${id}`, "PUT", input),
  deleteLayer: (id: number | string) =>
    writeJSON(`/api/v1/admin/topology/layers/${id}`, "DELETE"),
  createCategory: (input: CategoryInput) =>
    writeJSON("/api/v1/admin/topology/categories", "POST", input),
  updateCategory: (id: number | string, input: CategoryInput) =>
    writeJSON(`/api/v1/admin/topology/categories/${id}`, "PUT", input),
  deleteCategory: (id: number | string) =>
    writeJSON(`/api/v1/admin/topology/categories/${id}`, "DELETE"),
  createVendor: (categoryId: number | string, input: VendorInput) =>
    writeJSON("/api/v1/admin/topology/vendors", "POST", {
      ...input,
      category_id: categoryId,
    }),
  updateVendor: (id: number | string, input: VendorInput) =>
    writeJSON(`/api/v1/admin/topology/vendors/${id}`, "PUT", input),
  deleteVendor: (id: number | string) =>
    writeJSON(`/api/v1/admin/topology/vendors/${id}`, "DELETE"),
};

export type LocalUserQuery = {
  page?: number;
  pageSize?: number;
  q?: string;
  role?: "admin" | "user";
  active?: boolean;
};

export const usersApi = {
  async list(params: LocalUserQuery = {}): Promise<PageResult<LocalUserRow>> {
    const qs = new URLSearchParams();
    if (params.page) qs.set("page", String(params.page));
    if (params.pageSize) qs.set("page_size", String(params.pageSize));
    if (params.q) qs.set("q", params.q);
    if (params.role) qs.set("role", params.role);
    if (typeof params.active === "boolean")
      qs.set("active", String(params.active));
    const path = `/api/v1/admin/users/local${qs.toString() ? `?${qs.toString()}` : ""}`;
    const data = await apiJson<unknown>(path);
    return toPageResult<LocalUserDTO, LocalUserRow>(
      data,
      params.page ?? 1,
      params.pageSize ?? 20,
      toLocalUser,
    );
  },
  create: (input: {
    username: string;
    password: string;
    email?: string;
    role?: string;
  }) =>
    apiJson<LocalUserDTO>("/api/v1/admin/users/local", {
      method: "POST",
      body: JSON.stringify(input),
    }).then(toLocalUser),
  update: (
    id: number | string,
    input: { email?: string; role?: string; is_active?: boolean },
  ) =>
    apiJson<LocalUserDTO>(`/api/v1/admin/users/local/${id}`, {
      method: "PUT",
      body: JSON.stringify(input),
    }).then(toLocalUser),
  resetPassword: (id: number | string, password: string) =>
    writeJSON(`/api/v1/admin/users/local/${id}/password`, "PUT", { password }),
  remove: (id: number | string) =>
    writeJSON(`/api/v1/admin/users/local/${id}`, "DELETE"),
};

export const apiKeysApi = {
  async list(): Promise<ApiKeyRow[]> {
    return toApiKeys(await apiJson<ApiKeyDTO[]>("/api/v1/admin/api-keys"));
  },
  async create(name: string, expiresAt?: string): Promise<ApiKeyCreateResult> {
    const body: { name: string; expires_at?: string } = { name };
    if (expiresAt) body.expires_at = expiresAt;
    return toApiKeyCreated(
      await apiJson<ApiKeyCreateDTO>("/api/v1/admin/api-keys", {
        method: "POST",
        body: JSON.stringify(body),
      }),
    );
  },
  revoke: (id: number | string) =>
    writeJSON(`/api/v1/admin/api-keys/${id}`, "DELETE"),
};

export const settingsApi = {
  auth: () => apiJson<AuthSettingsDTO>("/api/v1/admin/auth-settings"),
  patchAuth: (input: {
    local_login_enabled?: boolean;
    ldap_enabled?: boolean;
    ldap_config?: unknown;
    oidc_enabled?: boolean;
    oidc_config?: unknown;
  }) =>
    apiJson<AuthSettingsDTO>("/api/v1/admin/auth-settings", {
      method: "PATCH",
      body: JSON.stringify(input),
    }),
  system: () => apiJson<SystemSettingsDTO>("/api/v1/admin/system-settings"),
  patchSystem: (input: Partial<SystemSettingsDTO>) =>
    apiJson<SystemSettingsDTO>("/api/v1/admin/system-settings", {
      method: "PATCH",
      body: JSON.stringify(input),
    }),
  config: () => apiJson<SystemSettingsDTO>("/api/v1/config"),
};

export const backupApi = {
  exportPath: "/api/v1/admin/backup/export",
  import: (backup: unknown, confirm = BACKUP_RESTORE_CONFIRM) =>
    apiJson<{ ok: boolean; message: string }>("/api/v1/admin/backup/import", {
      method: "POST",
      body: JSON.stringify({ confirm, backup }),
    }),
};
