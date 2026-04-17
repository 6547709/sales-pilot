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

/**
 * 页脚文字，可选。如 "企业销售赋能"。
 */
export const FOOTER_TEXT = process.env.NEXT_PUBLIC_FOOTER_TEXT;

/**
 * 备案号，可选。如 "京ICP备xxxx号"。
 */
export const FILING_NUMBER = process.env.NEXT_PUBLIC_FILING_NUMBER;
