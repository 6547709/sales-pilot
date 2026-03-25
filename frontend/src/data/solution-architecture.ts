/**
 * 全局解决方案架构图数据（与业务展示图对应）
 * keywords 用于与产品 category/name 做模糊匹配筛选
 */
import type { LucideIcon } from "lucide-react";
import {
  Activity,
  Bot,
  Box,
  Brain,
  Building2,
  Cable,
  Cloud,
  CloudCog,
  Database,
  FileStack,
  Fingerprint,
  Globe2,
  HardDrive,
  Layers3,
  LayoutGrid,
  LineChart,
  Lock,
  Network,
  Radar,
  Router,
  Scale,
  Server,
  ServerCog,
  Shield,
  ShieldCheck,
  Sparkles,
  ThermometerSnowflake,
  Workflow,
  Zap,
} from "lucide-react";

export type ArchCard = {
  id: string;
  label: string;
  icon: LucideIcon;
  /** 用于匹配后台产品分类/名称 */
  keywords: string[];
  hint?: string;
  domestic?: string[];
  foreign?: string[];
};

export type ArchLayer = {
  id: string;
  level: number;
  title: string;
  subtitle?: string;
  items: ArchCard[];
};

/** 左侧：安全体系 */
export const securityColumn: ArchCard[] = [
  {
    id: "sec-host",
    label: "主机安全",
    icon: ServerCog,
    keywords: ["主机", "EDR", "安全代理"],
  },
  {
    id: "sec-boundary",
    label: "边界安全",
    icon: Shield,
    keywords: ["边界", "防火墙", "WAF", "南北向"],
  },
  {
    id: "sec-secret",
    label: "机密管理",
    icon: Lock,
    keywords: ["机密", "Vault", "密钥", "证书"],
  },
  {
    id: "sec-identity",
    label: "身份认证",
    icon: Fingerprint,
    keywords: ["身份", "IAM", "SSO", "零信任"],
  },
  {
    id: "sec-zerotrust",
    label: "零信任",
    icon: ShieldCheck,
    keywords: ["零信任", "ZTNA", "微隔离"],
  },
  {
    id: "sec-encrypt",
    label: "数据加密",
    icon: Lock,
    keywords: ["加密", "KMS", "脱敏"],
  },
  {
    id: "sec-situation",
    label: "态势感知",
    icon: Radar,
    keywords: ["态势", "SOC", "SIEM", "XDR"],
  },
];

/** 右侧：运维体系 */
export const opsColumn: ArchCard[] = [
  {
    id: "ops-aiops",
    label: "AIOps",
    icon: Brain,
    keywords: ["AIOps", "智能运维", "根因"],
  },
  {
    id: "ops-cmdb",
    label: "CMDB",
    icon: Database,
    keywords: ["CMDB", "配置项", "资产"],
  },
  {
    id: "ops-probe",
    label: "拨测",
    icon: Activity,
    keywords: ["拨测", "可用性", "探测"],
  },
  {
    id: "ops-trace",
    label: "全链路监控",
    icon: LineChart,
    keywords: ["链路", "APM", "Tracing"],
  },
  {
    id: "ops-log",
    label: "日志分析",
    icon: FileStack,
    keywords: ["日志", "ELK", "可观测"],
  },
  {
    id: "ops-auto",
    label: "自动化运维",
    icon: Workflow,
    keywords: ["自动化", "编排", "Ansible"],
  },
  {
    id: "ops-dr",
    label: "灾备",
    icon: CloudCog,
    keywords: ["灾备", "RTO", "RPO", "双活"],
  },
];

