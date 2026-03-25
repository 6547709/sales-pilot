"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Textarea } from "@/components/ui/textarea";
import {
  apiFetch,
  fetchTopology,
  fetchSolutionCategories,
  type TopologyCategoryNode,
  type TopologyResponse,
} from "@/lib/api";

type TopologyLayer = {
  id: number;
  level: number;
  title: string;
  subtitle: string;
  sort_order: number;
};

type AdminCategory = {
  id: number;
  slug: string;
  label: string;
  icon_key: string;
  column_type: string;
  layer_id?: number | null;
  sort_order: number;
  keywords?: unknown;
  hint: string;
  is_active: boolean;
};

type VendorRow = {
  id: number;
  name: string;
  market: string;
  category_id: number;
  category_label: string;
};

function flattenVendors(t: TopologyResponse): VendorRow[] {
  const rows: VendorRow[] = [];
  const push = (node: TopologyCategoryNode) => {
    for (const v of node.vendors.domestic) {
      rows.push({
        id: v.id,
        name: v.name,
        market: "domestic",
        category_id: node.id,
        category_label: node.label,
      });
    }
    for (const v of node.vendors.foreign) {
      rows.push({
        id: v.id,
        name: v.name,
        market: "foreign",
        category_id: node.id,
        category_label: node.label,
      });
    }
  };
  t.security.forEach(push);
  t.ops.forEach(push);
  t.central_layers.forEach((b) => b.categories.forEach(push));
  return rows;
}

function keywordsToText(k: unknown): string {
  if (Array.isArray(k)) return JSON.stringify(k, null, 2);
  if (typeof k === "string") return k;
  return "[]";
}

