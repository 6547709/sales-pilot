"use client";

import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";

type Props = {
  children: string;
  className?: string;
};

/** 统一 Markdown 渲染（含 GFM 表格/任务列表），样式见 globals.css `.markdown-body` */
export function MarkdownContent({ children, className }: Props) {
  return (
    <div className={className ?? "markdown-body"}>
      <ReactMarkdown remarkPlugins={[remarkGfm]}>{children || ""}</ReactMarkdown>
    </div>
  );
}
