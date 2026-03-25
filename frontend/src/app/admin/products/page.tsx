"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { apiFetch, type Product } from "@/lib/api";

export default function AdminProductsPage() {
  const router = useRouter();
  const [list, setList] = useState<Product[]>([]);

  useEffect(() => {
    apiFetch("/api/v1/products").then(async (r) => {
      if (r.ok) setList(await r.json());
    });
  }, []);

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
      }),
    });
    if (res.ok) {
      const p = await res.json();
      router.push(`/admin/products/${p.id}`);
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-4">
        <h1 className="text-2xl font-bold text-primary">产品维护</h1>
        <Button onClick={createDraft}>新建产品</Button>
      </div>
      <ul className="grid gap-4 sm:grid-cols-2">
        {list.map((p) => (
          <li key={p.id}>
            <Link href={`/admin/products/${p.id}`}>
              <Card className="h-full transition hover:border-primary/40">
                <CardHeader>
                  <CardTitle className="text-lg">{p.name}</CardTitle>
                  <CardDescription>{p.category || "未分类"}</CardDescription>
                </CardHeader>
              </Card>
            </Link>
          </li>
        ))}
      </ul>
    </div>
  );
}
