package seed

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/6547709/sales-pilot/backend/internal/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func jsonList(items []string) datatypes.JSON {
	raw, _ := json.Marshal(items)
	return datatypes.JSON(raw)
}

func defaultServices() []models.Service {
	return []models.Service{
		{
			Slug: "v2v-migration", Name: "V2V 迁移服务", Kind: "v2v", SortOrder: 10, IsActive: true,
			Summary:      "异构虚拟化平台之间的批量迁移与校验。",
			Highlights:   jsonList([]string{"迁移评估", "分批割接", "结果校验"}),
			Deliverables: jsonList([]string{"迁移评估报告", "割接方案", "迁移执行与回滚预案"}),
			ValueProp:    "业务不中断完成平台替换，迁移结果可验证。",
			Features:     jsonList([]string{"批量迁移", "在线校验", "回滚预案"}),
			IconKey:      "arrow-left-right",
		},
		{
			Slug: "backup-service", Name: "备份服务", Kind: "backup", SortOrder: 20, IsActive: true,
			Summary:      "核心数据库与虚拟机的备份策略设计与演练。",
			Highlights:   jsonList([]string{"策略设计", "窗口优化", "恢复演练"}),
			Deliverables: jsonList([]string{"备份策略文档", "恢复演练报告"}),
			ValueProp:    "交付可量化的 RTO / RPO，而不是只交付一套工具。",
			Features:     jsonList([]string{"多库支持", "窗口评估", "季度演练"}),
			IconKey:      "database-backup",
		},
		{
			Slug: "dr-orchestration", Name: "容灾编排服务", Kind: "dr", SortOrder: 30, IsActive: true,
			Summary:      "双活与切换编排，按 RTO/RPO 分级落地。",
			Highlights:   jsonList([]string{"分级设计", "切换编排", "演练验证"}),
			Deliverables: jsonList([]string{"容灾设计", "切换手册", "演练记录"}),
			ValueProp:    "故障时敢切、切得回来。",
			Features:     jsonList([]string{"双活", "一键切换", "回切验证"}),
			IconKey:      "shield-check",
		},
		{
			Slug: "virtualization-upgrade", Name: "虚拟化升级服务", Kind: "upgrade", SortOrder: 40, IsActive: true,
			Summary:      "VMware 等平台的版本升级与兼容性验证。",
			Highlights:   jsonList([]string{"兼容性清单", "升级窗口", "回滚方案"}),
			Deliverables: jsonList([]string{"升级方案", "兼容性验证报告"}),
			ValueProp:    "在维护窗口内完成升级，失败可回滚。",
			Features:     jsonList([]string{"版本评估", "分批升级", "回滚"}),
			IconKey:      "arrow-up-circle",
		},
		{
			Slug: "cloud-migration", Name: "上云迁移服务", Kind: "migration", SortOrder: 50, IsActive: true,
			Summary:      "私有云与公有云之间的工作负载迁移。",
			Highlights:   jsonList([]string{"应用盘点", "依赖梳理", "迁移验证"}),
			Deliverables: jsonList([]string{"迁移评估", "迁移方案", "验收报告"}),
			ValueProp:    "按应用分批迁移，风险可控。",
			Features:     jsonList([]string{"应用盘点", "分批迁移", "验收"}),
			IconKey:      "cloud",
		},
		{
			Slug: "managed-ops", Name: "运维托管服务", Kind: "ops", SortOrder: 60, IsActive: true,
			Summary:      "日常巡检、容量与告警的持续托管。",
			Highlights:   jsonList([]string{"巡检", "容量", "告警响应"}),
			Deliverables: jsonList([]string{"月度巡检报告", "容量趋势分析"}),
			ValueProp:    "把运维负担交给专业团队，问题在发生前被发现。",
			Features:     jsonList([]string{"定期巡检", "容量规划", "响应 SLA"}),
			IconKey:      "activity",
		},
	}
}

// EnsureServices 按 slug 幂等写入默认服务，不覆盖后台已修改的内容。
func EnsureServices(db *gorm.DB) {
	for _, item := range defaultServices() {
		var existing models.Service
		err := db.Where("slug = ?", item.Slug).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.Create(&item).Error; err != nil {
				log.Printf("seed service %s: %v", item.Slug, err)
			}
		}
	}
}
