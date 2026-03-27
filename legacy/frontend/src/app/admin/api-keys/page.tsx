"use client";

import { useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { apiFetch } from "@/lib/api";
import { KeyRound, Trash2 } from "lucide-react";

type ApiKeyRow = {
  id: number;
  user_id: number;
  name: string;
  prefix: string;
  expires_at?: string | null;
  last_used_at?: string | null;
  is_active: boolean;
  created_at: string;
};

export default function AdminApiKeysPage() {
  const [list, setList] = useState<ApiKeyRow[]>([]);
  const [msg, setMsg] = useState("");
  const [name, setName] = useState("");
  const [expiresAt, setExpiresAt] = useState("");
  const [onceKey, setOnceKey] = useState<string | null>(null);

  const refresh = useCallback(async () => {
    const r = await apiFetch("/api/v1/admin/api-keys");
    if (r.ok) setList(await r.json());
    else setList([]);
  }, []);

  useEffect(() => {
    refresh();
  }, [refresh]);

  async function createKey() {
    setMsg("");
    setOnceKey(null);
    if (!name.trim()) {
      setMsg("请填写名称");
      return;
    }
    const body: { name: string; expires_at?: string } = {
      name: name.trim(),
    };
    if (expiresAt.trim()) {
      const d = new Date(expiresAt);
      if (Number.isNaN(d.getTime())) {
        setMsg("过期时间格式无效");
        return;
      }
      body.expires_at = d.toISOString();
    }
    const res = await apiFetch("/api/v1/admin/api-keys", {
      method: "POST",
      body: JSON.stringify(body),
    });
    if (res.ok) {
      const data = (await res.json()) as { key?: string; warning?: string };
      setOnceKey(data.key ?? null);
      setMsg(data.warning || "已创建");
      setName("");
      setExpiresAt("");
      await refresh();
    } else {
      const e = await res.json().catch(() => ({}));
      setMsg(e.error || "创建失败");
    }
  }

  async function revoke(id: number) {
    if (!confirm("作废后该密钥立即失效，确定？")) return;
    setMsg("");
    const res = await apiFetch(`/api/v1/admin/api-keys/${id}`, {
      method: "DELETE",
    });
    if (res.ok) {
      await refresh();
      setMsg("已作废");
    } else setMsg("操作失败");
  }

  return (
    <div className="mx-auto max-w-3xl space-y-8">
      <div>
        <Link href="/admin" className="text-sm text-primary">
          ← 控制台
        </Link>
        <h1 className="mt-1 flex items-center gap-2 text-2xl font-bold text-primary">
          <KeyRound className="h-7 w-7" />
          API 密钥
        </h1>
        <p className="mt-2 text-sm text-muted-foreground">
          供 MCP、脚本等使用请求头{" "}
          <code className="rounded bg-muted px-1">X-API-Key</code>{" "}
          调用 <code className="rounded bg-muted px-1">/api/v1/admin/*</code>
          ，无需频繁登录。密钥明文<strong>仅在创建时显示一次</strong>。
        </p>
        <p className="mt-1 text-xs text-muted-foreground">
          JWT 登录令牌默认有效期为 24 小时（见后端{" "}
          <code className="rounded bg-muted px-1">auth.TokenTTL</code>
          ）；完整说明见仓库{" "}
          <code className="rounded bg-muted px-1">docs/api-ai-product-maintenance.md</code>
          。
        </p>
      </div>

      {onceKey ? (
        <div className="rounded-lg border-2 border-primary bg-primary/5 p-4">
          <p className="mb-2 text-sm font-semibold text-primary">
            请立即复制保存（不会再显示）：
          </p>
          <pre className="break-all rounded-md bg-background p-3 font-mono text-xs">
            {onceKey}
          </pre>
          <Button
            className="mt-2"
            type="button"
            size="sm"
            variant="secondary"
            onClick={() => void navigator.clipboard.writeText(onceKey)}
          >
            复制到剪贴板
          </Button>
        </div>
      ) : null}

      {msg ? (
        <p className="rounded-md border border-primary/30 bg-primary/5 px-3 py-2 text-sm">
          {msg}
        </p>
      ) : null}

      <section className="rounded-lg border bg-card p-4">
        <h2 className="mb-3 font-semibold">新建密钥</h2>
        <div className="grid gap-3 sm:grid-cols-2">
          <div className="space-y-1.5">
            <Label>名称（备注用途）</Label>
            <Input
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="例如 MCP-生产"
            />
          </div>
          <div className="space-y-1.5">
            <Label>过期时间（可选，留空=永久）</Label>
            <Input
              type="datetime-local"
              value={expiresAt}
              onChange={(e) => setExpiresAt(e.target.value)}
            />
          </div>
        </div>
        <Button className="mt-3" type="button" onClick={createKey}>
          生成密钥
        </Button>
      </section>

      <div className="overflow-x-auto rounded-lg border">
        <table className="w-full text-left text-sm">
          <thead className="border-b bg-muted/50">
            <tr>
              <th className="px-3 py-2">ID</th>
              <th className="px-3 py-2">名称</th>
              <th className="px-3 py-2">前缀</th>
              <th className="px-3 py-2">过期</th>
              <th className="px-3 py-2">最近使用</th>
              <th className="px-3 py-2">状态</th>
              <th className="px-3 py-2">操作</th>
            </tr>
          </thead>
          <tbody>
            {list.map((k) => (
              <tr key={k.id} className="border-b">
                <td className="px-3 py-2 text-muted-foreground">{k.id}</td>
                <td className="px-3 py-2">{k.name}</td>
                <td className="px-3 py-2 font-mono text-xs">{k.prefix}</td>
                <td className="px-3 py-2 text-muted-foreground">
                  {k.expires_at
                    ? new Date(k.expires_at).toLocaleString("zh-CN")
                    : "—"}
                </td>
                <td className="px-3 py-2 text-muted-foreground">
                  {k.last_used_at
                    ? new Date(k.last_used_at).toLocaleString("zh-CN")
                    : "—"}
                </td>
                <td className="px-3 py-2">
                  {k.is_active ? "有效" : "已作废"}
                </td>
                <td className="px-3 py-2">
                  {k.is_active ? (
                    <Button
                      type="button"
                      size="sm"
                      variant="destructive"
                      onClick={() => revoke(k.id)}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  ) : null}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {list.length === 0 ? (
          <p className="p-6 text-center text-sm text-muted-foreground">
            暂无密钥
          </p>
        ) : null}
      </div>
    </div>
  );
}
