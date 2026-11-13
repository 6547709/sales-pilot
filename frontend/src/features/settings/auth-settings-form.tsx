"use client";

import { useEffect, useState } from "react";
import { Alert, AlertDescription } from "@/components/reui/alert";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Frame, FrameHeader, FramePanel, FrameTitle } from "@/components/reui/frame";
import { ConfirmDialog } from "@/components/patterns/confirm-dialog";
import { FormField } from "@/components/patterns/form-field";
import { settingsApi, type AuthSettingsDTO } from "@/lib/api";

type AuthState = {
  localLoginEnabled: boolean;
  ldapEnabled: boolean;
  ldapConfig: string;
  oidcEnabled: boolean;
  oidcConfig: string;
};

function fromDTO(dto: AuthSettingsDTO): AuthState {
  return {
    localLoginEnabled: dto.local_login_enabled !== false,
    ldapEnabled: dto.ldap_enabled,
    ldapConfig:
      typeof dto.ldap_config === "string"
        ? dto.ldap_config
        : JSON.stringify(dto.ldap_config || {}, null, 2),
    oidcEnabled: dto.oidc_enabled,
    oidcConfig:
      typeof dto.oidc_config === "string"
        ? dto.oidc_config
        : JSON.stringify(dto.oidc_config || {}, null, 2),
  };
}

export function AuthSettingsForm() {
  const [st, setSt] = useState<AuthState | null>(null);
  const [msg, setMsg] = useState<{ kind: "ok" | "error"; text: string } | null>(null);
  const [confirmOpen, setConfirmOpen] = useState(false);
  const [busy, setBusy] = useState(false);

  useEffect(() => {
    settingsApi
      .auth()
      .then((dto) => setSt(fromDTO(dto)))
      .catch((e) => setMsg({ kind: "error", text: e instanceof Error ? e.message : "加载失败" }));
  }, []);

  async function save() {
    if (!st) return;
    setBusy(true);
    setMsg(null);
    try {
      const ldap = JSON.parse(st.ldapConfig || "{}");
      const oidc = JSON.parse(st.oidcConfig || "{}");
      const dto = await settingsApi.patchAuth({
        local_login_enabled: st.localLoginEnabled,
        ldap_enabled: st.ldapEnabled,
        ldap_config: ldap,
        oidc_enabled: st.oidcEnabled,
        oidc_config: oidc,
      });
      setSt(fromDTO(dto));
      setMsg({ kind: "ok", text: "认证配置已保存" });
    } catch (e) {
      setMsg({
        kind: "error",
        text: e instanceof SyntaxError ? "JSON 格式错误" : e instanceof Error ? e.message : "保存失败",
      });
    } finally {
      setBusy(false);
      setConfirmOpen(false);
    }
  }

  if (!st) return <p className="text-sm text-muted-foreground">加载中…</p>;

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
            <FrameTitle>本地账号</FrameTitle>
          </FrameHeader>
          <div className="flex flex-col gap-2 px-5 pb-5">
            <label className="flex items-center gap-2 text-sm">
              <input
                type="checkbox"
                checked={st.localLoginEnabled}
                onChange={(e) => setSt({ ...st, localLoginEnabled: e.target.checked })}
              />
              允许用户名密码登录
            </label>
            <p className="text-xs text-muted-foreground">
              关闭后只能通过已启用的 LDAP / OIDC 登录，请至少保留一种登录方式。
            </p>
          </div>
        </FramePanel>

        <FramePanel>
          <FrameHeader>
            <FrameTitle>LDAP</FrameTitle>
          </FrameHeader>
          <div className="flex flex-col gap-3 px-5 pb-5">
            <label className="flex items-center gap-2 text-sm">
              <input
                type="checkbox"
                checked={st.ldapEnabled}
                onChange={(e) => setSt({ ...st, ldapEnabled: e.target.checked })}
              />
              启用 LDAP
            </label>
            <Textarea
              rows={10}
              className="font-mono text-xs"
              value={st.ldapConfig}
              onChange={(e) => setSt({ ...st, ldapConfig: e.target.value })}
            />
          </div>
        </FramePanel>

        <FramePanel>
          <FrameHeader>
            <FrameTitle>OIDC (SSO)</FrameTitle>
          </FrameHeader>
          <div className="flex flex-col gap-3 px-5 pb-5">
            <label className="flex items-center gap-2 text-sm">
              <input
                type="checkbox"
                checked={st.oidcEnabled}
                onChange={(e) => setSt({ ...st, oidcEnabled: e.target.checked })}
              />
              启用 OIDC
            </label>
            <p className="text-xs text-muted-foreground">
              redirect_url 需与 IdP 登记一致，指向{" "}
              <code className="rounded bg-muted px-1">/api/v1/auth/oidc/callback</code>
            </p>
            <Textarea
              rows={12}
              className="font-mono text-xs"
              value={st.oidcConfig}
              onChange={(e) => setSt({ ...st, oidcConfig: e.target.value })}
            />
          </div>
        </FramePanel>
      </Frame>

      <div>
        <Button disabled={busy} onClick={() => setConfirmOpen(true)}>
          保存认证配置
        </Button>
      </div>

      <ConfirmDialog
        open={confirmOpen}
        onOpenChange={setConfirmOpen}
        title="保存认证配置"
        description="修改登录方式可能影响现有账号的登录。确认保存？"
        confirmLabel="保存"
        onConfirm={save}
      />
    </div>
  );
}
