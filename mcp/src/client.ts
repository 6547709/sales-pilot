const base = (process.env.SALES_PILOT_BASE_URL ?? "http://127.0.0.1:8080").replace(/\/$/, "");
const apiKey = process.env.SALES_PILOT_API_KEY ?? "";

export async function api(path: string, init: RequestInit = {}): Promise<unknown> {
  if (!apiKey) {
    throw new Error("SALES_PILOT_API_KEY 未配置");
  }
  const res = await fetch(`${base}${path}`, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      "X-API-Key": apiKey,
      ...(init.headers ?? {}),
    },
  });
  const text = await res.text();
  const body: unknown = text ? JSON.parse(text) : null;
  if (!res.ok) {
    const message =
      body && typeof body === "object" && "error" in body
        ? String((body as { error: unknown }).error)
        : `HTTP ${res.status}`;
    throw new Error(message);
  }
  return body;
}

export function text(value: unknown) {
  return {
    content: [{ type: "text" as const, text: JSON.stringify(value, null, 2) }],
  };
}
