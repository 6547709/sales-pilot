"use client";

import { usePathname, useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { PublicShell } from "@/components/marketing/public-shell";
import { fetchMe } from "@/lib/api";
import { getToken } from "@/lib/auth";

function isPublicMarketingPath(pathname: string) {
  return (
    pathname === "/login" || pathname.startsWith("/auth/callback")
  );
}

/**
 * 企业内部：除登录与 SSO 回调外，首页与方案库等均需已登录。
 */
export function MarketingAuthLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const pathname = usePathname();
  const router = useRouter();
  const isPublic = isPublicMarketingPath(pathname);
  const [ready, setReady] = useState(isPublic);

  useEffect(() => {
    if (isPublic) {
      setReady(true);
      return;
    }
    setReady(false);
    if (!getToken()) {
      router.replace("/login");
      return;
    }
    fetchMe().then((u) => {
      if (!u) router.replace("/login");
      else setReady(true);
    });
  }, [pathname, isPublic, router]);

  if (!ready) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-background text-muted-foreground">
        校验登录…
      </div>
    );
  }

  return <PublicShell>{children}</PublicShell>;
}
