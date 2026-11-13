import { PageHeader } from "@/components/patterns/page-header";
import { AdminScenariosView } from "@/features/scenarios/admin-scenarios-view";

export default function AdminScenariosPage() {
  return (
    <div className="flex flex-col gap-6">
      <PageHeader
        title="场景维护"
        description="两级分类：场景组下维护场景，并关联产品与服务"
      />
      <AdminScenariosView />
    </div>
  );
}