/** 中部：分层栈（自下而上展示时 level 从 1 到 5） */
export const centralLayers: ArchLayer[] = [
  {
    id: "layer-physical",
    level: 1,
    title: "物理环境层",
    subtitle: "机房基础设施",
    items: [
      {
        id: "phy-dc",
        label: "数据中心",
        icon: Building2,
        keywords: ["数据中心", "机房", "IDC"],
      },
      {
        id: "phy-ups",
        label: "UPS / PDU",
        icon: Zap,
        keywords: ["UPS", "PDU", "供电"],
      },
      {
        id: "phy-hvac",
        label: "精密空调",
        icon: ThermometerSnowflake,
        keywords: ["空调", "制冷", "PUE"],
      },
      {
        id: "phy-cable",
        label: "综合布线",
        icon: Cable,
        keywords: ["布线", "光纤", "铜缆"],
      },
    ],
  },
  {
    id: "layer-hardware",
    level: 2,
    title: "基础硬件层",
    subtitle: "算力与存储",
    items: [
      {
        id: "hw-compute",
        label: "计算",
        icon: Server,
        keywords: ["服务器", "计算", "x86"],
        domestic: ["H3C", "浪潮", "华为"],
        foreign: ["Dell", "HPE"],
      },
      {
        id: "hw-san",
        label: "SAN 存储",
        icon: HardDrive,
        keywords: ["SAN", "块存储"],
      },
      {
        id: "hw-obj",
        label: "文件 / 对象",
        icon: Layers3,
        keywords: ["对象存储", "NAS", "文件"],
      },
      {
        id: "hw-gpu",
        label: "GPU",
        icon: Sparkles,
        keywords: ["GPU", "算力", "推理", "训练"],
        domestic: ["华为", "寒武纪"],
        foreign: ["Nvidia", "AMD"],
      },
    ],
  },
  {
    id: "layer-network",
    level: 3,
    title: "网络与连接层",
    subtitle: "连通与调度",
    items: [
      {
        id: "net-core",
        label: "核心网络",
        icon: Network,
        keywords: ["核心网", "骨干", "路由"],
      },
      {
        id: "net-sdn",
        label: "SDN / SD-WAN",
        icon: Router,
        keywords: ["SDN", "SD-WAN", "广域网"],
      },
      {
        id: "net-lb",
        label: "负载均衡",
        icon: Scale,
        keywords: ["负载均衡", "F5", "Ingress"],
        foreign: ["F5"],
      },
      {
        id: "net-dns",
        label: "DNS / IPAM",
        icon: Globe2,
        keywords: ["DNS", "IPAM", "地址管理"],
      },
    ],
  },
  {
    id: "layer-cloud-infra",
    level: 4,
    title: "云平台层 · 基础设施",
    subtitle: "虚拟化与容器底座",
    items: [
      {
        id: "cloud-virt",
        label: "虚拟化",
        icon: Box,
        keywords: ["虚拟化", "VMware", "迁移", "KubeVirt"],
        domestic: ["SmartX", "华为"],
        foreign: ["VMware"],
      },
      {
        id: "cloud-k8s",
        label: "容器云",
        icon: LayoutGrid,
        keywords: ["容器", "Kubernetes", "OpenShift", "K8s"],
      },
      {
        id: "cloud-ai",
        label: "AI 云",
        icon: Bot,
        keywords: ["AI 云", "模型服务", "MLOps"],
      },
      {
        id: "cloud-multi",
        label: "多云管理",
        icon: Cloud,
        keywords: ["多云", "混合云", "CMP"],
      },
    ],
  },
  {
    id: "layer-cloud-app",
    level: 5,
    title: "云平台层 · 应用与服务",
    subtitle: "中间件与智能化",
    items: [
      {
        id: "app-web",
        label: "Web / 中间件",
        icon: Server,
        keywords: ["中间件", "Web", "应用服务器"],
      },
      {
        id: "app-db",
        label: "数据库",
        icon: Database,
        keywords: ["数据库", "PostgreSQL", "MySQL"],
      },
      {
        id: "app-ai-mw",
        label: "AI 中间件",
        icon: Brain,
        keywords: ["向量", "Embedding", "RAG"],
      },
      {
        id: "app-agent",
        label: "AI / Agent / MCP",
        icon: Sparkles,
        keywords: ["Agent", "MCP", "智能体", "工作流"],
      },
    ],
  },
];

/** 扁平化所有卡片，便于搜索高亮 */
export function allArchCards(): ArchCard[] {
  return [
    ...securityColumn,
    ...opsColumn,
    ...centralLayers.flatMap((l) => l.items),
  ];
}
