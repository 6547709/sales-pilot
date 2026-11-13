import { PageHeader } from "@/components/patterns/page-header";
import { ServicesDataGrid } from "@/features/services/services-data-grid";

export default function AdminServicesPage() {
  return (
    <div className="flex flex-col gap-6">
      <PageHeader
        title="服务维护"
        description="维护我方服务，并关联到首页架构方框与场景"
      />
      <ServicesDataGrid />
    </div>
  );
}
