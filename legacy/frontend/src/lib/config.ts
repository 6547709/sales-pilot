/**
 * 后端 API 基址。
 * - 本地开发默认 http://localhost:8080
 * - 与 Nginx 同域部署时设为空字符串，使用相对路径 /api/...
 */
function resolveApiBase(): string {
  const v = process.env.NEXT_PUBLIC_API_URL;
  if (v !== undefined) {
    return v.replace(/\/$/, "");
  }
  return "http://localhost:8080";
}

export const API_BASE = resolveApiBase();
