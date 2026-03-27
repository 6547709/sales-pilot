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
          "flex w-full flex-col items-center gap-2 p-3 text-center outline-none transition-colors hover:bg-muted/20 focus-visible:ring-2 focus-visible:ring-ring",
          compact && "p-2",
          hasProductBar ? "rounded-t-xl" : "rounded-xl",
        )}
      >
        <span className="flex h-11 w-11 items-center justify-center rounded-lg bg-primary/10 text-primary shadow-inner group-hover:bg-primary/15">
          <Icon className="h-5 w-5" strokeWidth={1.75} />
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
      {hasProductBar ? (
        <div
          className={cn(
            "flex w-full flex-col gap-1.5 rounded-b-xl border-t border-border/60 bg-muted/50 px-2 py-2 dark:bg-muted/30",
            compact && "gap-1 py-1.5",
          )}
        >
          {productChips!.domestic.length > 0 ? (
            <div className="flex flex-wrap items-center justify-center gap-1">
              <span className="text-[9px] font-bold uppercase tracking-wide text-red-700 dark:text-red-400">
                国内
              </span>
              {productChips!.domestic.map((chip) => (
                <Link
                  key={chip.href}
                  href={chip.href}
                  onClick={(e) => e.stopPropagation()}
                  title={chip.label}
                  className={cn(
                    "max-w-[10rem] truncate rounded-md border border-red-200 bg-red-50 px-1.5 py-0.5 text-[9px] font-medium text-red-900 underline-offset-2 hover:underline dark:border-red-900/50 dark:bg-red-950/40 dark:text-red-200",
                    compact && "max-w-[7rem]",
                  )}
                >
                  {chip.label}
                </Link>
              ))}
            </div>
          ) : null}
          {productChips!.foreign.length > 0 ? (
            <div className="flex flex-wrap items-center justify-center gap-1">
              <span className="text-[9px] font-bold uppercase tracking-wide text-slate-600 dark:text-slate-400">
                国外
              </span>
              {productChips!.foreign.map((chip) => (
                <Link
                  key={chip.href}
                  href={chip.href}
                  onClick={(e) => e.stopPropagation()}
                  title={chip.label}
                  className={cn(
                    "max-w-[10rem] truncate rounded-md border border-slate-200 bg-slate-100 px-1.5 py-0.5 text-[9px] font-medium text-slate-800 underline-offset-2 hover:underline dark:border-slate-600 dark:bg-slate-800 dark:text-slate-200",
                    compact && "max-w-[7rem]",
                  )}
                >
                  {chip.label}
                </Link>
              ))}
            </div>
          ) : null}
          {productChips!.neutral.length > 0 ? (
            <div className="flex flex-wrap items-center justify-center gap-1">
              <span className="text-[9px] font-bold uppercase tracking-wide text-muted-foreground">
                不区分
              </span>
              {productChips!.neutral.map((chip) => (
                <Link
                  key={chip.href}
                  href={chip.href}
                  onClick={(e) => e.stopPropagation()}
                  title={chip.label}
                  className={cn(
                    "max-w-[10rem] truncate rounded-md border border-border bg-background px-1.5 py-0.5 text-[9px] font-medium text-foreground underline-offset-2 hover:underline",
                    compact && "max-w-[7rem]",
                  )}
                >
                  {chip.label}
                </Link>
              ))}
            </div>
          ) : null}
        </div>
      ) : null}
    </motion.div>
  );
}
