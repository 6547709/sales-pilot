# UI Optimizations Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Extend login timeout to 14 days, remove mobile hamburger menu in favor of inline links, and implement smart sticky navbar that hides on scroll down.

**Architecture:** 
- Backend: Update JWT TTL constant in Golang.
- Frontend: Modify React `public-shell.tsx` component to adjust responsive classes (removing `<Sheet>`, tweaking flex layouts) and use a custom `useEffect` scroll event listener to toggle a `-translate-y-full` class on the `<header>`.

**Tech Stack:** React, Next.js, Tailwind CSS, Golang

---

### Task 1: Update Backend Token TTL

**Files:**
- Modify: `backend/internal/auth/jwt.go:18-20`

**Step 1: Write minimal implementation**

```go
// TokenTTL 访问令牌有效期
const TokenTTL = 14 * 24 * time.Hour
```

**Step 2: Commit**

```bash
git add backend/internal/auth/jwt.go
git commit -m "chore(auth): extend token TTL to 14 days"
```

---

### Task 2: Implement Smart Sticky Navbar & Uncollapse Mobile Nav

**Files:**
- Modify: `frontend/src/components/marketing/public-shell.tsx`

**Step 1: Add Scroll Hook and State**

Modify `public-shell.tsx` to include scroll direction tracking state:
```tsx
  const [scrollDir, setScrollDir] = useState("up");
  const [lastScrollY, setLastScrollY] = useState(0);

  useEffect(() => {
    const handleScroll = () => {
      const currentScrollY = window.scrollY;
      if (currentScrollY > lastScrollY && currentScrollY > 50) {
        setScrollDir("down");
      } else if (currentScrollY < lastScrollY) {
        setScrollDir("up");
      }
      setLastScrollY(currentScrollY);
    };
    window.addEventListener("scroll", handleScroll);
    return () => window.removeEventListener("scroll", handleScroll);
  }, [lastScrollY]);
```

**Step 2: Update Header Classes for Smart Sticky**

Update the `<header>` element classes:
```tsx
      <header
        className={cn(
          "sticky top-0 z-50 border-b border-border/80 bg-background/90 backdrop-blur-md transition-transform duration-300",
          scrollDir === "down" ? "-translate-y-full" : "translate-y-0"
        )}
      >
```

**Step 3: Modify Navigation Layout (Remove Sheet, Show Links)**

In the `public-shell.tsx`, change:
- Remove `<Sheet>` component and its imports (`useSheet`, `SheetTrigger`, `SheetContent`, `SheetTitle`, etc. and `Menu` icon).
- Remove `hidden md:flex` from the `<nav>` so it shows on all screens.
- Keep the `user` logic (username and logout button/login link) inline.
- Specifically, the nav container should flex row, wrap, or just fit compactly. Let's combine them into a single flex container.

```tsx
          <div className="flex items-center gap-4">
            <nav className="flex items-center gap-1">
              {publicLinks.map((l) => (
                <Link
                  key={l.href}
                  href={l.href}
                  className={cn(
                    "rounded-lg px-2 py-1 text-sm font-medium transition-colors sm:px-3 sm:py-2",
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
```

**Step 4: Commit**

```bash
git add frontend/src/components/marketing/public-shell.tsx
git commit -m "feat(ui): implement smart sticky navbar and inline mobile nav"
```
