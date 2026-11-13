"use client";

import { useEffect, useMemo, useState } from "react";
import { Compass } from "lucide-react";
import { Input } from "@/components/ui/input";
import { Frame, FramePanel } from "@/components/reui/frame";
import { PageHeader } from "@/components/patterns/page-header";
import { EmptyState } from "@/components/patterns/empty-state";
import { ErrorState } from "@/components/patterns/error-state";
import { LoadingState } from "@/components/patterns/loading-state";
import { ScenarioCard } from "@/components/domain/scenario-card";
import { scenariosApi, type ScenarioGroup } from "@/lib/api";
import { useViewMode } from "@/lib/view-mode";

export function ScenarioOverview() {
  const { mode } = useViewMode();
  const [groups, setGroups] = useState<ScenarioGroup[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [activeGroup, setActiveGroup] = useState("all");
  const [q, setQ] = useState("");

  useEffect(() => {
    let alive = true;
    setLoading(true);
    setError(null);
    scenariosApi
      .groups(mode)
      .then((list) => {
        if (alive) setGroups(list);
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
  }, [mode]);

  const visible = useMemo(() => {
    const kw = q.trim().toLowerCase();
    return groups
      .map((g) => ({
        ...g,
        scenarios: g.scenarios.filter((item) => {
          if (item.productCount + item.serviceCount === 0) return false;
          if (!kw) return true;
          return (
            item.scenario.name.toLowerCase().includes(kw) ||
            item.scenario.summary.toLowerCase().includes(kw)
          );
        }),
      }))
      .filter((g) => g.scenarios.length > 0)
      .filter((g) => activeGroup === "all" || g.slug === activeGroup);
  }, [groups, q, activeGroup]);

  return (
    <div className="mx-auto flex max-w-[1200px] flex-col gap-5 px-4 py-6">
      <PageHeader
        title="按场景找方案"
        description="从客户真实需要出发，找到对应的服务、产品与案例。"
        actions={
          <Input
            value={q}
            onChange={(e) => setQ(e.target.value)}
            placeholder="搜索场景或关键字"
            className="w-64"
          />
        }
      />

      <Frame>
        <FramePanel>
          <div className="flex flex-wrap gap-1.5">
            <button
              type="button"
              onClick={() => setActiveGroup("all")}
              className={
                activeGroup === "all"
                  ? "rounded-full bg-primary/10 px-3 py-1 text-xs font-semibold text-primary"
                  : "rounded-full px-3 py-1 text-xs text-muted-foreground hover:bg-muted"
              }
            >
              全部
            </button>
            {groups.map((g) => (
              <button
                key={g.id}
                type="button"
                onClick={() => setActiveGroup(g.slug)}
                className={
                  activeGroup === g.slug
                    ? "rounded-full bg-primary/10 px-3 py-1 text-xs font-semibold text-primary"
                    : "rounded-full px-3 py-1 text-xs text-muted-foreground hover:bg-muted"
                }
              >
                {g.name}
              </button>
            ))}
          </div>
        </FramePanel>
      </Frame>

      {loading ? <LoadingState rows={6} /> : null}
      {error ? <ErrorState message={error} onRetry={() => setActiveGroup(activeGroup)} /> : null}

      {!loading && !error && visible.length === 0 ? (
        <EmptyState
          icon={Compass}
          title="没有匹配的场景"
          description="换个关键字，或等待管理员在后台补充场景。"
        />
      ) : null}

      {visible.map((g) => (
        <section key={g.id} className="flex flex-col gap-2">
          <h2 className="text-sm font-semibold text-foreground">{g.name}</h2>
          {g.description ? (
            <p className="-mt-1 text-xs text-muted-foreground">{g.description}</p>
          ) : null}
          <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
            {g.scenarios.map((item) => (
              <ScenarioCard key={item.scenario.id} item={item} />
            ))}
          </div>
        </section>
      ))}
    </div>
  );
}
