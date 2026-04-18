"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { API_BASE } from "@/lib/config";
import { useFooterConfig } from "@/lib/footer-config";
import {
  fetchLoginOptions,
  fetchMe,
  loginLdap,
  loginLocal,
  type LoginOptions,
} from "@/lib/api";
import { getToken, setToken } from "@/lib/auth";

export default function LoginPage() {
  const router = useRouter();
  const { footerText, filingNumber } = useFooterConfig();
  const [u, setU] = useState("");
  const [p, setP] = useState("");
  const [err, setErr] = useState("");
  const [loading, setLoading] = useState(false);
  const [opts, setOpts] = useState<LoginOptions | null>(null);

  useEffect(() => {
    fetchLoginOptions().then(setOpts);
  }, []);

  // 已登录用户无需停留在登录页
  useEffect(() => {
    const t = getToken();
    if (!t) return;
    fetchMe().then((user) => {
      if (user) router.replace(user.role === "admin" ? "/admin" : "/");
    });
  }, [router]);

  async function submitLocal(e: React.FormEvent) {
    e.preventDefault();
    setErr("");
    setLoading(true);
    try {
      const { access_token, user } = await loginLocal(u, p);
      setToken(access_token);
      router.replace(user.role === "admin" ? "/admin" : "/");
    } catch (ex) {
      setErr(ex instanceof Error ? ex.message : "登录失败");
    } finally {
      setLoading(false);
    }
  }

  async function submitLdap(e: React.FormEvent) {
    e.preventDefault();
    setErr("");
    setLoading(true);
    try {
      const { access_token, user } = await loginLdap(u, p);
      setToken(access_token);
      router.replace(user.role === "admin" ? "/admin" : "/");
    } catch (ex) {
      setErr(ex instanceof Error ? ex.message : "LDAP 失败");
    } finally {
      setLoading(false);
    }
  }

  function startOidc() {
    window.location.href = `${API_BASE}/api/v1/auth/oidc/start`;
  }

  const showLocal = opts?.local_enabled ?? true;
  const showLdap = opts?.ldap_enabled ?? false;
  const showOidc = opts?.oidc_enabled ?? false;
  const passwordTabs = (showLocal ? 1 : 0) + (showLdap ? 1 : 0);

  if (!opts) {
    return (
      <div className="flex min-h-screen flex-col items-center justify-center bg-gradient-to-b from-slate-50 to-white p-4 text-muted-foreground">
        加载登录方式…
      </div>
    );
  }

  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-gradient-to-b from-slate-50 to-white p-4">
      <Card className="w-full max-w-md border-primary/20 shadow-lg">
        <CardHeader>
          <CardTitle className="text-primary">Sales-Pilot</CardTitle>
          <CardDescription>企业销售赋能系统 — 登录</CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          {passwordTabs > 0 ? (
            <Tabs defaultValue={showLocal ? "local" : "ldap"}>
              <TabsList
                className={`grid w-full ${passwordTabs === 2 ? "grid-cols-2" : "grid-cols-1"}`}
              >
                {showLocal ? (
                  <TabsTrigger value="local">本地账号</TabsTrigger>
                ) : null}
                {showLdap ? (
                  <TabsTrigger value="ldap">LDAP</TabsTrigger>
                ) : null}
              </TabsList>
              {showLocal ? (
                <TabsContent value="local">
                  <form onSubmit={submitLocal} className="space-y-4">
                    <div className="space-y-2">
                      <Label htmlFor="user">用户名</Label>
                      <Input
                        id="user"
                        value={u}
                        onChange={(e) => setU(e.target.value)}
                        autoComplete="username"
                      />
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="pass">密码</Label>
                      <Input
                        id="pass"
                        type="password"
                        value={p}
                        onChange={(e) => setP(e.target.value)}
                        autoComplete="current-password"
                      />
                    </div>
                    {err ? (
                      <p className="text-sm text-destructive">{err}</p>
                    ) : null}
                    <Button
                      type="submit"
                      className="w-full"
                      disabled={loading}
                    >
                      {loading ? "登录中…" : "登录"}
                    </Button>
                  </form>
                </TabsContent>
              ) : null}
              {showLdap ? (
                <TabsContent value="ldap">
                  <form onSubmit={submitLdap} className="space-y-4">
                    <div className="space-y-2">
                      <Label htmlFor="lu">用户名</Label>
                      <Input
                        id="lu"
                        value={u}
                        onChange={(e) => setU(e.target.value)}
                      />
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="lp">密码</Label>
                      <Input
                        id="lp"
                        type="password"
                        value={p}
                        onChange={(e) => setP(e.target.value)}
                      />
                    </div>
                    {err ? (
                      <p className="text-sm text-destructive">{err}</p>
                    ) : null}
                    <Button
                      type="submit"
                      className="w-full"
                      variant="secondary"
                      disabled={loading}
                    >
                      LDAP 登录
                    </Button>
                  </form>
                </TabsContent>
              ) : null}
            </Tabs>
          ) : null}

          {passwordTabs > 0 && showOidc ? (
            <div className="relative">
              <div className="absolute inset-0 flex items-center">
                <span className="w-full border-t" />
              </div>
              <div className="relative flex justify-center text-xs uppercase">
                <span className="bg-card px-2 text-muted-foreground">或</span>
              </div>
            </div>
          ) : null}

          {showOidc ? (
            <Button
              type="button"
              variant="outline"
              className="w-full"
              onClick={startOidc}
            >
              使用企业单点登录 (SSO)
            </Button>
          ) : null}

          {!showLocal && !showLdap && !showOidc ? (
            <p className="text-center text-sm text-destructive">
              未启用任何登录方式，请管理员在「认证配置」中至少开启一种方式。
            </p>
          ) : null}

          <p className="text-center text-xs text-muted-foreground">
            登录后可访问首页与方案库等企业内容。
          </p>
          {(footerText || filingNumber) && (
            <div className="mt-4 border-t pt-4 text-center text-xs text-muted-foreground">
              {footerText}
              {footerText && filingNumber && <span className="mx-2">|</span>}
              {filingNumber}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
