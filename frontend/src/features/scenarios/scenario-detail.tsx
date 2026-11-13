"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useParams } from "next/navigation";
import { ChevronRight } from "lucide-react";
import { Frame, FramePanel } from "@/components/reui/frame";
import { Badge } from "@/components/reui/badge";
import { PageHeader } from "@/components/patterns/page-header";
import { ErrorState } from "@/components/patterns/error-state";
import { LoadingState } from "@/components/patterns/loading-state";
import { ProductCard } from "@/components/domain/product-card";
import { scenariosApi, type ScenarioDetail } from "@/lib/api";
import { useViewMode } from "@/lib/view-mode";

export function ScenarioDetailView() {
  const params = useParams();
  const slug = String(params.slug ?? "");
  const { mode } = useViewMode();
  const [detail, setDetail] = useState<ScenarioDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [reloadKey, setReloadKey] = useState(0);

  useEffect(() => {
    let alive = true;
    setLoading(true);
    setError(null);
    scenariosApi
      .get(slug, mode)
      .then((d) => {
        if (alive) setDetail(d);
      })
      .catch((e) => {
        if (alive) setError(e instanceof Error ? e.message : "加载失败");
      })
      .finally(() => {
        if (alive) setLoading(false);
      });
    return () => {
      alive = false;
    };
  }, [slug, mode, reloadKey]);

  if (loading) {
    return (
      <div className="mx-auto max-w-[1200px] px-4 py-8">
        <LoadingState rows={6} />
      </div>
    );
  }
  if (error || !detail) {
    return (
      <div className="mx-auto max-w-[1200px] px-4 py-8">
        <ErrorState message={error ?? "未找到该场景"} onRetry={() => setReloadKey((k) => k + 1)} />
      </div>
    );
  }

  const { scenario, group, services, products, cases } = detail;
  const categories = new Map<number, string>();
  for (const p of products) {
    if (p.solutionCategory?.id) categories.set(p.solutionCategory.id, p.solutionCategory.label);
  }

  return (
    <div className="mx-auto flex max-w-[1200px] flex-col gap-5 px-4 py-6">
      <div className="flex flex-wrap items-center gap-1 text-xs text-muted-foreground">
        <Link href="/scenarios" className="hover:text-foreground">
          按场景
        </Link>
        <ChevronRight className="h-3 w-3" />
        {group ? <span>{group.name}</span> : null}
        <ChevronRight className="h-3 w-3" />
        <span className="text-foreground">{scenario.name}</span>
      </div>

      <PageHeader
        title={scenario.name}
        description={scenario.summary}
        actions={group ? <Badge variant="info-light">{group.name}</Badge> : undefined}
      />

      <div className="grid gap-4 md:grid-cols-2">
        <Frame>
          <FramePanel>
            <h2 className="text-sm font-semibold">客户痛点</h2>
            <p className="mt-2 whitespace-pre-wrap text-sm text-muted-foreground">
              {scenario.painPoint || "—"}
            </p>
          </FramePanel>
        </Frame>
        <Frame>
          <FramePanel>
            <h2 className="text-sm font-semibold">我方价值</h2>
            <p className="mt-2 whitespace-pre-wrap text-sm text-muted-foreground">
              {scenario.valueProp || "—"}
            </p>
          </FramePanel>
        </Frame>
      </div>

      <section className="flex flex-col gap-2">
        <h2 className="text-sm font-semibold">推荐服务 · {services.length}</h2>
        {services.length === 0 ? (
          <p className="text-xs text-muted-foreground">暂无关联服务。</p>
        ) : (
          <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
            {services.map((s) => (
              <Link key={s.id} href={`/services/${s.slug}`} className="block h-full">
                <Frame className="h-full transition-colors hover:border-primary/40">
                  <FramePanel>
                    <p className="font-medium">{s.name}</p>
                    <p className="mt-1 line-clamp-2 text-xs text-muted-foreground">{s.summary}</p>
                    <div className="mt-2 flex flex-wrap gap-1">
                      {s.highlights.slice(0, 3).map((h) => (
                        <Badge key={h} variant="outline">
                          {h}
                        </Badge>
                      ))}
                    </div>
                  </FramePanel>
                </Frame>
              </Link>
            ))}
          </div>
        )}
      </section>

      <section className="flex flex-col gap-2">
        <h2 className="text-sm font-semibold">推荐产品 · {products.length}</h2>
        {products.length === 0 ? (
          <p className="text-xs text-muted-foreground">暂无关联产品。</p>
        ) : (
          <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
            {products.map((p) => (
              <ProductCard key={p.id} product={p} href={`/products/${p.id}`} />
            ))}
          </div>
        )}
      </section>

      {cases.length > 0 ? (
        <section className="flex flex-col gap-2">
          <h2 className="text-sm font-semibold">行业案例</h2>
          <div className="grid gap-3 md:grid-cols-2">
            {cases.map((c) => (
              <Frame key={c.id}>
                <FramePanel>
                  <div className="flex flex-wrap items-center gap-2">
                    {c.industry ? <Badge variant="info-light">{c.industry}</Badge> : null}
                    <span className="font-medium">{c.clientName}</span>
                  </div>
                  <p className="mt-2 text-sm text-muted-foreground">{c.solution}</p>
                  <p className="mt-1 text-xs text-muted-foreground">{c.valueDelivered}</p>
                </FramePanel>
              </Frame>
            ))}
          </div>
        </section>
      ) : null}

      {categories.size > 0 ? (
        <section className="flex flex-col gap-2">
          <h2 className="text-sm font-semibold">相关解决方案</h2>
          <div className="flex flex-wrap gap-1.5">
            {Array.from(categories.entries()).map(([id, label]) => (
              <Link key={id} href={`/?solution=${id}`}>
                <Badge variant="secondary">{label}</Badge>
              </Link>
            ))}
          </div>
        </section>
      ) : null}

    </div>
  );
}
