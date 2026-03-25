"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useEffect, useState } from "react";
import { Button, buttonVariants } from "@/components/ui/button";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";
import { fetchMe, type User } from "@/lib/api";
import { clearToken, getToken } from "@/lib/auth";
import { cn } from "@/lib/utils";
import { Menu } from "lucide-react";

const publicLinks = [
  { href: "/", label: "首页" },
  { href: "/products", label: "方案库" },
];

export function PublicShell({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const [user, setUser] = useState<User | null | undefined>(undefined);

  useEffect(() => {
    const t = getToken();
    if (!t) {
      setUser(null);
      return;
    }
    fetchMe().then(setUser);
  }, [pathname]);

  return (
    <div className="flex min-h-screen flex-col bg-background">
      <header className="sticky top-0 z-50 border-b border-border/80 bg-background/90 backdrop-blur-md">
        <div className="mx-auto flex h-14 max-w-[1400px] items-center justify-between gap-4 px-4">
          <Link href="/" className="flex items-center gap-2 font-semibold text-primary">
            <span className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary text-sm text-primary-foreground">
              SP
            </span>
            <span className="hidden sm:inline">Sales-Pilot</span>
          </Link>

          <nav className="hidden items-center gap-1 md:flex">
            {publicLinks.map((l) => (
              <Link
                key={l.href}
                href={l.href}
                className={cn(
                  "rounded-lg px-3 py-2 text-sm font-medium transition-colors",
                  pathname === l.href
                    ? "bg-primary/10 text-primary"
                    : "text-muted-foreground hover:bg-muted hover:text-foreground",
                )}
              >
                {l.label}
              </Link>
            ))}
          </nav>

          <div className="flex items-center gap-2">
            {user === undefined ? null : user ? (
              <>
                <span className="hidden max-w-[120px] truncate text-xs text-muted-foreground lg:inline">
                  {user.username}
                </span>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => {
                    clearToken();
                    setUser(null);
                    window.location.href = "/login";
                  }}
                >
                  退出
                </Button>
              </>
            ) : (
              <Link
                href="/login"
                className={cn(buttonVariants({ size: "sm" }))}
              >
                登录
              </Link>
            )}

            <Sheet>
              <SheetTrigger
                className={cn(
                  buttonVariants({ variant: "ghost", size: "icon" }),
                  "md:hidden",
                )}
              >
                <Menu className="h-5 w-5" />
                <span className="sr-only">菜单</span>
              </SheetTrigger>
              <SheetContent side="right" className="w-[280px]">
                <SheetHeader>
                  <SheetTitle>导航</SheetTitle>
                </SheetHeader>
                <nav className="mt-6 flex flex-col gap-1">
                  {publicLinks.map((l) => (
                    <Link
                      key={l.href}
                      href={l.href}
                      className="rounded-lg px-3 py-3 text-base font-medium hover:bg-muted"
                    >
                      {l.label}
                    </Link>
                  ))}
                  {user ? (
                    <button
                      type="button"
                      className="rounded-lg px-3 py-3 text-left text-base font-medium hover:bg-muted"
                      onClick={() => {
                        clearToken();
                        setUser(null);
                        window.location.href = "/login";
                      }}
                    >
                      退出登录
                    </button>
                  ) : (
                    <Link
                      href="/login"
                      className="rounded-lg px-3 py-3 text-base font-medium hover:bg-muted"
                    >
                      登录
                    </Link>
                  )}
                </nav>
              </SheetContent>
            </Sheet>
          </div>
        </div>
      </header>

      <div className="flex-1">{children}</div>

      <footer className="border-t border-border/60 py-8 text-center text-xs text-muted-foreground">
        <p>Sales-Pilot · 企业销售赋能</p>
      </footer>
    </div>
  );
}
