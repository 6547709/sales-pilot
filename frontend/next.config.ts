import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "standalone",
  // 关闭开发环境左下角的 Next.js 悬浮标识
  devIndicators: false,
};

export default nextConfig;
