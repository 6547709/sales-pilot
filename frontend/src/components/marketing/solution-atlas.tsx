"use client";

import { motion } from "framer-motion";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useCallback, useEffect, useMemo, useState } from "react";
import {
  ArchCardButton,
  type ArchCardManufacturerChip,
  type ArchCardProductChips,
} from "@/components/marketing/arch-card";
import {
  apiFetch,
  fetchTopology,
  fetchProducts,
  type Product,
  type TopologyCategoryNode,
  type TopologyResponse,
} from "@/lib/api";
import { topologyIcon } from "@/lib/topology-icons";
import { Badge } from "@/components/ui/badge";
import { buttonVariants } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { cn } from "@/lib/utils";
import { ChevronRight, Compass, Sparkles } from "lucide-react";

type Props = {
  initialSolutionId?: string;
  initialKeyword?: string;
  /** 当前 URL 中的 market=domestic|foreign（与下方赋能区筛选联动） */
  initialMarket?: string;
  /** 由父组件同步 URL 与赋能区（首页传入）；未传时仍用内部 router.push */
  onSelectSolution?: (cat: TopologyCategoryNode | null) => void;
};

function flattenCategories(t: TopologyResponse): TopologyCategoryNode[] {
  const a: TopologyCategoryNode[] = [...t.security, ...t.ops];
  for (const block of t.central_layers) {
    a.push(...block.categories);
  }
  return a;
}

function findById(
  t: TopologyResponse | null,
  id: string,
): TopologyCategoryNode | null {
  if (!t || !id) return null;
  const n = parseInt(id, 10);
  if (Number.isNaN(n)) return null;
  return flattenCategories(t).find((c) => c.id === n) ?? null;
}

/** 同一方案、同一市场桶内按厂商聚合；无厂商名则逐条显示产品名并链到详情 */
function buildBucketChips(
  products: Product[],
  categoryId: number,
): ArchCardManufacturerChip[] {
  const groups = new Map<string, Product[]>();
  for (const p of products) {
    const raw = (p.manufacturer_name || "").trim();
    const dedupeKey = raw ? raw.toLowerCase() : `__n_${p.id}`;
    if (!groups.has(dedupeKey)) groups.set(dedupeKey, []);
    groups.get(dedupeKey)!.push(p);
  }
  const chips: ArchCardManufacturerChip[] = [];
  for (const [, group] of groups) {
    const first = group[0];
    const mname = (first.manufacturer_name || "").trim();
    if (mname) {
      const q = new URLSearchParams({
        solution_category_id: String(categoryId),
        manufacturer: mname,
      });
      chips.push({ label: mname, href: `/products?${q}` });
    } else {
      for (const p of group) {
        chips.push({ label: p.name, href: `/products/${p.id}` });
      }
    }
  }
  chips.sort((a, b) => a.label.localeCompare(b.label, "zh-CN"));
  return chips;
}

/** 按解决方案 ID 分组，再按国内/国外/不区分输出厂商标签 */
function buildProductChipMap(list: Product[]): Map<number, ArchCardProductChips> {
  const raw = new Map<
    number,
    { domestic: Product[]; foreign: Product[]; neutral: Product[] }
  >();
  for (const p of list) {
    if (p.is_draft) continue;
    const sid = p.solution_category_id;
    if (sid == null || sid === undefined) continue;
    if (!raw.has(sid)) {
      raw.set(sid, { domestic: [], foreign: [], neutral: [] });
    }
    const g = raw.get(sid)!;
    const vm = (p.vendor_market || "all").toLowerCase().trim();
    if (vm === "domestic") g.domestic.push(p);
    else if (vm === "foreign") g.foreign.push(p);
    else g.neutral.push(p);
  }
  const m = new Map<number, ArchCardProductChips>();
  for (const [sid, buckets] of raw) {
    m.set(sid, {
      domestic: buildBucketChips(buckets.domestic, sid),
      foreign: buildBucketChips(buckets.foreign, sid),
      neutral: buildBucketChips(buckets.neutral, sid),
    });
  }
  return m;
}

/** 无产品时勿传 productChips，避免空灰条 */
function chipsFor(
  map: Map<number, ArchCardProductChips>,
  categoryId: number,
): ArchCardProductChips | undefined {
  const v = map.get(categoryId);
  if (!v) return undefined;
  const n =
    v.domestic.length + v.foreign.length + v.neutral.length;
  return n > 0 ? v : undefined;
}

function findByKeyword(
  t: TopologyResponse | null,
  kw: string,
): TopologyCategoryNode | null {
  if (!t || !kw.trim()) return null;
  const k = kw.trim().toLowerCase();
  return (
    flattenCategories(t).find((c) =>
      c.keywords.some(
        (x) =>
          k.includes(String(x).toLowerCase()) ||
          String(x).toLowerCase().includes(k),
      ),
    ) ?? null
  );
}

