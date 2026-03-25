"use client";

import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { apiFetch } from "@/lib/api";

type AuthSettings = {
  id: number;
  /** 旧库可能缺省，按 true 处理 */
  local_login_enabled?: boolean;
  ldap_enabled: boolean;
  ldap_config: string;
  oidc_enabled: boolean;
  oidc_config: string;
};

export default function AdminSettingsPage() {
  const [st, setSt] = useState<AuthSettings | null>(null);
  const [ldapJson, setLdapJson] = useState("");
  const [oidcJson, setOidcJson] = useState("");
  const [msg, setMsg] = useState("");

  useEffect(() => {
    apiFetch("/api/v1/admin/auth-settings").then(async (r) => {
      if (!r.ok) return;
      const data = await r.json();
      setSt({
        ...data,
        local_login_enabled: data.local_login_enabled !== false,
      });
      setLdapJson(
        typeof data.ldap_config === "string"
          ? data.ldap_config
          : JSON.stringify(data.ldap_config || {}, null, 2),
      );
      setOidcJson(
        typeof data.oidc_config === "string"
          ? data.oidc_config
          : JSON.stringify(data.oidc_config || {}, null, 2),
      );
    });
  }, []);

  async function save() {
    if (!st) return;
    setMsg("");
    let ldapParsed: unknown;
    let oidcParsed: unknown;
    try {
      ldapParsed = JSON.parse(ldapJson || "{}");
      oidcParsed = JSON.parse(oidcJson || "{}");
    } catch {
      setMsg("JSON 格式错误");
      return;
    }
    const res = await apiFetch("/api/v1/admin/auth-settings", {
      method: "PATCH",
      body: JSON.stringify({
        local_login_enabled: st.local_login_enabled !== false,
        ldap_enabled: st?.ldap_enabled,
        ldap_config: ldapParsed,
        oidc_enabled: st?.oidc_enabled,
        oidc_config: oidcParsed,
      }),
    });
    if (res.ok) {
      setMsg("已保存");
      const data = await res.json();
      setSt(data);
    } else {
      const e = await res.json().catch(() => ({}));
      setMsg(e.error || "保存失败");
    }
  }

  if (!st) return <p className="text-muted-foreground">加载中…</p>;

  return (
    <div className="mx-auto max-w-3xl space-y-6">
      <h1 className="text-2xl font-bold text-primary">认证配置</h1>
      <Card>
        <CardHeader>
          <CardTitle>本地账号</CardTitle>
        </CardHeader>
        <CardContent className="space-y-2">
          <div className="flex items-center gap-2">
            <input
              type="checkbox"
              id="lle"
              checked={st.local_login_enabled}
              onChange={(e) =>
                setSt({ ...st, local_login_enabled: e.target.checked })
              }
            />
            <Label htmlFor="lle">允许用户名密码登录</Label>
          </div>
          <p className="text-xs text-muted-foreground">
            关闭后，将无法使用本地用户名密码登录；仅可通过已启用的 LDAP /
            OIDC 进入系统。请至少保留一种登录方式。
          </p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <CardTitle>LDAP</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-center gap-2">
            <input
              type="checkbox"
              id="le"
              checked={st.ldap_enabled}
              onChange={(e) =>
                setSt({ ...st, ldap_enabled: e.target.checked })
              }
            />
            <Label htmlFor="le">启用 LDAP</Label>
          </div>
          <Textarea
            rows={10}
            className="font-mono text-xs"
            value={ldapJson}
            onChange={(e) => setLdapJson(e.target.value)}
          />
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <CardTitle>OIDC (SSO)</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-center gap-2">
            <input
              type="checkbox"
              id="oe"
              checked={st.oidc_enabled}
              onChange={(e) =>
                setSt({ ...st, oidc_enabled: e.target.checked })
              }
            />
            <Label htmlFor="oe">启用 OIDC</Label>
          </div>
          <p className="text-xs text-muted-foreground">
            redirect_url 需与 IdP 登记一致，指向{" "}
            <code className="rounded bg-muted px-1">
              /api/v1/auth/oidc/callback
            </code>
          </p>
          <Textarea
            rows={12}
            className="font-mono text-xs"
            value={oidcJson}
            onChange={(e) => setOidcJson(e.target.value)}
          />
        </CardContent>
      </Card>
      {msg ? <p className="text-sm text-primary">{msg}</p> : null}
      <Button onClick={save}>保存</Button>
    </div>
  );
}
