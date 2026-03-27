import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Sales-Pilot",
  description: "企业级销售赋能系统",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="zh-CN" className="h-full antialiased">
      <body className="min-h-full flex flex-col font-sans">{children}</body>
    </html>
  );
}
