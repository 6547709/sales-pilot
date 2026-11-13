"use client";

import { useRef, useState } from "react";
import { Download, Upload } from "lucide-react";
import { Alert, AlertDescription, AlertTitle } from "@/components/reui/alert";
import { Frame, FrameHeader, FramePanel, FrameTitle } from "@/components/reui/frame";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { ConfirmDialog } from "@/components/patterns/confirm-dialog";
import { FormField } from "@/components/patterns/form-field";
import { apiFetch, backupApi, BACKUP_RESTORE_CONFIRM } from "@/lib/api";
import { clearToken } from "@/lib/auth";

type SectionKey = "content" | "config" | "users";

const SECTION_LABEL: Record<SectionKey, string> = {
  content: "内容数据（拓扑、产品、服务、场景）",
  config: "系统配置（认证、页脚）",
  users: "用户与 API 密钥",
};

export function BackupView() {
  const [msg, setMsg] = useState<{ kind: "ok" | "error"; text: string } | null>(null);
  const [busy, setBusy] = useState(false);
  const [token, setToken] = useState("");
  const [file, setFile] = useState<File | null>(null);
  const [sections, setSections] = useState<Record<SectionKey, boolean>>({
    content: true,
    config: true,
    users: true,
  });
  const [confirmOpen, setConfirmOpen] = useState(false);
  const fileRef = useRef<HTMLInputElement>(null);

  const selected = (Object.keys(sections) as SectionKey[]).filter((k) => sections[k]);

  async function exportBackup() {
    setMsg(null);
    setBusy(true);
    try {
      const res = await apiFetch(backupApi.exportZipPath);
      if (!res.ok) {
        setMsg({ kind: "error", text: "导出失败" });
        return;
      }
      const blob = await res.blob();
      const url = URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = `sales-pilot-backup-${new Date().toISOString().slice(0, 19).replace(/:/g, "-")}.zip`;
      a.click();
      URL.revokeObjectURL(url);
      setMsg({
        kind: "ok",
        text: "已下载 zip 备份（含 content.json / config.json / users.json）。",
      });
    } finally {
      setBusy(false);
    }
  }

  function startImport() {
    setMsg(null);
    if (!file) {
      setMsg({ kind: "error", text: "请选择备份 zip 文件" });
      return;
    }
    if (selected.length === 0) {
      setMsg({ kind: "error", text: "请至少选择一个恢复分区" });
      return;
    }
    if (token !== BACKUP_RESTORE_CONFIRM) {
      setMsg({ kind: "error", text: `请输入确认令牌：${BACKUP_RESTORE_CONFIRM}` });
      return;
    }
    setConfirmOpen(true);
  }

  async function doImport() {
    if (!file) return;
    setBusy(true);
    try {
      const report = await backupApi.importZip(file, selected, token);
      const countText = Object.entries(report.counts)
        .filter(([, n]) => n > 0)
        .map(([k, n]) => `${k} ${n}`)
        .join("、");
      if (sections.users) {
        clearToken();
        setMsg({
          kind: "ok",
          text: `恢复完成（${countText}）。用户数据已覆盖，正在跳转登录…`,
        });
        setTimeout(() => {
          window.location.href = "/login";
        }, 1500);
      } else {
        setMsg({ kind: "ok", text: `恢复完成（${countText}），用户与登录状态未受影响。` });
      }
      setFile(null);
      if (fileRef.current) fileRef.current.value = "";
    } catch (e) {
      setMsg({ kind: "error", text: e instanceof Error ? e.message : "恢复失败" });
    } finally {
      setBusy(false);
      setConfirmOpen(false);
    }
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
              导出 zip，内含三份独立 JSON：内容数据、系统配置、用户与密钥。
            </p>
            <div>
              <Button type="button" disabled={busy} onClick={exportBackup} className="gap-1">
                <Download className="h-4 w-4" />
                下载 zip 备份
              </Button>
            </div>
          </div>
        </FramePanel>

        <FramePanel>
          <FrameHeader>
            <FrameTitle>恢复备份</FrameTitle>
          </FrameHeader>
          <div className="flex flex-col gap-4 px-5 pb-5">
            <Alert variant="warning">
              <AlertTitle>按分区覆盖</AlertTitle>
              <AlertDescription>
                只覆盖勾选的分区，未勾选的数据保持不变。恢复用户会清空当前登录状态。
              </AlertDescription>
            </Alert>

            <FormField label="备份 zip 文件">
              <input
                ref={fileRef}
                type="file"
                accept=".zip,application/zip"
                className="block text-sm"
                onChange={(e) => setFile(e.target.files?.[0] ?? null)}
              />
            </FormField>

            <div className="flex flex-col gap-2">
              <p className="text-xs font-medium text-muted-foreground">选择恢复分区</p>
              {(Object.keys(SECTION_LABEL) as SectionKey[]).map((key) => (
                <label key={key} className="flex items-center gap-2 text-sm">
                  <input
                    type="checkbox"
                    checked={sections[key]}
                    onChange={(e) => setSections({ ...sections, [key]: e.target.checked })}
                  />
                  {SECTION_LABEL[key]}
                </label>
              ))}
            </div>

            <FormField
              label="确认令牌"
              description={`请输入 ${BACKUP_RESTORE_CONFIRM} 以继续`}
            >
              <Input
                value={token}
                onChange={(e) => setToken(e.target.value)}
                placeholder={BACKUP_RESTORE_CONFIRM}
              />
            </FormField>

            <div>
              <Button type="button" disabled={busy} onClick={startImport} className="gap-1">
                <Upload className="h-4 w-4" />
                执行恢复
              </Button>
            </div>
          </div>
        </FramePanel>
      </Frame>

      <ConfirmDialog
        open={confirmOpen}
        onOpenChange={setConfirmOpen}
        title="确认覆盖所选分区"
        description={
          selected.length > 0
            ? `将覆盖：${selected.map((k) => SECTION_LABEL[k]).join("；")}。此操作不可撤销。`
            : undefined
        }
        confirmLabel="覆盖并恢复"
        destructive
        onConfirm={doImport}
      />
    </div>
  );
}
