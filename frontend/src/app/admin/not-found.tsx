import Link from "next/link";
import { EmptyState } from "@/components/patterns/empty-state";
import { buttonVariants } from "@/components/ui/button";

export default function AdminNotFound() {
  return (
    <div className="mx-auto max-w-lg py-16">
      <EmptyState
        title="页面不存在"
        description="链接可能已失效，或该记录已被删除。"
        action={
          <Link className={buttonVariants()} href="/admin">
            返回控制台
          </Link>
        }
      />
    </div>
  );
}
