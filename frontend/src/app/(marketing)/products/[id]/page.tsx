"use client";

import { motion } from "framer-motion";
import Link from "next/link";
import { useCallback, useEffect, useRef, useState } from "react";
import { useParams } from "next/navigation";
import { MarkdownContent } from "@/components/ui/markdown-content";
import { Badge } from "@/components/ui/badge";
import { Button, buttonVariants } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { cn } from "@/lib/utils";
import {
  apiFetch,
  parseJsonArray,
  parsePersonas,
  type Case,
  type Product,
  type SalesScript,
} from "@/lib/api";
import {
  vendorMarketBadgeClass,
  vendorMarketBadgeText,
} from "@/lib/vendor-market";
import {
  BookOpen,
  Building2,
  ClipboardList,
  Copy,
  Lightbulb,
  Mail,
  MessageSquare,
  Phone,
  Sparkles,
  Target,
  TrendingUp,
  UserCircle,
  Users,
  Zap,
} from "lucide-react";

/** 生成 tel: 链接用：去掉空格与常见连字符 */
function normalizeTelHref(phone: string) {
  return phone.trim().replace(/[\s\u00a0-]/g, "");
}

function ContactPhoneLink({ value }: { value?: string | null }) {
  const display = value?.trim();
  if (!display) return <span className="text-muted-foreground">—</span>;
  return (
    <a
      href={`tel:${normalizeTelHref(display)}`}
      className="text-primary underline-offset-2 hover:underline"
    >
      {display}
    </a>
  );
}

function ContactEmailLink({ value }: { value?: string | null }) {
  const display = value?.trim();
  if (!display) return <span className="text-muted-foreground">—</span>;
  return (
    <a href={`mailto:${display}`} className="break-all text-primary underline-offset-2 hover:underline">
      {display}
    </a>
  );
}

/** 内联复制：icon = 单条旁小按钮；text = 「复制全部」类按钮 */
function CopyTextControl({
  text,
  mode,
  label = "复制全部",
  className,
}: {
  text: string;
  mode: "icon" | "text";
  label?: string;
  className?: string;
}) {
  const [done, setDone] = useState(false);
  const handle = () => {
    if (!text.trim()) return;
    void navigator.clipboard.writeText(text).then(() => {
      setDone(true);
      window.setTimeout(() => setDone(false), 2000);
    });
  };
  if (!text.trim()) return null;
  if (mode === "icon") {
    return (
      <Button
        type="button"
        variant="ghost"
        size="icon-sm"
        className={cn("shrink-0 text-muted-foreground hover:text-foreground", className)}
        title={done ? "已复制" : "复制"}
        aria-label={done ? "已复制" : "复制"}
        onClick={handle}
      >
        <Copy className="size-3.5" />
      </Button>
    );
  }
  return (
    <Button
      type="button"
      variant="outline"
      size="xs"
      className={cn("gap-1", className)}
      onClick={handle}
    >
      <Copy className="size-3" />
      {done ? "已复制" : label}
    </Button>
  );
}

const sectionNav = [
  { id: "overview", label: "概览" },
  { id: "contacts", label: "联络" },
  { id: "highlights", label: "亮点" },
  { id: "questions", label: "黄金三问" },
  { id: "personas", label: "画像" },
  { id: "triggers", label: "触发" },
  { id: "competitor", label: "竞品" },
  { id: "roi", label: "ROI" },
  { id: "scripts", label: "话术" },
  { id: "cases", label: "案例" },
] as const;

