"use client";

import { useEffect, useState } from "react";
import { Alert, AlertDescription } from "@/components/reui/alert";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Frame, FramePanel } from "@/components/reui/frame";
import { FormField } from "@/components/patterns/form-field";
import { settingsApi } from "@/lib/api";

export function SystemSettingsForm() {
  const [footer, setFooter] = useState("");
  const [filing, setFiling] = useState("");
  const [msg, setMsg] = useState<{ kind: "ok" | "error"; text: string } | null>(null);
  const [busy, setBusy] = useState(false);

  useEffect(() => {
    settingsApi
      .config()
      .then((dto) => {
        setFooter(dto.footer_text ?? "");
        setFiling(dto.filing_number ?? "");
      })
      .catch((e) =>
        setMsg({ kind: "error", text: e instanceof Error ? e.message : "加载失败" }),
      );
  }, []);

  async function save() {
    setBusy(true);
    setMsg(null);
    try {
      await settingsApi.patchSystem({ footer_text: footer, filing_number: filing });
      setMsg({ kind: "ok", text: "系统配置已保存" });
    } catch (e) {
      setMsg({ kind: "error", text: e instanceof Error ? e.message : "保存失败" });
    } finally {
      setBusy(false);
    }
  }

  return (
    <div className="flex flex-col gap-4">
      {msg ? (
        <Alert variant={msg.kind === "ok" ? "success" : "destructive"}>
          <AlertDescription>{msg.text}</AlertDescription>
        </Alert>
      ) : null}

      <Frame>
        <FramePanel>
          <div className="grid gap-4">
            <FormField label="页脚文字">
              <Input
                value={footer}
                onChange={(e) => setFooter(e.target.value)}
                placeholder="如：企业销售赋能"
              />
            </FormField>
            <FormField label="备案号">
              <Input
                value={filing}
                onChange={(e) => setFiling(e.target.value)}
                placeholder="如：京ICP备xxxx号"
              />
            </FormField>
            <p className="text-xs text-muted-foreground">
              将在首页、登录页与管理后台底部显示。
            </p>
          </div>
        </FramePanel>
      </Frame>

      <div>
        <Button disabled={busy} onClick={save}>
          保存系统配置
        </Button>
      </div>
    </div>
  );
}
