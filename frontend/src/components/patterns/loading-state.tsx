import { Skeleton } from "@/components/ui/skeleton";

type LoadingStateProps = {
  rows?: number;
};

export function LoadingState({ rows = 5 }: LoadingStateProps) {
  return (
    <div className="flex flex-col gap-2" aria-busy="true" aria-live="polite">
      {Array.from({ length: rows }).map((_, i) => (
        <Skeleton key={i} className="h-10 w-full" />
      ))}
    </div>
  );
}
