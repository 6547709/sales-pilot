"use client";

import { useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { apiFetch, type Product } from "@/lib/api";

export default function ProductsPage() {
  const [q, setQ] = useState("");
  const [list, setList] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);

  const load = useCallback(async () => {
    setLoading(true);
    const path = q.trim()
      ? `/api/v1/products/search?q=${encodeURIComponent(q.trim())}`
      : "/api/v1/products";
    const res = await apiFetch(path);
    if (res.ok) setList(await res.json());
    setLoading(false);
  }, [q]);

  useEffect(() => {
    load();
  }, [load]);

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-primary">产品方案</h1>
        <p className="text-sm text-muted-foreground">
          搜索赋能内容，点击进入详情与话术。
        </p>
      </div>
      <div className="flex gap-2">
        <Input
          placeholder="搜索名称、分类、描述…"
          value={q}
          onChange={(e) => setQ(e.target.value)}
          className="max-w-md"
        />
        <Button type="button" variant="secondary" onClick={load}>
          搜索
        </Button>
      </div>
      {loading ? (
        <p className="text-muted-foreground">加载中…</p>
      ) : (
        <ul className="grid gap-4 sm:grid-cols-2">
          {list.map((p) => (
            <li key={p.id}>
              <Link href={`/products/${p.id}`}>
                <Card className="h-full transition hover:border-primary/40">
                  <CardHeader>
                    <CardTitle className="text-lg">{p.name}</CardTitle>
                    <CardDescription>{p.category}</CardDescription>
                  </CardHeader>
                </Card>
              </Link>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
