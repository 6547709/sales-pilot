"use client";

import { useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { useParams, useRouter } from "next/navigation";
import { MarkdownContent } from "@/components/ui/markdown-content";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Textarea } from "@/components/ui/textarea";
import {
  apiFetch,
  fetchSolutionCategories,
  parseJsonArray,
  parsePersonas,
  type Case,
  type Product,
} from "@/lib/api";
import {
  mergePersonasFromForm,
  PERSONA_PRESET_ROLES,
  splitPersonasForForm,
  type PersonaExtraRow,
} from "@/lib/personas-presets";
import { ClipboardList } from "lucide-react";

function emptyPresetRecord(): Record<string, string> {
  return Object.fromEntries(PERSONA_PRESET_ROLES.map((r) => [r, ""]));
}

export default function AdminProductEditPage() {
  const params = useParams();
  const router = useRouter();
  const id = params.id as string;
  const [product, setProduct] = useState<Product | null>(null);
  const [h1, setH1] = useState("");
  const [h2, setH2] = useState("");
  const [h3, setH3] = useState("");
  const [q1, setQ1] = useState("");
  const [q2, setQ2] = useState("");
  const [q3, setQ3] = useState("");
  const [presetPersonas, setPresetPersonas] = useState<Record<string, string>>(
    () => emptyPresetRecord(),
  );
  const [extraPersonas, setExtraPersonas] = useState<PersonaExtraRow[]>([]);
  const [msg, setMsg] = useState("");
  const [solutionOptions, setSolutionOptions] = useState<
    { id: number; label: string; slug: string; is_active: boolean }[]
  >([]);
  const [cases, setCases] = useState<Case[]>([]);
  const [newCase, setNewCase] = useState({
    client_name: "",
    pain_points: "",
    solution: "",
    value_delivered: "",
  });

  useEffect(() => {
    fetchSolutionCategories(false).then(setSolutionOptions);
  }, []);

  const refreshCases = useCallback(async () => {
    const cr = await apiFetch(`/api/v1/admin/products/${id}/cases`);
    if (cr.ok) setCases(await cr.json());
    else setCases([]);
  }, [id]);

  const load = useCallback(async () => {
    const [pr, cr] = await Promise.all([
      apiFetch(`/api/v1/admin/products/${id}`),
      apiFetch(`/api/v1/admin/products/${id}/cases`),
    ]);
    if (!pr.ok) {
      setProduct(null);
      return;
    }
    const p: Product = await pr.json();
    setProduct(p);
    const hl = parseJsonArray(p.highlights);
    setH1(hl[0] || "");
    setH2(hl[1] || "");
    setH3(hl[2] || "");
    const qs = parseJsonArray(p.discovery_questions);
    setQ1(qs[0] || "");
    setQ2(qs[1] || "");
    setQ3(qs[2] || "");
    const { preset, extra } = splitPersonasForForm(
      parsePersonas(p.target_personas),
    );
    setPresetPersonas(preset);
    setExtraPersonas(extra);
    if (cr.ok) setCases(await cr.json());
    else setCases([]);
  }, [id]);

  useEffect(() => {
    load();
  }, [load]);

  async function save() {
    if (!product) return;
    setMsg("");
    const personas = mergePersonasFromForm(presetPersonas, extraPersonas);
    const { solution_category: _sc, ...clean } = product;
    const body = {
      ...clean,
      solution_category_id: product.solution_category_id ?? null,
      is_draft: product.is_draft ?? false,
      highlights: [h1, h2, h3].filter(Boolean),
      discovery_questions: [q1, q2, q3].filter(Boolean),
      target_personas: personas,
    };
    const res = await apiFetch(`/api/v1/admin/products/${id}`, {
      method: "PUT",
      body: JSON.stringify(body),
    });
    if (res.ok) {
      router.replace("/admin/products");
    } else {
      const e = await res.json().catch(() => ({}));
      setMsg(e.error || "保存失败");
    }
  }

  async function remove() {
    if (!confirm("确定删除该产品及下属话术、案例？")) return;
    const res = await apiFetch(`/api/v1/admin/products/${id}`, {
      method: "DELETE",
    });
    if (res.ok) router.replace("/admin/products");
  }

  if (!product) {
    return <p className="text-muted-foreground">加载中…</p>;
  }

  return (
    <div className="mx-auto max-w-4xl space-y-6">
      <div className="flex flex-wrap items-center gap-4">
        <Link href="/admin/products" className="text-sm text-primary">
          ← 列表
        </Link>
        <h1 className="text-xl font-bold text-primary">编辑产品</h1>
      </div>

      <details className="rounded-lg border border-primary/20 bg-primary/[0.04] p-4 text-sm">
        <summary className="cursor-pointer font-medium text-primary">
          填写说明（建议先看）
        </summary>
        <ul className="mt-3 list-disc space-y-1.5 pl-5 text-muted-foreground">
          <li>
            <strong className="text-foreground">所属解决方案</strong>
            与首页架构图节点一致；保存后「分类」字段会与方案名称自动对齐。
          </li>
          <li>
            <strong className="text-foreground">市场定位</strong>{" "}
            前台列表、全局架构卡与详情会显示「国内」或「国外」标签；选「不区分」则不显示地域标签。方案库仍可按国内/国外筛选。
          </li>
          <li>
            <strong className="text-foreground">概览 / 竞品</strong>
            支持 Markdown 与 GFM 表格；亮点与三问用短句，便于销售扫读。
          </li>
          <li>
            <strong className="text-foreground">目标画像</strong>{" "}
            按预设角色填写；亦可「添加其他角色」自定义条目。
          </li>
          <li>
            <strong className="text-foreground">客户案例</strong>{" "}
            在页面下方单独维护，与前台详情页「客户案例」区块一致。
          </li>
        </ul>
      </details>

      <div className="grid gap-4 sm:grid-cols-2">
        <div className="space-y-2">
          <Label>名称</Label>
          <Input
            value={product.name}
            onChange={(e) =>
              setProduct({ ...product, name: e.target.value })
            }
          />
        </div>
        <div className="space-y-2">
          <Label>所属解决方案</Label>
          <select
            className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
            value={product.solution_category_id ?? ""}
            onChange={(e) => {
              const raw = e.target.value;
              const idNum = raw ? Number(raw) : null;
              const hit = solutionOptions.find((c) => c.id === idNum);
              setProduct({
                ...product,
                solution_category_id: idNum,
                category: hit?.label ?? product.category,
              });
            }}
          >
            <option value="">（不关联首页方案）</option>
            {solutionOptions.map((c) => (
              <option key={c.id} value={c.id}>
                {c.label} ({c.slug})
                {!c.is_active ? " · 已停用" : ""}
              </option>
            ))}
          </select>
        </div>
      </div>

      <div className="space-y-2">
        <Label>市场定位（国内 / 国外标签）</Label>
        <select
          className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring sm:max-w-md"
          value={(product.vendor_market || "all").toLowerCase()}
          onChange={(e) =>
            setProduct({ ...product, vendor_market: e.target.value })
          }
        >
          <option value="all">不区分（不展示国内/国外标签）</option>
          <option value="domestic">国内</option>
          <option value="foreign">国外</option>
        </select>
        <p className="text-xs text-muted-foreground">
          例如 VMware、OpenShift 等可归为「国外」；与「所属解决方案」配合后，首页同类方案下可一眼区分国内与国外产品。
        </p>
      </div>

      <div className="grid gap-4 sm:grid-cols-2">
        <div className="space-y-2">
          <Label>分类（展示文案）</Label>
          <Input
            value={product.category}
            readOnly={!!product.solution_category_id}
            className={product.solution_category_id ? "bg-muted/60" : ""}
            onChange={(e) =>
              setProduct({ ...product, category: e.target.value })
            }
          />
          {product.solution_category_id ? (
            <p className="text-xs text-muted-foreground">
              已关联解决方案，保存时后台会用方案名称覆盖此字段。
            </p>
          ) : null}
        </div>
        <div className="space-y-2">
          <Label>厂商名称</Label>
          <Input
            value={product.manufacturer_name ?? ""}
            onChange={(e) =>
              setProduct({
                ...product,
                manufacturer_name: e.target.value,
              })
            }
          />
        </div>
      </div>

      <div className="flex items-center gap-2 rounded-lg border border-dashed border-border bg-muted/20 px-3 py-2">
        <input
          id="is-draft"
          type="checkbox"
          checked={product.is_draft ?? false}
          onChange={(e) =>
            setProduct({ ...product, is_draft: e.target.checked })
          }
        />
        <Label htmlFor="is-draft" className="cursor-pointer font-normal">
          草稿（前台不展示；可在列表中一键发布）
        </Label>
      </div>

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <div className="space-y-2">
          <Label>销售联系人</Label>
          <Input
            value={product.sales_contact_name ?? ""}
            onChange={(e) =>
              setProduct({
                ...product,
                sales_contact_name: e.target.value,
              })
            }
          />
        </div>
        <div className="space-y-2">
          <Label>销售电话</Label>
          <Input
            value={product.sales_contact_phone ?? ""}
            onChange={(e) =>
              setProduct({
                ...product,
                sales_contact_phone: e.target.value,
              })
            }
          />
        </div>
        <div className="space-y-2">
          <Label>销售邮箱</Label>
          <Input
            value={product.sales_contact_email ?? ""}
            onChange={(e) =>
              setProduct({
                ...product,
                sales_contact_email: e.target.value,
              })
            }
          />
        </div>
        <div className="space-y-2">
          <Label>售前联系人</Label>
          <Input
            value={product.presales_contact_name ?? ""}
            onChange={(e) =>
              setProduct({
                ...product,
                presales_contact_name: e.target.value,
              })
            }
          />
        </div>
        <div className="space-y-2">
          <Label>售前电话</Label>
          <Input
            value={product.presales_contact_phone ?? ""}
            onChange={(e) =>
              setProduct({
                ...product,
                presales_contact_phone: e.target.value,
              })
            }
          />
        </div>
        <div className="space-y-2">
          <Label>售前邮箱</Label>
          <Input
            value={product.presales_contact_email ?? ""}
            onChange={(e) =>
              setProduct({
                ...product,
                presales_contact_email: e.target.value,
              })
            }
          />
        </div>
      </div>

      <div className="space-y-2">
        <Label>描述（支持 Markdown）</Label>
        <Tabs defaultValue="edit">
          <TabsList>
            <TabsTrigger value="edit">编辑</TabsTrigger>
            <TabsTrigger value="preview">预览</TabsTrigger>
          </TabsList>
          <TabsContent value="edit">
            <Textarea
              rows={8}
              value={product.description}
              onChange={(e) =>
                setProduct({ ...product, description: e.target.value })
              }
            />
          </TabsContent>
          <TabsContent value="preview">
            <div className="rounded-md border p-4">
              <MarkdownContent>{product.description}</MarkdownContent>
            </div>
          </TabsContent>
        </Tabs>
      </div>

      <div className="space-y-3">
        <Label>三大亮点（一句话利益点）</Label>
        <Input placeholder="亮点 1" value={h1} onChange={(e) => setH1(e.target.value)} />
        <Input placeholder="亮点 2" value={h2} onChange={(e) => setH2(e.target.value)} />
        <Input placeholder="亮点 3" value={h3} onChange={(e) => setH3(e.target.value)} />
      </div>

      <div className="space-y-3">
        <Label>黄金三问</Label>
        <Textarea rows={2} value={q1} onChange={(e) => setQ1(e.target.value)} />
        <Textarea rows={2} value={q2} onChange={(e) => setQ2(e.target.value)} />
        <Textarea rows={2} value={q3} onChange={(e) => setQ3(e.target.value)} />
      </div>

      <div className="space-y-2">
        <Label>触发事件</Label>
        <Textarea
          rows={3}
          value={product.trigger_events}
          onChange={(e) =>
            setProduct({ ...product, trigger_events: e.target.value })
          }
        />
      </div>

      <div className="space-y-4 rounded-lg border border-border bg-muted/20 p-4">
        <div>
          <Label className="text-base">目标画像</Label>
          <p className="mt-1 text-xs text-muted-foreground">
            预设常见决策角色；留空则不展示该角色。可添加自定义角色行。
          </p>
        </div>
        <div className="grid gap-3 sm:grid-cols-2">
          {PERSONA_PRESET_ROLES.map((role) => (
            <div key={role} className="space-y-1.5">
              <Label className="text-xs font-medium text-muted-foreground">
                {role}
              </Label>
              <Textarea
                rows={2}
                placeholder="一句话描述该角色关注点…"
                value={presetPersonas[role] ?? ""}
                onChange={(e) =>
                  setPresetPersonas((prev) => ({
                    ...prev,
                    [role]: e.target.value,
                  }))
                }
              />
            </div>
          ))}
        </div>
        <div className="space-y-3 border-t border-border/80 pt-4">
          <Label className="text-sm">其他角色</Label>
          {extraPersonas.length === 0 ? (
            <p className="text-xs text-muted-foreground">暂无自定义角色</p>
          ) : (
            <ul className="space-y-3">
              {extraPersonas.map((row) => (
                <li
                  key={row.id}
                  className="grid gap-2 rounded-md border border-border/60 bg-background p-3 sm:grid-cols-[minmax(0,140px)_1fr_auto]"
                >
                  <Input
                    placeholder="角色名称"
                    value={row.key}
                    onChange={(e) =>
                      setExtraPersonas((prev) =>
                        prev.map((x) =>
                          x.id === row.id ? { ...x, key: e.target.value } : x,
                        ),
                      )
                    }
                  />
                  <Textarea
                    rows={2}
                    placeholder="描述"
                    value={row.value}
                    onChange={(e) =>
                      setExtraPersonas((prev) =>
                        prev.map((x) =>
                          x.id === row.id ? { ...x, value: e.target.value } : x,
                        ),
                      )
                    }
                  />
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    className="justify-self-end sm:justify-self-center"
                    onClick={() =>
                      setExtraPersonas((prev) =>
                        prev.filter((x) => x.id !== row.id),
                      )
                    }
                  >
                    移除
                  </Button>
                </li>
              ))}
            </ul>
          )}
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() =>
              setExtraPersonas((prev) => [
                ...prev,
                { id: Date.now(), key: "", value: "" },
              ])
            }
          >
            添加其他角色
          </Button>
        </div>
      </div>

      <div className="space-y-2">
        <Label>竞品分析（Markdown）</Label>
        <Tabs defaultValue="edit">
          <TabsList>
            <TabsTrigger value="edit">编辑</TabsTrigger>
            <TabsTrigger value="preview">预览</TabsTrigger>
          </TabsList>
          <TabsContent value="edit">
            <Textarea
              rows={12}
              value={product.competitor_analysis}
              onChange={(e) =>
                setProduct({
                  ...product,
                  competitor_analysis: e.target.value,
                })
              }
            />
          </TabsContent>
          <TabsContent value="preview">
            <div className="rounded-md border p-4">
              <MarkdownContent>{product.competitor_analysis}</MarkdownContent>
            </div>
          </TabsContent>
        </Tabs>
      </div>

      <div className="space-y-2">
        <Label>ROI 指标</Label>
        <Textarea
          rows={4}
          value={product.roi_metrics}
          onChange={(e) =>
            setProduct({ ...product, roi_metrics: e.target.value })
          }
        />
      </div>

      <Card className="border-primary/15">
        <CardHeader className="pb-3">
          <CardTitle className="flex items-center gap-2 text-lg">
            <ClipboardList className="h-5 w-5 text-primary" />
            客户案例
          </CardTitle>
          <p className="text-sm text-muted-foreground">
            每条案例包含客户名称、痛点、方案与价值；保存后立即作用于前台产品详情页。
          </p>
        </CardHeader>
        <CardContent className="space-y-6">
          {cases.map((c, idx) => (
            <div key={c.id}>
              {idx > 0 ? <Separator className="mb-6" /> : null}
              <div className="mb-2 flex flex-wrap items-center justify-between gap-2">
                <span className="text-sm font-medium text-muted-foreground">
                  案例 #{c.id}
                </span>
                <div className="flex flex-wrap gap-2">
                  <Button
                    type="button"
                    size="sm"
                    variant="secondary"
                    onClick={async () => {
                      setMsg("");
                      const res = await apiFetch(
                        `/api/v1/admin/cases/${c.id}`,
                        {
                          method: "PUT",
                          body: JSON.stringify(c),
                        },
                      );
                      if (res.ok) {
                        setMsg("案例已保存");
                        await refreshCases();
                      } else {
                        const e = await res.json().catch(() => ({}));
                        setMsg(e.error || "案例保存失败");
                      }
                    }}
                  >
                    保存本条
                  </Button>
                  <Button
                    type="button"
                    size="sm"
                    variant="destructive"
                    onClick={async () => {
                      if (!confirm("确定删除该案例？")) return;
                      setMsg("");
                      const res = await apiFetch(
                        `/api/v1/admin/cases/${c.id}`,
                        { method: "DELETE" },
                      );
                      if (res.ok) {
                        setMsg("已删除案例");
                        await refreshCases();
                      } else setMsg("删除失败");
                    }}
                  >
                    删除
                  </Button>
                </div>
              </div>
              <div className="grid gap-3 sm:grid-cols-2">
                <div className="space-y-1.5 sm:col-span-2">
                  <Label>客户名称</Label>
                  <Input
                    value={c.client_name}
                    onChange={(e) =>
                      setCases((prev) =>
                        prev.map((x) =>
                          x.id === c.id
                            ? { ...x, client_name: e.target.value }
                            : x,
                        ),
                      )
                    }
                  />
                </div>
                <div className="space-y-1.5 sm:col-span-2">
                  <Label>痛点</Label>
                  <Textarea
                    rows={3}
                    value={c.pain_points}
                    onChange={(e) =>
                      setCases((prev) =>
                        prev.map((x) =>
                          x.id === c.id
                            ? { ...x, pain_points: e.target.value }
                            : x,
                        ),
                      )
                    }
                  />
                </div>
                <div className="space-y-1.5 sm:col-span-2">
                  <Label>方案</Label>
                  <Textarea
                    rows={3}
                    value={c.solution}
                    onChange={(e) =>
                      setCases((prev) =>
                        prev.map((x) =>
                          x.id === c.id
                            ? { ...x, solution: e.target.value }
                            : x,
                        ),
                      )
                    }
                  />
                </div>
                <div className="space-y-1.5 sm:col-span-2">
                  <Label>交付价值</Label>
                  <Textarea
                    rows={3}
                    value={c.value_delivered}
                    onChange={(e) =>
                      setCases((prev) =>
                        prev.map((x) =>
                          x.id === c.id
                            ? { ...x, value_delivered: e.target.value }
                            : x,
                        ),
                      )
                    }
                  />
                </div>
              </div>
            </div>
          ))}

          <Separator />

          <div>
            <p className="mb-3 text-sm font-medium">新增案例</p>
            <div className="grid gap-3 sm:grid-cols-2">
              <div className="space-y-1.5 sm:col-span-2">
                <Label>客户名称</Label>
                <Input
                  value={newCase.client_name}
                  onChange={(e) =>
                    setNewCase({ ...newCase, client_name: e.target.value })
                  }
                />
              </div>
              <div className="space-y-1.5 sm:col-span-2">
                <Label>痛点</Label>
                <Textarea
                  rows={3}
                  value={newCase.pain_points}
                  onChange={(e) =>
                    setNewCase({ ...newCase, pain_points: e.target.value })
                  }
                />
              </div>
              <div className="space-y-1.5 sm:col-span-2">
                <Label>方案</Label>
                <Textarea
                  rows={3}
                  value={newCase.solution}
                  onChange={(e) =>
                    setNewCase({ ...newCase, solution: e.target.value })
                  }
                />
              </div>
              <div className="space-y-1.5 sm:col-span-2">
                <Label>交付价值</Label>
                <Textarea
                  rows={3}
                  value={newCase.value_delivered}
                  onChange={(e) =>
                    setNewCase({ ...newCase, value_delivered: e.target.value })
                  }
                />
              </div>
            </div>
            <Button
              className="mt-3"
              type="button"
              variant="outline"
              onClick={async () => {
                setMsg("");
                const res = await apiFetch(
                  `/api/v1/admin/products/${id}/cases`,
                  {
                    method: "POST",
                    body: JSON.stringify({
                      client_name: newCase.client_name.trim(),
                      pain_points: newCase.pain_points,
                      solution: newCase.solution,
                      value_delivered: newCase.value_delivered,
                    }),
                  },
                );
                if (res.ok) {
                  setNewCase({
                    client_name: "",
                    pain_points: "",
                    solution: "",
                    value_delivered: "",
                  });
                  setMsg("已新增案例");
                  await refreshCases();
                } else {
                  const e = await res.json().catch(() => ({}));
                  setMsg(e.error || "新增失败");
                }
              }}
            >
              添加案例
            </Button>
          </div>
        </CardContent>
      </Card>

      {msg ? <p className="text-sm text-primary">{msg}</p> : null}
      <div className="flex flex-wrap gap-2">
        <Button onClick={save}>保存</Button>
        <Button variant="destructive" type="button" onClick={remove}>
          删除
        </Button>
      </div>
    </div>
  );
}
