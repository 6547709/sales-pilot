package initialize

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/business"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"gorm.io/gorm"
)

// initBusinessData 业务数据初始化（拓扑、解决方案分类等）
type initBusinessData struct{}

func (i *initBusinessData) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	var count int64
	db.Model(&business.TopologyLayer{}).Count(&count)
	return count > 0
}

func (i *initBusinessData) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	// 1. 创建拓扑层级（5层）
	layers := []business.TopologyLayer{
		{Name: "physical", SortOrder: 0, Level: 1, Title: "物理环境层", Subtitle: "Physical Environment"},
		{Name: "iaas", SortOrder: 1, Level: 2, Title: "云平台层·基础设施", Subtitle: "Infrastructure as a Service"},
		{Name: "paas", SortOrder: 2, Level: 3, Title: "云平台层·数据服务", Subtitle: "Platform as a Service"},
		{Name: "saas", SortOrder: 3, Level: 4, Title: "应用服务层", Subtitle: "Software as a Service"},
		{Name: "business", SortOrder: 4, Level: 5, Title: "业务层", Subtitle: "Business Layer"},
	}
	if err = db.Create(&layers).Error; err != nil {
		return ctx, err
	}

	// 建立层级ID映射（按 level 映射）
	layerMap := make(map[int]uint)
	for _, l := range layers {
		layerMap[l.Level] = l.ID
	}

	// 2. 创建拓扑分类（34个）
	categories := []business.TopologyCategory{
		// 安全体系 (column_type=security)
		{ColumnType: "security", Name: "sec-host", Slug: "sec-host", Label: "主机安全", IconKey: "ServerCog", Keywords: `["主机","EDR","安全代理"]`, SortOrder: 0, IsActive: true},
		{ColumnType: "security", Name: "sec-boundary", Slug: "sec-boundary", Label: "边界安全", IconKey: "Shield", Keywords: `["边界","防火墙","WAF","南北向"]`, SortOrder: 1, IsActive: true},
		{ColumnType: "security", Name: "sec-secret", Slug: "sec-secret", Label: "机密管理", IconKey: "Lock", Keywords: `["机密","Vault","密钥","证书"]`, SortOrder: 2, IsActive: true},
		{ColumnType: "security", Name: "sec-identity", Slug: "sec-identity", Label: "身份认证", IconKey: "Fingerprint", Keywords: `["身份","IAM","SSO","零信任"]`, SortOrder: 3, IsActive: true},
		{ColumnType: "security", Name: "sec-zerotrust", Slug: "sec-zerotrust", Label: "零信任", IconKey: "ShieldCheck", Keywords: `["零信任","ZTNA","微隔离"]`, SortOrder: 4, IsActive: true},
		{ColumnType: "security", Name: "sec-encrypt", Slug: "sec-encrypt", Label: "数据加密", IconKey: "Lock", Keywords: `["加密","KMS","脱敏"]`, SortOrder: 5, IsActive: true},
		{ColumnType: "security", Name: "sec-situation", Slug: "sec-situation", Label: "态势感知", IconKey: "Radar", Keywords: `["态势","SOC","SIEM","XDR"]`, SortOrder: 6, IsActive: true},
		// 运维体系 (column_type=ops)
		{ColumnType: "ops", Name: "ops-aiops", Slug: "ops-aiops", Label: "AIOps", IconKey: "Brain", Keywords: `["AIOps","智能运维","根因"]`, SortOrder: 0, IsActive: true},
		{ColumnType: "ops", Name: "ops-cmdb", Slug: "ops-cmdb", Label: "CMDB", IconKey: "Database", Keywords: `["CMDB","配置项","资产"]`, SortOrder: 1, IsActive: true},
		{ColumnType: "ops", Name: "ops-probe", Slug: "ops-probe", Label: "拨测", IconKey: "Activity", Keywords: `["拨测","可用性","探测"]`, SortOrder: 2, IsActive: true},
		{ColumnType: "ops", Name: "ops-trace", Slug: "ops-trace", Label: "全链路监控", IconKey: "LineChart", Keywords: `["链路","APM","Tracing"]`, SortOrder: 3, IsActive: true},
		{ColumnType: "ops", Name: "ops-log", Slug: "ops-log", Label: "日志分析", IconKey: "FileStack", Keywords: `["日志","ELK","可观测"]`, SortOrder: 4, IsActive: true},
		{ColumnType: "ops", Name: "ops-auto", Slug: "ops-auto", Label: "自动化运维", IconKey: "Workflow", Keywords: `["自动化","编排","Ansible"]`, SortOrder: 5, IsActive: true},
		{ColumnType: "ops", Name: "ops-dr", Slug: "ops-dr", Label: "灾备", IconKey: "CloudCog", Keywords: `["灾备","RTO","RPO","双活"]`, SortOrder: 6, IsActive: true},
		// 中心层级 (column_type=central) 按 layer_id 分组
		// Level 5 - 云平台层·应用与服务 (layer_id=layerMap[5])
		{LayerID: layerMap[5], ColumnType: "central", Name: "app-web", Slug: "app-web", Label: "Web / 中间件", IconKey: "Server", Keywords: `["中间件","Web","应用服务器"]`, SortOrder: 0, IsActive: true},
		{LayerID: layerMap[5], ColumnType: "central", Name: "app-db", Slug: "app-db", Label: "数据库", IconKey: "Database", Keywords: `["数据库","PostgreSQL","MySQL"]`, SortOrder: 1, IsActive: true},
		{LayerID: layerMap[5], ColumnType: "central", Name: "app-ai-mw", Slug: "app-ai-mw", Label: "AI 中间件", IconKey: "Brain", Keywords: `["向量","Embedding","RAG"]`, SortOrder: 2, IsActive: true},
		{LayerID: layerMap[5], ColumnType: "central", Name: "app-agent", Slug: "app-agent", Label: "AI / Agent / MCP", IconKey: "Sparkles", Keywords: `["Agent","MCP","Dify","智能体"]`, SortOrder: 3, IsActive: true},
		// Level 4 - 云平台层·基础设施 (layer_id=layerMap[4])
		{LayerID: layerMap[4], ColumnType: "central", Name: "cloud-virt", Slug: "cloud-virt", Label: "虚拟化", IconKey: "Box", Keywords: `["虚拟化","VMware","迁移","KubeVirt"]`, SortOrder: 0, IsActive: true},
		{LayerID: layerMap[4], ColumnType: "central", Name: "cloud-k8s", Slug: "cloud-k8s", Label: "容器云", IconKey: "LayoutGrid", Keywords: `["容器","Kubernetes","OpenShift","K8s"]`, SortOrder: 1, IsActive: true},
		{LayerID: layerMap[4], ColumnType: "central", Name: "cloud-storage", Slug: "cloud-storage", Label: "存储", IconKey: "HardDrive", Keywords: `["存储","Ceph","SDS","分布式存储"]`, SortOrder: 2, IsActive: true},
		{LayerID: layerMap[4], ColumnType: "central", Name: "cloud-network", Slug: "cloud-network", Label: "网络", IconKey: "Network", Keywords: `["网络","SDN","VPC","Underlay"]`, SortOrder: 3, IsActive: true},
		// Level 3 - 云平台层·数据服务 (layer_id=layerMap[3])
		{LayerID: layerMap[3], ColumnType: "central", Name: "data-bigdata", Slug: "data-bigdata", Label: "大数据平台", IconKey: "BarChart", Keywords: `["大数据","Hadoop","Spark","Flink"]`, SortOrder: 0, IsActive: true},
		{LayerID: layerMap[3], ColumnType: "central", Name: "data-dw", Slug: "data-dw", Label: "数据仓库", IconKey: "Database", Keywords: `["数据仓库","Hive","ClickHouse","StarRocks"]`, SortOrder: 1, IsActive: true},
		{LayerID: layerMap[3], ColumnType: "central", Name: "data-ml", Slug: "data-ml", Label: "机器学习", IconKey: "Brain", Keywords: `["机器学习","TensorFlow","PyTorch","MLOps"]`, SortOrder: 2, IsActive: true},
		// Level 2 - 支撑层 (layer_id=layerMap[2])
		{LayerID: layerMap[2], ColumnType: "central", Name: "support-api", Slug: "support-api", Label: "API 管理", IconKey: "Gateway", Keywords: `["API","网关","Kong","APISIX"]`, SortOrder: 0, IsActive: true},
		{LayerID: layerMap[2], ColumnType: "central", Name: "support-devops", Slug: "support-devops", Label: "DevOps 平台", IconKey: "Workflow", Keywords: `["DevOps","CI/CD","Jenkins","GitLab"]`, SortOrder: 1, IsActive: true},
		{LayerID: layerMap[2], ColumnType: "central", Name: "support-test", Slug: "support-test", Label: "测试管理", IconKey: "CheckCircle", Keywords: `["测试","QA","自动化测试","Jest"]`, SortOrder: 2, IsActive: true},
		// Level 1 - 接入层 (layer_id=layerMap[1])
		{LayerID: layerMap[1], ColumnType: "central", Name: "access-cdn", Slug: "access-cdn", Label: "CDN", IconKey: "Globe", Keywords: `["CDN","加速","边缘计算","DNS"]`, SortOrder: 0, IsActive: true},
		{LayerID: layerMap[1], ColumnType: "central", Name: "access-dns", Slug: "access-dns", Label: "DNS", IconKey: "Globe", Keywords: `["DNS","域名解析","权威DNS","HTTPDNS"]`, SortOrder: 1, IsActive: true},
		{LayerID: layerMap[1], ColumnType: "central", Name: "access-gateway", Slug: "access-gateway", Label: "负载均衡", IconKey: "Server", Keywords: `["负载均衡","SLB","Nginx","Envoy"]`, SortOrder: 2, IsActive: true},
	}
	if err = db.Create(&categories).Error; err != nil {
		return ctx, err
	}

	global.GVA_LOG.Info("业务数据初始化成功")
	next = context.WithValue(ctx, "initBusinessData", true)
	return next, nil
}

// InitBusinessData 自动初始化业务数据（如果不存在）
func InitBusinessData() {
	if global.GVA_DB == nil {
		return
	}
	ctx := context.WithValue(context.Background(), "db", global.GVA_DB)
	init := &initBusinessData{}
	if !init.DataInserted(ctx) {
		global.GVA_LOG.Info("开始初始化业务数据...")
		init.InitializeData(ctx)
	}
}
