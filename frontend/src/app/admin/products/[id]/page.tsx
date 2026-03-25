"use client";

import { useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { useParams, useRouter } from "next/navigation";
import ReactMarkdown from "react-markdown";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Textarea } from "@/components/ui/textarea";
import {
  apiFetch,
  parseJsonArray,
  parsePersonas,
  type Product,
} from "@/lib/api";

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
  const [personasJson, setPersonasJson] = useState("{}");
  const [msg, setMsg] = useState("");

  const load = useCallback(async () => {
    const res = await apiFetch(`/api/v1/products/${id}`);
    if (!res.ok) {
      setProduct(null);
      return;
    }
    const p: Product = await res.json();
    setProduct(p);
    const hl = parseJsonArray(p.highlights);
    setH1(hl[0] || "");
    setH2(hl[1] || "");
    setH3(hl[2] || "");
    const qs = parseJsonArray(p.discovery_questions);
    setQ1(qs[0] || "");
    setQ2(qs[1] || "");
    setQ3(qs[2] || "");
    setPersonasJson(JSON.stringify(parsePersonas(p.target_personas), null, 2));
  }, [id]);

  useEffect(() => {
    load();
  }, [load]);

  async function save() {
    if (!product) return;
    setMsg("");
    let personas: Record<string, string>;
    try {
      personas = JSON.parse(personasJson) as Record<string, string>;
    } catch {
      setMsg("目标画像 JSON 无效");
      return;
    }
    const body = {
      ...product,
      highlights: [h1, h2, h3].filter(Boolean),
      discovery_questions: [q1, q2, q3].filter(Boolean),
      target_personas: personas,
    };
    const res = await apiFetch(`/api/v1/admin/products/${id}`, {
      method: "PUT",
      body: JSON.stringify(body),
    });
    if (res.ok) {
      setMsg("已保存");
      setProduct(await res.json());
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
          <Label>分类</Label>
          <Input
            value={product.category}
            onChange={(e) =>
              setProduct({ ...product, category: e.target.value })
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
            <div className="prose prose-sm max-w-none rounded-md border p-4 dark:prose-invert">
              <ReactMarkdown>{product.description}</ReactMarkdown>
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

      <div className="space-y-2">
        <Label>目标画像（JSON 对象）</Label>
        <Textarea
          rows={6}
          className="font-mono text-xs"
          value={personasJson}
          onChange={(e) => setPersonasJson(e.target.value)}
        />
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
            <div className="prose prose-sm max-w-none rounded-md border p-4 dark:prose-invert">
              <ReactMarkdown>{product.competitor_analysis}</ReactMarkdown>
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