export function SolutionAtlas({
  initialSolutionId,
  initialKeyword,
  initialMarket = "",
  onSelectSolution,
}: Props) {
  const router = useRouter();
  const [topology, setTopology] = useState<TopologyResponse | null>(null);
  const [loadErr, setLoadErr] = useState<string | null>(null);
  const [active, setActive] = useState<TopologyCategoryNode | null>(null);
  const [productChipMap, setProductChipMap] = useState<
    Map<number, ArchCardProductChips>
  >(() => new Map());

  useEffect(() => {
    // 并行加载拓扑和产品数据，利用缓存加速
    Promise.all([
      fetchTopology()
        .then(setTopology)
        .catch(() => setLoadErr("无法加载拓扑，请确认后端已启动并已执行数据库迁移。")),
      fetchProducts()
        .then((list) =>
          setProductChipMap(buildProductChipMap(list)),
        )
        .catch(() => setProductChipMap(new Map())),
    ]);
  }, []);

  useEffect(() => {
    if (!topology) return;
    const bySol = findById(topology, initialSolutionId ?? "");
    if (bySol) {
      setActive(bySol);
      return;
    }
    const byKw = findByKeyword(topology, initialKeyword ?? "");
    setActive(byKw);
  }, [topology, initialSolutionId, initialKeyword]);

  const hasMarket =
    initialMarket === "domestic" || initialMarket === "foreign";

  const applyFilter = useCallback(
    (cat: TopologyCategoryNode | null) => {
      setActive(cat);
      if (onSelectSolution) {
        onSelectSolution(cat);
        return;
      }
      const sp = new URLSearchParams();
      if (cat) sp.set("solution", String(cat.id));
      const qs = sp.toString();
      router.push(qs ? `/?${qs}` : "/", { scroll: false });
      requestAnimationFrame(() => {
        document.getElementById("solution-spotlight")?.scrollIntoView({
          behavior: "smooth",
          block: "start",
        });
      });
    },
    [router, onSelectSolution],
  );

  const archMainClick = useCallback(
    (c: TopologyCategoryNode) => {
      if (active?.id === c.id) {
        if (hasMarket) {
          applyFilter(c);
        } else {
          applyFilter(null);
        }
      } else {
        applyFilter(c);
      }
    },
    [active, hasMarket, applyFilter],
  );

  const centralSorted = useMemo(() => {
    if (!topology) return [];
    return [...topology.central_layers].sort(
      (a, b) => b.layer.level - a.layer.level,
    );
  }, [topology]);

  if (loadErr) {
    return (
      <div className="mx-auto max-w-lg px-4 py-16 text-center text-sm text-destructive">
        {loadErr}
      </div>
    );
  }

  if (!topology) {
    return (
      <div className="relative overflow-hidden">
        <div className="pointer-events-none absolute inset-0 bg-[radial-gradient(ellipse_80%_50%_at_50%_-20%,rgba(0,82,204,0.12),transparent)]" />
        <div className="relative mx-auto max-w-[1400px] px-4 py-8 sm:py-12">
          <div className="mb-8 text-center sm:mb-10">
            <Skeleton className="mb-3 h-6 w-24 mx-auto" />
            <Skeleton className="h-10 w-80 mx-auto sm:w-96" />
          </div>
          <div className="grid gap-6 lg:grid-cols-[minmax(0,240px)_1fr_minmax(0,240px)] lg:gap-4 xl:gap-6">
            <div className="rounded-2xl border border-primary/15 bg-gradient-to-b from-slate-50/90 to-white p-4 shadow-sm dark:from-slate-950/40 dark:to-background">
              <Skeleton className="mb-1 h-5 w-20" />
              <Skeleton className="mb-4 h-3 w-36" />
              <div className="grid grid-cols-2 gap-2 sm:grid-cols-1">
                {[...Array(4)].map((_, i) => (
                  <Skeleton key={i} className="h-10 w-full" />
                ))}
              </div>
            </div>
            <div className="space-y-3">
              {[...Array(2)].map((_, i) => (
                <Skeleton key={i} className="h-28 w-full rounded-2xl" />
              ))}
            </div>
            <div className="rounded-2xl border border-primary/15 bg-gradient-to-b from-slate-50/90 to-white p-4 shadow-sm dark:from-slate-950/40 dark:to-background">
              <Skeleton className="mb-1 h-5 w-20" />
              <Skeleton className="mb-4 h-3 w-36" />
              <div className="grid grid-cols-2 gap-2 sm:grid-cols-1">
                {[...Array(4)].map((_, i) => (
                  <Skeleton key={i} className="h-10 w-full" />
                ))}
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="relative overflow-hidden">
      <div className="pointer-events-none absolute inset-0 bg-[radial-gradient(ellipse_80%_50%_at_50%_-20%,rgba(0,82,204,0.12),transparent)]" />
      <div className="relative mx-auto max-w-[1400px] px-4 py-8 sm:py-12">
        <motion.div
          initial={{ opacity: 0, y: 12 }}
          animate={{ opacity: 1, y: 0 }}
          className="mb-8 text-center sm:mb-10"
        >
          <Badge variant="secondary" className="mb-3 gap-1">
            <Compass className="h-3 w-3" />
            全局视角
          </Badge>
          <h1 className="text-balance text-3xl font-bold tracking-tight text-foreground sm:text-4xl md:text-5xl">
            按架构层级与体系
            <span className="text-primary"> 选择解决方案</span>
          </h1>
        </motion.div>

        <div className="grid gap-6 lg:grid-cols-[minmax(0,240px)_1fr_minmax(0,240px)] lg:gap-4 xl:gap-6">
          <motion.aside
            initial={{ opacity: 0, x: -16 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ delay: 0.05 }}
            className="rounded-2xl border border-primary/15 bg-gradient-to-b from-slate-50/90 to-white p-4 shadow-sm dark:from-slate-950/40 dark:to-background"
          >
            <h2 className="mb-1 flex items-center gap-2 text-sm font-bold text-primary">
              <span className="h-2 w-2 rounded-full bg-primary" />
              安全体系
            </h2>
            <p className="mb-4 text-xs text-muted-foreground">
              身份、边界与数据全生命周期
            </p>
            <ul className="grid grid-cols-2 gap-2 sm:grid-cols-1">
              {topology.security.map((c) => (
                <li key={c.id}>
                  <ArchCardButton
                    icon={topologyIcon(c.icon_key)}
                    label={c.label}
                    compact
                    selected={active?.id === c.id}
                    hint={c.hint || null}
                    onClick={() => archMainClick(c)}
                    productChips={chipsFor(productChipMap, c.id)}
                  />
                </li>
              ))}
            </ul>
          </motion.aside>

          <div className="space-y-3">
            {centralSorted.map((block, idx) => (
              <motion.section
                key={block.layer.id}
                initial={{ opacity: 0, y: 16 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.08 + idx * 0.04 }}
                className="rounded-2xl border border-blue-200/60 bg-gradient-to-r from-[#0052cc]/[0.07] via-background to-[#0052cc]/[0.05] p-4 shadow-sm dark:border-blue-900/40"
              >
                <div className="mb-3 flex flex-wrap items-end justify-between gap-2 border-b border-primary/10 pb-2">
                  <div>
                    <h3 className="text-base font-bold text-foreground sm:text-lg">
                      {block.layer.title}
                    </h3>
                    {block.layer.subtitle ? (
                      <p className="text-xs text-muted-foreground">
                        {block.layer.subtitle}
                      </p>
                    ) : null}
                  </div>
                  <Sparkles className="hidden h-5 w-5 text-primary/40 sm:block" />
                </div>
                <div className="grid grid-cols-2 gap-2 sm:grid-cols-4">
                  {block.categories.map((c) => (
                    <ArchCardButton
                      key={c.id}
                      icon={topologyIcon(c.icon_key)}
                      label={c.label}
                      selected={active?.id === c.id}
                      hint={c.hint || null}
                      onClick={() => archMainClick(c)}
                      productChips={chipsFor(productChipMap, c.id)}
                    />
                  ))}
                </div>
              </motion.section>
            ))}
            <div className="flex justify-center pt-2">
              <Link
                href="/products"
                className={cn(
                  buttonVariants({ variant: "outline", size: "sm" }),
                  "inline-flex gap-1",
                )}
              >
                全部方案库
                <ChevronRight className="h-4 w-4" />
              </Link>
            </div>
          </div>

          <motion.aside
            initial={{ opacity: 0, x: 16 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ delay: 0.05 }}
            className="rounded-2xl border border-primary/15 bg-gradient-to-b from-slate-50/90 to-white p-4 shadow-sm dark:from-slate-950/40 dark:to-background"
          >
            <h2 className="mb-1 flex items-center gap-2 text-sm font-bold text-primary">
              <span className="h-2 w-2 rounded-full bg-primary" />
              运维体系
            </h2>
            <p className="mb-4 text-xs text-muted-foreground">
              可观测、自动化与连续性
            </p>
            <ul className="grid grid-cols-2 gap-2 sm:grid-cols-1">
              {topology.ops.map((c) => (
                <li key={c.id}>
                  <ArchCardButton
                    icon={topologyIcon(c.icon_key)}
                    label={c.label}
                    compact
                    selected={active?.id === c.id}
                    hint={c.hint || null}
                    onClick={() => archMainClick(c)}
                    productChips={chipsFor(productChipMap, c.id)}
                  />
                </li>
              ))}
            </ul>
          </motion.aside>
        </div>
      </div>
    </div>
  );
}
