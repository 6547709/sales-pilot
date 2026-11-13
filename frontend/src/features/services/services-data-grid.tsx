"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { useTable, type ColumnDef } from "@tanstack/react-table";
import { Pencil, Plus, Trash2, Wrench } from "lucide-react";
import {
  DataGrid,
  DataGridContainer,
  dataGridFeatures,
  type DataGridFeatures,
} from "@/components/reui/data-grid/data-grid";
import { DataGridTable } from "@/components/reui/data-grid/data-grid-table";
import { DataGridPagination } from "@/components/reui/data-grid/data-grid-pagination";
import { Badge } from "@/components/reui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { ConfirmDialog } from "@/components/patterns/confirm-dialog";
import { ErrorState } from "@/components/patterns/error-state";
import { FormField } from "@/components/patterns/form-field";
import {
  adminServicesApi,
  fetchSolutionCategories,
  type Service,
  type ServiceInput,
} from "@/lib/api";

const KINDS: { value: string; label: string }[] = [
  { value: "v2v", label: "V2V 迁移" },
  { value: "backup", label: "备份" },
  { value: "dr", label: "容灾编排" },
  { value: "migration", label: "迁移" },
  { value: "upgrade", label: "升级" },
  { value: "ops", label: "运维" },
];

function kindLabel(kind: string) {
  return KINDS.find((k) => k.value === kind)?.label ?? kind;
}

type FormState = {
  id?: number;
  slug: string;
  name: string;
  kind: string;
  summary: string;
  highlights: string;
  deliverables: string;
  valueProp: string;
  features: string;
  iconKey: string;
  sortOrder: number;
  isActive: boolean;
  categoryIds: number[];
};

const emptyForm: FormState = {
  slug: "",
  name: "",
  kind: "backup",
  summary: "",
  highlights: "",
  deliverables: "",
  valueProp: "",
  features: "",
  iconKey: "",
  sortOrder: 0,
  isActive: true,
  categoryIds: [],
};

function linesToArray(v: string) {
  return v
    .split("\n")
    .map((s) => s.trim())
    .filter(Boolean);
}

