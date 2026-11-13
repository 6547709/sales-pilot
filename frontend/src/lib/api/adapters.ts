import type {
  ApiKeyCreateDTO,
  ApiKeyDTO,
  CaseDTO,
  LocalUserDTO,
  LoginOptionsDTO,
  ProductDTO,
  SalesScriptDTO,
  SolutionCategoryDTO,
  TopologyCategoryNodeDTO,
  TopologyDTO,
  TopologyLayerBlockDTO,
  TopologyVendorDTO,
  UserDTO,
} from "./types";

export type User = {
  id: number;
  username: string;
  email: string;
  role: string;
  authProvider: string;
};

export type LocalUserRow = {
  id: number;
  username: string;
  email: string;
  role: string;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
};

export type SolutionCategoryRef = {
  id: number;
  slug: string;
  label: string;
  columnType?: string;
};

export type Product = {
  id: number;
  name: string;
  category: string;
  vendorMarket?: string;
  isDraft?: boolean;
  solutionCategoryId?: number | null;
  solutionCategory?: SolutionCategoryRef | null;
  manufacturerName?: string;
  manufacturerLogo?: string;
  salesContactName?: string;
  salesContactPhone?: string;
  salesContactEmail?: string;
  presalesContactName?: string;
  presalesContactPhone?: string;
  presalesContactEmail?: string;
  description: string;
  highlights: string[];
  targetPersonas: Record<string, string>;
  triggerEvents: string;
  discoveryQuestions: string[];
  competitorAnalysis: string;
  roiMetrics: string;
  updatedAt: string;
  createdAt?: string;
};

export type SalesScript = {
  id: number;
  productId: number;
  scenario: string;
  content: string;
};

export type Case = {
  id: number;
  productId: number;
  clientName: string;
  painPoints: string;
  solution: string;
  valueDelivered: string;
};

export type ApiKeyRow = {
  id: number;
  userId: number;
  name: string;
  prefix: string;
  expiresAt?: string | null;
  lastUsedAt?: string | null;
  isActive: boolean;
  createdAt: string;
};

export type ApiKeyCreateResult = {
  id: number;
  name: string;
  prefix: string;
  expiresAt?: string | null;
  key: string;
  warning: string;
};

export type LoginOptions = {
  localEnabled: boolean;
  ldapEnabled: boolean;
  oidcEnabled: boolean;
};

/** 拓扑与解决方案分类暂时沿用后端形状，迁移在后续任务处理。 */
export type SolutionCategory = SolutionCategoryDTO;
export type TopologyVendor = TopologyVendorDTO;
export type TopologyCategoryNode = TopologyCategoryNodeDTO;
export type TopologyLayerBlock = TopologyLayerBlockDTO;
export type TopologyResponse = TopologyDTO;

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

export function toUser(dto: UserDTO): User {
  return {
    id: dto.id,
    username: dto.username,
    email: dto.email,
    role: dto.role,
    authProvider: dto.auth_provider,
  };
}

export function toLocalUser(dto: LocalUserDTO): LocalUserRow {
  return {
    id: dto.id,
    username: dto.username,
    email: dto.email,
    role: dto.role,
    isActive: dto.is_active,
    createdAt: dto.created_at,
    updatedAt: dto.updated_at,
  };
}

export function toLocalUsers(dtos: LocalUserDTO[]): LocalUserRow[] {
  return dtos.map(toLocalUser);
}

export function toProduct(dto: ProductDTO): Product {
  return {
    id: dto.id,
    name: dto.name,
    category: dto.category,
    vendorMarket: dto.vendor_market,
    isDraft: dto.is_draft,
    solutionCategoryId: dto.solution_category_id,
    solutionCategory: dto.solution_category
      ? {
          id: dto.solution_category.id,
          slug: dto.solution_category.slug,
          label: dto.solution_category.label,
          columnType: dto.solution_category.column_type,
        }
      : null,
    manufacturerName: dto.manufacturer_name,
    manufacturerLogo: dto.manufacturer_logo,
    salesContactName: dto.sales_contact_name,
    salesContactPhone: dto.sales_contact_phone,
    salesContactEmail: dto.sales_contact_email,
    presalesContactName: dto.presales_contact_name,
    presalesContactPhone: dto.presales_contact_phone,
    presalesContactEmail: dto.presales_contact_email,
    description: dto.description,
    highlights: parseJsonArray(dto.highlights),
    targetPersonas: parsePersonas(dto.target_personas),
    triggerEvents: dto.trigger_events,
    discoveryQuestions: parseJsonArray(dto.discovery_questions),
    competitorAnalysis: dto.competitor_analysis,
    roiMetrics: dto.roi_metrics,
    updatedAt: dto.updated_at,
    createdAt: dto.created_at,
  };
}

