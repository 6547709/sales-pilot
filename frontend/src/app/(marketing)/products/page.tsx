"use client";

import { motion } from "framer-motion";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { Suspense, useCallback, useEffect, useRef, useState } from "react";
import { useSearchParams } from "next/navigation";
import { Badge } from "@/components/ui/badge";
import { buttonVariants } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { apiFetch, fetchSolutionCategories, parseJsonArray, type Product } from "@/lib/api";
import {
  vendorMarketBadgeClass,
  vendorMarketBadgeText,
} from "@/lib/vendor-market";
import { LayoutGrid, Search } from "lucide-react";

function ProductsInner() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const isFirstLoad = useRef(true);
  const initialCat = searchParams.get("solution_category_id") ?? "";
  const initialQ = searchParams.get("q") ?? searchParams.get("kw") ?? "";
  const initialMarketRaw = searchParams.get("market") ?? "";
  const initialMarket =
    initialMarketRaw === "domestic" || initialMarketRaw === "foreign"
      ? initialMarketRaw
      : "";
  const initialMfg = searchParams.get("manufacturer") ?? "";
  const [solutionFilter, setSolutionFilter] = useState(initialCat);
  const [q, setQ] = useState(initialQ);
  const [marketFilter, setMarketFilter] = useState(initialMarket);
  const [mfrFilter, setMfrFilter] = useState(initialMfg);
  const [list, setList] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [categories, setCategories] = useState<
    { id: number; label: string; slug: string; is_active: boolean }[]
  >([]);

  useEffect(() => {
    setSolutionFilter(initialCat);
    setQ(initialQ);
    setMarketFilter(initialMarket);
    setMfrFilter(initialMfg);
  }, [initialCat, initialQ, initialMarket, initialMfg]);

  useEffect(() => {
    fetchSolutionCategories(false).then(setCategories);
  }, []);

  const load = useCallback(async () => {
    setLoading(true);
    const sid = solutionFilter.trim();
    const kw = q.trim();
    let path: string;
    const mfr = mfrFilter.trim();
    if (kw) {
      const params = new URLSearchParams({ q: kw });
      if (sid) params.set("solution_category_id", sid);
      if (marketFilter) params.set("vendor_market", marketFilter);
      if (mfr) params.set("manufacturer", mfr);
      path = `/api/v1/products/search?${params.toString()}`;
    } else if (sid) {
      const params = new URLSearchParams({
        solution_category_id: sid,
      });
      if (marketFilter) params.set("vendor_market", marketFilter);
      if (mfr) params.set("manufacturer", mfr);
      path = `/api/v1/products?${params.toString()}`;
    } else {
      const params = new URLSearchParams();
      if (marketFilter) params.set("vendor_market", marketFilter);
      if (mfr) params.set("manufacturer", mfr);
      path =
        params.toString() === ""
          ? "/api/v1/products"
          : `/api/v1/products?${params.toString()}`;
    }
    const res = await apiFetch(path);
    if (res.ok) {
      const data = await res.json();
      // 如果只有1个结果，自动跳转到详情页
      if (data.length === 1 && !isFirstLoad.current) {
        isFirstLoad.current = true;
        router.replace(`/products/${data[0].id}`);
        setLoading(false);
        return;
      }
      setList(data);
    } else {
      setList([]);
    }
    isFirstLoad.current = false;
    setLoading(false);
  }, [q, solutionFilter, marketFilter, mfrFilter]);

  useEffect(() => {
    load();
  }, [load]);

  return (
    <div className="mx-auto max-w-[1400px] px-4 py-8 sm:py-10">
      <div className="mb-8 flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <Badge variant="secondary" className="mb-2 gap-1">
            <LayoutGrid className="h-3 w-3" />
            方案库
          </Badge>
          <h1 className="text-3xl font-bold tracking-tight">解决方案目录</h1>
          <p className="mt-1 max-w-xl text-sm text-muted-foreground">
            可按首页所选「解决方案」、厂商名称、市场定位与关键词组合筛选。从架构卡点击厂商标签进入时会自动带上方案与厂商条件。
          </p>
          {mfrFilter.trim() ? (
            <p className="mt-2 text-sm text-primary">
              当前厂商：<span className="font-medium">{mfrFilter.trim()}</span>
            </p>
          ) : null}
        </div>
        <Link
          href="/"
          className={cn(buttonVariants({ variant: "outline", size: "sm" }))}
        >
          返回架构总览
        </Link>
      </div>

      <div className="mb-8 flex flex-col gap-3 lg:flex-row lg:items-center">
        <div className="grid flex-1 gap-3 sm:grid-cols-2">
          <div className="space-y-1.5">
            <label className="text-xs font-medium text-muted-foreground">
              所属解决方案
            </label>
            <select
              className="flex h-11 w-full rounded-md border border-input bg-background px-3 text-sm shadow-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
              value={solutionFilter}
              onChange={(e) => setSolutionFilter(e.target.value)}
            >
              <option value="">全部</option>
              {categories.map((c) => (
                <option key={c.id} value={String(c.id)}>
                  {c.label}
                  {!c.is_active ? "（已停用）" : ""}
                </option>
              ))}
            </select>
          </div>
          <div className="space-y-1.5">
            <label className="text-xs font-medium text-muted-foreground">
              市场定位
            </label>
            <select
              className="flex h-11 w-full rounded-md border border-input bg-background px-3 text-sm shadow-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
              value={marketFilter}
              onChange={(e) => setMarketFilter(e.target.value)}
            >
              <option value="">全部（含不区分）</option>
              <option value="domestic">国内</option>
              <option value="foreign">国外</option>
            </select>
          </div>
          <div className="relative sm:col-span-2">
            <label className="mb-1.5 block text-xs font-medium text-muted-foreground">
              关键词
            </label>
            <Search className="pointer-events-none absolute bottom-3 left-3 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="搜索名称、分类、描述…"
              value={q}
              onChange={(e) => setQ(e.target.value)}
              onKeyDown={(e) => e.key === "Enter" && load()}
              className="h-11 pl-9"
            />
          </div>
        </div>
        <button
          type="button"
          className={cn(
            buttonVariants({ size: "default" }),
            "h-11 w-full shrink-0 lg:mt-6 lg:w-28",
          )}
          onClick={load}
        >
          搜索
        </button>
      </div>

      {loading ? (
        <p className="text-sm text-muted-foreground">加载中…</p>
      ) : (
        <ul className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {list.map((p, i) => (
            <motion.li
              key={p.id}
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: Math.min(i * 0.03, 0.3) }}
            >
              <Link href={`/products/${p.id}`} className="block h-full">
                <Card className="h-full flex flex-col border-border/80 transition hover:border-primary/40 hover:shadow-lg overflow-hidden">
                  <CardHeader className={(solutionFilter || mfrFilter || q || marketFilter) ? "pb-4" : ""}>
                    <div className="flex flex-wrap items-start justify-between gap-2">
                      <div className="flex items-center gap-3">
                        {p.manufacturer_logo ? (
                          <img src={p.manufacturer_logo} alt={p.manufacturer_name ?? "logo"} className="max-h-6 max-w-[80px] object-contain" />
                        ) : null}
                        <CardTitle className="text-lg leading-snug">{p.name}</CardTitle>
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
                    <CardDescription>
                      {p.manufacturer_name
                        ? `${p.manufacturer_name} · ${p.category || "方案"}`
                        : p.category || "未分类"}
                    </CardDescription>
                  </CardHeader>
                  {(solutionFilter || mfrFilter || q || marketFilter) ? (
                    (() => {
                      const hl = parseJsonArray(p.highlights);
                      const qs = parseJsonArray(p.discovery_questions);
                      if (hl.length === 0 && qs.length === 0) return null;
                      return (
                        <CardContent className="mt-auto px-6 py-4 border-t bg-muted/10 text-xs flex-1">
                          <div className="space-y-4">
                            {hl.length > 0 && (
                              <div>
                                <div className="font-semibold text-primary/80 mb-1.5 flex items-center gap-1.5">
                                  <span className="text-[10px]">✨</span> 亮点优势
                                </div>
                                <ul className="space-y-1 text-muted-foreground list-disc pl-4 marker:text-primary/40">
                                  {hl.map((h, i) => <li key={i} className="line-clamp-2">{h}</li>)}
                                </ul>
                              </div>
                            )}
                            {qs.length > 0 && (
                              <div>
                                <div className="font-semibold text-primary/80 mb-1.5 flex items-center gap-1.5">
                                  <span className="text-[10px]">💡</span> 黄金三问
                                </div>
                                <ul className="space-y-1 text-muted-foreground list-disc pl-4 marker:text-primary/40">
                                  {qs.map((q, i) => <li key={i} className="line-clamp-2">{q}</li>)}
                                </ul>
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
      {!loading && list.length === 0 ? (
        <p className="rounded-xl border border-dashed p-10 text-center text-sm text-muted-foreground">
          无匹配结果，请调整筛选或关键词。
        </p>
      ) : null}
    </div>
  );
}

export default function ProductsPage() {
  return (
    <Suspense
      fallback={
        <p className="p-8 text-center text-muted-foreground">加载中…</p>
      }
    >
      <ProductsInner />
    </Suspense>
  );
}