export default function ProductDetailPage() {
  const params = useParams();
  const id = params.id as string;
  const [product, setProduct] = useState<Product | null>(null);
  const [scripts, setScripts] = useState<SalesScript[]>([]);
  const [cases, setCases] = useState<Case[]>([]);
  const [activeSection, setActiveSection] = useState<string>("overview");
  const observerRef = useRef<IntersectionObserver | null>(null);

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

  useEffect(() => {
    const sections = sectionNav.map((s) => document.getElementById(s.id)).filter(Boolean);
    observerRef.current?.disconnect();
    observerRef.current = new IntersectionObserver(
      (entries) => {
        const visible = entries
          .filter((e) => e.isIntersecting)
          .sort((a, b) => b.intersectionRatio - a.intersectionRatio);
        if (visible[0]?.target.id) setActiveSection(visible[0].target.id);
      },
      { rootMargin: "-20% 0px -55% 0px", threshold: [0, 0.2, 0.5] },
    );
    sections.forEach((el) => el && observerRef.current?.observe(el));
    return () => observerRef.current?.disconnect();
  }, [product]);

  if (!product) {
    return (
      <div className="flex min-h-[50vh] items-center justify-center text-muted-foreground">
        加载方案…
      </div>
    );
  }

  const prod = product;
  const highlights = parseJsonArray(prod.highlights);
  const questions = parseJsonArray(prod.discovery_questions);
  const personas = parsePersonas(prod.target_personas);
  const highlightsJoined = highlights.join("\n\n");
  const questionsJoined = questions.join("\n\n");
  const scriptsJoined = scripts
    .map((s) => `【${s.scenario}】\n${s.content}`)
    .join("\n\n");
  const personasJoined = Object.entries(personas)
    .map(([k, v]) => `${k}\n${v}`)
    .join("\n\n");

  function scrollToSection(sid: string) {
    document.getElementById(sid)?.scrollIntoView({ behavior: "smooth", block: "start" });
  }

  function formatCaseText(c: Case) {
    return [
      c.client_name,
      `痛点：${c.pain_points}`,
      `方案：${c.solution}`,
      `价值：${c.value_delivered}`,
    ].join("\n");
  }

  const casesJoined = cases.map((c) => formatCaseText(c)).join("\n\n---\n\n");

  const cardMotion = {
    initial: { opacity: 0, y: 12 },
    animate: { opacity: 1, y: 0 },
  };

  return (
    <div className="relative pb-10">
      <div className="mx-auto max-w-[1200px] px-4 py-6 sm:py-8">
        <div className="mb-6 flex flex-wrap items-center gap-3">
          <Link
            href="/products"
            className={cn(
              buttonVariants({ variant: "ghost", size: "sm" }),
              "text-muted-foreground",
            )}
          >
            ← 方案库
          </Link>
          <Link
            href="/"
            className={cn(
              buttonVariants({ variant: "ghost", size: "sm" }),
              "text-muted-foreground",
            )}
          >
            架构总览
          </Link>
        </div>

        {/* 移动端：横向章节导航 */}
        <div className="sticky top-14 z-30 -mx-4 mb-6 border-y bg-background/95 py-2 backdrop-blur lg:hidden">
          <div className="flex gap-1 overflow-x-auto px-4 pb-1">
            {sectionNav.map((s) => (
              <button
                key={s.id}
                type="button"
                onClick={() => scrollToSection(s.id)}
                className={cn(
                  "shrink-0 rounded-full px-3 py-1.5 text-xs font-medium transition-colors",
                  activeSection === s.id
                    ? "bg-primary text-primary-foreground"
                    : "bg-muted text-muted-foreground",
                )}
              >
                {s.label}
              </button>
            ))}
          </div>
        </div>

        <div className="lg:grid lg:grid-cols-[200px_1fr] lg:gap-10">
          <aside className="mb-8 hidden lg:block">
            <nav className="sticky top-24 space-y-0.5 text-sm">
              {sectionNav.map((s) => (
                <button
                  key={s.id}
                  type="button"
                  onClick={() => scrollToSection(s.id)}
                  className={cn(
                    "flex w-full rounded-lg px-3 py-2 text-left transition-colors",
                    activeSection === s.id
                      ? "bg-primary/10 font-medium text-primary"
                      : "text-muted-foreground hover:bg-muted",
                  )}
                >
                  {s.label}
                </button>
              ))}
            </nav>
          </aside>

          <div className="grid auto-rows-auto gap-4 sm:gap-5">
            {/* 概览 */}
            <motion.section id="overview" {...cardMotion}>
              <Card className="overflow-hidden border-primary/20 bg-gradient-to-br from-primary/[0.06] to-transparent">
                <CardHeader className="pb-2">
                  <div className="flex flex-wrap items-center gap-2">
                    <Badge>{prod.category || "方案"}</Badge>
                    {vendorMarketBadgeText(prod.vendor_market) ? (
                      <Badge
                        variant="outline"
                        className={`font-medium ${vendorMarketBadgeClass(prod.vendor_market)}`}
                      >
                        {vendorMarketBadgeText(prod.vendor_market)}
                      </Badge>
                    ) : null}
                    <Badge variant="outline" className="gap-1">
                      <BookOpen className="h-3 w-3" />
                      赋能包
                    </Badge>
                  </div>
                  <CardTitle className="text-2xl leading-tight sm:text-3xl">
                    {prod.name}
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <MarkdownContent className="markdown-body text-sm text-foreground/90">
                    {prod.description}
                  </MarkdownContent>
                </CardContent>
              </Card>
            </motion.section>

            <motion.section id="contacts" {...cardMotion}>
              <Card>
                <CardHeader className="pb-2">
                  <CardTitle className="flex items-center gap-2 text-lg">
                    <Building2 className="h-5 w-5 text-primary" />
                    厂商与联络方式
                  </CardTitle>
                </CardHeader>
                <CardContent className="grid gap-4 sm:grid-cols-2">
                  <div className="rounded-xl border bg-muted/20 p-4">
                    <p className="mb-2 text-xs font-semibold uppercase tracking-wide text-muted-foreground">
                      厂商
                    </p>
                    <p className="text-sm font-medium">
                      {prod.manufacturer_name?.trim() || "—"}
                    </p>
                  </div>
                  <div className="rounded-xl border bg-muted/20 p-4">
                    <p className="mb-3 text-xs font-semibold uppercase tracking-wide text-muted-foreground">
                      销售
                    </p>
                    <ul className="space-y-2 text-sm">
                      <li className="flex items-start gap-2">
                        <UserCircle className="mt-0.5 h-4 w-4 shrink-0 text-primary" />
                        <span>
                          {prod.sales_contact_name?.trim() || "—"}
                        </span>
                      </li>
                      <li className="flex items-start gap-2">
                        <Phone className="mt-0.5 h-4 w-4 shrink-0 text-primary" />
                        <ContactPhoneLink value={prod.sales_contact_phone} />
                      </li>
                      <li className="flex items-start gap-2">
                        <Mail className="mt-0.5 h-4 w-4 shrink-0 text-primary" />
                        <ContactEmailLink value={prod.sales_contact_email} />
                      </li>
                    </ul>
                  </div>
                  <div className="rounded-xl border bg-muted/20 p-4 sm:col-span-2">
                    <p className="mb-3 text-xs font-semibold uppercase tracking-wide text-muted-foreground">
                      售前
                    </p>
                    <ul className="grid gap-3 sm:grid-cols-3">
                      <li className="flex items-start gap-2 text-sm">
                        <UserCircle className="mt-0.5 h-4 w-4 shrink-0 text-primary" />
                        <span>
                          {prod.presales_contact_name?.trim() || "—"}
                        </span>
                      </li>
                      <li className="flex items-start gap-2 text-sm">
                        <Phone className="mt-0.5 h-4 w-4 shrink-0 text-primary" />
                        <ContactPhoneLink value={prod.presales_contact_phone} />
                      </li>
                      <li className="flex items-start gap-2 text-sm">
                        <Mail className="mt-0.5 h-4 w-4 shrink-0 text-primary" />
                        <ContactEmailLink value={prod.presales_contact_email} />
                      </li>
                    </ul>
                  </div>
                </CardContent>
              </Card>
            </motion.section>

            <div className="grid gap-4 md:grid-cols-2">
              {/* 三大亮点 — 三张小卡 */}
              <motion.section id="highlights" {...cardMotion} className="md:col-span-1">
                <Card className="h-full">
                  <CardHeader className="pb-2">
                    <div className="flex flex-wrap items-start justify-between gap-2">
                      <CardTitle className="flex items-center gap-2 text-lg">
                        <Sparkles className="h-5 w-5 text-primary" />
                        三大亮点
                      </CardTitle>
                      <CopyTextControl text={highlightsJoined} mode="text" label="复制三条" />
                    </div>
                  </CardHeader>
                  <CardContent className="grid gap-2">
                    {highlights.map((h, i) => (
                      <div
                        key={i}
                        className="flex gap-1 rounded-lg border border-border/80 bg-muted/30 p-3 text-sm leading-relaxed"
                      >
                        <div className="min-w-0 flex-1">
                          <span className="font-semibold text-primary">{i + 1}. </span>
                          {h}
                        </div>
                        <CopyTextControl text={h} mode="icon" />
                      </div>
                    ))}
                  </CardContent>
                </Card>
              </motion.section>

              {/* 黄金三问 */}
              <motion.section id="questions" {...cardMotion} className="md:col-span-1">
                <Card className="h-full">
                  <CardHeader className="pb-2">
                    <div className="flex flex-wrap items-start justify-between gap-2">
                      <CardTitle className="flex items-center gap-2 text-lg">
                        <Lightbulb className="h-5 w-5 text-primary" />
                        黄金三问
                      </CardTitle>
                      <CopyTextControl text={questionsJoined} mode="text" label="复制三问" />
                    </div>
                  </CardHeader>
                  <CardContent className="space-y-3">
                    {questions.map((q, i) => (
                      <div
                        key={i}
                        className="flex gap-2 rounded-lg border p-3 text-sm"
                      >
                        <span className="flex h-7 w-7 shrink-0 items-center justify-center rounded-full bg-primary/10 text-xs font-bold text-primary">
                          {i + 1}
                        </span>
                        <p className="min-w-0 flex-1 leading-relaxed text-foreground/90">{q}</p>
                        <CopyTextControl text={q} mode="icon" />
                      </div>
                    ))}
                  </CardContent>
                </Card>
              </motion.section>
            </div>

            <div className="grid gap-4 md:grid-cols-2">
              <motion.section id="personas" {...cardMotion}>
                <Card className="h-full">
                  <CardHeader className="pb-2">
                    <div className="flex flex-wrap items-start justify-between gap-2">
                      <CardTitle className="flex items-center gap-2 text-lg">
                        <Users className="h-5 w-5 text-primary" />
                        目标画像
                      </CardTitle>
                      <CopyTextControl text={personasJoined} mode="text" label="复制全部画像" />
                    </div>
                  </CardHeader>
                  <CardContent className="space-y-2">
                    {Object.entries(personas).map(([k, v]) => {
                      const block = `${k}\n${v}`;
                      return (
                        <div
                          key={k}
                          className="flex gap-1 rounded-lg border border-dashed p-3 text-sm"
                        >
                          <div className="min-w-0 flex-1">
                            <span className="font-semibold text-primary">{k}</span>
                            <Separator className="my-2" />
                            <p className="text-muted-foreground">{v}</p>
                          </div>
                          <CopyTextControl text={block} mode="icon" />
                        </div>
                      );
                    })}
                  </CardContent>
                </Card>
              </motion.section>

              <motion.section id="triggers" {...cardMotion}>
                <Card className="h-full">
                  <CardHeader className="pb-2">
                    <div className="flex flex-wrap items-start justify-between gap-2">
                      <CardTitle className="flex items-center gap-2 text-lg">
                        <Zap className="h-5 w-5 text-primary" />
                        触发事件
                      </CardTitle>
                      <CopyTextControl
                        text={prod.trigger_events || ""}
                        mode="text"
                        label="复制本段"
                      />
                    </div>
                  </CardHeader>
                  <CardContent>
                    <p className="whitespace-pre-wrap text-sm leading-relaxed text-foreground/90">
                      {prod.trigger_events}
                    </p>
                  </CardContent>
                </Card>
              </motion.section>
            </div>

            <motion.section id="competitor" {...cardMotion}>
              <Card>
                <CardHeader className="pb-2">
                  <div className="flex flex-wrap items-start justify-between gap-2">
                    <CardTitle className="flex items-center gap-2 text-lg">
                      <Target className="h-5 w-5 text-primary" />
                      竞品分析
                    </CardTitle>
                    <CopyTextControl
                      text={prod.competitor_analysis || ""}
                      mode="text"
                      label="复制原文"
                    />
                  </div>
                </CardHeader>
                <CardContent>
                  <div className="rounded-xl border bg-muted/20 p-4">
                    <MarkdownContent className="markdown-body text-sm">
                      {prod.competitor_analysis}
                    </MarkdownContent>
                  </div>
                </CardContent>
              </Card>
            </motion.section>

            <motion.section id="roi" {...cardMotion}>
              <Card>
                <CardHeader className="pb-2">
                  <div className="flex flex-wrap items-start justify-between gap-2">
                    <CardTitle className="flex items-center gap-2 text-lg">
                      <TrendingUp className="h-5 w-5 text-primary" />
                      ROI 指标
                    </CardTitle>
                    <CopyTextControl
                      text={prod.roi_metrics || ""}
                      mode="text"
                      label="复制本段"
                    />
                  </div>
                </CardHeader>
                <CardContent>
                  <p className="whitespace-pre-wrap text-sm leading-relaxed">
                    {prod.roi_metrics}
                  </p>
                </CardContent>
              </Card>
            </motion.section>

            <motion.section id="scripts" {...cardMotion}>
              <Card>
                <CardHeader className="pb-2">
                  <div className="flex flex-wrap items-start justify-between gap-2">
                    <CardTitle className="flex items-center gap-2 text-lg">
                      <MessageSquare className="h-5 w-5 text-primary" />
                      销售话术
                    </CardTitle>
                    <CopyTextControl text={scriptsJoined} mode="text" label="复制全部话术" />
                  </div>
                </CardHeader>
                <CardContent>
                  <ul className="grid gap-3 sm:grid-cols-2">
                    {scripts.map((s) => (
                      <li
                        key={s.id}
                        className="rounded-xl border bg-card p-4 shadow-sm"
                      >
                        <div className="flex items-start justify-between gap-2">
                          <p className="min-w-0 flex-1 text-sm font-semibold text-primary">
                            {s.scenario}
                          </p>
                          <CopyTextControl
                            text={`【${s.scenario}】\n${s.content}`}
                            mode="icon"
                          />
                        </div>
                        <Separator className="my-2" />
                        <p className="text-sm leading-relaxed text-foreground/90">
                          {s.content}
                        </p>
                      </li>
                    ))}
                  </ul>
                </CardContent>
              </Card>
            </motion.section>

            <motion.section id="cases" {...cardMotion}>
              <Card>
                <CardHeader className="pb-2">
                  <div className="flex flex-wrap items-start justify-between gap-2">
                    <CardTitle className="flex items-center gap-2 text-lg">
                      <ClipboardList className="h-5 w-5 text-primary" />
                      客户案例
                    </CardTitle>
                    <CopyTextControl text={casesJoined} mode="text" label="复制全部案例" />
                  </div>
                </CardHeader>
                <CardContent>
                  <ul className="grid gap-4 md:grid-cols-2">
                    {cases.map((c) => (
                      <li
                        key={c.id}
                        className="rounded-xl border border-primary/10 bg-gradient-to-br from-muted/40 to-background p-4"
                      >
                        <div className="flex items-start justify-between gap-2">
                          <p className="min-w-0 flex-1 font-semibold">{c.client_name}</p>
                          <CopyTextControl text={formatCaseText(c)} mode="icon" />
                        </div>
                        <p className="mt-2 text-xs text-muted-foreground">
                          痛点
                        </p>
                        <p className="text-sm">{c.pain_points}</p>
                        <p className="mt-2 text-xs text-muted-foreground">方案</p>
                        <p className="text-sm">{c.solution}</p>
                        <p className="mt-2 text-xs text-muted-foreground">价值</p>
                        <p className="text-sm">{c.value_delivered}</p>
                      </li>
                    ))}
                  </ul>
                </CardContent>
              </Card>
            </motion.section>
          </div>
        </div>
      </div>
    </div>
  );
}
