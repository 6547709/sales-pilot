import { API_BASE } from "./config";
import { clearToken, getToken } from "./auth";

export type User = {
  id: number;
  username: string;
  email: string;
  role: string;
  auth_provider: string;
};

export type Product = {
  id: number;
  name: string;
  category: string;
  description: string;
  highlights: string | string[];
  target_personas: Record<string, string> | string;
  trigger_events: string;
  discovery_questions: string | string[];
  competitor_analysis: string;
  roi_metrics: string;
  updated_at: string;
};

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
