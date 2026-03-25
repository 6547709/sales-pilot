"use client";

import { SiteHeader } from "@/components/site-header";
import type { User } from "@/lib/api";

export function AppFrame({
  user,
  children,
}: {
  user: User | null;
  children: React.ReactNode;
}) {
  return (
    <>
      <SiteHeader user={user} />
      <main className="mx-auto w-full max-w-6xl flex-1 px-4 py-6">
        {children}
      </main>
    </>
  );
}
