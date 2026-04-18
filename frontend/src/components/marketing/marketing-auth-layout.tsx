"use client";

import { usePathname, useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { PublicShell } from "@/components/marketing/public-shell";
import { FooterConfigProvider } from "@/lib/footer-config";
import { fetchMe } from "@/lib/api";
import { getToken } from "@/lib/auth";

function isPublicMarketingPath(pathname: string) {
  return (
    pathname === "/login" || pathname.startsWith("/auth/callback")
  );
}

/**
 * 企业内部：除登录与 SSO 回调外，首页与方案库等均需已登录。
 *
 * 优化：token 存在时先显示内容，后台异步验证；token 不存在时才跳转登录页。
 */
export function MarketingAuthLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const pathname = usePathname();
  const router = useRouter();
  const isPublic = isPublicMarketingPath(pathname);

  // ready 状态：null=待验证, false=验证中/无token, true=已验证
  const [ready, setReady] = useState<boolean | null>(isPublic ? true : null);

  useEffect(() => {
    if (isPublic) {
      setReady(true);
      return;
    }

    const token = getToken();
    if (!token) {
      // 无 token，立即跳转登录页
      setReady(false);
      router.replace("/login");
      return;
    }

    // token 存在，先标记为就绪显示内容，后台异步验证
    setReady(true);

    // 后台验证 token 有效性，如果失效则跳转登录页
    fetchMe()
      .then((u) => {
        if (!u) {
          router.replace("/login");
        }
      })
      .catch(() => {
        router.replace("/login");
      });
  }, [pathname, isPublic, router]);

  if (ready !== true) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-background text-muted-foreground">
        校验登录…
      </div>
    );
  }

  return (
    <FooterConfigProvider>
      <PublicShell>{children}</PublicShell>
    </FooterConfigProvider>
  );
}
