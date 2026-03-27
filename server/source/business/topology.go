package business

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/model/business"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderBusiness = system.InitOrderInternal + 1

type initBusiness struct{}

// auto run
func init() {
	system.RegisterInit(initOrderBusiness, &initBusiness{})
}

func (i *initBusiness) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(
		&business.SolutionCategory{},
		&business.Product{},
		&business.SalesScript{},
		&business.Case{},
		&business.TopologyLayer{},
		&business.TopologyCategory{},
		&business.TopologyVendor{},
	)
}

func (i *initBusiness) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&business.TopologyLayer{})
}

func (i *initBusiness) InitializerName() string {
	return "sys_business_topology"
}

func (i *initBusiness) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	// 检查是否已有数据
	var count int64
	db.Model(&business.TopologyLayer{}).Count(&count)
	if count > 0 {
		return ctx, nil
	}

	// 1. 创建拓扑层级（5层，从上到下 level 5→1）
	layers := []business.TopologyLayer{
		{Level: 5, Title: "云平台层 · 应用与服务", Subtitle: "中间件与智能化", SortOrder: 0},
		{Level: 4, Title: "云平台层 · 基础设施", Subtitle: "虚拟化与容器底座", SortOrder: 1},
		{Level: 3, Title: "网络与连接层", Subtitle: "连通与调度", SortOrder: 2},
		{Level: 2, Title: "基础硬件层", Subtitle: "算力与存储", SortOrder: 3},
		{Level: 1, Title: "物理环境层", Subtitle: "机房基础设施", SortOrder: 4},
	}

	if err = db.Create(&layers).Error; err != nil {
		return ctx, errors.Wrap(err, "拓扑层级初始化失败")
	}

	// 建立层级ID映射（按 level 映射）
	layerMap := make(map[int]uint)
	for _, l := range layers {
		layerMap[l.Level] = l.ID
	}

	// 2. 创建拓扑分类（34个）
	categories := []business.TopologyCategory{
		// 安全体系 (column_type=security) id 1-7
		{ColumnType: "security", Name: "sec-host", Slug: "sec-host", Label: "主机安全", IconKey: "ServerCog", Keywords: "[\"主机\",\"EDR\",\"安全代理\"]", SortOrder: 0, IsActive: true},
		{ColumnType: "security", Name: "sec-boundary", Slug: "sec-boundary", Label: "边界安全", IconKey: "Shield", Keywords: "[\"边界\",\"防火墙\",\"WAF\",\"南北向\"]", SortOrder: 1, IsActive: true},
		{ColumnType: "security", Name: "sec-secret", Slug: "sec-secret", Label: "机密管理", IconKey: "Lock", Keywords: "[\"机密\",\"Vault\",\"密钥\",\"证书\"]", SortOrder: 2, IsActive: true},
		{ColumnType: "security", Name: "sec-identity", Slug: "sec-identity", Label: "身份认证", IconKey: "Fingerprint", Keywords: "[\"身份\",\"IAM\",\"SSO\",\"零信任\"]", SortOrder: 3, IsActive: true},
		{ColumnType: "security", Name: "sec-zerotrust", Slug: "sec-zerotrust", Label: "零信任", IconKey: "ShieldCheck", Keywords: "[\"零信任\",\"ZTNA\",\"微隔离\"]", SortOrder: 4, IsActive: true},
		{ColumnType: "security", Name: "sec-encrypt", Slug: "sec-encrypt", Label: "数据加密", IconKey: "Lock", Keywords: "[\"加密\",\"KMS\",\"脱敏\"]", SortOrder: 5, IsActive: true},
		{ColumnType: "security", Name: "sec-situation", Slug: "sec-situation", Label: "态势感知", IconKey: "Radar", Keywords: "[\"态势\",\"SOC\",\"SIEM\",\"XDR\"]", SortOrder: 6, IsActive: true},

		// 运维体系 (column_type=ops) id 8-14
		{ColumnType: "ops", Name: "ops-aiops", Slug: "ops-aiops", Label: "AIOps", IconKey: "Brain", Keywords: "[\"AIOps\",\"智能运维\",\"根因\"]", SortOrder: 0, IsActive: true},
		{ColumnType: "ops", Name: "ops-cmdb", Slug: "ops-cmdb", Label: "CMDB", IconKey: "Database", Keywords: "[\"CMDB\",\"配置项\",\"资产\"]", SortOrder: 1, IsActive: true},
		{ColumnType: "ops", Name: "ops-probe", Slug: "ops-probe", Label: "拨测", IconKey: "Activity", Keywords: "[\"拨测\",\"可用性\",\"探测\"]", SortOrder: 2, IsActive: true},
		{ColumnType: "ops", Name: "ops-trace", Slug: "ops-trace", Label: "全链路监控", IconKey: "LineChart", Keywords: "[\"链路\",\"APM\",\"Tracing\"]", SortOrder: 3, IsActive: true},
		{ColumnType: "ops", Name: "ops-log", Slug: "ops-log", Label: "日志分析", IconKey: "FileStack", Keywords: "[\"日志\",\"ELK\",\"可观测\"]", SortOrder: 4, IsActive: true},
		{ColumnType: "ops", Name: "ops-auto", Slug: "ops-auto", Label: "自动化运维", IconKey: "Workflow", Keywords: "[\"自动化\",\"编排\",\"Ansible\"]", SortOrder: 5, IsActive: true},
		{ColumnType: "ops", Name: "ops-dr", Slug: "ops-dr", Label: "灾备", IconKey: "CloudCog", Keywords: "[\"灾备\",\"RTO\",\"RPO\",\"双活\"]", SortOrder: 6, IsActive: true},

		// 中心层级 (column_type=central) id 15-34，按 layer_id 分组
		// Level 5 - 云平台层·应用与服务 (layer_id=1)
		{LayerID: layerMap[5], ColumnType: "central", Name: "app-web", Slug: "app-web", Label: "Web / 中间件", IconKey: "Server", Keywords: "[\"中间件\",\"Web\",\"应用服务器\"]", SortOrder: 0, IsActive: true},
		{LayerID: layerMap[5], ColumnType: "central", Name: "app-db", Slug: "app-db", Label: "数据库", IconKey: "Database", Keywords: "[\"数据库\",\"PostgreSQL\",\"MySQL\"]", SortOrder: 1, IsActive: true},
		{LayerID: layerMap[5], ColumnType: "central", Name: "app-ai-mw", Slug: "app-ai-mw", Label: "AI 中间件", IconKey: "Brain", Keywords: "[\"向量\",\"Embedding\",\"RAG\"]", SortOrder: 2, IsActive: true},
		{LayerID: layerMap[5], ColumnType: "central", Name: "app-agent", Slug: "app-agent", Label: "AI / Agent / MCP", IconKey: "Sparkles", Keywords: "[\"Agent\",\"MCP\",\"Dify\",\"智能体\"]", SortOrder: 3, IsActive: true},

		// Level 4 - 云平台层·基础设施 (layer_id=2)
		{LayerID: layerMap[4], ColumnType: "central", Name: "cloud-virt", Slug: "cloud-virt", Label: "虚拟化", IconKey: "Box", Keywords: "[\"虚拟化\",\"VMware\",\"迁移\",\"KubeVirt\"]", SortOrder: 0, IsActive: true},
		{LayerID: layerMap[4], ColumnType: "central", Name: "cloud-k8s", Slug: "cloud-k8s", Label: "容器云", IconKey: "LayoutGrid", Keywords: "[\"容器\",\"Kubernetes\",\"OpenShift\",\"K8s\"]", SortOrder: 1, IsActive: true},
		{LayerID: layerMap[4], ColumnType: "central", Name: "cloud-ai", Slug: "cloud-ai", Label: "AI 云", IconKey: "Bot", Keywords: "[\"AI 云\",\"模型服务\",\"MLOps\"]", SortOrder: 2, IsActive: true},
		{LayerID: layerMap[4], ColumnType: "central", Name: "cloud-multi", Slug: "cloud-multi", Label: "多云管理", IconKey: "Cloud", Keywords: "[\"多云\",\"混合云\",\"CMP\"]", SortOrder: 3, IsActive: true},

		// Level 3 - 网络与连接层 (layer_id=3)
		{LayerID: layerMap[3], ColumnType: "central", Name: "net-core", Slug: "net-core", Label: "核心网络", IconKey: "Network", Keywords: "[\"核心网\",\"骨干\",\"路由\"]", SortOrder: 0, IsActive: true},
		{LayerID: layerMap[3], ColumnType: "central", Name: "net-sdn", Slug: "net-sdn", Label: "SDN / SD-WAN", IconKey: "Router", Keywords: "[\"SDN\",\"SD-WAN\",\"广域网\"]", SortOrder: 1, IsActive: true},
		{LayerID: layerMap[3], ColumnType: "central", Name: "net-lb", Slug: "net-lb", Label: "负载均衡", IconKey: "Scale", Keywords: "[\"负载均衡\",\"F5\",\"Ingress\"]", SortOrder: 2, IsActive: true},
		{LayerID: layerMap[3], ColumnType: "central", Name: "net-dns", Slug: "net-dns", Label: "DNS / IPAM", IconKey: "Globe2", Keywords: "[\"DNS\",\"IPAM\",\"地址管理\"]", SortOrder: 3, IsActive: true},

		// Level 2 - 基础硬件层 (layer_id=4)
		{LayerID: layerMap[2], ColumnType: "central", Name: "hw-compute", Slug: "hw-compute", Label: "计算", IconKey: "Server", Keywords: "[\"服务器\",\"计算\",\"x86\"]", SortOrder: 0, IsActive: true},
		{LayerID: layerMap[2], ColumnType: "central", Name: "hw-san", Slug: "hw-san", Label: "SAN 存储", IconKey: "HardDrive", Keywords: "[\"SAN\",\"块存储\"]", SortOrder: 1, IsActive: true},
		{LayerID: layerMap[2], ColumnType: "central", Name: "hw-obj", Slug: "hw-obj", Label: "文件 / 对象", IconKey: "Layers3", Keywords: "[\"对象存储\",\"NAS\",\"文件\"]", SortOrder: 2, IsActive: true},
		{LayerID: layerMap[2], ColumnType: "central", Name: "hw-gpu", Slug: "hw-gpu", Label: "GPU", IconKey: "Sparkles", Keywords: "[\"GPU\",\"算力\",\"推理\",\"训练\"]", SortOrder: 3, IsActive: true},

		// Level 1 - 物理环境层 (layer_id=5)
		{LayerID: layerMap[1], ColumnType: "central", Name: "phy-dc", Slug: "phy-dc", Label: "数据中心", IconKey: "Building2", Keywords: "[\"数据中心\",\"机房\",\"IDC\"]", SortOrder: 0, IsActive: true},
		{LayerID: layerMap[1], ColumnType: "central", Name: "phy-ups", Slug: "phy-ups", Label: "UPS / PDU", IconKey: "Zap", Keywords: "[\"UPS\",\"PDU\",\"供电\"]", SortOrder: 1, IsActive: true},
		{LayerID: layerMap[1], ColumnType: "central", Name: "phy-hvac", Slug: "phy-hvac", Label: "精密空调", IconKey: "ThermometerSnowflake", Keywords: "[\"空调\",\"制冷\",\"PUE\"]", SortOrder: 2, IsActive: true},
		{LayerID: layerMap[1], ColumnType: "central", Name: "phy-cable", Slug: "phy-cable", Label: "综合布线", IconKey: "Cable", Keywords: "[\"布线\",\"光纤\",\"铜缆\"]", SortOrder: 3, IsActive: true},
	}

	if err = db.Create(&categories).Error; err != nil {
		return ctx, errors.Wrap(err, "拓扑分类初始化失败")
	}

	// 3. 创建解决方案分类
	solutionCategories := []business.SolutionCategory{
		{Label: "安全合规"},
		{Label: "云计算与基础设施"},
		{Label: "数据管理与分析"},
		{Label: "应用现代化"},
		{Label: "DevOps 与自动化"},
		{Label: "运维可观测性"},
	}

	if err = db.Create(&solutionCategories).Error; err != nil {
		return ctx, errors.Wrap(err, "解决方案分类初始化失败")
	}

	next = context.WithValue(ctx, i.InitializerName(), true)
	return next, nil
}

func (i *initBusiness) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	var count int64
	db.Model(&business.TopologyLayer{}).Count(&count)
	return count > 0
}
