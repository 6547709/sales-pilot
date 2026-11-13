package seed

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/6547709/sales-pilot/backend/internal/models"
	"gorm.io/gorm"
)

type scenarioSeed struct {
	GroupSlug string
	Slug      string
	Name      string
	Summary   string
	PainPoint string
	ValueProp string
	IconKey   string
	SortOrder int
}

func defaultGroups() []models.ScenarioGroup {
	return []models.ScenarioGroup{
		{Slug: "data-protection", Name: "数据保护", Description: "备份、归档与恢复", SortOrder: 10},
		{Slug: "infrastructure-upgrade", Name: "基础架构升级", Description: "版本、容量与系统升级", SortOrder: 20},
		{Slug: "virtualization-replacement", Name: "虚拟化替代", Description: "平台替换与国产化", SortOrder: 30},
		{Slug: "dr-migration", Name: "容灾与迁移", Description: "迁移、双活与切换", SortOrder: 40},
	}
}

func defaultScenarios() []scenarioSeed {
	return []scenarioSeed{
		{"data-protection", "core-db-backup", "核心数据库备份", "备份窗口不足，恢复结果不可验证。", "全量备份经常超出窗口，恢复演练缺少标准流程。", "备份窗口评估 + 策略调优 + 季度恢复演练，交付可量化的 RTO / RPO。", "database-backup", 10},
		{"data-protection", "vm-backup", "虚拟机备份", "多平台统一备份与快速恢复。", "虚拟化平台分散，备份策略不统一。", "统一策略覆盖多平台，恢复时间可预期。", "server", 20},
		{"data-protection", "file-archive", "文件与对象归档", "长期留存、合规与低成本存储。", "归档数据占满主存储，合规留存难管理。", "冷热分层归档，满足留存要求并降低成本。", "archive", 30},
		{"infrastructure-upgrade", "vmware-upgrade", "VMware 升级", "版本升级、兼容性验证与回滚。", "版本过旧，升级风险高且缺少回滚预案。", "兼容性清单 + 分批升级 + 回滚方案。", "arrow-up-circle", 10},
		{"infrastructure-upgrade", "storage-expand", "存储扩容", "容量评估、扩容与数据再平衡。", "容量接近上限，扩容影响业务。", "在线扩容与数据再平衡，业务无感。", "hard-drive", 20},
		{"infrastructure-upgrade", "os-upgrade", "操作系统升级", "批量升级与兼容性清单。", "操作系统停止维护，补丁无法获取。", "批量升级路径与兼容性验证。", "terminal", 30},
		{"virtualization-replacement", "vmware-to-domestic", "VMware → 国产虚拟化", "平滑替换，保留运维习惯。", "授权成本高、国产化要求。", "业务不中断完成平台替换，运维习惯可延续。", "server", 10},
		{"virtualization-replacement", "vmware-to-hci", "VMware → 超融合", "一体化架构替换与容量规划。", "架构复杂、扩容不灵活。", "一体化架构与容量规划，降低运维复杂度。", "layers", 20},
		{"dr-migration", "v2v-migration", "V2V 迁移", "异构平台批量迁移与校验。", "平台替换缺少迁移工具与经验。", "分批迁移、在线校验、失败可回滚。", "arrow-left-right", 10},
		{"dr-migration", "active-active", "双活容灾", "按 RTO/RPO 分级设计切换。", "单中心运行，故障时无法快速切换。", "分级双活设计与切换演练。", "shield-check", 20},
	}
}

// EnsureScenarios 幂等写入默认场景组与场景。
func EnsureScenarios(db *gorm.DB) {
	groupIDs := make(map[string]uint, len(defaultGroups()))
	for _, g := range defaultGroups() {
		var existing models.ScenarioGroup
		err := db.Where("slug = ?", g.Slug).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.Create(&g).Error; err != nil {
				log.Printf("seed scenario group %s: %v", g.Slug, err)
				continue
			}
			groupIDs[g.Slug] = g.ID
			continue
		}
		if err == nil {
			groupIDs[g.Slug] = existing.ID
		}
	}
	for _, sc := range defaultScenarios() {
		gid, ok := groupIDs[sc.GroupSlug]
		if !ok {
			continue
		}
		var existing models.Scenario
		err := db.Where("slug = ?", sc.Slug).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			item := models.Scenario{
				GroupID: gid, Slug: sc.Slug, Name: sc.Name, Summary: sc.Summary,
				PainPoint: sc.PainPoint, ValueProp: sc.ValueProp, IconKey: sc.IconKey,
				SortOrder: sc.SortOrder, IsActive: true,
			}
			if err := db.Create(&item).Error; err != nil {
				log.Printf("seed scenario %s: %v", sc.Slug, err)
			}
		}
	}
	if os.Getenv("SEED_DEMO") == "1" {
		linkDemoScenarioContent(db)
	}
}

// linkDemoScenarioContent 仅在演示环境建立示例关联，生产由管理员维护。
func linkDemoScenarioContent(db *gorm.DB) {
	services := make(map[string]uint)
	var svcList []models.Service
	db.Find(&svcList)
	for _, s := range svcList {
		services[s.Slug] = s.ID
	}
	serviceLinks := map[string][]string{
		"core-db-backup":     {"backup-service", "dr-orchestration", "managed-ops"},
		"vm-backup":          {"backup-service", "managed-ops"},
		"file-archive":       {"backup-service"},
		"vmware-upgrade":     {"virtualization-upgrade", "managed-ops"},
		"storage-expand":     {"managed-ops"},
		"os-upgrade":         {"managed-ops"},
		"vmware-to-domestic": {"v2v-migration", "virtualization-upgrade"},
		"vmware-to-hci":      {"v2v-migration", "virtualization-upgrade"},
		"v2v-migration":      {"v2v-migration", "cloud-migration"},
		"active-active":      {"dr-orchestration", "managed-ops"},
	}
	productLinks := map[string]string{
		"core-db-backup":     "备份",
		"vm-backup":          "备份",
		"vmware-upgrade":     "虚拟化",
		"vmware-to-domestic": "虚拟化",
		"vmware-to-hci":      "虚拟化",
		"v2v-migration":      "虚拟化",
		"active-active":      "容灾",
	}
	for slug, svcSlugs := range serviceLinks {
		var scenario models.Scenario
		if err := db.Where("slug = ?", slug).First(&scenario).Error; err != nil {
			continue
		}
		for _, ss := range svcSlugs {
			if id, ok := services[ss]; ok {
				db.Where("scenario_id = ? AND service_id = ?", scenario.ID, id).
					FirstOrCreate(&models.ScenarioService{ScenarioID: scenario.ID, ServiceID: id})
			}
		}
	}
	for slug, keyword := range productLinks {
		var scenario models.Scenario
		if err := db.Where("slug = ?", slug).First(&scenario).Error; err != nil {
			continue
		}
		var products []models.Product
		db.Where("LOWER(category) LIKE ? AND is_draft = ?", "%"+strings.ToLower(keyword)+"%", false).
			Limit(3).Find(&products)
		for _, p := range products {
			db.Where("scenario_id = ? AND product_id = ?", scenario.ID, p.ID).
				FirstOrCreate(&models.ScenarioProduct{ScenarioID: scenario.ID, ProductID: p.ID})
		}
	}
}
