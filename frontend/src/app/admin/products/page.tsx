import { PageHeader } from "@/components/patterns/page-header";
import { AdminProductsList } from "@/features/products/admin-products-list";

export default function AdminProductsPage() {
  return (
    <div className="flex flex-col gap-6">
      <PageHeader title="产品维护" description="编辑方案内容、发布状态与关联信息" />
      <AdminProductsList />
    </div>
  );
}
