package seed

import (
	"encoding/json"
	"log"

	"github.com/6547709/sales-pilot/backend/internal/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type layerDef struct {
	Level    int
	Title    string
	Subtitle string
	Ord      int
}

type catDef struct {
	Slug       string
	Label      string
	Icon       string
	Column     string // security | ops | central
	LayerLevel *int   // 与 TopologyLayer.Level 对应，侧栏为 nil
	Ord        int
	Keywords   []string
	Hint       string
	Vendors    []vendDef
}

func intPtr(v int) *int { return &v }

type vendDef struct {
	Market string // domestic | foreign
	Name   string
	Ord    int
}

// EnsureTopology 若尚无解决方案分类则写入默认拓扑（与初版前端静态数据一致）
func EnsureTopology(db *gorm.DB) {
	var n int64
	db.Model(&models.SolutionCategory{}).Count(&n)
	if n > 0 {
		return
	}
	log.Println("写入默认拓扑与厂商数据…")

	layers := []layerDef{
		{5, "云平台层 · 应用与服务", "中间件与智能化", 0},
		{4, "云平台层 · 基础设施", "虚拟化与容器底座", 1},
		{3, "网络与连接层", "连通与调度", 2},
		{2, "基础硬件层", "算力与存储", 3},
		{1, "物理环境层", "机房基础设施", 4},
	}
	layerIDByLevel := map[int]uint{}
	for _, ld := range layers {
		l := models.TopologyLayer{
			Level: ld.Level, Title: ld.Title, Subtitle: ld.Subtitle, SortOrder: ld.Ord,
		}
		if err := db.Create(&l).Error; err != nil {
			log.Fatal(err)
		}
		layerIDByLevel[ld.Level] = l.ID
	}

	cats := []catDef{
		// 安全
		{"sec-host", "主机安全", "ServerCog", "security", nil, 0, []string{"主机", "EDR", "安全代理"}, "", nil},
		{"sec-boundary", "边界安全", "Shield", "security", nil, 1, []string{"边界", "防火墙", "WAF", "南北向"}, "", nil},
		{"sec-secret", "机密管理", "Lock", "security", nil, 2, []string{"机密", "Vault", "密钥", "证书"}, "", nil},
		{"sec-identity", "身份认证", "Fingerprint", "security", nil, 3, []string{"身份", "IAM", "SSO", "零信任"}, "", nil},
		{"sec-zerotrust", "零信任", "ShieldCheck", "security", nil, 4, []string{"零信任", "ZTNA", "微隔离"}, "", nil},
		{"sec-encrypt", "数据加密", "Lock", "security", nil, 5, []string{"加密", "KMS", "脱敏"}, "", nil},
		{"sec-situation", "态势感知", "Radar", "security", nil, 6, []string{"态势", "SOC", "SIEM", "XDR"}, "", nil},
		// 运维
		{"ops-aiops", "AIOps", "Brain", "ops", nil, 0, []string{"AIOps", "智能运维", "根因"}, "", nil},
		{"ops-cmdb", "CMDB", "Database", "ops", nil, 1, []string{"CMDB", "配置项", "资产"}, "", nil},
		{"ops-probe", "拨测", "Activity", "ops", nil, 2, []string{"拨测", "可用性", "探测"}, "", nil},
		{"ops-trace", "全链路监控", "LineChart", "ops", nil, 3, []string{"链路", "APM", "Tracing"}, "", nil},
		{"ops-log", "日志分析", "FileStack", "ops", nil, 4, []string{"日志", "ELK", "可观测"}, "", nil},
		{"ops-auto", "自动化运维", "Workflow", "ops", nil, 5, []string{"自动化", "编排", "Ansible"}, "", nil},
		{"ops-dr", "灾备", "CloudCog", "ops", nil, 6, []string{"灾备", "RTO", "RPO", "双活"}, "", nil},
		// 中部 L5
		{"app-web", "Web / 中间件", "Server", "central", intPtr(5), 0, []string{"中间件", "Web", "应用服务器"}, "", nil},
		{"app-db", "数据库", "Database", "central", intPtr(5), 1, []string{"数据库", "PostgreSQL", "MySQL"}, "", nil},
		{"app-ai-mw", "AI 中间件", "Brain", "central", intPtr(5), 2, []string{"向量", "Embedding", "RAG"}, "", nil},
		// 提示文案留空，由管理员在「拓扑维护」中按需填写，避免未录入时出现示例性描述
		{"app-agent", "AI / Agent / MCP", "Sparkles", "central", intPtr(5), 3, []string{"Agent", "MCP", "智能体", "工作流"}, "", nil},
		// L4
		{"cloud-virt", "虚拟化", "Box", "central", intPtr(4), 0, []string{"虚拟化", "VMware", "迁移", "KubeVirt"}, "", []vendDef{
			{"domestic", "SmartX", 0}, {"domestic", "华为", 1}, {"foreign", "VMware", 2},
		}},
		{"cloud-k8s", "容器云", "LayoutGrid", "central", intPtr(4), 1, []string{"容器", "Kubernetes", "OpenShift", "K8s"}, "", nil},
		{"cloud-ai", "AI 云", "Bot", "central", intPtr(4), 2, []string{"AI 云", "模型服务", "MLOps"}, "", nil},
		{"cloud-multi", "多云管理", "Cloud", "central", intPtr(4), 3, []string{"多云", "混合云", "CMP"}, "", nil},
		// L3
		{"net-core", "核心网络", "Network", "central", intPtr(3), 0, []string{"核心网", "骨干", "路由"}, "", nil},
		{"net-sdn", "SDN / SD-WAN", "Router", "central", intPtr(3), 1, []string{"SDN", "SD-WAN", "广域网"}, "", nil},
		{"net-lb", "负载均衡", "Scale", "central", intPtr(3), 2, []string{"负载均衡", "F5", "Ingress"}, "", []vendDef{{"foreign", "F5", 0}}},
		{"net-dns", "DNS / IPAM", "Globe2", "central", intPtr(3), 3, []string{"DNS", "IPAM", "地址管理"}, "", nil},
		// L2
		{"hw-compute", "计算", "Server", "central", intPtr(2), 0, []string{"服务器", "计算", "x86"}, "", []vendDef{
			{"domestic", "H3C", 0}, {"domestic", "浪潮", 1}, {"domestic", "华为", 2},
			{"foreign", "Dell", 3}, {"foreign", "HPE", 4},
		}},
		{"hw-san", "SAN 存储", "HardDrive", "central", intPtr(2), 1, []string{"SAN", "块存储"}, "", nil},
		{"hw-obj", "文件 / 对象", "Layers3", "central", intPtr(2), 2, []string{"对象存储", "NAS", "文件"}, "", nil},
		{"hw-gpu", "GPU", "Sparkles", "central", intPtr(2), 3, []string{"GPU", "算力", "推理", "训练"}, "", []vendDef{
			{"domestic", "华为", 0}, {"domestic", "寒武纪", 1}, {"foreign", "Nvidia", 2}, {"foreign", "AMD", 3},
		}},
		// L1
		{"phy-dc", "数据中心", "Building2", "central", intPtr(1), 0, []string{"数据中心", "机房", "IDC"}, "", nil},
		{"phy-ups", "UPS / PDU", "Zap", "central", intPtr(1), 1, []string{"UPS", "PDU", "供电"}, "", nil},
		{"phy-hvac", "精密空调", "ThermometerSnowflake", "central", intPtr(1), 2, []string{"空调", "制冷", "PUE"}, "", nil},
		{"phy-cable", "综合布线", "Cable", "central", intPtr(1), 3, []string{"布线", "光纤", "铜缆"}, "", nil},
	}

	for _, cd := range cats {
		kw, _ := json.Marshal(cd.Keywords)
		cat := models.SolutionCategory{
			Slug: cd.Slug, Label: cd.Label, IconKey: cd.Icon, ColumnType: cd.Column,
			SortOrder: cd.Ord, Keywords: datatypes.JSON(kw), Hint: cd.Hint, IsActive: true,
		}
		if cd.LayerLevel != nil {
			lid, ok := layerIDByLevel[*cd.LayerLevel]
			if ok {
				cat.LayerID = &lid
			}
		}
		if err := db.Create(&cat).Error; err != nil {
			log.Fatal(err)
		}
		for _, vd := range cd.Vendors {
			v := models.SolutionVendor{
				CategoryID: cat.ID, Market: vd.Market, Name: vd.Name, SortOrder: vd.Ord,
			}
			if err := db.Create(&v).Error; err != nil {
				log.Fatal(err)
			}
		}
	}
	log.Println("拓扑种子完成")
}

// FixLegacyTopologyHints 清除历史种子中的示例性 hint（如已录入真实内容则不会误伤）
func FixLegacyTopologyHints(db *gorm.DB) {
	for _, h := range []string{"如 Dify 等", "如Dify等"} {
		_ = db.Model(&models.SolutionCategory{}).
			Where("slug = ? AND trim(hint) = ?", "app-agent", h).
			Update("hint", "").Error
	}
}
