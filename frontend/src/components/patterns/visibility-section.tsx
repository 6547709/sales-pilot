import { cn } from "@/lib/utils";
import { VisibilityBadge } from "@/components/patterns/visibility-badge";

type VisibilitySectionProps = {
  kind: "customer" | "internal";
  title: string;
  children: React.ReactNode;
  className?: string;
};

export function VisibilitySection({ kind, title, children, className }: VisibilitySectionProps) {
  return (
    <section
      className={cn(
        "overflow-hidden rounded-xl border border-border bg-card",
        kind === "customer" ? "border-l-4 border-l-primary" : "border-l-4 border-l-warning",
        className,
      )}
    >
      <header className="flex items-center justify-between gap-2 border-b border-border px-4 py-3">
        <h3 className="text-sm font-semibold">{title}</h3>
        <VisibilityBadge kind={kind} />
      </header>
      <div className="p-4">{children}</div>
    </section>
  );
}
