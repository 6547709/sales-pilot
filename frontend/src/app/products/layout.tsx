"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { AppFrame } from "@/components/app-frame";
import { fetchMe, type User } from "@/lib/api";
import { getToken } from "@/lib/auth";

export default function ProductsLayout({
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
      else {
        setUser(u);
        setReady(true);
      }
    });
  }, [router]);

  if (!ready) {
    return (
      <p className="p-8 text-center text-muted-foreground">加载中…</p>
    );
  }

  return <AppFrame user={user}>{children}</AppFrame>;
}
