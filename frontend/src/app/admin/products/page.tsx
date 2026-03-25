"use client";

import { useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { Button, buttonVariants } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { apiFetch, type Product } from "@/lib/api";
import { FileEdit, Pencil, Trash2 } from "lucide-react";

function formatTime(iso: string) {
  try {
    return new Date(iso).toLocaleString("zh-CN", {
      dateStyle: "short",
      timeStyle: "short",
    });
  } catch {
    return iso;
  }
}

export default function AdminProductsPage() {
  const router = useRouter();
  const [list, setList] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [msg, setMsg] = useState("");

  const refresh = useCallback(async () => {
    setLoading(true);
    setMsg("");
    const r = await apiFetch("/api/v1/admin/products");
    if (r.ok) setList(await r.json());
    else setList([]);
    setLoading(false);
  }, []);

  useEffect(() => {
    refresh();
  }, [refresh]);

  async function createDraft() {
    const res = await apiFetch("/api/v1/admin/products", {
      method: "POST",
      body: JSON.stringify({
        name: "未命名方案",
        category: "",
        description: "",
        highlights: [],
        target_personas: {},
        trigger_events: "",
        discovery_questions: [],
        competitor_analysis: "",
        roi_metrics: "",
        is_draft: true,
        vendor_market: "all",
      }),
    });
    if (res.ok) {
      const p = await res.json();
      router.push(`/admin/products/${p.id}`);
    }
  }

  async function toggleDraft(p: Product) {
    setMsg("");
    const next = !p.is_draft;
    const res = await apiFetch(`/api/v1/admin/products/${p.id}/draft`, {
      method: "PATCH",
      body: JSON.stringify({ is_draft: next }),
    });
    if (res.ok) await refresh();
    else {
      const e = await res.json().catch(() => ({}));
      setMsg(e.error || "切换草稿状态失败");
    }
  }

  async function delProduct(id: number, name: string) {
    if (!confirm(`确定删除产品「${name}」及其话术与案例？`)) return;
    setMsg("");
    const res = await apiFetch(`/api/v1/admin/products/${id}`, {
      method: "DELETE",
    });
    if (res.ok) await refresh();
    else setMsg("删除失败");
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-4">
        <h1 className="text-2xl font-bold text-primary">产品维护</h1>
        <Button onClick={createDraft}>新建产品</Button>
      </div>

      {msg ? (
        <p className="rounded-md border border-destructive/30 bg-destructive/5 px-3 py-2 text-sm text-destructive">
          {msg}
        </p>
      ) : null}

      <div className="overflow-x-auto rounded-lg border border-border bg-card">
        <table className="w-full min-w-[960px] border-collapse text-left text-sm">
          <thead>
            <tr className="border-b bg-muted/50">
              <th className="px-3 py-3 font-semibold">ID</th>
              <th className="px-3 py-3 font-semibold">名称</th>
              <th className="px-3 py-3 font-semibold">状态</th>
              <th className="px-3 py-3 font-semibold">分类</th>
              <th className="px-3 py-3 font-semibold">厂商</th>
              <th className="px-3 py-3 font-semibold">市场</th>
              <th className="px-3 py-3 font-semibold">方案 ID</th>
              <th className="px-3 py-3 font-semibold">更新时间</th>
              <th className="px-3 py-3 font-semibold text-right">操作</th>
            </tr>
          </thead>
          <tbody>
            {loading ? (
              <tr>
                <td colSpan={9} className="px-3 py-10 text-center text-muted-foreground">
                  加载中…
                </td>
              </tr>
            ) : list.length === 0 ? (
              <tr>
                <td colSpan={9} className="px-3 py-10 text-center text-muted-foreground">
                  暂无产品，点击「新建产品」开始。
                </td>
              </tr>
            ) : (
              list.map((p) => (
                <tr
                  key={p.id}
                  className="border-b border-border/80 transition-colors hover:bg-muted/30"
                >
                  <td className="px-3 py-2.5 tabular-nums text-muted-foreground">
                    {p.id}
                  </td>
                  <td className="max-w-[200px] px-3 py-2.5 font-medium">
                    <span className="line-clamp-2">{p.name}</span>
                  </td>
                  <td className="whitespace-nowrap px-3 py-2.5">
                    {p.is_draft ? (
                      <span className="rounded-full bg-amber-100 px-2 py-0.5 text-xs font-medium text-amber-900 dark:bg-amber-950/50 dark:text-amber-200">
                        草稿
                      </span>
                    ) : (
                      <span className="rounded-full bg-emerald-100 px-2 py-0.5 text-xs font-medium text-emerald-900 dark:bg-emerald-950/50 dark:text-emerald-200">
                        已发布
                      </span>
                    )}
                  </td>
                  <td className="max-w-[120px] px-3 py-2.5 text-muted-foreground">
                    <span className="line-clamp-2">
                      {p.category || "—"}
                    </span>
                  </td>
                  <td className="max-w-[100px] px-3 py-2.5 text-muted-foreground">
                    <span className="line-clamp-1">
                      {p.manufacturer_name?.trim() || "—"}
                    </span>
                  </td>
                  <td className="whitespace-nowrap px-3 py-2.5 text-muted-foreground">
                    {p.vendor_market === "foreign"
                      ? "国外"
                      : p.vendor_market === "domestic"
                        ? "国内"
                        : "不区分"}
                  </td>
                  <td className="px-3 py-2.5 tabular-nums text-muted-foreground">
                    {p.solution_category_id ?? "—"}
                  </td>
                  <td className="whitespace-nowrap px-3 py-2.5 text-muted-foreground">
                    {formatTime(p.updated_at)}
                  </td>
                  <td className="px-3 py-2.5 text-right">
                    <div className="flex flex-wrap justify-end gap-1">
                      <Link
                        href={`/admin/products/${p.id}`}
                        className={cn(
                          buttonVariants({ variant: "default", size: "sm" }),
                          "inline-flex gap-1",
                        )}
                      >
                        <Pencil className="h-3.5 w-3.5" />
                        编辑
                      </Link>
                      <Button
                        type="button"
                        size="sm"
                        variant="secondary"
                        className="gap-1"
                        onClick={() => toggleDraft(p)}
                      >
                        <FileEdit className="h-3.5 w-3.5" />
                        {p.is_draft ? "发布" : "转草稿"}
                      </Button>
                      <Button
                        type="button"
                        size="sm"
                        variant="destructive"
                        className="gap-1"
                        onClick={() => delProduct(p.id, p.name)}
                      >
                        <Trash2 className="h-3.5 w-3.5" />
                        删除
                      </Button>
                    </div>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
