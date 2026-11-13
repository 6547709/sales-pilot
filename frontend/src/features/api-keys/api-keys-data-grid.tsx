"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { useTable, type ColumnDef } from "@tanstack/react-table";
import { KeyRound, Trash2 } from "lucide-react";
import {
  DataGrid,
  DataGridContainer,
  dataGridFeatures,
  type DataGridFeatures,
} from "@/components/reui/data-grid/data-grid";
import { DataGridTable } from "@/components/reui/data-grid/data-grid-table";
import { DataGridPagination } from "@/components/reui/data-grid/data-grid-pagination";
import { Alert, AlertDescription, AlertTitle } from "@/components/reui/alert";
import { Badge } from "@/components/reui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { ConfirmDialog } from "@/components/patterns/confirm-dialog";
import { ErrorState } from "@/components/patterns/error-state";
import { FormField } from "@/components/patterns/form-field";
import { apiKeysApi, type ApiKeyRow } from "@/lib/api";

function formatTime(iso?: string | null) {
  if (!iso) return "—";
  try {
    return new Date(iso).toLocaleString("zh-CN", { dateStyle: "short", timeStyle: "short" });
  } catch {
    return iso;
  }
}

function statusOf(k: ApiKeyRow) {
  if (!k.isActive) return { label: "已作废", variant: "outline" as const };
  if (k.expiresAt && new Date(k.expiresAt).getTime() < Date.now())
    return { label: "已过期", variant: "warning-light" as const };
  return { label: "有效", variant: "success-light" as const };
}

export function ApiKeysDataGrid() {
  const [list, setList] = useState<ApiKeyRow[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [name, setName] = useState("");
  const [expiresAt, setExpiresAt] = useState("");
  const [onceKey, setOnceKey] = useState<string | null>(null);
  const [revokeTarget, setRevokeTarget] = useState<ApiKeyRow | null>(null);
  const [busy, setBusy] = useState(false);

  const load = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      setList(await apiKeysApi.list());
    } catch (e) {
      setError(e instanceof Error ? e.message : "加载失败");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  async function createKey() {
    if (!name.trim()) {
      setError("请填写名称");
      return;
    }
    setBusy(true);
    setError(null);
    setOnceKey(null);
    try {
      const iso = expiresAt ? new Date(expiresAt).toISOString() : undefined;
      const created = await apiKeysApi.create(name.trim(), iso);
      setOnceKey(created.key);
      setName("");
      setExpiresAt("");
      await load();
    } catch (e) {
      setError(e instanceof Error ? e.message : "创建失败");
    } finally {
      setBusy(false);
    }
  }

  const columns = useMemo<ColumnDef<DataGridFeatures, ApiKeyRow>[]>(
    () => [
      {
        accessorKey: "name",
        header: "名称",
        cell: ({ row }) => <span className="font-medium">{row.original.name}</span>,
      },
      {
        accessorKey: "prefix",
        header: "前缀",
        cell: ({ row }) => (
          <span className="font-mono text-xs text-muted-foreground">{row.original.prefix}</span>
        ),
      },
      {
        id: "status",
        header: "状态",
        cell: ({ row }) => {
          const s = statusOf(row.original);
          return <Badge variant={s.variant}>{s.label}</Badge>;
        },
      },
      {
        id: "expires",
        header: "过期时间",
        cell: ({ row }) => (
          <span className="whitespace-nowrap text-muted-foreground">
            {formatTime(row.original.expiresAt)}
          </span>
        ),
      },
      {
        id: "lastUsed",
        header: "最后使用",
        cell: ({ row }) => (
          <span className="whitespace-nowrap text-muted-foreground">
            {formatTime(row.original.lastUsedAt)}
          </span>
        ),
      },
      {
        id: "actions",
        header: () => <span className="text-right">操作</span>,
        cell: ({ row }) => (
          <div className="flex justify-end">
            <Button
              size="sm"
              variant="destructive"
              className="gap-1"
              disabled={!row.original.isActive}
              onClick={() => setRevokeTarget(row.original)}
            >
              <Trash2 className="h-3.5 w-3.5" />
              作废
            </Button>
          </div>
        ),
      },
    ],
    [],
  );

  const table = useTable({
    features: dataGridFeatures,
    data: list,
    columns,
    getRowId: (row) => String(row.id),
    initialState: { pagination: { pageIndex: 0, pageSize: 10 } },
  });

  return (
    <div className="flex flex-col gap-4">
      <div className="rounded-lg border border-border bg-card p-3">
        <div className="flex flex-wrap items-end gap-2">
          <FormField label="名称" className="min-w-56 flex-1">
            <Input
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="例如：mcp-server"
            />
          </FormField>
          <FormField label="过期时间（可选）" className="min-w-56">
            <Input
              type="datetime-local"
              value={expiresAt}
              onChange={(e) => setExpiresAt(e.target.value)}
            />
          </FormField>
          <Button disabled={busy} onClick={createKey} className="gap-1">
            <KeyRound className="h-4 w-4" />
            创建密钥
          </Button>
        </div>
      </div>

      {onceKey ? (
        <Alert variant="warning">
          <AlertTitle>完整密钥仅显示一次</AlertTitle>
          <AlertDescription>
            <p className="break-all font-mono text-xs">{onceKey}</p>
            <p className="mt-1">请立即保存；丢失后只能作废并重新创建。</p>
          </AlertDescription>
        </Alert>
      ) : null}

      {error ? <ErrorState message={error} onRetry={load} /> : null}

      <DataGrid
        table={table}
        recordCount={list.length}
        isLoading={loading}
        emptyMessage="暂无 API 密钥"
        tableLayout={{ dense: true, headerBackground: true, stripped: true }}
      >
        <DataGridContainer>
          <DataGridTable />
        </DataGridContainer>
        <DataGridPagination sizes={[10, 20, 50]} rowsPerPageLabel="每页" />
      </DataGrid>

      <ConfirmDialog
        open={!!revokeTarget}
        onOpenChange={(open) => {
          if (!open) setRevokeTarget(null);
        }}
        title="作废 API 密钥"
        description={revokeTarget ? `作废「${revokeTarget.name}」后，使用该密钥的自动化与 MCP 将立即失效。` : undefined}
        confirmLabel="作废"
        destructive
        onConfirm={async () => {
          if (!revokeTarget) return;
          await apiKeysApi.revoke(revokeTarget.id);
          setRevokeTarget(null);
          await load();
        }}
      />
    </div>
  );
}
