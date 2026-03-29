"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useEffect, useState } from "react";
import { Button, buttonVariants } from "@/components/ui/button";
import { fetchMe, type User } from "@/lib/api";
import { clearToken, getToken } from "@/lib/auth";
import { cn } from "@/lib/utils";

const publicLinks = [
  { href: "/", label: "首页" },
  { href: "/products", label: "方案库" },
];

export function PublicShell({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const [user, setUser] = useState<User | null | undefined>(undefined);
  const [scrollDir, setScrollDir] = useState("up");

  useEffect(() => {
    let lastScroll = window.scrollY;
    const handleScroll = () => {
      const currentScrollY = window.scrollY;
      if (currentScrollY > lastScroll && currentScrollY > 50) {
        setScrollDir("down");
      } else if (currentScrollY < lastScroll) {
        setScrollDir("up");
      }
      lastScroll = currentScrollY;
    };
    window.addEventListener("scroll", handleScroll, { passive: true });
    return () => window.removeEventListener("scroll", handleScroll);
  }, []);

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
      <header
        className={cn(
          "sticky top-0 z-50 border-b border-border/80 bg-background/90 backdrop-blur-md transition-transform duration-300",
          scrollDir === "down" ? "-translate-y-full" : "translate-y-0"
        )}
      >
        <div className="mx-auto flex h-14 max-w-[1400px] items-center justify-between gap-4 px-4">
          <Link href="/" className="flex items-center gap-2 font-semibold text-primary">
            <span className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary text-sm text-primary-foreground">
              SP
            </span>
            <span className="hidden sm:inline">Sales-Pilot</span>
          </Link>

          <div className="flex items-center gap-1 sm:gap-4">
            <nav className="flex items-center gap-1">
              {publicLinks.map((l) => (
                <Link
                  key={l.href}
                  href={l.href}
                  className={cn(
                    "rounded-lg px-2 py-1 text-xs font-medium transition-colors sm:px-3 sm:py-2 sm:text-sm",
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
                    className="h-8 px-2 text-xs sm:h-9 sm:px-3 sm:text-sm"
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
                  className={cn(buttonVariants({ size: "sm" }), "h-8 px-2 text-xs sm:h-9 sm:px-3 sm:text-sm")}
                >
                  登录
                </Link>
              )}
            </div>
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
