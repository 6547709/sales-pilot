"use client";

import { useRef, useState } from "react";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { apiFetch, BACKUP_RESTORE_CONFIRM } from "@/lib/api";
import { clearToken } from "@/lib/auth";
import { Download, Upload } from "lucide-react";

export default function AdminBackupPage() {
  const [msg, setMsg] = useState("");
  const [busy, setBusy] = useState(false);
  const [restoreToken, setRestoreToken] = useState("");
  const [jsonText, setJsonText] = useState("");
  const fileRef = useRef<HTMLInputElement>(null);

  async function exportBackup() {
    setMsg("");
    setBusy(true);
    try {
      const res = await apiFetch("/api/v1/admin/backup/export");
      if (!res.ok) {
        setMsg("导出失败");
        return;
      }
      const blob = await res.blob();
      const url = URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = `sales-pilot-backup-${new Date().toISOString().slice(0, 19).replace(/:/g, "-")}.json`;
      a.click();
      URL.revokeObjectURL(url);
      setMsg("已下载备份文件，请妥善保管（含密码哈希）。");
    } finally {
      setBusy(false);
    }
  }

  async function importBackup() {
    setMsg("");
    if (restoreToken !== BACKUP_RESTORE_CONFIRM) {
      setMsg(`请在确认框中完整输入：${BACKUP_RESTORE_CONFIRM}`);
      return;
    }
    let backup: unknown;
    try {
      backup = JSON.parse(jsonText || "{}");
    } catch {
      setMsg("JSON 解析失败");
      return;
    }
    if (
      !backup ||
      typeof backup !== "object" ||
      (backup as { version?: number }).version !== 1
    ) {
      setMsg("备份文件格式无效（需 version: 1）");
      return;
    }
    if (
      !window.confirm(
        "将清空当前数据库并写入备份中的全部数据。所有用户需重新登录。确定继续？",
      )
    ) {
      return;
    }
    setBusy(true);
    try {
      const res = await apiFetch("/api/v1/admin/backup/import", {
        method: "POST",
        body: JSON.stringify({
          confirm: BACKUP_RESTORE_CONFIRM,
          backup,
        }),
      });
      const data = await res.json().catch(() => ({}));
      if (res.ok) {
        clearToken();
        setMsg(data.message || "恢复成功，即将跳转登录…");
        setTimeout(() => {
          window.location.href = "/login";
        }, 1500);
      } else {
        setMsg(data.error || "恢复失败");
      }
    } finally {
      setBusy(false);
    }
  }

  function onPickFile(f: File | null) {
    if (!f) return;
    const reader = new FileReader();
    reader.onload = () => {
      setJsonText(String(reader.result || ""));
      setMsg(`已读取文件：${f.name}`);
    };
    reader.readAsText(f);
  }

  return (
    <div className="mx-auto max-w-3xl space-y-8">
      <div>
        <Link href="/admin" className="text-sm text-primary">
          ← 控制台
        </Link>
        <h1 className="mt-1 text-2xl font-bold text-primary">备份与恢复</h1>
        <p className="mt-2 text-sm text-muted-foreground">
          导出为 JSON，包含拓扑、产品、话术、案例、用户（含本地密码哈希）与认证配置。
          恢复将<strong>整库覆盖</strong>，请先在安全环境验证备份文件。
        </p>
      </div>

      {msg ? (
        <p className="rounded-md border border-primary/30 bg-primary/5 px-3 py-2 text-sm">
          {msg}
        </p>
      ) : null}

      <section className="rounded-lg border bg-card p-4">
        <h2 className="mb-3 flex items-center gap-2 font-semibold">
          <Download className="h-5 w-5 text-primary" />
          导出备份
        </h2>
        <Button type="button" disabled={busy} onClick={exportBackup}>
          下载 JSON 备份
        </Button>
      </section>

      <section className="rounded-lg border border-destructive/30 bg-destructive/[0.03] p-4">
        <h2 className="mb-3 flex items-center gap-2 font-semibold text-destructive">
          <Upload className="h-5 w-5" />
          恢复备份（危险操作）
        </h2>
        <div className="space-y-3">
          <div>
            <Label className="mb-2 block">选择备份文件或粘贴 JSON</Label>
            <input
              ref={fileRef}
              type="file"
              accept=".json,application/json"
              className="mb-2 block text-sm"
              onChange={(e) => onPickFile(e.target.files?.[0] ?? null)}
            />
            <Textarea
              rows={8}
              className="font-mono text-xs"
              placeholder="或在此粘贴 sales-pilot-backup.json 全文…"
              value={jsonText}
              onChange={(e) => setJsonText(e.target.value)}
            />
          </div>
          <div>
            <Label>
              确认令牌（输入：<code className="rounded bg-muted px-1">{BACKUP_RESTORE_CONFIRM}</code>）
            </Label>
            <Input
              className="mt-1"
              value={restoreToken}
              onChange={(e) => setRestoreToken(e.target.value)}
              placeholder={BACKUP_RESTORE_CONFIRM}
            />
          </div>
          <Button
            type="button"
            variant="destructive"
            disabled={busy}
            onClick={importBackup}
          >
            执行恢复
          </Button>
        </div>
      </section>
    </div>
  );
}
