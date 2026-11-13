"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import {
  useTable,
  type ColumnDef,
  type PaginationState,
} from "@tanstack/react-table";
import { KeyRound, Pencil, Trash2, Users } from "lucide-react";
import {
  DataGrid,
  DataGridContainer,
  dataGridFeatures,
  type DataGridFeatures,
} from "@/components/reui/data-grid/data-grid";
import { DataGridTable } from "@/components/reui/data-grid/data-grid-table";
import { DataGridPagination } from "@/components/reui/data-grid/data-grid-pagination";
import { Badge } from "@/components/reui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { ConfirmDialog } from "@/components/patterns/confirm-dialog";
import { ErrorState } from "@/components/patterns/error-state";
import { FormField } from "@/components/patterns/form-field";
import { usersApi, type LocalUserRow } from "@/lib/api";

function formatTime(iso: string) {
  try {
    return new Date(iso).toLocaleString("zh-CN", {
      dateStyle: "short",
      timeStyle: "short",
    });
  } catch {
    return iso;
  }
}

export function UsersDataGrid() {
  const [data, setData] = useState<LocalUserRow[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [pagination, setPagination] = useState<PaginationState>({
    pageIndex: 0,
    pageSize: 20,
  });
  const [q, setQ] = useState("");
  const [role, setRole] = useState("all");
  const [active, setActive] = useState("all");

  const [createOpen, setCreateOpen] = useState(false);
  const [newUser, setNewUser] = useState({
    username: "",
    password: "",
    email: "",
    role: "user",
  });
  const [editing, setEditing] = useState<LocalUserRow | null>(null);
  const [editEmail, setEditEmail] = useState("");
  const [editRole, setEditRole] = useState("user");
  const [pwdTarget, setPwdTarget] = useState<LocalUserRow | null>(null);
  const [pwd, setPwd] = useState("");
  const [deleteTarget, setDeleteTarget] = useState<LocalUserRow | null>(null);
  const [toggleTarget, setToggleTarget] = useState<LocalUserRow | null>(null);
  const [busy, setBusy] = useState(false);

  const load = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const r = await usersApi.list({
        page: pagination.pageIndex + 1,
        pageSize: pagination.pageSize,
        q: q.trim() || undefined,
        role: role === "all" ? undefined : (role as "admin" | "user"),
        active: active === "all" ? undefined : active === "true",
      });
      setData(r.items);
      setTotal(r.total);
    } catch (e) {
      setError(e instanceof Error ? e.message : "加载失败");
    } finally {
      setLoading(false);
    }
  }, [pagination, q, role, active]);

  useEffect(() => {
    load();
  }, [load]);

  function resetPage() {
    setPagination((p) => ({ ...p, pageIndex: 0 }));
  }

  async function run(action: () => Promise<void>) {
    setBusy(true);
    setError(null);
    try {
      await action();
      await load();
    } catch (e) {
      setError(e instanceof Error ? e.message : "操作失败");
    } finally {
      setBusy(false);
    }
  }

  const columns = useMemo<ColumnDef<DataGridFeatures, LocalUserRow>[]>(
    () => [
      {
        accessorKey: "username",
        header: "用户名",
        cell: ({ row }) => <span className="font-medium">{row.original.username}</span>,
      },
      {
        accessorKey: "email",
        header: "邮箱",
        cell: ({ row }) => (
          <span className="line-clamp-1 text-muted-foreground">{row.original.email || "—"}</span>
        ),
      },
      {
        id: "role",
        header: "角色",
        cell: ({ row }) => (
          <Badge variant={row.original.role === "admin" ? "info-light" : "secondary"}>
            {row.original.role === "admin" ? "管理员" : "用户"}
          </Badge>
        ),
      },
      {
        id: "status",
        header: "状态",
        cell: ({ row }) => (
          <Badge variant={row.original.isActive ? "success-light" : "outline"}>
            {row.original.isActive ? "启用" : "停用"}
          </Badge>
        ),
      },
      {
        accessorKey: "createdAt",
        header: "创建时间",
        cell: ({ row }) => (
          <span className="whitespace-nowrap text-muted-foreground">
            {formatTime(row.original.createdAt)}
          </span>
        ),
      },
      {
        id: "actions",
        header: () => <span className="text-right">操作</span>,
        cell: ({ row }) => {
          const u = row.original;
          return (
            <div className="flex flex-wrap justify-end gap-1">
              <Button
                size="sm"
                variant="secondary"
                className="gap-1"
                onClick={() => {
                  setEditing(u);
                  setEditEmail(u.email);
                  setEditRole(u.role);
                }}
              >
                <Pencil className="h-3.5 w-3.5" />
                编辑
              </Button>
              <Button
                size="sm"
                variant="outline"
                className="gap-1"
                onClick={() => {
                  setPwdTarget(u);
                  setPwd("");
                }}
              >
                <KeyRound className="h-3.5 w-3.5" />
                重置密码
              </Button>
              <Button
                size="sm"
                variant="outline"
                disabled={busy}
                onClick={() => setToggleTarget(u)}
              >
                {u.isActive ? "停用" : "启用"}
              </Button>
              <Button
                size="sm"
                variant="destructive"
                className="gap-1"
                onClick={() => setDeleteTarget(u)}
              >
                <Trash2 className="h-3.5 w-3.5" />
                删除
              </Button>
            </div>
          );
        },
      },
    ],
    [busy],
  );

  const table = useTable({
    features: dataGridFeatures,
    data,
    columns,
    state: { pagination },
    onPaginationChange: setPagination,
    manualPagination: true,
    rowCount: total,
    getRowId: (row) => String(row.id),
  });

  return (
    <div className="flex flex-col gap-4">
      <div className="flex flex-wrap items-center gap-2">
        <Input
          value={q}
          onChange={(e) => {
            setQ(e.target.value);
            resetPage();
          }}
          placeholder="搜索用户名或邮箱"
          className="w-56"
        />
        <Select
          value={role}
          onValueChange={(v) => {
            setRole(String(v));
            resetPage();
          }}
        >
          <SelectTrigger className="w-32">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">全部角色</SelectItem>
            <SelectItem value="admin">管理员</SelectItem>
            <SelectItem value="user">用户</SelectItem>
          </SelectContent>
        </Select>
        <Select
          value={active}
          onValueChange={(v) => {
            setActive(String(v));
            resetPage();
          }}
        >
          <SelectTrigger className="w-32">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">全部状态</SelectItem>
            <SelectItem value="true">启用</SelectItem>
            <SelectItem value="false">停用</SelectItem>
          </SelectContent>
        </Select>
        <Button variant="outline" onClick={() => load()}>
          刷新
        </Button>
        <Button className="ml-auto" onClick={() => setCreateOpen(true)}>
          新建用户
        </Button>
      </div>

      {error ? <ErrorState message={error} onRetry={load} /> : null}

      <DataGrid
        table={table}
        recordCount={total}
        isLoading={loading}
        emptyMessage="暂无本地用户"
        tableLayout={{ dense: true, headerBackground: true, stripped: true }}
      >
        <DataGridContainer>
          <DataGridTable />
        </DataGridContainer>
        <DataGridPagination
          sizes={[10, 20, 50]}
          rowsPerPageLabel="每页"
          previousPageLabel="上一页"
          nextPageLabel="下一页"
        />
      </DataGrid>

      <p className="flex items-center gap-1.5 text-xs text-muted-foreground">
        <Users className="h-3.5 w-3.5" />
        仅显示 auth_provider=local 的账号，共 {total} 个
      </p>

      <Dialog open={createOpen} onOpenChange={setCreateOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>新建本地用户</DialogTitle>
          </DialogHeader>
          <div className="grid gap-3">
            <FormField label="用户名">
              <Input
                value={newUser.username}
                onChange={(e) => setNewUser({ ...newUser, username: e.target.value })}
              />
            </FormField>
            <FormField label="初始密码" description="至少 6 位">
              <Input
                type="password"
                value={newUser.password}
                onChange={(e) => setNewUser({ ...newUser, password: e.target.value })}
              />
            </FormField>
            <FormField label="邮箱（可选）">
              <Input
                value={newUser.email}
                onChange={(e) => setNewUser({ ...newUser, email: e.target.value })}
              />
            </FormField>
            <FormField label="角色">
              <Select
                value={newUser.role}
                onValueChange={(v) => setNewUser({ ...newUser, role: String(v) })}
              >
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="user">用户</SelectItem>
                  <SelectItem value="admin">管理员</SelectItem>
                </SelectContent>
              </Select>
            </FormField>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setCreateOpen(false)}>
              取消
            </Button>
            <Button
              disabled={busy}
              onClick={() =>
                run(async () => {
                  await usersApi.create({
                    username: newUser.username.trim(),
                    password: newUser.password,
                    email: newUser.email.trim(),
                    role: newUser.role,
                  });
                  setCreateOpen(false);
                  setNewUser({ username: "", password: "", email: "", role: "user" });
                })
              }
            >
              创建
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={!!editing} onOpenChange={(open) => !open && setEditing(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>编辑用户 · {editing?.username}</DialogTitle>
          </DialogHeader>
          <div className="grid gap-3">
            <FormField label="邮箱">
              <Input value={editEmail} onChange={(e) => setEditEmail(e.target.value)} />
            </FormField>
            <FormField label="角色">
              <Select value={editRole} onValueChange={(v) => setEditRole(String(v))}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="user">用户</SelectItem>
                  <SelectItem value="admin">管理员</SelectItem>
                </SelectContent>
              </Select>
            </FormField>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setEditing(null)}>
              取消
            </Button>
            <Button
              disabled={busy}
              onClick={() =>
                run(async () => {
                  if (!editing) return;
                  await usersApi.update(editing.id, { email: editEmail, role: editRole });
                  setEditing(null);
                })
              }
            >
              保存
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={!!pwdTarget} onOpenChange={(open) => !open && setPwdTarget(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>重置密码 · {pwdTarget?.username}</DialogTitle>
          </DialogHeader>
          <FormField label="新密码" description="至少 6 位">
            <Input type="password" value={pwd} onChange={(e) => setPwd(e.target.value)} />
          </FormField>
          <DialogFooter>
            <Button variant="outline" onClick={() => setPwdTarget(null)}>
              取消
            </Button>
            <Button
              disabled={busy || pwd.length < 6}
              onClick={() =>
                run(async () => {
                  if (!pwdTarget) return;
                  await usersApi.resetPassword(pwdTarget.id, pwd);
                  setPwdTarget(null);
                })
              }
            >
              重置
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <ConfirmDialog
        open={!!toggleTarget}
        onOpenChange={(open) => {
          if (!open) setToggleTarget(null);
        }}
        title={toggleTarget?.isActive ? "停用用户" : "启用用户"}
        description={
          toggleTarget
            ? `${toggleTarget.isActive ? "停用后该账号将无法登录。" : "启用后该账号可重新登录。"}`
            : undefined
        }
        confirmLabel={toggleTarget?.isActive ? "停用" : "启用"}
        destructive={!!toggleTarget?.isActive}
        onConfirm={async () => {
          if (!toggleTarget) return;
          await run(async () => {
            await usersApi.update(toggleTarget.id, { is_active: !toggleTarget.isActive });
            setToggleTarget(null);
          });
        }}
      />

      <ConfirmDialog
        open={!!deleteTarget}
        onOpenChange={(open) => {
          if (!open) setDeleteTarget(null);
        }}
        title="删除用户"
        description={deleteTarget ? `确定删除「${deleteTarget.username}」？该操作不可撤销。` : undefined}
        confirmLabel="删除"
        destructive
        onConfirm={async () => {
          if (!deleteTarget) return;
          await run(async () => {
            await usersApi.remove(deleteTarget.id);
            setDeleteTarget(null);
          });
        }}
      />
    </div>
  );
}
