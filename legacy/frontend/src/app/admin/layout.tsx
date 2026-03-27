"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { AdminShell } from "@/components/admin/admin-shell";
import { fetchMe, type User } from "@/lib/api";
import { getToken } from "@/lib/auth";

export default function AdminLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();
  const [user, setUser] = useState<User | null>(null);
  const [ready, setReady] = useState(false);

  useEffect(() => {
    if (!getToken()) {
      router.replace("/login");
      return;
    }
    fetchMe().then((u) => {
      if (!u) router.replace("/login");
      else if (u.role !== "admin") router.replace("/");
      else {
        setUser(u);
        setReady(true);
      }
    });
  }, [router]);

  if (!ready || !user) {
    return (
      <div className="flex min-h-screen items-center justify-center text-muted-foreground">
        校验权限…
      </div>
    );
  }

  return <AdminShell user={user}>{children}</AdminShell>;
}