export default function AdminTopologyPage() {
  const [msg, setMsg] = useState("");
  const [layers, setLayers] = useState<TopologyLayer[]>([]);
  const [categories, setCategories] = useState<AdminCategory[]>([]);
  const [vendorRows, setVendorRows] = useState<VendorRow[]>([]);
  const [loading, setLoading] = useState(true);

  const [layerForm, setLayerForm] = useState({
    level: 3,
    title: "",
    subtitle: "",
    sort_order: 0,
  });

  const [catForm, setCatForm] = useState({
    slug: "",
    label: "",
    icon_key: "layers",
    column_type: "central",
    layer_id: "" as string,
    sort_order: 0,
    keywords: "[]",
    hint: "",
    is_active: true,
  });

  const [editCat, setEditCat] = useState<AdminCategory | null>(null);
  const [editKeywords, setEditKeywords] = useState("[]");

  const [vendorForm, setVendorForm] = useState({
    category_id: "",
    market: "domestic",
    name: "",
    sort_order: 0,
  });

  const refresh = useCallback(async () => {
    setLoading(true);
    setMsg("");
    try {
      const [lr, cr, topo] = await Promise.all([
        apiFetch("/api/v1/admin/topology/layers"),
        fetchSolutionCategories(false),
        fetchTopology().catch(() => null),
      ]);
      if (lr.ok) setLayers(await lr.json());
      setCategories(cr as AdminCategory[]);
      if (topo) setVendorRows(flattenVendors(topo));
      else setVendorRows([]);
    } catch {
      setMsg("加载失败，请检查网络与登录状态");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    refresh();
  }, [refresh]);

  useEffect(() => {
    if (editCat) setEditKeywords(keywordsToText(editCat.keywords));
  }, [editCat]);

  const layerOptions = useMemo(
    () =>
      [...layers].sort((a, b) => b.level - a.level || a.sort_order - b.sort_order),
    [layers],
  );

  async function createLayer() {
    setMsg("");
    const res = await apiFetch("/api/v1/admin/topology/layers", {
      method: "POST",
      body: JSON.stringify({
        level: layerForm.level,
        title: layerForm.title.trim(),
        subtitle: layerForm.subtitle.trim(),
        sort_order: layerForm.sort_order,
      }),
    });
    if (res.ok) {
      setLayerForm({ level: 3, title: "", subtitle: "", sort_order: 0 });
      await refresh();
      setMsg("已新增分层");
    } else {
      const e = await res.json().catch(() => ({}));
      setMsg(e.error || "创建分层失败");
    }
  }

  async function deleteLayer(id: number) {
    if (!confirm("确定删除该分层？其下不能有解决方案。")) return;
    const res = await apiFetch(`/api/v1/admin/topology/layers/${id}`, {
      method: "DELETE",
    });
    if (res.ok) {
      await refresh();
      setMsg("已删除分层");
    } else {
      const e = await res.json().catch(() => ({}));
      setMsg(e.error || "删除失败");
    }
  }

  async function createCategory() {
    setMsg("");
    let keywordsJson: unknown = [];
    try {
      keywordsJson = JSON.parse(catForm.keywords || "[]");
      if (!Array.isArray(keywordsJson)) throw new Error("not array");
    } catch {
      setMsg("关键词须为 JSON 数组，例如 [\"虚拟化\",\"K8s\"]");
      return;
    }
    const lid =
      catForm.column_type === "central" && catForm.layer_id
        ? Number(catForm.layer_id)
        : null;
    const res = await apiFetch("/api/v1/admin/topology/categories", {
      method: "POST",
      body: JSON.stringify({
        slug: catForm.slug.trim(),
        label: catForm.label.trim(),
        icon_key: catForm.icon_key.trim() || "layers",
        column_type: catForm.column_type,
        layer_id: lid,
        sort_order: catForm.sort_order,
        keywords: keywordsJson,
        hint: catForm.hint.trim(),
        is_active: catForm.is_active,
      }),
    });
    if (res.ok) {
      setCatForm({
        slug: "",
        label: "",
        icon_key: "layers",
        column_type: "central",
        layer_id: "",
        sort_order: 0,
        keywords: "[]",
        hint: "",
        is_active: true,
      });
      await refresh();
      setMsg("已新增解决方案节点");
    } else {
      const e = await res.json().catch(() => ({}));
      setMsg(e.error || "创建失败");
    }
  }

  async function saveCategory() {
    if (!editCat) return;
    setMsg("");
    let keywordsJson: unknown = [];
    try {
      keywordsJson = JSON.parse(editKeywords || "[]");
      if (!Array.isArray(keywordsJson)) throw new Error("not array");
    } catch {
      setMsg("关键词须为 JSON 数组");
      return;
    }
    const res = await apiFetch(
      `/api/v1/admin/topology/categories/${editCat.id}`,
      {
        method: "PUT",
        body: JSON.stringify({
          ...editCat,
          keywords: keywordsJson,
        }),
      },
    );
    if (res.ok) {
      setEditCat(null);
      await refresh();
      setMsg("已保存解决方案");
    } else {
      const e = await res.json().catch(() => ({}));
      setMsg(e.error || "保存失败");
    }
  }

  async function deleteCategory(id: number) {
    if (!confirm("确定删除？若仍有产品关联将失败。")) return;
    const res = await apiFetch(`/api/v1/admin/topology/categories/${id}`, {
      method: "DELETE",
    });
    if (res.ok) {
      if (editCat?.id === id) setEditCat(null);
      await refresh();
      setMsg("已删除解决方案");
    } else {
      const e = await res.json().catch(() => ({}));
      setMsg(e.error || "删除失败");
    }
  }

  async function createVendor() {
    setMsg("");
    const cid = Number(vendorForm.category_id);
    if (!cid || !vendorForm.name.trim()) {
      setMsg("请选择解决方案并填写厂商名称");
      return;
    }
    const res = await apiFetch("/api/v1/admin/topology/vendors", {
      method: "POST",
      body: JSON.stringify({
        category_id: cid,
        market: vendorForm.market,
        name: vendorForm.name.trim(),
        sort_order: vendorForm.sort_order,
      }),
    });
    if (res.ok) {
      setVendorForm({
        category_id: vendorForm.category_id,
        market: "domestic",
        name: "",
        sort_order: 0,
      });
      await refresh();
      setMsg("已添加厂商");
    } else {
      const e = await res.json().catch(() => ({}));
      setMsg(e.error || "添加失败");
    }
  }

  async function deleteVendor(id: number) {
    if (!confirm("确定删除该厂商？")) return;
    const res = await apiFetch(`/api/v1/admin/topology/vendors/${id}`, {
      method: "DELETE",
    });
    if (res.ok) {
      await refresh();
      setMsg("已删除厂商");
    } else {
      setMsg("删除失败");
    }
  }

  return (
    <div className="mx-auto max-w-5xl space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <Link href="/admin" className="text-sm text-primary">
            ← 控制台
          </Link>
          <h1 className="mt-1 text-xl font-bold text-primary">拓扑维护</h1>
          <p className="text-sm text-muted-foreground">
            管理首页架构分层、解决方案卡片与厂商（国内 / 国外）。停用方案仍可在后台编辑，但首页不展示。
          </p>
        </div>
        <Button variant="outline" type="button" onClick={() => refresh()}>
          刷新
        </Button>
      </div>

      {msg ? (
        <p className="rounded-md border border-primary/30 bg-primary/5 px-3 py-2 text-sm">
          {msg}
        </p>
      ) : null}

      {loading ? (
        <p className="text-sm text-muted-foreground">加载中…</p>
      ) : (
        <Tabs defaultValue="layers">
          <TabsList className="flex-wrap">
            <TabsTrigger value="layers">分层</TabsTrigger>
            <TabsTrigger value="categories">解决方案</TabsTrigger>
            <TabsTrigger value="vendors">厂商</TabsTrigger>
          </TabsList>

          <TabsContent value="layers" className="space-y-6 pt-4">
            <div className="rounded-lg border p-4">
              <h2 className="mb-3 font-semibold">新增分层（中部栏）</h2>
              <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
                <div className="space-y-1.5">
                  <Label>level（越大越靠上）</Label>
                  <Input
                    type="number"
                    value={layerForm.level}
                    onChange={(e) =>
                      setLayerForm({
                        ...layerForm,
                        level: Number(e.target.value) || 0,
                      })
                    }
                  />
                </div>
                <div className="space-y-1.5 sm:col-span-2">
                  <Label>标题</Label>
                  <Input
                    value={layerForm.title}
                    onChange={(e) =>
                      setLayerForm({ ...layerForm, title: e.target.value })
                    }
                  />
                </div>
                <div className="space-y-1.5">
                  <Label>排序</Label>
                  <Input
                    type="number"
                    value={layerForm.sort_order}
                    onChange={(e) =>
                      setLayerForm({
                        ...layerForm,
                        sort_order: Number(e.target.value) || 0,
                      })
                    }
                  />
                </div>
                <div className="space-y-1.5 sm:col-span-2 lg:col-span-4">
                  <Label>副标题</Label>
                  <Input
                    value={layerForm.subtitle}
                    onChange={(e) =>
                      setLayerForm({ ...layerForm, subtitle: e.target.value })
                    }
                  />
                </div>
              </div>
              <Button className="mt-3" type="button" onClick={createLayer}>
                创建分层
              </Button>
            </div>

            <div className="overflow-x-auto rounded-lg border">
              <table className="w-full text-left text-sm">
                <thead className="border-b bg-muted/40">
                  <tr>
                    <th className="p-3">ID</th>
                    <th className="p-3">level</th>
                    <th className="p-3">标题</th>
                    <th className="p-3">副标题</th>
                    <th className="p-3">排序</th>
                    <th className="p-3">操作</th>
                  </tr>
                </thead>
                <tbody>
                  {layers.map((L) => (
                    <tr key={L.id} className="border-b last:border-0">
                      <td className="p-3">{L.id}</td>
                      <td className="p-3">{L.level}</td>
                      <td className="p-3 font-medium">{L.title}</td>
                      <td className="p-3 text-muted-foreground">
                        {L.subtitle}
                      </td>
                      <td className="p-3">{L.sort_order}</td>
                      <td className="p-3">
                        <Button
                          variant="destructive"
                          size="sm"
                          type="button"
                          onClick={() => deleteLayer(L.id)}
                        >
                          删除
                        </Button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </TabsContent>

          <TabsContent value="categories" className="space-y-6 pt-4">
            <div className="rounded-lg border p-4">
              <h2 className="mb-3 font-semibold">新增解决方案节点</h2>
              <div className="grid gap-3 sm:grid-cols-2">
                <div className="space-y-1.5">
                  <Label>slug（英文唯一）</Label>
                  <Input
                    value={catForm.slug}
                    onChange={(e) =>
                      setCatForm({ ...catForm, slug: e.target.value })
                    }
                  />
                </div>
                <div className="space-y-1.5">
                  <Label>显示名称</Label>
                  <Input
                    value={catForm.label}
                    onChange={(e) =>
                      setCatForm({ ...catForm, label: e.target.value })
                    }
                  />
                </div>
                <div className="space-y-1.5">
                  <Label>图标 key（Lucide 名，如 server、shield）</Label>
                  <Input
                    value={catForm.icon_key}
                    onChange={(e) =>
                      setCatForm({ ...catForm, icon_key: e.target.value })
                    }
                  />
                </div>
                <div className="space-y-1.5">
                  <Label>栏目</Label>
                  <select
                    className="flex h-10 w-full rounded-md border border-input bg-background px-3 text-sm"
                    value={catForm.column_type}
                    onChange={(e) =>
                      setCatForm({ ...catForm, column_type: e.target.value })
                    }
                  >
                    <option value="security">安全（左栏）</option>
                    <option value="ops">运维（右栏）</option>
                    <option value="central">中部（需选分层）</option>
                  </select>
                </div>
                {catForm.column_type === "central" ? (
                  <div className="space-y-1.5">
                    <Label>所属分层</Label>
                    <select
                      className="flex h-10 w-full rounded-md border border-input bg-background px-3 text-sm"
                      value={catForm.layer_id}
                      onChange={(e) =>
                        setCatForm({ ...catForm, layer_id: e.target.value })
                      }
                    >
                      <option value="">请选择</option>
                      {layerOptions.map((L) => (
                        <option key={L.id} value={L.id}>
                          L{L.level} · {L.title}
                        </option>
                      ))}
                    </select>
                  </div>
                ) : null}
                <div className="space-y-1.5">
                  <Label>排序</Label>
                  <Input
                    type="number"
                    value={catForm.sort_order}
                    onChange={(e) =>
                      setCatForm({
                        ...catForm,
                        sort_order: Number(e.target.value) || 0,
                      })
                    }
                  />
                </div>
                <div className="flex items-center gap-2 sm:col-span-2">
                  <input
                    id="cat-active"
                    type="checkbox"
                    checked={catForm.is_active}
                    onChange={(e) =>
                      setCatForm({ ...catForm, is_active: e.target.checked })
                    }
                  />
                  <Label htmlFor="cat-active" className="font-normal">
                    启用（首页展示）
                  </Label>
                </div>
                <div className="space-y-1.5 sm:col-span-2">
                  <Label>关键词（JSON 数组，用于首页模糊匹配）</Label>
                  <Textarea
                    rows={2}
                    className="font-mono text-xs"
                    value={catForm.keywords}
                    onChange={(e) =>
                      setCatForm({ ...catForm, keywords: e.target.value })
                    }
                  />
                </div>
                <div className="space-y-1.5 sm:col-span-2">
                  <Label>悬停提示 hint</Label>
                  <Input
                    value={catForm.hint}
                    onChange={(e) =>
                      setCatForm({ ...catForm, hint: e.target.value })
                    }
                  />
                </div>
              </div>
              <Button className="mt-3" type="button" onClick={createCategory}>
                创建
              </Button>
            </div>

            <div className="overflow-x-auto rounded-lg border">
              <table className="w-full text-left text-sm">
                <thead className="border-b bg-muted/40">
                  <tr>
                    <th className="p-3">ID</th>
                    <th className="p-3">名称</th>
                    <th className="p-3">slug</th>
                    <th className="p-3">栏目</th>
                    <th className="p-3">分层</th>
                    <th className="p-3">状态</th>
                    <th className="p-3">操作</th>
                  </tr>
                </thead>
                <tbody>
                  {categories.map((c) => (
                    <tr key={c.id} className="border-b last:border-0">
                      <td className="p-3">{c.id}</td>
                      <td className="p-3 font-medium">{c.label}</td>
                      <td className="p-3 text-muted-foreground">{c.slug}</td>
                      <td className="p-3">{c.column_type}</td>
                      <td className="p-3">{c.layer_id ?? "—"}</td>
                      <td className="p-3">
                        {c.is_active ? "启用" : "停用"}
                      </td>
                      <td className="p-3">
                        <Button
                          variant="outline"
                          size="sm"
                          type="button"
                          className="mr-2"
                          onClick={() => setEditCat(c)}
                        >
                          编辑
                        </Button>
                        <Button
                          variant="destructive"
                          size="sm"
                          type="button"
                          onClick={() => deleteCategory(c.id)}
                        >
                          删
                        </Button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            {editCat ? (
              <div className="rounded-lg border border-primary/30 bg-primary/[0.03] p-4">
                <h3 className="mb-3 font-semibold">
                  编辑 #{editCat.id} {editCat.label}
                </h3>
                <div className="grid gap-3 sm:grid-cols-2">
                  <div className="space-y-1.5">
                    <Label>slug</Label>
                    <Input
                      value={editCat.slug}
                      onChange={(e) =>
                        setEditCat({ ...editCat, slug: e.target.value })
                      }
                    />
                  </div>
                  <div className="space-y-1.5">
                    <Label>名称</Label>
                    <Input
                      value={editCat.label}
                      onChange={(e) =>
                        setEditCat({ ...editCat, label: e.target.value })
                      }
                    />
                  </div>
                  <div className="space-y-1.5">
                    <Label>icon_key</Label>
                    <Input
                      value={editCat.icon_key}
                      onChange={(e) =>
                        setEditCat({ ...editCat, icon_key: e.target.value })
                      }
                    />
                  </div>
                  <div className="space-y-1.5">
                    <Label>栏目</Label>
                    <select
                      className="flex h-10 w-full rounded-md border border-input bg-background px-3 text-sm"
                      value={editCat.column_type}
                      onChange={(e) =>
                        setEditCat({
                          ...editCat,
                          column_type: e.target.value,
                        })
                      }
                    >
                      <option value="security">security</option>
                      <option value="ops">ops</option>
                      <option value="central">central</option>
                    </select>
                  </div>
                  <div className="space-y-1.5">
                    <Label>layer_id（central 时填分层 ID）</Label>
                    <Input
                      type="number"
                      value={editCat.layer_id ?? ""}
                      onChange={(e) => {
                        const v = e.target.value;
                        setEditCat({
                          ...editCat,
                          layer_id: v ? Number(v) : null,
                        });
                      }}
                    />
                  </div>
                  <div className="space-y-1.5">
                    <Label>排序</Label>
                    <Input
                      type="number"
                      value={editCat.sort_order}
                      onChange={(e) =>
                        setEditCat({
                          ...editCat,
                          sort_order: Number(e.target.value) || 0,
                        })
                      }
                    />
                  </div>
                  <div className="flex items-center gap-2 sm:col-span-2">
                    <input
                      id="edit-active"
                      type="checkbox"
                      checked={editCat.is_active}
                      onChange={(e) =>
                        setEditCat({
                          ...editCat,
                          is_active: e.target.checked,
                        })
                      }
                    />
                    <Label htmlFor="edit-active" className="font-normal">
                      启用
                    </Label>
                  </div>
                  <div className="space-y-1.5 sm:col-span-2">
                    <Label>关键词 JSON</Label>
                    <Textarea
                      rows={3}
                      className="font-mono text-xs"
                      value={editKeywords}
                      onChange={(e) => setEditKeywords(e.target.value)}
                    />
                  </div>
                  <div className="space-y-1.5 sm:col-span-2">
                    <Label>hint</Label>
                    <Input
                      value={editCat.hint}
                      onChange={(e) =>
                        setEditCat({ ...editCat, hint: e.target.value })
                      }
                    />
                  </div>
                </div>
                <div className="mt-3 flex flex-wrap gap-2">
                  <Button type="button" onClick={saveCategory}>
                    保存
                  </Button>
                  <Button
                    variant="outline"
                    type="button"
                    onClick={() => setEditCat(null)}
                  >
                    取消
                  </Button>
                </div>
              </div>
            ) : null}
          </TabsContent>

          <TabsContent value="vendors" className="space-y-6 pt-4">
            <p className="text-sm text-muted-foreground">
              下列表来自当前<strong>已启用</strong>方案的公开拓扑；新增厂商后请点「刷新」。若方案已停用，仍可在此选择 ID
              添加，刷新后可能暂不在下表出现。
            </p>
            <div className="rounded-lg border p-4">
              <h2 className="mb-3 font-semibold">添加厂商</h2>
              <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
                <div className="space-y-1.5 sm:col-span-2">
                  <Label>解决方案</Label>
                  <select
                    className="flex h-10 w-full rounded-md border border-input bg-background px-3 text-sm"
                    value={vendorForm.category_id}
                    onChange={(e) =>
                      setVendorForm({
                        ...vendorForm,
                        category_id: e.target.value,
                      })
                    }
                  >
                    <option value="">请选择</option>
                    {categories.map((c) => (
                      <option key={c.id} value={c.id}>
                        {c.label} (#{c.id})
                      </option>
                    ))}
                  </select>
                </div>
                <div className="space-y-1.5">
                  <Label>市场</Label>
                  <select
                    className="flex h-10 w-full rounded-md border border-input bg-background px-3 text-sm"
                    value={vendorForm.market}
                    onChange={(e) =>
                      setVendorForm({ ...vendorForm, market: e.target.value })
                    }
                  >
                    <option value="domestic">国内</option>
                    <option value="foreign">国外</option>
                  </select>
                </div>
                <div className="space-y-1.5">
                  <Label>排序</Label>
                  <Input
                    type="number"
                    value={vendorForm.sort_order}
                    onChange={(e) =>
                      setVendorForm({
                        ...vendorForm,
                        sort_order: Number(e.target.value) || 0,
                      })
                    }
                  />
                </div>
                <div className="space-y-1.5 sm:col-span-2 lg:col-span-4">
                  <Label>厂商名称</Label>
                  <Input
                    value={vendorForm.name}
                    onChange={(e) =>
                      setVendorForm({ ...vendorForm, name: e.target.value })
                    }
                  />
                </div>
              </div>
              <Button className="mt-3" type="button" onClick={createVendor}>
                添加
              </Button>
            </div>

            <div className="overflow-x-auto rounded-lg border">
              <table className="w-full text-left text-sm">
                <thead className="border-b bg-muted/40">
                  <tr>
                    <th className="p-3">ID</th>
                    <th className="p-3">方案</th>
                    <th className="p-3">市场</th>
                    <th className="p-3">名称</th>
                    <th className="p-3">操作</th>
                  </tr>
                </thead>
                <tbody>
                  {vendorRows.map((v) => (
                    <tr key={v.id} className="border-b last:border-0">
                      <td className="p-3">{v.id}</td>
                      <td className="p-3">
                        {v.category_label}{" "}
                        <span className="text-muted-foreground">
                          (#{v.category_id})
                        </span>
                      </td>
                      <td className="p-3">
                        {v.market === "foreign" ? "国外" : "国内"}
                      </td>
                      <td className="p-3 font-medium">{v.name}</td>
                      <td className="p-3">
                        <Button
                          variant="destructive"
                          size="sm"
                          type="button"
                          onClick={() => deleteVendor(v.id)}
                        >
                          删除
                        </Button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
              {vendorRows.length === 0 ? (
                <p className="p-6 text-center text-sm text-muted-foreground">
                  暂无数据或拓扑加载失败
                </p>
              ) : null}
            </div>
          </TabsContent>
        </Tabs>
      )}
    </div>
  );
}
