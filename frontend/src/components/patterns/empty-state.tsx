import type { LucideIcon } from "lucide-react";
import { IconStack } from "@/components/reui/icon-stack";
import { IconTile } from "@/components/reui/icon-tile";

type EmptyStateProps = {
  icon?: LucideIcon;
  title: string;
  description?: string;
  action?: React.ReactNode;
};

export function EmptyState({ icon: Icon, title, description, action }: EmptyStateProps) {
  return (
    <div className="flex flex-col items-center justify-center gap-3 px-6 py-12 text-center">
      <IconStack>
        <IconTile variant="elevated" size="lg">
          {Icon ? <Icon className="h-5 w-5 text-muted-foreground" /> : null}
        </IconTile>
      </IconStack>
      <div className="space-y-1">
        <p className="text-sm font-medium">{title}</p>
        {description ? (
          <p className="text-xs text-muted-foreground">{description}</p>
        ) : null}
      </div>
      {action}
    </div>
  );
}