export function toProducts(dtos: ProductDTO[]): Product[] {
  return dtos.map(toProduct);
}

/** 写回后端时恢复 snake_case；associations 不参与写入。 */
export function toProductDTO(p: Product): ProductDTO {
  return {
    id: p.id,
    name: p.name,
    category: p.category,
    vendor_market: p.vendorMarket,
    is_draft: p.isDraft,
    solution_category_id: p.solutionCategoryId ?? null,
    manufacturer_name: p.manufacturerName,
    manufacturer_logo: p.manufacturerLogo,
    sales_contact_name: p.salesContactName,
    sales_contact_phone: p.salesContactPhone,
    sales_contact_email: p.salesContactEmail,
    presales_contact_name: p.presalesContactName,
    presales_contact_phone: p.presalesContactPhone,
    presales_contact_email: p.presalesContactEmail,
    description: p.description,
    highlights: p.highlights,
    target_personas: p.targetPersonas,
    trigger_events: p.triggerEvents,
    discovery_questions: p.discoveryQuestions,
    competitor_analysis: p.competitorAnalysis,
    roi_metrics: p.roiMetrics,
    updated_at: p.updatedAt,
    created_at: p.createdAt,
  };
}

export function toSalesScript(dto: SalesScriptDTO): SalesScript {
  return {
    id: dto.id,
    productId: dto.product_id,
    scenario: dto.scenario,
    content: dto.content,
  };
}

export function toSalesScripts(dtos: SalesScriptDTO[]): SalesScript[] {
  return dtos.map(toSalesScript);
}

export function toCase(dto: CaseDTO): Case {
  return {
    id: dto.id,
    productId: dto.product_id,
    clientName: dto.client_name,
    painPoints: dto.pain_points,
    solution: dto.solution,
    valueDelivered: dto.value_delivered,
  };
}

export function toCases(dtos: CaseDTO[]): Case[] {
  return dtos.map(toCase);
}

export function toApiKey(dto: ApiKeyDTO): ApiKeyRow {
  return {
    id: dto.id,
    userId: dto.user_id,
    name: dto.name,
    prefix: dto.prefix,
    expiresAt: dto.expires_at,
    lastUsedAt: dto.last_used_at,
    isActive: dto.is_active,
    createdAt: dto.created_at,
  };
}

export function toApiKeys(dtos: ApiKeyDTO[]): ApiKeyRow[] {
  return dtos.map(toApiKey);
}

export function toApiKeyCreated(dto: ApiKeyCreateDTO): ApiKeyCreateResult {
  return {
    id: dto.id,
    name: dto.name,
    prefix: dto.prefix,
    expiresAt: dto.expires_at,
    key: dto.key,
    warning: dto.warning,
  };
}

export function toLoginOptions(dto: LoginOptionsDTO): LoginOptions {
  return {
    localEnabled: dto.local_enabled,
    ldapEnabled: dto.ldap_enabled,
    oidcEnabled: dto.oidc_enabled,
  };
}

export type PageResult<T> = {
  items: T[];
  total: number;
  page: number;
  pageSize: number;
};

/** 后端带 page 时返回 {items,...}，否则返回裸数组；统一成 PageResult。 */
export function toPageResult<T, M>(
  data: unknown,
  page: number,
  pageSize: number,
  map: (dto: T) => M,
): PageResult<M> {
  if (Array.isArray(data)) {
    const items = (data as T[]).map(map);
    return { items, total: items.length, page: 1, pageSize: items.length };
  }
  const p = data as { items?: T[]; total?: number; page?: number; page_size?: number };
  const items = (p.items ?? []).map(map);
  return {
    items,
    total: p.total ?? items.length,
    page: p.page ?? page,
    pageSize: p.page_size ?? pageSize,
  };
}
