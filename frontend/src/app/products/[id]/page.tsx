"use client";

import { useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { useParams } from "next/navigation";
import { Button, buttonVariants } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import {
  apiFetch,
  parseJsonArray,
  parsePersonas,
  type Case,
  type Product,
  type SalesScript,
} from "@/lib/api";

const sections = [
  { id: "highlights", label: "三大亮点" },
  { id: "personas", label: "目标画像" },
  { id: "triggers", label: "触发事件" },
  { id: "questions", label: "黄金三问" },
  { id: "competitor", label: "竞品分析" },
  { id: "roi", label: "ROI 指标" },
  { id: "scripts", label: "销售话术" },
  { id: "cases", label: "客户案例" },
] as const;

export default function ProductDetailPage() {
  const params = useParams();
  const id = params.id as string;
  const [product, setProduct] = useState<Product | null>(null);
  const [scripts, setScripts] = useState<SalesScript[]>([]);
  const [cases, setCases] = useState<Case[]>([]);

  const load = useCallback(async () => {
    const [pr, sr, cr] = await Promise.all([
      apiFetch(`/api/v1/products/${id}`),
      apiFetch(`/api/v1/products/${id}/scripts`),
      apiFetch(`/api/v1/products/${id}/cases`),
    ]);
    if (pr.ok) setProduct(await pr.json());
    if (sr.ok) setScripts(await sr.json());
    if (cr.ok) setCases(await cr.json());
  }, [id]);

  useEffect(() => {
    load();
  }, [load]);

  if (!product) {
    return <p className="text-muted-foreground">加载中…</p>;
  }

  const prod = product;
  const highlights = parseJsonArray(prod.highlights);
  const questions = parseJsonArray(prod.discovery_questions);
  const personas = parsePersonas(prod.target_personas);
  const scriptText = scripts.map((s) => `【${s.scenario}】\n${s.content}`).join("\n\n");

  function copyScripts() {
    const text = scriptText || prod.description || "";
    void navigator.clipboard.writeText(text);
  }

  return (
    <div className="relative pb-24 lg:pb-8">
      <div className="mb-4">
        <Link
          href="/products"
          className="text-sm text-primary hover:underline"
        >
          ← 返回列表
        </Link>
      </div>

      <div className="lg:grid lg:grid-cols-[200px_1fr] lg:gap-10">
        <aside className="mb-6 hidden lg:block">
          <nav className="sticky top-24 space-y-1 text-sm">
            {sections.map((s) => (
              <a
                key={s.id}
                href={`#${s.id}`}
                className="block rounded-md px-2 py-1.5 text-muted-foreground hover:bg-muted hover:text-foreground"
              >
                {s.label}
              </a>
            ))}
          </nav>
        </aside>

        <article className="min-w-0 space-y-10">
          <header>
            <h1 className="text-2xl font-bold text-primary md:text-3xl">
              {prod.name}
            </h1>
            <p className="mt-1 text-sm text-muted-foreground">
              {prod.category}
            </p>
            <p className="mt-4 whitespace-pre-wrap text-foreground/90">
              {prod.description}
            </p>
          </header>

          <section id="highlights">
            <h2 className="mb-3 text-lg font-semibold">三大亮点</h2>
            <ul className="list-inside list-disc space-y-2 text-foreground/90">
              {highlights.map((h, i) => (
                <li key={i}>{h}</li>
              ))}
            </ul>
          </section>

          <section id="personas">
            <h2 className="mb-3 text-lg font-semibold">目标画像</h2>
            <ul className="space-y-2">
              {Object.entries(personas).map(([k, v]) => (
                <li key={k}>
                  <span className="font-medium text-primary">{k}</span>：{v}
                </li>
              ))}
            </ul>
          </section>

          <section id="triggers">
            <h2 className="mb-3 text-lg font-semibold">触发事件</h2>
            <p className="whitespace-pre-wrap text-foreground/90">
              {prod.trigger_events}
            </p>
          </section>

          <section id="questions">
            <h2 className="mb-3 text-lg font-semibold">黄金三问</h2>
            <ol className="list-decimal space-y-3 pl-5 text-foreground/90">
              {questions.map((q, i) => (
                <li key={i}>{q}</li>
              ))}
            </ol>
          </section>

          <section id="competitor">
            <h2 className="mb-3 text-lg font-semibold">竞品分析</h2>
            <ScrollArea className="max-h-[480px] rounded-md border p-4">
              <div className="prose prose-sm max-w-none dark:prose-invert">
                <pre className="whitespace-pre-wrap font-sans text-sm">
                  {prod.competitor_analysis}
                </pre>
              </div>
            </ScrollArea>
          </section>

          <section id="roi">
            <h2 className="mb-3 text-lg font-semibold">ROI 指标</h2>
            <p className="whitespace-pre-wrap text-foreground/90">
              {prod.roi_metrics}
            </p>
          </section>

          <section id="scripts">
            <h2 className="mb-3 text-lg font-semibold">销售话术</h2>
            <div className="space-y-4">
              {scripts.map((s) => (
                <div key={s.id} className="rounded-lg border p-4">
                  <p className="font-medium text-primary">{s.scenario}</p>
                  <Separator className="my-2" />
                  <p className="whitespace-pre-wrap text-sm">{s.content}</p>
                </div>
              ))}
            </div>
          </section>

          <section id="cases">
            <h2 className="mb-3 text-lg font-semibold">客户案例</h2>
            <div className="space-y-4">
              {cases.map((c) => (
                <div key={c.id} className="rounded-lg border p-4">
                  <p className="font-medium">{c.client_name}</p>
                  <p className="mt-2 text-sm text-muted-foreground">
                    痛点：{c.pain_points}
                  </p>
                  <p className="mt-1 text-sm">方案：{c.solution}</p>
                  <p className="mt-1 text-sm">价值：{c.value_delivered}</p>
                </div>
              ))}
            </div>
          </section>
        </article>
      </div>

      {/* 移动端底部悬浮条 */}
      <div className="fixed inset-x-0 bottom-0 z-50 border-t bg-card/95 p-3 backdrop-blur lg:hidden">
        <div className="mx-auto flex max-w-lg gap-2">
          <Button className="flex-1" onClick={copyScripts}>
            复制话术
          </Button>
          <a
            href={
              process.env.NEXT_PUBLIC_DIAL_NUMBER
                ? `tel:${process.env.NEXT_PUBLIC_DIAL_NUMBER}`
                : "#"
            }
            onClick={(e) => {
              if (!process.env.NEXT_PUBLIC_DIAL_NUMBER) {
                e.preventDefault();
                alert(
                  "可在环境变量 NEXT_PUBLIC_DIAL_NUMBER 中配置热线号码。",
                );
              }
            }}
            className={cn(
              buttonVariants({ variant: "secondary" }),
              "inline-flex flex-1 items-center justify-center",
            )}
          >
            一键拨号
          </a>
        </div>
      </div>
    </div>
  );
}
