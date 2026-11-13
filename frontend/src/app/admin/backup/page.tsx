import { PageHeader } from "@/components/patterns/page-header";
import { BackupView } from "@/features/backup/backup-view";

export default function AdminBackupPage() {
  return (
    <div className="mx-auto flex max-w-3xl flex-col gap-6">
      <PageHeader
        title="备份与恢复"
        description="导出 JSON 备份；恢复会整库覆盖，请先在安全环境验证文件"
      />
      <BackupView />
    </div>
  );
}
