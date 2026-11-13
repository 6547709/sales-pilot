"use client";

import { useRef, useState } from "react";
import { Download, Upload } from "lucide-react";
import { Alert, AlertDescription, AlertTitle } from "@/components/reui/alert";
import { Frame, FrameHeader, FramePanel, FrameTitle } from "@/components/reui/frame";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { ConfirmDialog } from "@/components/patterns/confirm-dialog";
import { FormField } from "@/components/patterns/form-field";
import { apiFetch, BACKUP_RESTORE_CONFIRM } from "@/lib/api";
import { clearToken } from "@/lib/auth";

type BackupSummary = {
  version: number;
  products: number;
  users: number;
  scripts: number;
  cases: number;
  layers: number;
};

function summarize(value: unknown): BackupSummary | null {
  if (!value || typeof value !== "object") return null;
  const v = value as Record<string, unknown>;
  const n = (key: string) => (Array.isArray(v[key]) ? (v[key] as unknown[]).length : 0);
  const version = typeof v.version === "number" ? v.version : 0;
  return {
    version,
    products: n("products"),
    users: n("users"),
    scripts: n("sales_scripts"),
    cases: n("cases"),
    layers: n("topology_layers"),
  };
}

export function BackupView() {
  const [msg, setMsg] = useState<{ kind: "ok" | "error"; text: string } | null>(null);
  const [busy, setBusy] = useState(false);
  const [restoreToken, setRestoreToken] = useState("");
  const [jsonText, setJsonText] = useState("");
  const [pending, setPending] = useState<{ backup: unknown; summary: BackupSummary } | null>(null);
  const fileRef = useRef<HTMLInputElement>(null);

  async function exportBackup() {
    setMsg(null);
    setBusy(true);
    try {
      const res = await apiFetch("/api/v1/admin/backup/export");
      if (!res.ok) {
        setMsg({ kind: "error", text: "导出失败" });
        return;
      }
      const blob = await res.blob();
      const url = URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = `sales-pilot-backup-${new Date().toISOString().slice(0, 19).replace(/:/g, "-")}.json`;
      a.click();
      URL.revokeObjectURL(url);
      setMsg({ kind: "ok", text: "已下载备份文件，请妥善保管（含密码哈希）。" });
    } finally {
      setBusy(false);
    }
  }

  function startImport() {
    setMsg(null);
    if (restoreToken !== BACKUP_RESTORE_CONFIRM) {
      setMsg({ kind: "error", text: `请在确认框中完整输入：${BACKUP_RESTORE_CONFIRM}` });
      return;
    }
    let backup: unknown;
    try {
      backup = JSON.parse(jsonText || "{}");
    } catch {
      setMsg({ kind: "error", text: "JSON 解析失败" });
      return;
    }
    const summary = summarize(backup);
    if (!summary || (summary.version !== 1 && summary.version !== 2)) {
      setMsg({ kind: "error", text: "备份文件格式无效（需 version: 1 或 2）" });
      return;
    }
    setPending({ backup, summary });
  }

  async function doImport() {
    if (!pending) return;
    setBusy(true);
    try {
      const data = await apiFetch("/api/v1/admin/backup/import", {
        method: "POST",
        body: JSON.stringify({ confirm: BACKUP_RESTORE_CONFIRM, backup: pending.backup }),
      });
      const body = await data.json().catch(() => ({}));
      if (data.ok) {
        clearToken();
        setMsg({ kind: "ok", text: body.message || "恢复成功，即将跳转登录…" });
        setTimeout(() => {
          window.location.href = "/login";
        }, 1500);
      } else {
        setMsg({ kind: "error", text: body.error || "恢复失败" });
      }
    } finally {
      setBusy(false);
      setPending(null);
    }
  }

  function onPickFile(f: File | null) {
    if (!f) return;
    const reader = new FileReader();
    reader.onload = () => {
      setJsonText(String(reader.result || ""));
      setMsg({ kind: "ok", text: `已读取文件：${f.name}` });
    };
    reader.readAsText(f);
  }

  return (
    <div className="flex flex-col gap-4">
      {msg ? (
        <Alert variant={msg.kind === "ok" ? "success" : "destructive"}>
          <AlertDescription>{msg.text}</AlertDescription>
        </Alert>
      ) : null}

      <Frame stacked>
        <FramePanel>
          <FrameHeader>
            <FrameTitle>导出备份</FrameTitle>
          </FrameHeader>
          <div className="flex flex-col gap-3 px-5 pb-5">
            <p className="text-xs text-muted-foreground">
              包含拓扑、产品、话术、案例、用户（含本地密码哈希）与认证配置。
            </p>
            <div>
              <Button type="button" disabled={busy} onClick={exportBackup} className="gap-1">
                <Download className="h-4 w-4" />
                下载 JSON 备份
              </Button>
            </div>
          </div>
        </FramePanel>

        <FramePanel>
          <FrameHeader>
            <FrameTitle>恢复备份（危险操作）</FrameTitle>
          </FrameHeader>
          <div className="flex flex-col gap-3 px-5 pb-5">
            <Alert variant="warning">
              <AlertTitle>整库覆盖</AlertTitle>
              <AlertDescription>
                恢复会清空当前数据库并写入备份中的全部数据，所有用户需重新登录。
              </AlertDescription>
            </Alert>
            <FormField label="选择备份文件或粘贴 JSON">
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
            </FormField>
            <FormField
              label="确认令牌"
              description={`请输入 ${BACKUP_RESTORE_CONFIRM} 以继续`}
            >
              <Input
                value={restoreToken}
                onChange={(e) => setRestoreToken(e.target.value)}
                placeholder={BACKUP_RESTORE_CONFIRM}
              />
            </FormField>
            <div>
              <Button type="button" variant="destructive" disabled={busy} onClick={startImport} className="gap-1">
                <Upload className="h-4 w-4" />
                执行恢复
              </Button>
            </div>
          </div>
        </FramePanel>
      </Frame>

      <ConfirmDialog
        open={!!pending}
        onOpenChange={(open) => {
          if (!open) setPending(null);
        }}
        title="确认整库覆盖"
        description={
          pending
            ? `将写入 ${pending.summary.products} 个产品、${pending.summary.users} 个用户、${pending.summary.scripts} 条话术、${pending.summary.cases} 个案例。此操作不可撤销。`
            : undefined
        }
        confirmLabel="覆盖并恢复"
        destructive
        onConfirm={doImport}
      />
    </div>
  );
}
