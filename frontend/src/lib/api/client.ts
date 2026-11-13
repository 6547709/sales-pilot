import { API_BASE } from "@/lib/config";
import { clearToken, getToken } from "@/lib/auth";
import { problemMessage } from "./problem";

export class ApiError extends Error {
  status: number;
  code?: string;

  constructor(status: number, message: string, code?: string) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.code = code;
  }
}

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

export async function apiJson<T>(
  path: string,
  init: RequestInit = {},
  opts?: { token?: string | null },
): Promise<T> {
  const res = await apiFetch(path, init, opts);
  if (!res.ok) {
    if (res.status === 401) clearToken();
    const body: unknown = await res.json().catch(() => null);
    const code =
      body && typeof body === "object" && typeof (body as { code?: unknown }).code === "string"
        ? (body as { code: string }).code
        : undefined;
    throw new ApiError(
      res.status,
      problemMessage(body, `请求失败（${res.status}）`),
      code,
    );
  }
  return (await res.json()) as T;
}
