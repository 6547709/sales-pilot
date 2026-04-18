"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { Button, buttonVariants } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { useFooterConfig } from "@/lib/footer-config";
import type { User } from "@/lib/api";
import { clearToken } from "@/lib/auth";
import {
  DatabaseBackup,
  KeyRound,
  LayoutDashboard,
  Network,
  Package,
  Settings,
  Shield,
  Users,
} from "lucide-react";

const adminNav = [
  { href: "/admin", label: "控制台", icon: LayoutDashboard },
  { href: "/admin/products", label: "产品维护", icon: Package },
  { href: "/admin/topology", label: "拓扑维护", icon: Network },
  { href: "/admin/api-keys", label: "API 密钥", icon: KeyRound },
  { href: "/admin/users", label: "本地用户", icon: Users },
  { href: "/admin/backup", label: "备份恢复", icon: DatabaseBackup },
  { href: "/admin/settings", label: "认证配置", icon: Settings },
];

export function AdminShell({
  user,
  children,
}: {
  user: User;
  children: React.ReactNode;
}) {
  const pathname = usePathname();
  const { footerText, filingNumber } = useFooterConfig();

  return (
    <div className="min-h-screen bg-muted/30">
      <div className="flex min-h-screen">
        <aside className="hidden w-56 shrink-0 border-r border-border bg-card md:flex md:flex-col">
          <div className="flex h-14 items-center gap-2 border-b px-4">
            <Shield className="h-5 w-5 text-primary" />
            <span className="font-semibold text-primary">管理后台</span>
          </div>
          <nav className="flex flex-1 flex-col gap-0.5 p-2">
            {adminNav.map((item) => {
              const active =
                pathname === item.href ||
                (item.href !== "/admin" && pathname.startsWith(item.href));
              const Icon = item.icon;
              return (
                <Link
                  key={item.href}
                  href={item.href}
                  className={cn(
                    "flex items-center gap-2 rounded-lg px-3 py-2.5 text-sm font-medium transition-colors",
                    active
                      ? "bg-primary/10 text-primary"
                      : "text-muted-foreground hover:bg-muted hover:text-foreground",
                  )}
                >
                  <Icon className="h-4 w-4 shrink-0" />
                  {item.label}
                </Link>
              );
            })}
          </nav>
          <div className="border-t p-3">
            <p className="truncate text-xs text-muted-foreground">{user.username}</p>
            <Button
              variant="outline"
              size="sm"
              className="mt-2 w-full"
              onClick={() => {
                clearToken();
                window.location.href = "/login";
              }}
            >
              退出
            </Button>
          </div>
        </aside>

        <div className="flex min-w-0 flex-1 flex-col">
          <header className="sticky top-0 z-40 flex h-14 items-center justify-between gap-4 border-b bg-background/95 px-4 backdrop-blur md:px-6">
            <div className="flex items-center gap-3 md:hidden">
              <Shield className="h-5 w-5 text-primary" />
              <span className="font-semibold">管理</span>
            </div>
            <div className="hidden md:block" />
            <div className="flex items-center gap-2">
              <span className="hidden text-xs text-muted-foreground sm:inline">
                {user.username}
              </span>
              <Link
                href="/"
                className={cn(buttonVariants({ variant: "ghost", size: "sm" }))}
              >
                返回前台
              </Link>
            </div>
          </header>

          <main className="flex-1 p-4 md:p-6">{children}</main>
        </div>
      </div>

      {(footerText || filingNumber) && (
        <div className="fixed bottom-16 left-0 right-0 z-40 border-t bg-muted/50 py-2 text-center text-xs text-muted-foreground md:hidden">
          {footerText}
          {footerText && filingNumber && <span className="mx-2">|</span>}
          {filingNumber}
        </div>
      )}
      {/* 移动端底部导航 */}
      <nav className="fixed bottom-0 left-0 right-0 z-40 flex border-t bg-card/95 backdrop-blur md:hidden">
        {adminNav.map((item) => {
          const active =
            pathname === item.href ||
            (item.href !== "/admin" && pathname.startsWith(item.href));
          const Icon = item.icon;
          return (
            <Link
              key={item.href}
              href={item.href}
              className={cn(
                "flex flex-1 flex-col items-center gap-0.5 py-2 text-[10px] font-medium",
                active ? "text-primary" : "text-muted-foreground",
              )}
            >
              <Icon className="h-5 w-5" />
              {item.label}
            </Link>
          );
        })}
      </nav>
      <div className="h-16 md:hidden" />
    </div>
  );
}
