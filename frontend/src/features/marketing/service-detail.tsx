"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useParams } from "next/navigation";
import { Frame, FramePanel } from "@/components/reui/frame";
import { Badge } from "@/components/reui/badge";
import { PageHeader } from "@/components/patterns/page-header";
import { ErrorState } from "@/components/patterns/error-state";
import { LoadingState } from "@/components/patterns/loading-state";
import { servicesApi, type Scenario, type Service } from "@/lib/api";
import { useViewMode } from "@/lib/view-mode";

const KIND_LABEL: Record<string, string> = {
  v2v: "V2V 迁移",
  backup: "备份",
  dr: "容灾编排",
  migration: "迁移",
  upgrade: "升级",
  ops: "运维",
};

export function ServiceDetail() {
  const params = useParams();
  const slug = String(params.slug ?? "");
  const { mode } = useViewMode();
  const [service, setService] = useState<Service | null>(null);
  const [scenarios, setScenarios] = useState<Scenario[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [reloadKey, setReloadKey] = useState(0);

  useEffect(() => {
    let alive = true;
    setLoading(true);
    setError(null);
    Promise.all([
      servicesApi.get(slug, mode),
      servicesApi.scenarios(slug, mode).catch(() => []),
    ])
      .then(([s, sc]) => {
        if (!alive) return;
        setService(s);
        setScenarios(sc);
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
      <div className="mx-auto max-w-3xl px-4 py-8">
        <LoadingState rows={5} />
      </div>
    );
  }
  if (error || !service) {
    return (
      <div className="mx-auto max-w-3xl px-4 py-8">
        <ErrorState message={error ?? "未找到该服务"} onRetry={() => setReloadKey((k) => k + 1)} />
      </div>
    );
  }

  return (
    <div className="mx-auto flex max-w-3xl flex-col gap-5 px-4 py-6">
      <PageHeader
        title={service.name}
        description={service.summary}
        actions={<Badge variant="info-light">{KIND_LABEL[service.kind] ?? service.kind}</Badge>}
      />

      {service.valueProp ? (
        <Frame>
          <FramePanel>
            <h2 className="text-sm font-semibold">服务价值</h2>
            <p className="mt-2 whitespace-pre-wrap text-sm text-muted-foreground">
              {service.valueProp}
            </p>
          </FramePanel>
        </Frame>
      ) : null}

      <div className="grid gap-4 md:grid-cols-2">
        <Frame>
          <FramePanel>
            <h2 className="text-sm font-semibold">服务特点</h2>
            {service.features.length > 0 ? (
              <ul className="mt-2 list-disc space-y-1 pl-5 text-sm text-muted-foreground">
                {service.features.map((f) => (
                  <li key={f}>{f}</li>
                ))}
              </ul>
            ) : (
              <p className="mt-2 text-sm text-muted-foreground">—</p>
            )}
          </FramePanel>
        </Frame>
        <Frame>
          <FramePanel>
            <h2 className="text-sm font-semibold">交付物</h2>
            {service.deliverables.length > 0 ? (
              <ul className="mt-2 list-disc space-y-1 pl-5 text-sm text-muted-foreground">
                {service.deliverables.map((d) => (
                  <li key={d}>{d}</li>
                ))}
              </ul>
            ) : (
              <p className="mt-2 text-sm text-muted-foreground">—</p>
            )}
          </FramePanel>
        </Frame>
      </div>

      {scenarios.length > 0 ? (
        <div className="flex flex-col gap-2">
          <h2 className="text-sm font-semibold">适用场景</h2>
          <div className="flex flex-wrap gap-1.5">
            {scenarios.map((sc) => (
              <Link key={sc.id} href={`/scenarios/${sc.slug}`}>
                <Badge variant="secondary">{sc.name}</Badge>
              </Link>
            ))}
          </div>
        </div>
      ) : null}

      {service.highlights.length > 0 ? (
        <div className="flex flex-wrap gap-1.5">
          {service.highlights.map((h) => (
            <Badge key={h} variant="outline">
              {h}
            </Badge>
          ))}
        </div>
      ) : null}
    </div>
  );
}