export function ServicesDataGrid() {
  const [list, setList] = useState<Service[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [form, setForm] = useState<FormState | null>(null);
  const [deleteTarget, setDeleteTarget] = useState<Service | null>(null);
  const [categories, setCategories] = useState<{ id: number; label: string }[]>([]);
  const [busy, setBusy] = useState(false);

  const load = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      setList(await adminServicesApi.list());
    } catch (e) {
      setError(e instanceof Error ? e.message : "加载失败");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    load();
    fetchSolutionCategories(false).then((cs) =>
      setCategories(cs.map((c) => ({ id: c.id, label: c.label }))),
    );
  }, [load]);

  async function save() {
    if (!form) return;
    setBusy(true);
    setError(null);
    const payload: ServiceInput = {
      slug: form.slug.trim(),
      name: form.name.trim(),
      kind: form.kind,
      summary: form.summary,
      highlights: linesToArray(form.highlights),
      deliverables: linesToArray(form.deliverables),
      value_prop: form.valueProp,
      features: linesToArray(form.features),
      icon_key: form.iconKey,
      sort_order: form.sortOrder,
      is_active: form.isActive,
      category_ids: form.categoryIds,
    };
    try {
      if (form.id) await adminServicesApi.update(form.id, payload);
      else await adminServicesApi.create(payload);
      setForm(null);
      await load();
    } catch (e) {
      setError(e instanceof Error ? e.message : "保存失败");
    } finally {
      setBusy(false);
    }
  }

  const columns = useMemo<ColumnDef<DataGridFeatures, Service>[]>(
    () => [
      {
        accessorKey: "name",
        header: "名称",
        cell: ({ row }) => (
          <div className="flex flex-col">
            <span className="font-medium">{row.original.name}</span>
            <span className="text-xs text-muted-foreground">{row.original.slug}</span>
          </div>
        ),
      },
      {
        id: "kind",
        header: "类型",
        cell: ({ row }) => <Badge variant="secondary">{kindLabel(row.original.kind)}</Badge>,
      },
      {
        id: "summary",
        header: "摘要",
        cell: ({ row }) => (
          <span className="line-clamp-2 max-w-[320px] text-muted-foreground">
            {row.original.summary}
          </span>
        ),
      },
      { id: "sort", header: "排序", cell: ({ row }) => row.original.sortOrder },
      {
        id: "visible",
        header: "可见性",
        cell: () => <Badge variant="info-light">客户可见</Badge>,
      },
      {
        id: "status",
        header: "状态",
        cell: ({ row }) => (
          <Badge variant={row.original.isActive ? "success-light" : "outline"}>
            {row.original.isActive ? "启用" : "停用"}
          </Badge>
        ),
      },
      {
        id: "actions",
        header: () => <span className="text-right">操作</span>,
        cell: ({ row }) => (
          <div className="flex justify-end gap-1">
            <Button
              size="sm"
              variant="secondary"
              className="gap-1"
              onClick={() =>
                setForm({
                  id: row.original.id,
                  slug: row.original.slug,
                  name: row.original.name,
                  kind: row.original.kind,
                  summary: row.original.summary,
                  highlights: row.original.highlights.join("\n"),
                  deliverables: row.original.deliverables.join("\n"),
                  valueProp: row.original.valueProp ?? "",
                  features: row.original.features.join("\n"),
                  iconKey: row.original.iconKey ?? "",
                  sortOrder: row.original.sortOrder,
                  isActive: row.original.isActive,
                  categoryIds: [],
                })
              }
            >
              <Pencil className="h-3.5 w-3.5" />
              编辑
            </Button>
            <Button
              size="sm"
              variant="destructive"
              className="gap-1"
              onClick={() => setDeleteTarget(row.original)}
            >
              <Trash2 className="h-3.5 w-3.5" />
              删除
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
    initialState: { pagination: { pageIndex: 0, pageSize: 20 } },
  });

  return (
    <div className="flex flex-col gap-4">
      <div className="flex items-center gap-2">
        <Button variant="outline" onClick={() => load()}>
          刷新
        </Button>
        <Button className="ml-auto gap-1" onClick={() => setForm({ ...emptyForm })}>
          <Plus className="h-4 w-4" />
          新建服务
        </Button>
      </div>

      {error ? <ErrorState message={error} onRetry={load} /> : null}

      <DataGrid
        table={table}
        recordCount={list.length}
        isLoading={loading}
        emptyMessage="暂无服务"
        tableLayout={{ dense: true, headerBackground: true, stripped: true }}
      >
        <DataGridContainer>
          <DataGridTable />
        </DataGridContainer>
        <DataGridPagination sizes={[10, 20, 50]} rowsPerPageLabel="每页" />
      </DataGrid>

      <p className="flex items-center gap-1.5 text-xs text-muted-foreground">
        <Wrench className="h-3.5 w-3.5" />
        共 {list.length} 项服务，全部对客户可见
      </p>

      <Dialog open={!!form} onOpenChange={(open) => !open && setForm(null)}>
        <DialogContent className="max-h-[85vh] overflow-y-auto sm:max-w-2xl">
          <DialogHeader>
            <DialogTitle>{form?.id ? "编辑服务" : "新建服务"}</DialogTitle>
          </DialogHeader>
          {form ? (
            <div className="grid gap-3 sm:grid-cols-2">
              <FormField label="名称">
                <Input value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} />
              </FormField>
              <FormField label="slug">
                <Input value={form.slug} onChange={(e) => setForm({ ...form, slug: e.target.value })} />
              </FormField>
              <FormField label="类型">
                <Select value={form.kind} onValueChange={(v) => setForm({ ...form, kind: String(v) })}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {KINDS.map((k) => (
                      <SelectItem key={k.value} value={k.value}>
                        {k.label}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </FormField>
              <FormField label="排序">
                <Input
                  type="number"
                  value={form.sortOrder}
                  onChange={(e) => setForm({ ...form, sortOrder: Number(e.target.value) || 0 })}
                />
              </FormField>
              <FormField label="摘要" className="sm:col-span-2">
                <Textarea
                  rows={2}
                  value={form.summary}
                  onChange={(e) => setForm({ ...form, summary: e.target.value })}
                />
              </FormField>
              <FormField label="亮点（每行一条）">
                <Textarea
                  rows={3}
                  value={form.highlights}
                  onChange={(e) => setForm({ ...form, highlights: e.target.value })}
                />
              </FormField>
              <FormField label="交付物（每行一条）">
                <Textarea
                  rows={3}
                  value={form.deliverables}
                  onChange={(e) => setForm({ ...form, deliverables: e.target.value })}
                />
              </FormField>
              <FormField label="客户价值" className="sm:col-span-2">
                <Textarea
                  rows={2}
                  value={form.valueProp}
                  onChange={(e) => setForm({ ...form, valueProp: e.target.value })}
                />
              </FormField>
              <FormField label="特点（每行一条）" className="sm:col-span-2">
                <Textarea
                  rows={2}
                  value={form.features}
                  onChange={(e) => setForm({ ...form, features: e.target.value })}
                />
              </FormField>
              <FormField label="图标 key">
                <Input value={form.iconKey} onChange={(e) => setForm({ ...form, iconKey: e.target.value })} />
              </FormField>
              <FormField label="状态">
                <Select
                  value={form.isActive ? "on" : "off"}
                  onValueChange={(v) => setForm({ ...form, isActive: String(v) === "on" })}
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="on">启用</SelectItem>
                    <SelectItem value="off">停用</SelectItem>
                  </SelectContent>
                </Select>
              </FormField>
              <FormField label="所属解决方案（架构方框）" className="sm:col-span-2">
                <div className="flex flex-wrap gap-1.5">
                  {categories.map((c) => {
                    const checked = form.categoryIds.includes(c.id);
                    return (
                      <button
                        key={c.id}
                        type="button"
                        onClick={() =>
                          setForm({
                            ...form,
                            categoryIds: checked
                              ? form.categoryIds.filter((id) => id !== c.id)
                              : [...form.categoryIds, c.id],
                          })
                        }
                        className={
                          checked
                            ? "rounded-full bg-primary/10 px-2.5 py-1 text-xs font-medium text-primary"
                            : "rounded-full border border-border px-2.5 py-1 text-xs text-muted-foreground"
                        }
                      >
                        {c.label}
                      </button>
                    );
                  })}
                </div>
              </FormField>
            </div>
          ) : null}
          <DialogFooter>
            <Button variant="outline" onClick={() => setForm(null)}>
              取消
            </Button>
            <Button disabled={busy} onClick={save}>
              保存
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <ConfirmDialog
        open={!!deleteTarget}
        onOpenChange={(open) => {
          if (!open) setDeleteTarget(null);
        }}
        title="删除服务"
        description={deleteTarget ? `确定删除「${deleteTarget.name}」？相关场景关联会一并移除。` : undefined}
        confirmLabel="删除"
        destructive
        onConfirm={async () => {
          if (!deleteTarget) return;
          await adminServicesApi.remove(deleteTarget.id);
          setDeleteTarget(null);
          await load();
        }}
      />
    </div>
  );
}
