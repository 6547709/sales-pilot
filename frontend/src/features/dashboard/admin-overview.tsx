"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import {
  DatabaseBackup,
  KeyRound,
  Network,
  Package,
  Settings,
  Users,
} from "lucide-react";
import { PageHeader } from "@/components/patterns/page-header";
import { StatCard } from "@/components/patterns/stat-card";
import { Frame, FramePanel } from "@/components/reui/frame";
import { adminProductsApi, apiKeysApi, usersApi } from "@/lib/api";

type Load<T> =
  | { status: "loading" }
  | { status: "ok"; value: T }
  | { status: "error" };

const shortcuts = [
  { href: "/admin/products", label: "产品维护", description: "编辑方案内容与发布状态", icon: Package },
  { href: "/admin/topology", label: "拓扑维护", description: "层、分类与厂商", icon: Network },
  { href: "/admin/users", label: "本地用户", description: "账号、角色与启停", icon: Users },
  { href: "/admin/api-keys", label: "API 密钥", description: "自动化与 MCP 接入", icon: KeyRound },
  { href: "/admin/backup", label: "备份恢复", description: "导出与全量恢复", icon: DatabaseBackup },
  { href: "/admin/settings", label: "认证配置", description: "LDAP / OIDC / 页脚", icon: Settings },
];

function valueOf(load: Load<number>) {
  if (load.status === "loading") return "…";
  if (load.status === "error") return "—";
  return load.value;
}

export function AdminOverview() {
  const [products, setProducts] = useState<Load<number>>({ status: "loading" });
  const [users, setUsers] = useState<Load<number>>({ status: "loading" });
  const [keys, setKeys] = useState<Load<number>>({ status: "loading" });

  useEffect(() => {
    let alive = true;
    Promise.allSettled([
      adminProductsApi.list({ page: 1, pageSize: 1 }),
      usersApi.list({ page: 1, pageSize: 1 }),
      apiKeysApi.list(),
    ]).then(([p, u, k]) => {
      if (!alive) return;
      setProducts(
        p.status === "fulfilled"
          ? { status: "ok", value: p.value.total }
          : { status: "error" },
      );
      setUsers(
        u.status === "fulfilled"
          ? { status: "ok", value: u.value.total }
          : { status: "error" },
      );
      if (k.status === "fulfilled") {
        const now = Date.now();
        const active = k.value.filter(
          (x) => x.isActive && (!x.expiresAt || new Date(x.expiresAt).getTime() > now),
        ).length;
        setKeys({ status: "ok", value: active });
      } else {
        setKeys({ status: "error" });
      }
    });
    return () => {
      alive = false;
    };
  }, []);

  return (
    <div className="flex flex-col gap-6">
      <PageHeader title="控制台" description="内容、账号与接入状态的快速概览" />

      <div className="grid gap-4 sm:grid-cols-3">
        <StatCard label="产品总数" value={valueOf(products)} hint="含草稿" icon={Package} />
        <StatCard label="本地用户" value={valueOf(users)} hint="本地账号" icon={Users} />
        <StatCard label="有效密钥" value={valueOf(keys)} hint="未过期且启用" icon={KeyRound} />
      </div>

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {shortcuts.map((item) => {
          const Icon = item.icon;
          return (
            <Link key={item.href} href={item.href} className="block">
              <Frame className="h-full transition-colors hover:border-primary/40">
                <FramePanel>
                  <div className="flex items-start gap-3">
                    <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary">
                      <Icon className="h-4 w-4" />
                    </span>
                    <div className="min-w-0">
                      <p className="font-medium">{item.label}</p>
                      <p className="mt-0.5 text-xs text-muted-foreground">
                        {item.description}
                      </p>
                    </div>
                  </div>
                </FramePanel>
              </Frame>
            </Link>
          );
        })}
      </div>
    </div>
  );
}
