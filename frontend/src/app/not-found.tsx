import Link from "next/link";
import { EmptyState } from "@/components/patterns/empty-state";
import { buttonVariants } from "@/components/ui/button";

export default function NotFound() {
  return (
    <div className="mx-auto max-w-lg py-16">
      <EmptyState
        title="页面不存在"
        description="链接可能已失效。"
        action={
          <Link className={buttonVariants()} href="/">
            返回首页
          </Link>
        }
      />
    </div>
  );
}
