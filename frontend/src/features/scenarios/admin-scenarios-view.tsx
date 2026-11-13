"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { useTable, type ColumnDef } from "@tanstack/react-table";
import { Pencil, Plus, Trash2, Link2 } from "lucide-react";
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
  adminProductsApi,
  adminScenariosApi,
  adminServicesApi,
  type AdminScenarioGroupRow,
  type Product,
  type Service,
  type ScenarioSummary,
} from "@/lib/api";

type ScenarioForm = {
  id?: number;
  groupId: number;
  slug: string;
  name: string;
  summary: string;
  painPoint: string;
  valueProp: string;
  iconKey: string;
  sortOrder: number;
  isActive: boolean;
};

const MAX_GROUPS = 6;

export function AdminScenariosView() {
  const [groups, setGroups] = useState<AdminScenarioGroupRow[]>([]);
  const [all, setAll] = useState<{ groupId: number; items: ScenarioSummary[] }[]>([]);
  const [activeId, setActiveId] = useState<number | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [groupForm, setGroupForm] = useState<{
    id?: number;
    name: string;
    slug: string;
    description: string;
    sortOrder: number;
  } | null>(null);
  const [scenarioForm, setScenarioForm] = useState<ScenarioForm | null>(null);
  const [deleteScenario, setDeleteScenario] = useState<ScenarioSummary | null>(null);
  const [deleteGroup, setDeleteGroup] = useState<AdminScenarioGroupRow | null>(null);
  const [linksTarget, setLinksTarget] = useState<ScenarioSummary | null>(null);
  const [linkProducts, setLinkProducts] = useState<Product[]>([]);
  const [linkServices, setLinkServices] = useState<Service[]>([]);
  const [pickedProducts, setPickedProducts] = useState<number[]>([]);
  const [pickedServices, setPickedServices] = useState<number[]>([]);
  const [busy, setBusy] = useState(false);

  const load = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const [groupRows, list] = await Promise.all([
        adminScenariosApi.groups(),
        adminScenariosApi.list(),
      ]);
      setGroups(groupRows);
      setAll(
        list.map((g) => ({
          groupId: g.id,
          items: g.scenarios,
        })),
      );
      setActiveId((prev) => prev ?? groupRows[0]?.id ?? null);
    } catch (e) {
      setError(e instanceof Error ? e.message : "加载失败");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  const activeGroup = groups.find((g) => g.id === activeId) ?? null;
  const scenarios = all.find((g) => g.groupId === activeId)?.items ?? [];

  async function openLinks(item: ScenarioSummary) {
    setLinksTarget(item);
    setError(null);
    try {
      const [detail, products, services] = await Promise.all([
        adminScenariosApi.get(item.scenario.id),
        adminProductsApi.list({ page: 1, pageSize: 100 }),
        adminServicesApi.list(),
      ]);
      setPickedProducts(detail.products.map((p) => p.id));
      setPickedServices(detail.services.map((s) => s.id));
      setLinkProducts(products.items);
      setLinkServices(services);
    } catch (e) {
      setError(e instanceof Error ? e.message : "加载关联失败");
    }
  }

  async function saveLinks() {
    if (!linksTarget) return;
    setBusy(true);
    try {
      await adminScenariosApi.setLinks(linksTarget.scenario.id, pickedProducts, pickedServices);
      setLinksTarget(null);
      await load();
    } catch (e) {
      setError(e instanceof Error ? e.message : "保存关联失败");
    } finally {
      setBusy(false);
    }
  }

  async function saveGroup() {
    if (!groupForm) return;
    setBusy(true);
    try {
      if (groupForm.id) {
        await adminScenariosApi.updateGroup(groupForm.id, {
          name: groupForm.name,
          slug: groupForm.slug,
          description: groupForm.description,
          sort_order: groupForm.sortOrder,
        });
      } else {
        await adminScenariosApi.createGroup({
          name: groupForm.name,
          slug: groupForm.slug,
          description: groupForm.description,
        });
      }
      setGroupForm(null);
      await load();
    } catch (e) {
      setError(e instanceof Error ? e.message : "保存失败");
    } finally {
      setBusy(false);
    }
  }

  async function saveScenario() {
    if (!scenarioForm) return;
    setBusy(true);
    const payload = {
      group_id: scenarioForm.groupId,
      slug: scenarioForm.slug.trim(),
      name: scenarioForm.name.trim(),
      summary: scenarioForm.summary,
      pain_point: scenarioForm.painPoint,
      value_prop: scenarioForm.valueProp,
      icon_key: scenarioForm.iconKey,
      sort_order: scenarioForm.sortOrder,
      is_active: scenarioForm.isActive,
    };
    try {
      if (scenarioForm.id) await adminScenariosApi.update(scenarioForm.id, payload);
      else await adminScenariosApi.create(payload);
      setScenarioForm(null);
      await load();
    } catch (e) {
      setError(e instanceof Error ? e.message : "保存失败");
    } finally {
      setBusy(false);
    }
  }

  const columns = useMemo<ColumnDef<DataGridFeatures, ScenarioSummary>[]>(
    () => [
      {
        id: "name",
        header: "场景",
        cell: ({ row }) => (
          <div className="flex flex-col">
            <span className="font-medium">{row.original.scenario.name}</span>
            <span className="text-xs text-muted-foreground">{row.original.scenario.slug}</span>
          </div>
        ),
      },
      {
        id: "summary",
        header: "摘要",
        cell: ({ row }) => (
          <span className="line-clamp-2 max-w-[340px] text-muted-foreground">
            {row.original.scenario.summary}
          </span>
        ),
      },
      { id: "products", header: "产品", cell: ({ row }) => row.original.productCount },
      { id: "services", header: "服务", cell: ({ row }) => row.original.serviceCount },
      {
        id: "status",
        header: "状态",
        cell: ({ row }) => {
          const s = row.original;
          if (!s.scenario.isActive) return <Badge variant="outline">停用</Badge>;
          if (s.productCount + s.serviceCount === 0)
            return <Badge variant="warning-light">未关联</Badge>;
          return <Badge variant="success-light">启用</Badge>;
        },
      },
      {
        id: "actions",
        header: () => <span className="text-right">操作</span>,
        cell: ({ row }) => (
          <div className="flex flex-wrap justify-end gap-1">
            <Button
              size="sm"
              variant="secondary"
              className="gap-1"
              onClick={() =>
                setScenarioForm({
                  id: row.original.scenario.id,
                  groupId: row.original.scenario.groupId,
                  slug: row.original.scenario.slug,
                  name: row.original.scenario.name,
                  summary: row.original.scenario.summary,
                  painPoint: row.original.scenario.painPoint ?? "",
                  valueProp: row.original.scenario.valueProp ?? "",
                  iconKey: row.original.scenario.iconKey ?? "",
                  sortOrder: row.original.scenario.sortOrder,
                  isActive: row.original.scenario.isActive,
                })
              }
            >
              <Pencil className="h-3.5 w-3.5" />
              编辑
            </Button>
            <Button size="sm" variant="outline" className="gap-1" onClick={() => openLinks(row.original)}>
              <Link2 className="h-3.5 w-3.5" />
              关联
            </Button>
            <Button
              size="sm"
              variant="destructive"
              className="gap-1"
              onClick={() => setDeleteScenario(row.original)}
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
    data: scenarios,
    columns,
    getRowId: (row) => String(row.scenario.id),
    initialState: { pagination: { pageIndex: 0, pageSize: 10 } },
  });

  return (
    <div className="flex flex-col gap-4">
      {error ? <ErrorState message={error} onRetry={load} /> : null}

      <div className="grid gap-4 lg:grid-cols-[240px_1fr]">
        <aside className="rounded-xl border border-border bg-card p-2">
          <div className="flex flex-col gap-1">
            {groups.map((g) => (
              <button
                key={g.id}
                type="button"
                onClick={() => setActiveId(g.id)}
                className={
                  g.id === activeId
                    ? "flex items-center justify-between rounded-lg bg-primary/10 px-3 py-2 text-left text-sm font-semibold text-primary"
                    : "flex items-center justify-between rounded-lg px-3 py-2 text-left text-sm text-muted-foreground hover:bg-muted"
                }
              >
                <span>{g.name}</span>
                <span className="text-xs">{g.scenario_count}</span>
              </button>
            ))}
          </div>
          <Button
            variant="outline"
            size="sm"
            className="mt-2 w-full"
            disabled={groups.length >= MAX_GROUPS}
            onClick={() => setGroupForm({ name: "", slug: "", description: "", sortOrder: 0 })}
          >
            <Plus className="mr-1 h-3.5 w-3.5" />
            新建场景组{groups.length >= MAX_GROUPS ? "（已达上限）" : ""}
          </Button>
          {activeGroup ? (
            <div className="mt-2 flex flex-col gap-1 border-t border-border pt-2">
              <Button
                variant="ghost"
                size="sm"
                onClick={() =>
                  setGroupForm({
                    id: activeGroup.id,
                    name: activeGroup.name,
                    slug: activeGroup.slug,
                    description: activeGroup.description,
                    sortOrder: activeGroup.sort_order,
                  })
                }
              >
                编辑当前组
              </Button>
              <Button variant="ghost" size="sm" onClick={() => setDeleteGroup(activeGroup)}>
                删除当前组
              </Button>
            </div>
          ) : null}
        </aside>

        <div className="flex min-w-0 flex-col gap-3">
          <div className="flex items-center gap-2">
            <Button variant="outline" size="sm" onClick={() => load()}>
              刷新
            </Button>
            <Button
              size="sm"
              className="ml-auto"
              disabled={!activeGroup}
              onClick={() =>
                setScenarioForm({
                  groupId: activeGroup?.id ?? 0,
                  slug: "",
                  name: "",
                  summary: "",
                  painPoint: "",
                  valueProp: "",
                  iconKey: "",
                  sortOrder: 0,
                  isActive: true,
                })
              }
            >
              新建场景
            </Button>
          </div>

          <DataGrid
            table={table}
            recordCount={scenarios.length}
            isLoading={loading}
            emptyMessage="该组还没有场景"
            tableLayout={{ dense: true, headerBackground: true, stripped: true }}
          >
            <DataGridContainer>
              <DataGridTable />
            </DataGridContainer>
            <DataGridPagination sizes={[10, 20, 50]} rowsPerPageLabel="每页" />
          </DataGrid>
        </div>
      </div>

      <Dialog open={!!groupForm} onOpenChange={(open) => !open && setGroupForm(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{groupForm?.id ? "编辑场景组" : "新建场景组"}</DialogTitle>
          </DialogHeader>
          {groupForm ? (
            <div className="grid gap-3">
              <FormField label="名称">
                <Input value={groupForm.name} onChange={(e) => setGroupForm({ ...groupForm, name: e.target.value })} />
              </FormField>
              <FormField label="slug">
                <Input value={groupForm.slug} onChange={(e) => setGroupForm({ ...groupForm, slug: e.target.value })} />
              </FormField>
              <FormField label="描述">
                <Input
                  value={groupForm.description}
                  onChange={(e) => setGroupForm({ ...groupForm, description: e.target.value })}
                />
              </FormField>
              <FormField label="排序">
                <Input
                  type="number"
                  value={groupForm.sortOrder}
                  onChange={(e) =>
                    setGroupForm({ ...groupForm, sortOrder: Number(e.target.value) || 0 })
                  }
                />
              </FormField>
            </div>
          ) : null}
          <DialogFooter>
            <Button variant="outline" onClick={() => setGroupForm(null)}>
              取消
            </Button>
            <Button disabled={busy} onClick={saveGroup}>
              保存
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={!!scenarioForm} onOpenChange={(open) => !open && setScenarioForm(null)}>
        <DialogContent className="max-h-[85vh] overflow-y-auto sm:max-w-xl">
          <DialogHeader>
            <DialogTitle>{scenarioForm?.id ? "编辑场景" : "新建场景"}</DialogTitle>
          </DialogHeader>
          {scenarioForm ? (
            <div className="grid gap-3 sm:grid-cols-2">
              <FormField label="名称">
                <Input
                  value={scenarioForm.name}
                  onChange={(e) => setScenarioForm({ ...scenarioForm, name: e.target.value })}
                />
              </FormField>
              <FormField label="slug">
                <Input
                  value={scenarioForm.slug}
                  onChange={(e) => setScenarioForm({ ...scenarioForm, slug: e.target.value })}
                />
              </FormField>
              <FormField label="所属组">
                <Select
                  value={String(scenarioForm.groupId)}
                  onValueChange={(v) => setScenarioForm({ ...scenarioForm, groupId: Number(v) })}
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {groups.map((g) => (
                      <SelectItem key={g.id} value={String(g.id)}>
                        {g.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </FormField>
              <FormField label="排序">
                <Input
                  type="number"
                  value={scenarioForm.sortOrder}
                  onChange={(e) =>
                    setScenarioForm({ ...scenarioForm, sortOrder: Number(e.target.value) || 0 })
                  }
                />
              </FormField>
              <FormField label="一句话摘要" className="sm:col-span-2">
                <Input
                  value={scenarioForm.summary}
                  onChange={(e) => setScenarioForm({ ...scenarioForm, summary: e.target.value })}
                />
              </FormField>
              <FormField label="客户痛点" className="sm:col-span-2">
                <Textarea
                  rows={2}
                  value={scenarioForm.painPoint}
                  onChange={(e) => setScenarioForm({ ...scenarioForm, painPoint: e.target.value })}
                />
              </FormField>
              <FormField label="我方价值" className="sm:col-span-2">
                <Textarea
                  rows={2}
                  value={scenarioForm.valueProp}
                  onChange={(e) => setScenarioForm({ ...scenarioForm, valueProp: e.target.value })}
                />
              </FormField>
              <FormField label="状态">
                <Select
                  value={scenarioForm.isActive ? "on" : "off"}
                  onValueChange={(v) => setScenarioForm({ ...scenarioForm, isActive: String(v) === "on" })}
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
            </div>
          ) : null}
          <DialogFooter>
            <Button variant="outline" onClick={() => setScenarioForm(null)}>
              取消
            </Button>
            <Button disabled={busy} onClick={saveScenario}>
              保存
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={!!linksTarget} onOpenChange={(open) => !open && setLinksTarget(null)}>
        <DialogContent className="max-h-[85vh] overflow-y-auto sm:max-w-2xl">
          <DialogHeader>
            <DialogTitle>关联 · {linksTarget?.scenario.name}</DialogTitle>
          </DialogHeader>
          <div className="grid gap-4 sm:grid-cols-2">
            <div>
              <p className="mb-2 text-xs font-semibold text-muted-foreground">产品</p>
              <div className="flex max-h-64 flex-col gap-1 overflow-y-auto rounded-lg border border-border p-2">
                {linkProducts.map((p) => (
                  <label key={p.id} className="flex items-center gap-2 text-xs">
                    <input
                      type="checkbox"
                      checked={pickedProducts.includes(p.id)}
                      onChange={(e) =>
                        setPickedProducts((prev) =>
                          e.target.checked ? [...prev, p.id] : prev.filter((id) => id !== p.id),
                        )
                      }
                    />
                    <span className="truncate">{p.name}</span>
                  </label>
                ))}
              </div>
            </div>
            <div>
              <p className="mb-2 text-xs font-semibold text-muted-foreground">服务</p>
              <div className="flex max-h-64 flex-col gap-1 overflow-y-auto rounded-lg border border-border p-2">
                {linkServices.map((s) => (
                  <label key={s.id} className="flex items-center gap-2 text-xs">
                    <input
                      type="checkbox"
                      checked={pickedServices.includes(s.id)}
                      onChange={(e) =>
                        setPickedServices((prev) =>
                          e.target.checked ? [...prev, s.id] : prev.filter((id) => id !== s.id),
                        )
                      }
                    />
                    <span className="truncate">{s.name}</span>
                  </label>
                ))}
              </div>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setLinksTarget(null)}>
              取消
            </Button>
            <Button disabled={busy} onClick={saveLinks}>
              保存关联
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <ConfirmDialog
        open={!!deleteScenario}
        onOpenChange={(open) => {
          if (!open) setDeleteScenario(null);
        }}
        title="删除场景"
        description={deleteScenario ? `确定删除「${deleteScenario.scenario.name}」？` : undefined}
        confirmLabel="删除"
        destructive
        onConfirm={async () => {
          if (!deleteScenario) return;
          await adminScenariosApi.remove(deleteScenario.scenario.id);
          setDeleteScenario(null);
          await load();
        }}
      />

      <ConfirmDialog
        open={!!deleteGroup}
        onOpenChange={(open) => {
          if (!open) setDeleteGroup(null);
        }}
        title="删除场景组"
        description={
          deleteGroup
            ? `确定删除「${deleteGroup.name}」？组内仍有场景时会被拒绝。`
            : undefined
        }
        confirmLabel="删除"
        destructive
        onConfirm={async () => {
          if (!deleteGroup) return;
          await adminScenariosApi.removeGroup(deleteGroup.id);
          setDeleteGroup(null);
          setActiveId(null);
          await load();
        }}
      />
    </div>
  );
}
