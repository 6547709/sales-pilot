import Link from "next/link";
import { Card, CardHeader, CardTitle } from "@/components/ui/card";

export default function AdminHome() {
  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-primary">管理后台</h1>
      <div className="grid gap-4 sm:grid-cols-2">
        <Link href="/admin/products">
          <Card className="h-full transition hover:border-primary/40">
            <CardHeader>
              <CardTitle>产品编辑</CardTitle>
              <p className="text-sm text-muted-foreground">
                Markdown 预览、维护赋能内容
              </p>
            </CardHeader>
          </Card>
        </Link>
        <Link href="/admin/settings">
          <Card className="h-full transition hover:border-primary/40">
            <CardHeader>
              <CardTitle>认证配置</CardTitle>
              <p className="text-sm text-muted-foreground">
                LDAP / OIDC 参数
              </p>
            </CardHeader>
          </Card>
        </Link>
      </div>
    </div>
  );
}
