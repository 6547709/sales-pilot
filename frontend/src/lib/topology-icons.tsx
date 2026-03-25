import type { LucideIcon } from "lucide-react";
import * as Icons from "lucide-react";
import { LayoutGrid } from "lucide-react";

/** 将后台维护的 icon_key 映射为 Lucide 图标 */
export function topologyIcon(key: string): LucideIcon {
  if (!key) return LayoutGrid;
  const Icon = (Icons as unknown as Record<string, LucideIcon | undefined>)[
    key
  ];
  return Icon ?? LayoutGrid;
}
