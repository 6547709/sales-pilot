export type ViewMode = "sales" | "customer";

export function withView(path: string, mode: ViewMode): string {
  const sep = path.includes("?") ? "&" : "?";
  return `${path}${sep}view=${mode}`;
}
