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

	// 1. 创建拓扑层级
	layers := []business.TopologyLayer{
		{SortOrder: 1, Level: 1, Name: "身份与访问", Title: "身份与访问", Subtitle: "身份认证、权限管理、合规审计"},
		{SortOrder: 2, Level: 2, Name: "安全边界", Title: "安全边界", Subtitle: "防火墙、入侵检测、零信任"},
		{SortOrder: 3, Level: 3, Name: "计算资源", Title: "计算资源", Subtitle: "服务器、容器、Serverless"},
		{SortOrder: 4, Level: 4, Name: "存储资源", Title: "存储资源", Subtitle: "对象存储、块存储、文件存储"},
		{SortOrder: 5, Level: 5, Name: "数据层", Title: "数据层", Subtitle: "数据库、大数据、数据湖"},
		{SortOrder: 6, Level: 6, Name: "应用层", Title: "应用层", Subtitle: "应用服务、中间件、API"},
		{SortOrder: 7, Level: 7, Name: "运维体系", Title: "运维体系", Subtitle: "监控、自动化、运维工具"},
	}

	if err = db.Create(&layers).Error; err != nil {
		return ctx, errors.Wrap(err, "拓扑层级初始化失败")
	}

	// 建立层级ID映射
	layerMap := make(map[string]uint)
	for _, l := range layers {
		layerMap[l.Name] = l.ID
	}

	// 2. 创建拓扑分类
	categories := []business.TopologyCategory{
		// 身份与访问
		{LayerID: layerMap["身份与访问"], Name: "IAM", Slug: "iam", Label: "身份与访问管理", IconKey: "security", Keywords: "[\"身份\",\"权限\",\"SSO\",\"MFA\"]", Hint: "统一身份认证与权限管理"},
		{LayerID: layerMap["身份与访问"], Name: "IDaaS", Slug: "idaas", Label: "云身份服务", IconKey: "cloud", Keywords: "[\"IDaaS\",\"云身份\",\"OIDC\",\"SAML\"]", Hint: "云原生身份即服务"},
		{LayerID: layerMap["身份与访问"], Name: "PAM", Slug: "pam", Label: "特权身份管理", IconKey: "key", Keywords: "[\"特权账号\",\"堡垒机\",\"密钥\"]", Hint: "特权账号与密钥管理"},

		// 安全边界
		{LayerID: layerMap["安全边界"], Name: "Firewall", Slug: "firewall", Label: "防火墙", IconKey: "shield", Keywords: "[\"防火墙\",\"WAF\",\"NGFW\"]", Hint: "网络与应用防火墙"},
		{LayerID: layerMap["安全边界"], Name: "IDS/IPS", Slug: "ids-ips", Label: "入侵检测/防御", IconKey: "warning", Keywords: "[\"IDS\",\"IPS\",\"入侵检测\",\"威胁防护\"]", Hint: "入侵检测与防御系统"},
		{LayerID: layerMap["安全边界"], Name: "SASE", Slug: "sase", Label: "SASE", IconKey: "connection", Keywords: "[\"SASE\",\"零信任\",\"SD-WAN\"]", Hint: "安全访问服务边缘"},

		// 计算资源
		{LayerID: layerMap["计算资源"], Name: "Server", Slug: "server", Label: "服务器", IconKey: "server", Keywords: "[\"服务器\",\"物理机\",\"虚拟化\"]", Hint: "x86服务器与虚拟化平台"},
		{LayerID: layerMap["计算资源"], Name: "Container", Slug: "container", Label: "容器平台", IconKey: "box", Keywords: "[\"容器\",\"K8s\",\"Docker\",\"编排\"]", Hint: "Kubernetes与容器编排"},
		{LayerID: layerMap["计算资源"], Name: "Function", Slug: "function", Label: "Serverless", IconKey: "cloud", Keywords: "[\"Serverless\",\"函数计算\",\"FaaS\"]", Hint: "函数即服务"},

		// 存储资源
		{LayerID: layerMap["存储资源"], Name: "ObjectStorage", Slug: "object", Label: "对象存储", IconKey: "folder", Keywords: "[\"对象存储\",\"OSS\",\"S3\",\"MinIO\"]", Hint: "非结构化数据存储"},
		{LayerID: layerMap["存储资源"], Name: "BlockStorage", Slug: "block", Label: "块存储", IconKey: "database", Keywords: "[\"块存储\",\"云盘\",\"SAN\"]", Hint: "块设备与云硬盘"},
		{LayerID: layerMap["存储资源"], Name: "FileStorage", Slug: "file", Label: "文件存储", IconKey: "document", Keywords: "[\"文件存储\",\"NAS\",\"HDFS\"]", Hint: "文件存储与共享"},

		// 数据层
		{LayerID: layerMap["数据层"], Name: "RDBMS", Slug: "rdbms", Label: "关系数据库", IconKey: "database", Keywords: "[\"MySQL\",\"PostgreSQL\",\"Oracle\",\"SQL Server\"]", Hint: "关系型数据库"},
		{LayerID: layerMap["数据层"], Name: "NoSQL", Slug: "nosql", Label: "NoSQL数据库", IconKey: "collection", Keywords: "[\"MongoDB\",\"Redis\",\"Elasticsearch\",\"NoSQL\"]", Hint: "非关系型数据库"},
		{LayerID: layerMap["数据层"], Name: "DataLake", Slug: "datalake", Label: "数据湖", IconKey: "folder-opened", Keywords: "[\"数据湖\",\"Delta Lake\",\"Iceberg\",\"仓\"]", Hint: "大数据存储与分析"},

		// 应用层
		{LayerID: layerMap["应用层"], Name: "APIGateway", Slug: "api-gateway", Label: "API网关", IconKey: "connection", Keywords: "[\"API网关\",\"Kong\",\"网关\",\"路由\"]", Hint: "API管理与流量控制"},
		{LayerID: layerMap["应用层"], Name: "MQ", Slug: "mq", Label: "消息队列", IconKey: "chat-line-square", Keywords: "[\"Kafka\",\"RabbitMQ\",\"RocketMQ\",\"消息队列\"]", Hint: "异步消息与事件驱动"},
		{LayerID: layerMap["应用层"], Name: "JobScheduler", Slug: "scheduler", Label: "任务调度", IconKey: "timer", Keywords: "[\"定时任务\",\"调度\",\"Crontab\",\"XXL-Job\"]", Hint: "任务调度与定时任务"},

		// 运维体系
		{LayerID: layerMap["运维体系"], Name: "Monitor", Slug: "monitor", Label: "监控告警", IconKey: "monitor", Keywords: "[\"Prometheus\",\"Grafana\",\"Zabbix\",\"监控\"]", Hint: "监控与可视化"},
		{LayerID: layerMap["运维体系"], Name: "Log", Slug: "log", Label: "日志分析", IconKey: "document", Keywords: "[\"日志\",\"ELK\",\"Loki\",\"日志分析\"]", Hint: "日志采集与分析"},
		{LayerID: layerMap["运维体系"], Name: "CI/CD", Slug: "cicd", Label: "持续交付", IconKey: "refresh", Keywords: "[\"CI/CD\",\"Jenkins\",\"GitLab\",\"ArgoCD\"]", Hint: "持续集成与部署"},
		{LayerID: layerMap["运维体系"], Name: "IaC", Slug: "iac", Label: "基础设施即代码", IconKey: "setting", Keywords: "[\"Terraform\",\"Ansible\",\"IaC\",\"基础设施\"]", Hint: "基础设施代码化"},
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
