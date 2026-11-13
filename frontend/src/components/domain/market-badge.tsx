import { Badge } from "@/components/reui/badge";

export function MarketBadge({ market }: { market?: string }) {
  if (market === "domestic") return <Badge variant="info-light">国内</Badge>;
  if (market === "foreign") return <Badge variant="secondary">国外</Badge>;
  return <Badge variant="outline">不区分</Badge>;
}
