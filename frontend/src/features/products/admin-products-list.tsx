"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import {
  useTable,
  type ColumnDef,
  type PaginationState,
  type SortingState,
} from "@tanstack/react-table";
import { Package, Pencil, FileEdit, Trash2 } from "lucide-react";
import {
  DataGrid,
  DataGridContainer,
  dataGridFeatures,
  type DataGridFeatures,
} from "@/components/reui/data-grid/data-grid";
import { DataGridTable } from "@/components/reui/data-grid/data-grid-table";
import { DataGridPagination } from "@/components/reui/data-grid/data-grid-pagination";
import { DataGridColumnHeader } from "@/components/reui/data-grid/data-grid-column-header";
import { Badge } from "@/components/reui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { ConfirmDialog } from "@/components/patterns/confirm-dialog";
import { ErrorState } from "@/components/patterns/error-state";
import { MarketBadge } from "@/components/domain/market-badge";
import { PublishStatusBadge } from "@/components/domain/publish-status-badge";
import { adminProductsApi, type Product } from "@/lib/api";

const SORT_MAP: Record<string, string> = {
  id: "id",
  name: "name",
  updatedAt: "updated_at",
};

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

export function AdminProductsList() {
  const router = useRouter();
  const [data, setData] = useState<Product[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [pagination, setPagination] = useState<PaginationState>({
    pageIndex: 0,
    pageSize: 20,
  });
  const [sorting, setSorting] = useState<SortingState>([
    { id: "updatedAt", desc: true },
  ]);
  const [q, setQ] = useState("");
  const [status, setStatus] = useState("all");
  const [market, setMarket] = useState("all");
  const [pendingDelete, setPendingDelete] = useState<Product | null>(null);

  const load = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const sort = sorting[0]?.id ? SORT_MAP[sorting[0].id] : undefined;
      const r = await adminProductsApi.list({
        page: pagination.pageIndex + 1,
        pageSize: pagination.pageSize,
        sort,
        order: sorting[0]?.desc ? "desc" : "asc",
        q: q.trim() || undefined,
        status: status === "all" ? undefined : (status as "draft" | "published"),
        market: market === "all" ? undefined : market,
      });
      setData(r.items);
      setTotal(r.total);
    } catch (e) {
      setError(e instanceof Error ? e.message : "加载失败");
    } finally {
      setLoading(false);
    }
  }, [pagination, sorting, q, status, market]);

  useEffect(() => {
    load();
  }, [load]);

  function resetPage() {
    setPagination((p) => ({ ...p, pageIndex: 0 }));
  }

  async function toggleDraft(p: Product) {
    await adminProductsApi.setDraft(p.id, !p.isDraft);
    await load();
  }

  const columns = useMemo<ColumnDef<DataGridFeatures, Product>[]>(
    () => [
      {
        accessorKey: "id",
        header: ({ column }) => <DataGridColumnHeader column={column} title="ID" />,
        cell: ({ row }) => (
          <span className="tabular-nums text-muted-foreground">{row.original.id}</span>
        ),
        size: 64,
      },
      {
        accessorKey: "name",
        header: ({ column }) => <DataGridColumnHeader column={column} title="名称" />,
        cell: ({ row }) => (
          <span className="line-clamp-2 max-w-[260px] font-medium">{row.original.name}</span>
        ),
      },
      {
        id: "status",
        header: "状态",
        cell: ({ row }) => <PublishStatusBadge isDraft={!!row.original.isDraft} />,
      },
      {
        id: "category",
        header: "分类",
        cell: ({ row }) => (
          <span className="line-clamp-1 text-muted-foreground">
            {row.original.category || "—"}
          </span>
        ),
      },
      {
        id: "manufacturer",
        header: "厂商",
        cell: ({ row }) => (
          <span className="line-clamp-1 text-muted-foreground">
            {row.original.manufacturerName?.trim() || "—"}
          </span>
        ),
      },
      {
        id: "market",
        header: "市场",
        cell: ({ row }) => <MarketBadge market={row.original.vendorMarket} />,
      },
      {
        accessorKey: "updatedAt",
        header: ({ column }) => <DataGridColumnHeader column={column} title="更新时间" />,
        cell: ({ row }) => (
          <span className="whitespace-nowrap text-muted-foreground">
            {formatTime(row.original.updatedAt)}
          </span>
        ),
      },
      {
        id: "actions",
        header: () => <span className="text-right">操作</span>,
        cell: ({ row }) => {
          const p = row.original;
          return (
            <div className="flex flex-wrap justify-end gap-1">
              <Link href={`/admin/products/${p.id}`}>
                <Button size="sm" className="gap-1">
                  <Pencil className="h-3.5 w-3.5" />
                  编辑
                </Button>
              </Link>
              <Button
                size="sm"
                variant="secondary"
                className="gap-1"
                onClick={() => toggleDraft(p)}
              >
                <FileEdit className="h-3.5 w-3.5" />
                {p.isDraft ? "发布" : "转草稿"}
              </Button>
              <Button
                size="sm"
                variant="destructive"
                className="gap-1"
                onClick={() => setPendingDelete(p)}
              >
                <Trash2 className="h-3.5 w-3.5" />
                删除
              </Button>
            </div>
          );
        },
      },
    ],
    [],
  );

  const table = useTable({
    features: dataGridFeatures,
    data,
    columns,
    state: { pagination, sorting },
    onPaginationChange: setPagination,
    onSortingChange: (updater) => {
      setSorting(updater);
      resetPage();
    },
    manualPagination: true,
    manualSorting: true,
    rowCount: total,
    getRowId: (row) => String(row.id),
  });

  return (
    <div className="flex flex-col gap-4">
      <div className="flex flex-wrap items-center gap-2">
        <Input
          value={q}
          onChange={(e) => {
            setQ(e.target.value);
            resetPage();
          }}
          placeholder="搜索名称或分类"
          className="w-56"
        />
        <Select
          value={status}
          onValueChange={(v) => {
            setStatus(String(v));
            resetPage();
          }}
        >
          <SelectTrigger className="w-32">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">全部状态</SelectItem>
            <SelectItem value="published">已发布</SelectItem>
            <SelectItem value="draft">草稿</SelectItem>
          </SelectContent>
        </Select>
        <Select
          value={market}
          onValueChange={(v) => {
            setMarket(String(v));
            resetPage();
          }}
        >
          <SelectTrigger className="w-32">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">全部市场</SelectItem>
            <SelectItem value="domestic">国内</SelectItem>
            <SelectItem value="foreign">国外</SelectItem>
          </SelectContent>
        </Select>
        <Button variant="outline" onClick={() => load()}>
          刷新
        </Button>
        <Button
          className="ml-auto"
          onClick={async () => {
            const created = await adminProductsApi.create({
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
            });
            router.push(`/admin/products/${created.id}`);
          }}
        >
          新建产品
        </Button>
      </div>

      {error ? <ErrorState message={error} onRetry={load} /> : null}

      <DataGrid
        table={table}
        recordCount={total}
        isLoading={loading}
        emptyMessage="暂无产品"
        tableLayout={{ dense: true, headerBackground: true }}
      >
        <DataGridContainer>
          <DataGridTable />
        </DataGridContainer>
        <DataGridPagination
          sizes={[10, 20, 50]}
          rowsPerPageLabel="每页"
          previousPageLabel="上一页"
          nextPageLabel="下一页"
        />
      </DataGrid>

      <ConfirmDialog
        open={!!pendingDelete}
        onOpenChange={(open) => {
          if (!open) setPendingDelete(null);
        }}
        title="删除产品"
        description={
          pendingDelete
            ? `确定删除「${pendingDelete.name}」及其话术与案例？该操作不可撤销。`
            : undefined
        }
        confirmLabel="删除"
        destructive
        onConfirm={async () => {
          if (!pendingDelete) return;
          await adminProductsApi.remove(pendingDelete.id);
          setPendingDelete(null);
          await load();
        }}
      />

      <p className="flex items-center gap-1.5 text-xs text-muted-foreground">
        <Package className="h-3.5 w-3.5" />
        共 {total} 个产品
      </p>
    </div>
  );
}
