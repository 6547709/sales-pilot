import { PageHeader } from "@/components/patterns/page-header";
import { UsersDataGrid } from "@/features/users/users-data-grid";

export default function AdminUsersPage() {
  return (
    <div className="flex flex-col gap-6">
      <PageHeader
        title="本地用户"
        description="仅管理 auth_provider=local 的账号；LDAP/OIDC 用户不在此列表"
      />
      <UsersDataGrid />
    </div>
  );
}
