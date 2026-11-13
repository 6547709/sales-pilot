import Link from "next/link";
import { ArrowRight } from "lucide-react";
import { Frame, FramePanel } from "@/components/reui/frame";
import type { ScenarioSummary } from "@/lib/api";

export function ScenarioCard({ item }: { item: ScenarioSummary }) {
  return (
    <Link href={`/scenarios/${item.scenario.slug}`} className="block h-full">
      <Frame className="h-full transition-colors hover:border-primary/40">
        <FramePanel>
          <div className="flex h-full flex-col gap-2">
            <p className="font-medium">{item.scenario.name}</p>
            <p className="line-clamp-2 min-h-8 text-xs text-muted-foreground">
              {item.scenario.summary}
            </p>
            <div className="mt-auto flex items-center justify-between gap-2 text-xs text-muted-foreground">
              <span>
                服务 {item.serviceCount} · 产品 {item.productCount}
              </span>
              <span className="flex items-center gap-1 font-medium text-primary">
                查看方案
                <ArrowRight className="h-3.5 w-3.5" />
              </span>
            </div>
          </div>
        </FramePanel>
      </Frame>
    </Link>
  );
}
