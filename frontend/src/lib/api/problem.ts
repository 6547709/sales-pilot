export type Problem = {
  error: string;
  code?: string;
};

export function isProblem(v: unknown): v is Problem {
  return (
    !!v &&
    typeof v === "object" &&
    typeof (v as { error?: unknown }).error === "string"
  );
}

export function problemMessage(v: unknown, fallback: string): string {
  return isProblem(v) ? v.error : fallback;
}
