"use client";

import { useEffect, useState, Suspense } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { fetchMe } from "@/lib/api";
import { setToken } from "@/lib/auth";

function CallbackInner() {
  const router = useRouter();
  const sp = useSearchParams();
  const [msg, setMsg] = useState("正在完成登录…");

  useEffect(() => {
    const token = sp.get("token");
    if (!token) {
      setMsg("缺少令牌，请重试 SSO。");
      return;
    }
    setToken(token);
    fetchMe().then((u) => {
      router.replace(u?.role === "admin" ? "/admin" : "/");
    });
  }, [sp, router]);

  return (
    <p className="p-8 text-center text-muted-foreground">{msg}</p>
  );
}

export default function AuthCallbackPage() {
  return (
    <Suspense
      fallback={
        <p className="p-8 text-center text-muted-foreground">加载中…</p>
      }
    >
      <CallbackInner />
    </Suspense>
  );
}
