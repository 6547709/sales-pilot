"use client";

import { useState } from "react";
import Link from "next/link";
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
import { loginLdap, loginLocal } from "@/lib/api";
import { setToken } from "@/lib/auth";

export default function LoginPage() {
  const router = useRouter();
  const [u, setU] = useState("");
  const [p, setP] = useState("");
  const [err, setErr] = useState("");
  const [loading, setLoading] = useState(false);

  async function submitLocal(e: React.FormEvent) {
    e.preventDefault();
    setErr("");
    setLoading(true);
    try {
      const { access_token } = await loginLocal(u, p);
      setToken(access_token);
      router.replace("/products");
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
      const { access_token } = await loginLdap(u, p);
      setToken(access_token);
      router.replace("/products");
    } catch (ex) {
      setErr(ex instanceof Error ? ex.message : "LDAP 失败");
    } finally {
      setLoading(false);
    }
  }

  function startOidc() {
    window.location.href = `${API_BASE}/api/v1/auth/oidc/start`;
  }

  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-gradient-to-b from-slate-50 to-white p-4">
      <Card className="w-full max-w-md border-primary/20 shadow-lg">
        <CardHeader>
          <CardTitle className="text-primary">Sales-Pilot</CardTitle>
          <CardDescription>企业销售赋能系统 — 登录</CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          <Tabs defaultValue="local">
            <TabsList className="grid w-full grid-cols-2">
              <TabsTrigger value="local">本地账号</TabsTrigger>
              <TabsTrigger value="ldap">LDAP</TabsTrigger>
            </TabsList>
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
          </Tabs>

          <div className="relative">
            <div className="absolute inset-0 flex items-center">
              <span className="w-full border-t" />
            </div>
            <div className="relative flex justify-center text-xs uppercase">
              <span className="bg-card px-2 text-muted-foreground">
                或
              </span>
            </div>
          </div>

          <Button
            type="button"
            variant="outline"
            className="w-full"
            onClick={startOidc}
          >
            使用企业单点登录 (SSO)
          </Button>
          <p className="text-center text-xs text-muted-foreground">
            演示账号：<code className="rounded bg-muted px-1">admin</code> /{" "}
            <code className="rounded bg-muted px-1">admin123</code>
            <br />
            <Link href="/products" className="text-primary underline">
              返回首页
            </Link>
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
