/** 后端 JSON 形状（snake_case）。页面不得直接使用这些类型，统一经 adapters 转成视图模型。 */

export type UserDTO = {
  id: number;
  username: string;
  email: string;
  role: string;
  auth_provider: string;
};

export type LocalUserDTO = {
  id: number;
  username: string;
  email: string;
  role: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
};

export type SolutionCategoryRefDTO = {
  id: number;
  slug: string;
  label: string;
  column_type?: string;
};

export type SolutionCategoryDTO = {
  id: number;
  slug: string;
  label: string;
  column_type: string;
  layer_id?: number | null;
  is_active: boolean;
};

export type ProductDTO = {
  id: number;
  name: string;
  category: string;
  vendor_market?: string;
  is_draft?: boolean;
  solution_category_id?: number | null;
  solution_category?: SolutionCategoryRefDTO | null;
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
  customer_value?: string;
  customer_features?: string | string[];
  target_personas: Record<string, string> | string;
  trigger_events: string;
  discovery_questions: string | string[];
  competitor_analysis: string;
  roi_metrics: string;
  updated_at: string;
  created_at?: string;
};

export type SalesScriptDTO = {
  id: number;
  product_id: number;
  scenario: string;
  content: string;
};

export type CaseDTO = {
  id: number;
  product_id: number;
  client_name: string;
  pain_points: string;
  solution: string;
  value_delivered: string;
};

export type TopologyVendorDTO = {
  id: number;
  name: string;
  sort_order: number;
};

export type TopologyCategoryNodeDTO = {
  id: number;
  slug: string;
  label: string;
  icon_key: string;
  keywords: string[];
  hint: string;
  vendors: { domestic: TopologyVendorDTO[]; foreign: TopologyVendorDTO[] };
  services?: ServiceDTO[];
};

export type TopologyLayerBlockDTO = {
  layer: {
    id: number;
    level: number;
    title: string;
    subtitle: string;
    sort_order: number;
  };
  categories: TopologyCategoryNodeDTO[];
};

export type TopologyDTO = {
  security: TopologyCategoryNodeDTO[];
  ops: TopologyCategoryNodeDTO[];
  central_layers: TopologyLayerBlockDTO[];
};

export type ApiKeyDTO = {
  id: number;
  user_id: number;
  name: string;
  prefix: string;
  expires_at?: string | null;
  last_used_at?: string | null;
  is_active: boolean;
  created_at: string;
};

export type ApiKeyCreateDTO = {
  id: number;
  name: string;
  prefix: string;
  expires_at?: string | null;
  key: string;
  warning: string;
};

export type LoginOptionsDTO = {
  local_enabled: boolean;
  ldap_enabled: boolean;
  oidc_enabled: boolean;
};

export type AuthSettingsDTO = {
  id: number;
  local_login_enabled?: boolean;
  ldap_enabled: boolean;
  ldap_config: string | Record<string, unknown>;
  oidc_enabled: boolean;
  oidc_config: string | Record<string, unknown>;
};

export type SystemSettingsDTO = {
  footer_text?: string;
  filing_number?: string;
};

export type SearchDTO = { q?: string; keyword?: string; results?: ProductDTO[] } | ProductDTO[];

export type PageDTO<T> = {
  items: T[];
  total: number;
  page: number;
  page_size: number;
};

export type ServiceDTO = {
  id: number;
  slug: string;
  name: string;
  kind: string;
  summary: string;
  highlights: string | string[];
  deliverables: string | string[];
  value_prop?: string;
  features?: string | string[];
  icon_key?: string;
  sort_order: number;
  is_active: boolean;
  updated_at: string;
};

export type ScenarioDTO = {
  id: number;
  group_id: number;
  slug: string;
  name: string;
  summary: string;
  pain_point?: string;
  value_prop?: string;
  icon_key?: string;
  sort_order: number;
  is_active: boolean;
};

export type ScenarioSummaryDTO = ScenarioDTO & {
  product_count: number;
  service_count: number;
};

export type ScenarioGroupDTO = {
  id: number;
  slug: string;
  name: string;
  description: string;
  sort_order: number;
  scenarios: ScenarioSummaryDTO[];
};

export type ScenarioDetailDTO = {
  scenario: ScenarioDTO;
  group?: { id: number; slug: string; name: string; description: string; sort_order: number } | null;
  services: ServiceDTO[];
  products: ProductDTO[];
  cases: CaseDTO[];
};
