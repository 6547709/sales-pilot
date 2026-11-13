"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { ErrorState } from "@/components/patterns/error-state";
import { ApiError } from "@/lib/api";

export default function AdminError({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  const router = useRouter();
  const status = error instanceof ApiError ? error.status : undefined;

  useEffect(() => {
    if (status === 401) router.replace("/login");
  }, [status, router]);

  if (status === 401) {
    return <p className="p-6 text-sm text-muted-foreground">登录已过期，正在跳转…</p>;
  }
  if (status === 403) {
    return (
      <ErrorState title="没有访问权限" message="当前账号无权查看该页面。" onRetry={reset} />
    );
  }
  return <ErrorState message={error.message || "请求失败，请稍后重试。"} onRetry={reset} />;
}
