import { Lock } from "lucide-react";
import { Badge } from "@/components/reui/badge";

export function VisibilityBadge({ kind }: { kind: "customer" | "internal" }) {
  if (kind === "customer") return <Badge variant="info-light">客户可见</Badge>;
  return (
    <Badge variant="warning-light">
      <Lock className="h-3 w-3" />
      仅销售
    </Badge>
  );
}
