"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { Button, buttonVariants } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { mainNav } from "@/config/nav";
import type { User } from "@/lib/api";
import { clearToken } from "@/lib/auth";

export function SiteHeader({ user }: { user: User | null }) {
  const pathname = usePathname();
  const items = mainNav.filter(
    (n) => !n.adminOnly || user?.role === "admin",
  );

  return (
    <header className="sticky top-0 z-40 border-b bg-card/95 backdrop-blur">
      <div className="mx-auto flex max-w-6xl items-center gap-4 px-4 py-3">
        <Link
          href="/products"
          className="text-lg font-semibold text-primary"
        >
          Sales-Pilot
        </Link>
        <nav className="flex flex-1 flex-wrap gap-1 text-sm">
          {items.map((n) => (
            <Link
              key={n.href}
              href={n.href}
              className={
                pathname === n.href || pathname.startsWith(n.href + "/")
                  ? "rounded-md bg-primary/10 px-3 py-1.5 font-medium text-primary"
                  : "rounded-md px-3 py-1.5 text-muted-foreground hover:bg-muted"
              }
            >
              {n.label}
            </Link>
          ))}
        </nav>
        <div className="flex items-center gap-2">
          {user ? (
            <>
              <span className="hidden text-xs text-muted-foreground sm:inline">
                {user.username}
              </span>
              <Button
                variant="outline"
                size="sm"
                onClick={() => {
                  clearToken();
                  window.location.href = "/login";
                }}
              >
                退出
              </Button>
            </>
          ) : (
            <Link
              href="/login"
              className={cn(buttonVariants({ variant: "default", size: "sm" }))}
            >
              登录
            </Link>
          )}
        </div>
      </div>
    </header>
  );
}
