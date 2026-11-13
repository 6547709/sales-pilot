import { LoadingState } from "@/components/patterns/loading-state";

export default function AdminLoading() {
  return (
    <div className="flex flex-col gap-4">
      <LoadingState rows={6} />
    </div>
  );
}
