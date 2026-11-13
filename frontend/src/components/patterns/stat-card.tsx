import type { LucideIcon } from "lucide-react";
import { Frame, FramePanel } from "@/components/reui/frame";
import { IconTile } from "@/components/reui/icon-tile";

type StatCardProps = {
  label: string;
  value: React.ReactNode;
  hint?: string;
  icon: LucideIcon;
};

export function StatCard({ label, value, hint, icon: Icon }: StatCardProps) {
  return (
    <Frame>
      <FramePanel>
        <div className="flex items-start justify-between gap-3">
          <div className="min-w-0">
            <p className="text-xs text-muted-foreground">{label}</p>
            <p className="mt-1 text-2xl font-semibold tabular-nums">{value}</p>
            {hint ? <p className="mt-1 text-xs text-muted-foreground">{hint}</p> : null}
          </div>
          <IconTile variant="elevated" size="lg">
            <Icon className="h-5 w-5 text-primary" />
          </IconTile>
        </div>
      </FramePanel>
    </Frame>
  );
}
