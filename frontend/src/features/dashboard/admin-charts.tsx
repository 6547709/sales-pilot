"use client";

import {
  Bar,
  BarChart,
  LabelList,
  Pie,
  PieChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
  Cell,
} from "recharts";
import { Frame, FrameHeader, FramePanel, FrameTitle } from "@/components/reui/frame";
import type { AdminStats, LabelCount } from "@/lib/api";

const TOOLTIP_STYLE = {
  background: "var(--popover)",
  border: "1px solid var(--border)",
  borderRadius: 8,
  fontSize: 12,
  color: "var(--popover-foreground)",
};

function MarketChart({ data }: { data: LabelCount[] }) {
  if (data.length === 0) {
    return <p className="py-10 text-center text-xs text-muted-foreground">暂无产品数据</p>;
  }
  return (
    <div className="flex flex-col items-center">
      <div className="h-52 w-full">
        <ResponsiveContainer width="100%" height="100%">
          <PieChart>
            <Pie
              data={data}
              dataKey="count"
              nameKey="label"
              innerRadius="55%"
              outerRadius="80%"
              paddingAngle={2}
              stroke="var(--background)"
              strokeWidth={2}
              isAnimationActive={false}
            >
              {data.map((_, i) => (
                <Cell key={i} fill={`var(--chart-${(i % 5) + 1})`} />
              ))}
            </Pie>
            <Tooltip contentStyle={TOOLTIP_STYLE} />
          </PieChart>
        </ResponsiveContainer>
      </div>
      <div className="mt-1 flex flex-wrap justify-center gap-3 text-xs text-muted-foreground">
        {data.map((item, i) => (
          <span key={item.label} className="flex items-center gap-1.5">
            <span
              className="h-2 w-2 rounded-full"
              style={{ background: `var(--chart-${(i % 5) + 1})` }}
            />
            {item.label} {item.count}
          </span>
        ))}
      </div>
    </div>
  );
}

function CategoryChart({ data }: { data: LabelCount[] }) {
  const top = data.slice(0, 8);
  if (top.length === 0) {
    return <p className="py-10 text-center text-xs text-muted-foreground">暂无分类数据</p>;
  }
  return (
    <div className="h-64 w-full">
      <ResponsiveContainer width="100%" height="100%">
        <BarChart data={top} layout="vertical" margin={{ top: 4, right: 32, bottom: 4, left: 4 }}>
          <XAxis type="number" hide />
          <YAxis
            type="category"
            dataKey="label"
            width={92}
            axisLine={false}
            tickLine={false}
            tick={{ fill: "var(--muted-foreground)", fontSize: 11 }}
          />
          <Tooltip cursor={{ fill: "var(--muted)" }} contentStyle={TOOLTIP_STYLE} />
          <Bar dataKey="count" fill="var(--chart-1)" radius={[0, 4, 4, 0]} barSize={14}>
            <LabelList
              dataKey="count"
              position="right"
              fill="var(--muted-foreground)"
              fontSize={11}
            />
          </Bar>
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
}

export function AdminCharts({ stats }: { stats: AdminStats | null }) {
  const market = stats?.products.by_market ?? [];
  const categories = stats?.products.by_category ?? [];

  return (
    <div className="grid gap-4 lg:grid-cols-2">
      <Frame>
        <FramePanel>
          <FrameHeader>
            <FrameTitle>产品市场分布</FrameTitle>
          </FrameHeader>
          <div className="px-4 pb-4">
            <MarketChart data={market} />
          </div>
        </FramePanel>
      </Frame>
      <Frame>
        <FramePanel>
          <FrameHeader>
            <FrameTitle>产品分类分布（前 8）</FrameTitle>
          </FrameHeader>
          <div className="px-4 pb-4">
            <CategoryChart data={categories} />
          </div>
        </FramePanel>
      </Frame>
    </div>
  );
}
