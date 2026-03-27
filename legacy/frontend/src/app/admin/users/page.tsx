"use client";

import { useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { Button, buttonVariants } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { apiFetch, type LocalUserRow } from "@/lib/api";
import { cn } from "@/lib/utils";
import { Trash2, KeyRound } from "lucide-react";

export default function AdminUsersPage() {
  const [list, setList] = useState<LocalUserRow[]>([]);
  const [msg, setMsg] = useState("");
  const [loading, setLoading] = useState(true);

  const [newUser, setNewUser] = useState({
    username: "",
    password: "",
    email: "",
    role: "user",
  });

  const [pwdRow, setPwdRow] = useState<{ id: number; v: string } | null>(null);

  const refresh = useCallback(async () => {
    setLoading(true);
    const r = await apiFetch("/api/v1/admin/users/local");
    if (r.ok) setList(await r.json());
    else setList([]);
    setLoading(false);
  }, []);

  useEffect(() => {
    refresh();
  }, [refresh]);

  async function createUser() {
    setMsg("");
    const res = await apiFetch("/api/v1/admin/users/local", {
      method: "POST",
      body: JSON.stringify({
        username: newUser.username.trim(),
        password: newUser.password,
        email: newUser.email.trim(),
        role: newUser.role,
      }),
    });
    if (res.ok) {
      setNewUser({ username: "", password: "", email: "", role: "user" });
      setMsg("已创建用户");
      await refresh();
    } else {
      const e = await res.json().catch(() => ({}));
      setMsg(e.error || "创建失败");
    }
  }

  async function saveRow(id: number) {
    const u = list.find((x) => x.id === id);
    if (!u) return;
    setMsg("");
    const res = await apiFetch(`/api/v1/admin/users/local/${u.id}`, {
      method: "PUT",
      body: JSON.stringify({
        email: u.email,
        role: u.role,
        is_active: u.is_active,
      }),
    });
    if (res.ok) {
      setMsg("已保存");
      await refresh();
    } else {
      const e = await res.json().catch(() => ({}));
      setMsg(e.error || "保存失败");
    }
  }

  async function resetPassword(id: number) {
    if (!pwdRow || pwdRow.id !== id || pwdRow.v.length < 6) {
      setMsg("请填写至少 6 位新密码");
      return;
    }
    setMsg("");
    const res = await apiFetch(`/api/v1/admin/users/local/${id}/password`, {
      method: "PUT",
      body: JSON.stringify({ password: pwdRow.v }),
    });
    if (res.ok) {
      setPwdRow(null);
      setMsg("密码已更新");
    } else {
      const e = await res.json().catch(() => ({}));
      setMsg(e.error || "重置失败");
    }
  }

  async function removeUser(id: number) {
    if (!confirm("确定删除该用户？")) return;
    setMsg("");
    const res = await apiFetch(`/api/v1/admin/users/local/${id}`, {
      method: "DELETE",
    });
    if (res.ok) {
      setMsg("已删除");
      await refresh();
    } else {
      const e = await res.json().catch(() => ({}));
      setMsg(e.error || "删除失败");
    }
  }

  return (
    <div className="space-y-8">
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <Link href="/admin" className="text-sm text-primary">
            ← 控制台
          </Link>
          <h1 className="mt-1 text-2xl font-bold text-primary">本地用户</h1>
          <p className="mt-1 text-sm text-muted-foreground">
            仅管理 auth_provider=local 的账号；LDAP/OIDC 用户不在此列表。
          </p>
        </div>
      </div>

      {msg ? (
        <p className="rounded-md border border-primary/30 bg-primary/5 px-3 py-2 text-sm">
          {msg}
        </p>
      ) : null}

      <section className="rounded-lg border bg-card p-4">
        <h2 className="mb-3 font-semibold">新建本地用户</h2>
        <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
          <div className="space-y-1.5">
            <Label>用户名</Label>
            <Input
              value={newUser.username}
              onChange={(e) =>
                setNewUser({ ...newUser, username: e.target.value })
              }
            />
          </div>
          <div className="space-y-1.5">
            <Label>初始密码（≥6 位）</Label>
            <Input
              type="password"
              autoComplete="new-password"
              value={newUser.password}
              onChange={(e) =>
                setNewUser({ ...newUser, password: e.target.value })
              }
            />
          </div>
          <div className="space-y-1.5">
            <Label>邮箱</Label>
            <Input
              placeholder="可空，默认 用户名@local"
              value={newUser.email}
              onChange={(e) =>
                setNewUser({ ...newUser, email: e.target.value })
              }
            />
          </div>
          <div className="space-y-1.5">
            <Label>角色</Label>
            <select
              className="flex h-10 w-full rounded-md border border-input bg-background px-3 text-sm"
              value={newUser.role}
              onChange={(e) =>
                setNewUser({ ...newUser, role: e.target.value })
              }
            >
              <option value="user">user</option>
              <option value="admin">admin</option>
            </select>
          </div>
        </div>
        <Button className="mt-3" type="button" onClick={createUser}>
          创建
        </Button>
      </section>

      <div className="overflow-x-auto rounded-lg border border-border bg-card">
        <table className="w-full min-w-[900px] border-collapse text-left text-sm">
          <thead>
            <tr className="border-b bg-muted/50">
              <th className="px-3 py-3 font-semibold">ID</th>
              <th className="px-3 py-3 font-semibold">用户名</th>
              <th className="px-3 py-3 font-semibold">邮箱</th>
              <th className="px-3 py-3 font-semibold">角色</th>
              <th className="px-3 py-3 font-semibold">启用</th>
              <th className="px-3 py-3 font-semibold">操作</th>
            </tr>
          </thead>
          <tbody>
            {loading ? (
              <tr>
                <td colSpan={6} className="px-3 py-8 text-center text-muted-foreground">
                  加载中…
                </td>
              </tr>
            ) : (
              list.map((u) => (
                <tr key={u.id} className="border-b border-border/80">
                  <td className="px-3 py-2 tabular-nums text-muted-foreground">
                    {u.id}
                  </td>
                  <td className="px-3 py-2 font-medium">{u.username}</td>
                  <td className="px-3 py-2">
                    <Input
                      className="h-9"
                      value={u.email}
                      onChange={(e) =>
                        setList((prev) =>
                          prev.map((x) =>
                            x.id === u.id ? { ...x, email: e.target.value } : x,
                          ),
                        )
                      }
                    />
                  </td>
                  <td className="px-3 py-2">
                    <select
                      className="flex h-9 w-full rounded-md border border-input bg-background px-2 text-sm"
                      value={u.role}
                      onChange={(e) =>
                        setList((prev) =>
                          prev.map((x) =>
                            x.id === u.id ? { ...x, role: e.target.value } : x,
                          ),
                        )
                      }
                    >
                      <option value="user">user</option>
                      <option value="admin">admin</option>
                    </select>
                  </td>
                  <td className="px-3 py-2">
                    <input
                      type="checkbox"
                      checked={u.is_active}
                      onChange={(e) =>
                        setList((prev) =>
                          prev.map((x) =>
                            x.id === u.id
                              ? { ...x, is_active: e.target.checked }
                              : x,
                          ),
                        )
                      }
                    />
                  </td>
                  <td className="px-3 py-2">
                    <div className="flex flex-wrap gap-1">
                      <Button
                        type="button"
                        size="sm"
                        variant="secondary"
                        onClick={() => saveRow(u.id)}
                      >
                        保存
                      </Button>
                      <Button
                        type="button"
                        size="sm"
                        variant="outline"
                        onClick={() =>
                          setPwdRow(
                            pwdRow?.id === u.id
                              ? null
                              : { id: u.id, v: "" },
                          )
                        }
                      >
                        <KeyRound className="mr-1 h-3.5 w-3.5" />
                        改密
                      </Button>
                      <button
                        type="button"
                        className={cn(
                          buttonVariants({
                            variant: "ghost",
                            size: "sm",
                          }),
                          "text-destructive",
                        )}
                        onClick={() => removeUser(u.id)}
                      >
                        <Trash2 className="h-4 w-4" />
                      </button>
                    </div>
                    {pwdRow?.id === u.id ? (
                      <div className="mt-2 flex flex-wrap items-end gap-2">
                        <div className="space-y-1">
                          <Label className="text-xs">新密码</Label>
                          <Input
                            type="password"
                            className="h-9 w-40"
                            value={pwdRow.v}
                            onChange={(e) =>
                              setPwdRow({ id: u.id, v: e.target.value })
                            }
                          />
                        </div>
                        <Button
                          type="button"
                          size="sm"
                          onClick={() => resetPassword(u.id)}
                        >
                          确认重置
                        </Button>
                      </div>
                    ) : null}
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
