import { PageHeader } from "@/components/patterns/page-header";
import { TopologyEditor } from "@/features/topology/topology-editor";

export default function AdminTopologyPage() {
  return (
    <div className="mx-auto flex max-w-5xl flex-col gap-6">
      <PageHeader
        title="拓扑维护"
        description="管理首页架构分层、解决方案卡片与厂商（国内 / 国外）。停用方案仍可在后台编辑，但首页不展示。"
      />
      <TopologyEditor />
    </div>
  );
}
