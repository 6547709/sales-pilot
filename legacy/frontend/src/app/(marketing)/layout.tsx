import { MarketingAuthLayout } from "@/components/marketing/marketing-auth-layout";

/** 确保首页查询参数（solution 等）参与渲染，避免赋能方案区不随卡片切换更新 */
export const dynamic = "force-dynamic";

export default function MarketingLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <MarketingAuthLayout>{children}</MarketingAuthLayout>;
}
