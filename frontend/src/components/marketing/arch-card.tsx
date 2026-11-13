"use client";

import { motion } from "framer-motion";
import type { LucideIcon } from "lucide-react";
import Link from "next/link";
import { cn } from "@/lib/utils";

/** 卡片底部：按厂商聚合（标签为厂商名，点击进入方案库该厂商下全部产品） */
export type ArchCardManufacturerChip = {
  label: string;
  href: string;
};

export type ArchCardProductChips = {
  domestic: ArchCardManufacturerChip[];
  foreign: ArchCardManufacturerChip[];
  neutral: ArchCardManufacturerChip[];
};

type Props = {
  icon: LucideIcon;
  label: string;
  selected?: boolean;
  /** 点击卡片主体：选择/取消该解决方案 */
  onClick?: () => void;
  className?: string;
  compact?: boolean;
  hint?: string | null;
  /** 有已发布且关联本方案的产品时才传入；无产品则不展示底部区域 */
  productChips?: ArchCardProductChips;
  /** 本方案下的我方服务，展示在产品芯片之上 */
  services?: ArchCardManufacturerChip[];
  /** 产品超过 4 个时的「查看全部」链接 */
  moreHref?: string;
};

/** 架构图卡片：仅展示真实录入的产品；无关联产品时无底部内容 */
export function ArchCardButton({
  icon: Icon,
  label,
  selected,
  onClick,
  className,
  compact,
  hint,
  productChips,
  services,
  moreHref,
}: Props) {
  const rawHint = hint?.trim() ?? "";
  const legacyPlaceholders = new Set(["如 Dify 等", "如Dify等", "如 dify 等"]);
  const hintText =
    rawHint && !legacyPlaceholders.has(rawHint) ? rawHint : "";
  const hasProductBar =
    !!productChips &&
    (productChips.domestic.length > 0 ||
      productChips.foreign.length > 0 ||
      productChips.neutral.length > 0);
  const productEntries = productChips
    ? [
        ...productChips.domestic.map((chip) => ({ chip, market: "domestic" as const })),
        ...productChips.foreign.map((chip) => ({ chip, market: "foreign" as const })),
        ...productChips.neutral.map((chip) => ({ chip, market: "neutral" as const })),
      ]
    : [];
  const shownProducts = productEntries.slice(0, 4);
  const hiddenProducts = productEntries.length - shownProducts.length;
  const hasServiceBar = !!services && services.length > 0;
  const hasFooter = hasProductBar || hasServiceBar;

  return (
    <motion.div
      initial={false}
      className={cn(
        "group relative flex w-full flex-col gap-0 rounded-xl border bg-gradient-to-b from-card to-muted/30 text-center shadow-sm transition-colors",
        "border-border/80 hover:border-primary/40 hover:shadow-md",
        selected && "border-primary ring-2 ring-primary/25",
        className,
      )}
    >
      <button
        type="button"
        onClick={onClick}
        className={cn(
          "flex w-full flex-col items-center gap-1.5 p-2.5 text-center outline-none transition-colors hover:bg-muted/20 focus-visible:ring-2 focus-visible:ring-ring",
          compact && "p-1.5",
          hasFooter ? "rounded-t-xl" : "rounded-xl",
        )}
      >
        <span className="flex h-9 w-9 items-center justify-center rounded-lg bg-primary/10 text-primary shadow-inner group-hover:bg-primary/15">
          <Icon className="h-4 w-4" strokeWidth={1.75} />
        </span>
        <span className="text-xs font-semibold leading-tight text-foreground sm:text-sm">
          {label}
        </span>
        {hintText ? (
          <span className="line-clamp-2 text-[10px] text-muted-foreground">
            {hintText}
          </span>
        ) : null}
      </button>
      {hasFooter ? (
        <div
          className={cn(
            "flex w-full flex-col gap-1 rounded-b-xl border-t border-border/60 bg-muted/50 px-2 py-1.5 dark:bg-muted/30",
            compact && "gap-1 py-1",
          )}
        >
          {hasServiceBar ? (
            <div className="flex flex-wrap items-center justify-center gap-1">
              <span className="text-[9px] font-bold uppercase tracking-wide text-primary">
                服务
              </span>
              {services!.map((chip) => (
                <Link
                  key={chip.href}
                  href={chip.href}
                  onClick={(e) => e.stopPropagation()}
                  title={chip.label}
                  className={cn(
                    "max-w-[10rem] truncate rounded-md border border-primary/30 bg-primary/10 px-1.5 py-0.5 text-[9px] font-semibold text-primary underline-offset-2 hover:underline",
                    compact && "max-w-[7rem]",
                  )}
                >
                  {chip.label}
                </Link>
              ))}
            </div>
          ) : null}
          {hasProductBar ? (
            <div className="flex flex-wrap items-center justify-center gap-1">
              <span className="text-[9px] font-bold uppercase tracking-wide text-muted-foreground">
                产品
              </span>
              {shownProducts.map(({ chip, market }) => (
                <Link
                  key={chip.href}
                  href={chip.href}
                  onClick={(e) => e.stopPropagation()}
                  title={chip.label}
                  className={cn(
                    "max-w-[10rem] truncate rounded-md border px-1.5 py-0.5 text-[9px] font-medium underline-offset-2 hover:underline",
                    market === "domestic" &&
                      "border-red-200 bg-red-50 text-red-900 dark:border-red-900/50 dark:bg-red-950/40 dark:text-red-200",
                    market === "foreign" &&
                      "border-slate-200 bg-slate-100 text-slate-800 dark:border-slate-600 dark:bg-slate-800 dark:text-slate-200",
                    market === "neutral" && "border-border bg-background text-foreground",
                    compact && "max-w-[7rem]",
                  )}
                >
                  {chip.label}
                </Link>
              ))}
              {hiddenProducts > 0 ? (
                <Link
                  href={moreHref ?? "/products"}
                  onClick={(e) => e.stopPropagation()}
                  className="rounded-md border border-dashed border-border px-1.5 py-0.5 text-[9px] font-medium text-muted-foreground hover:text-foreground"
                >
                  +{hiddenProducts}
                </Link>
              ) : null}
            </div>
          ) : null}
        </div>
      ) : null}
    </motion.div>
  );
}
