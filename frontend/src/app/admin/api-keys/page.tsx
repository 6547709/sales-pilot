import { PageHeader } from "@/components/patterns/page-header";
import { ApiKeysDataGrid } from "@/features/api-keys/api-keys-data-grid";

export default function AdminApiKeysPage() {
  return (
    <div className="flex flex-col gap-6">
      <PageHeader
        title="API 密钥"
        description="供自动化脚本与 MCP 使用；密钥等价于管理员权限，请妥善保管"
      />
      <ApiKeysDataGrid />
    </div>
  );
}
