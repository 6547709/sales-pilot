/** 后台编辑用的预设角色（与 JSON 存储互转） */
export const PERSONA_PRESET_ROLES = [
  "CEO",
  "CIO",
  "CTO",
  "CFO",
  "主任",
  "工程师",
  "技术总监",
  "业务负责人",
  "运维负责人",
  "安全负责人",
  "采购负责人",
] as const;

export type PersonaPresetRole = (typeof PERSONA_PRESET_ROLES)[number];

const PRESET_SET = new Set<string>(PERSONA_PRESET_ROLES);

export type PersonaExtraRow = { id: number; key: string; value: string };

/** 从存储的画像对象拆成「预设字段 + 自定义行」 */
export function splitPersonasForForm(raw: Record<string, string>): {
  preset: Record<string, string>;
  extra: PersonaExtraRow[];
} {
  const preset = Object.fromEntries(
    PERSONA_PRESET_ROLES.map((r) => [r, raw[r] ?? ""]),
  ) as Record<string, string>;
  const extra: PersonaExtraRow[] = [];
  let idx = 0;
  for (const [k, v] of Object.entries(raw)) {
    if (!PRESET_SET.has(k))
      extra.push({ id: idx++, key: k, value: v ?? "" });
  }
  return { preset, extra };
}

/** 合并表单为写入接口的对象（空值省略） */
export function mergePersonasFromForm(
  preset: Record<string, string>,
  extra: PersonaExtraRow[],
): Record<string, string> {
  const out: Record<string, string> = {};
  for (const r of PERSONA_PRESET_ROLES) {
    const t = (preset[r] ?? "").trim();
    if (t) out[r] = t;
  }
  for (const row of extra) {
    const k = row.key.trim();
    const v = row.value.trim();
    if (k && v) out[k] = v;
  }
  return out;
}
