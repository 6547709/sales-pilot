/** 产品 vendor_market 字段对应前台展示文案（domestic/foreign 以外不展示标签） */
export function vendorMarketBadgeText(
  market: string | undefined | null,
): string | null {
  const m = String(market || "")
    .toLowerCase()
    .trim();
  if (m === "domestic") return "国内";
  if (m === "foreign") return "国外";
  return null;
}

/** 用于 Badge 的配色区分 */
export function vendorMarketBadgeClass(market: string | undefined | null): string {
  const m = String(market || "")
    .toLowerCase()
    .trim();
  if (m === "domestic") {
    return "border-red-200 bg-red-50 text-red-900 dark:border-red-900/50 dark:bg-red-950/40 dark:text-red-200";
  }
  if (m === "foreign") {
    return "border-slate-300 bg-slate-100 text-slate-800 dark:border-slate-600 dark:bg-slate-800 dark:text-slate-200";
  }
  return "";
}
