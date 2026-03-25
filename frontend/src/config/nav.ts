/** 动态导航配置：后续可扩展「招投标工具」等模块 */
export type NavItem = {
  href: string;
  label: string;
  adminOnly?: boolean;
};

export const mainNav: NavItem[] = [
  { href: "/products", label: "产品方案" },
  { href: "/admin", label: "管理后台", adminOnly: true },
  { href: "/admin/products", label: "产品维护", adminOnly: true },
  { href: "/admin/settings", label: "认证配置", adminOnly: true },
];
