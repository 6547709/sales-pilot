import { PageHeader } from "@/components/patterns/page-header";
import { AuthSettingsForm } from "@/features/settings/auth-settings-form";
import { SystemSettingsForm } from "@/features/settings/system-settings-form";

export default function AdminSettingsPage() {
  return (
    <div className="mx-auto flex max-w-3xl flex-col gap-6">
      <PageHeader title="认证配置" description="登录方式与站点页脚信息" />
      <AuthSettingsForm />
      <SystemSettingsForm />
    </div>
  );
}
