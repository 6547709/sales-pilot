"use client";

import { motion } from "framer-motion";
import Link from "next/link";
import { useCallback, useEffect, useState } from "react";
import { Badge } from "@/components/ui/badge";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { apiFetch, fetchTopology, parseJsonArray, type Product } from "@/lib/api";
import {
  vendorMarketBadgeClass,
  vendorMarketBadgeText,
} from "@/lib/vendor-market";
import { ArrowRight, Layers } from "lucide-react";

type Props = {
  /** 与首页架构卡片联动：按解决方案 ID 精确筛选 */
  solutionId?: string;
  /** 展示用标题（解决方案名称） */
  solutionLabel?: string;
  /** 兼容旧参数：关键词搜索 */
  keyword?: string;
  /** 与架构卡「国内/国外」一致，对应后端 vendor_market */
  vendorMarket?: "domestic" | "foreign";
};

export function ProductSpotlight({
  solutionId,
  solutionLabel,
  keyword = "",
  vendorMarket,
}: Props) {
  const [list, setList] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [resolvedLabel, setResolvedLabel] = useState(solutionLabel ?? "");

  useEffect(() => {
    if (solutionLabel) {
      setResolvedLabel(solutionLabel);
      return;
    }
    const sid = solutionId?.trim();
    if (!sid) {
      setResolvedLabel("");
      return;
    }
    fetchTopology()
      .then((t) => {
        const all = [
          ...t.security,
          ...t.ops,
          ...t.central_layers.flatMap((b) => b.categories),
        ];
        const hit = all.find((c) => String(c.id) === sid);
        setResolvedLabel(hit?.label ?? `解决方案 #${sid}`);
      })
      .catch(() => setResolvedLabel(`解决方案 #${sid}`));
  }, [solutionId, solutionLabel]);

  const load = useCallback(async () => {
    setLoading(true);
    let path = "/api/v1/products";
    const sid = solutionId?.trim();
    const kw = keyword.trim();
    if (sid) {
      const q = new URLSearchParams({
        solution_category_id: sid,
      });
      if (vendorMarket) q.set("vendor_market", vendorMarket);
      path = `/api/v1/products?${q.toString()}`;
    } else if (kw) {
      const q = new URLSearchParams({ q: kw });
      if (vendorMarket) q.set("vendor_market", vendorMarket);
      path = `/api/v1/products/search?${q.toString()}`;
    }
    const res = await apiFetch(path);
    if (res.ok) {
      const data: Product[] = await res.json();
      setList(data.slice(0, 12));
    } else setList([]);
    setLoading(false);
  }, [solutionId, keyword, vendorMarket]);

  useEffect(() => {
    load();
  }, [load]);

  const marketSuffix =
    vendorMarket === "domestic"
      ? "（国内）"
      : vendorMarket === "foreign"
        ? "（国外）"
        : "";

  const titleLabel =
    resolvedLabel || (keyword.trim() ? keyword.trim() : null);

  const listHref = (() => {
    if (solutionId?.trim()) {
      const q = new URLSearchParams({
        solution_category_id: solutionId.trim(),
      });
      if (vendorMarket) q.set("market", vendorMarket);
      return `/products?${q.toString()}`;
    }
    if (keyword.trim()) {
      return `/products?q=${encodeURIComponent(keyword.trim())}`;
    }
    return "/products";
  })();

  return (
    <section
      id="solution-spotlight"
      className="scroll-mt-24 border-t border-border/60 bg-muted/20 py-10 sm:py-14"
    >
      <div className="mx-auto max-w-[1400px] px-4">
        <div className="mb-6 flex flex-wrap items-end justify-between gap-4">
          <div>
            <Badge variant="outline" className="mb-2 gap-1">
              <Layers className="h-3 w-3" />
              赋能方案
            </Badge>
            <h2 className="text-2xl font-bold tracking-tight sm:text-3xl">
              {titleLabel ? (
                <>
                  与「<span className="text-primary">{titleLabel}</span>」相关的方案
                  {marketSuffix ? (
                    <span className="text-muted-foreground">{marketSuffix}</span>
                  ) : null}
                </>
              ) : (
                "精选方案"
              )}
            </h2>
            <p className="mt-1 text-sm text-muted-foreground">
              已通过「所属解决方案」关联。列表中的「国内 / 国外」来自后台市场定位；全局架构卡仅展示您已录入的真实产品。亦可在方案库全文搜索或按地域筛选。
            </p>
          </div>
          <Link
            href={listHref}
            className="inline-flex items-center gap-1 text-sm font-medium text-primary hover:underline"
          >
            打开完整列表
            <ArrowRight className="h-4 w-4" />
          </Link>
        </div>

        {loading ? (
          <ul className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {[...Array(6)].map((_, i) => (
              <li key={i}>
                <Card className="h-full">
                  <CardHeader>
                    <div className="flex flex-wrap items-start justify-between gap-2">
                      <Skeleton className="h-6 w-3/4" />
                      <Skeleton className="h-5 w-16" />
                    </div>
                    <Skeleton className="mt-2 h-4 w-1/2" />
                  </CardHeader>
                </Card>
              </li>
            ))}
          </ul>
        ) : list.length === 0 ? (
          <p className="rounded-xl border border-dashed bg-background/80 p-8 text-center text-sm text-muted-foreground">
            暂无匹配方案。请尝试其他解决方案，或到方案库搜索。
          </p>
        ) : (
          <ul className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {list.map((p, i) => (
              <motion.li
                key={p.id}
                initial={{ opacity: 0, y: 8 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: i * 0.04 }}
              >
                <Link href={`/products/${p.id}`} className="block h-full">
                  <Card className="h-full flex flex-col border-border/80 transition hover:border-primary/35 hover:shadow-md overflow-hidden">
                    <CardHeader className={(solutionId || keyword || vendorMarket) ? "pb-4" : ""}>
                      <div className="flex flex-wrap items-start justify-between gap-2">
                        <div className="flex items-center gap-3">
                          {p.manufacturer_logo ? (
                            <img src={p.manufacturer_logo} alt={p.manufacturer_name ?? "logo"} className="max-h-6 max-w-[80px] object-contain" />
                          ) : null}
                          <CardTitle className="text-lg leading-snug">
                            {p.name}
                          </CardTitle>
                        </div>
                        {vendorMarketBadgeText(p.vendor_market) ? (
                          <Badge
                            variant="outline"
                            className={`shrink-0 text-xs font-medium ${vendorMarketBadgeClass(p.vendor_market)}`}
                          >
                            {vendorMarketBadgeText(p.vendor_market)}
                          </Badge>
                        ) : null}
                      </div>
                      <CardDescription className="line-clamp-2">
                        {p.manufacturer_name
                          ? `${p.manufacturer_name} · ${p.category || "方案"}`
                          : p.category || "未分类"}
                      </CardDescription>
                    </CardHeader>
                    {(solutionId || keyword || vendorMarket) ? (
                      (() => {
                        const hl = parseJsonArray(p.highlights);
                        const qs = parseJsonArray(p.discovery_questions);
                        if (hl.length === 0 && qs.length === 0) return null;
                        return (
                          <CardContent className="mt-auto px-6 py-4 border-t bg-muted/10 text-xs flex-1">
                            <div className="space-y-5">
                              {hl.length > 0 && (
                                <div>
                                  <div className="font-semibold text-primary/80 mb-2 flex items-center gap-1.5">
                                    <span className="text-[11px]">✨</span> 亮点优势
                                  </div>
                                  <div className="space-y-2.5 text-muted-foreground">
                                    {hl.map((h, i) => (
                                      <div key={i} className="flex items-start gap-2">
                                        <span className="flex h-4 w-4 shrink-0 items-center justify-center rounded-full bg-primary/15 text-[10px] font-bold text-primary mt-px">{i + 1}</span>
                                        <span className="line-clamp-2 leading-relaxed">{h}</span>
                                      </div>
                                    ))}
                                  </div>
                                </div>
                              )}
                              {qs.length > 0 && (
                                <div>
                                  <div className="font-semibold text-primary/80 mb-2 flex items-center gap-1.5">
                                    <span className="text-[11px]">💡</span> 黄金三问
                                  </div>
                                  <div className="space-y-2.5 text-muted-foreground">
                                    {qs.map((q, i) => (
                                      <div key={i} className="flex items-start gap-2">
                                        <span className="flex h-4 w-4.5 shrink-0 items-center justify-center rounded-[4px] bg-primary/15 px-1 text-[9px] font-bold text-primary mt-px">Q{i + 1}</span>
                                        <span className="line-clamp-2 leading-relaxed">{q}</span>
                                      </div>
                                    ))}
                                  </div>
                                </div>
                              )}
                            </div>
                          </CardContent>
                        );
                      })()
                    ) : null}
                  </Card>
                </Link>
              </motion.li>
            ))}
          </ul>
        )}
      </div>
    </section>
  );
}
