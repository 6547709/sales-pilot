package seed

import (
	"encoding/json"
	"log"

	"github.com/6547709/sales-pilot/backend/internal/auth"
	"github.com/6547709/sales-pilot/backend/internal/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Run 若库为空则写入默认管理员、认证配置占位与示例产品
func Run(db *gorm.DB) {
	var n int64
	db.Model(&models.User{}).Count(&n)
	if n > 0 {
		return
	}
	log.Println("执行初始数据种子...")

	hash, err := auth.HashPassword("admin123")
	if err != nil {
		log.Fatal(err)
	}
	admin := models.User{
		Username: "admin", Email: "admin@local",
		PasswordHash: hash, AuthProvider: "local", Role: "admin", IsActive: true,
	}
	if err := db.Create(&admin).Error; err != nil {
		log.Fatal(err)
	}

	ldapEx, _ := json.Marshal(map[string]any{
		"host":             "ldap.example.com",
		"port":             389,
		"use_tls":          false,
		"bind_dn":          "cn=readonly,dc=example,dc=com",
		"bind_password":    "",
		"user_search_base": "ou=users,dc=example,dc=com",
		"user_filter":      "(uid=%s)",
	})
	oidcEx, _ := json.Marshal(map[string]any{
		"issuer":        "https://keycloak.example/realms/corp",
		"client_id":     "",
		"client_secret": "",
		"redirect_url":  "http://localhost:8080/api/v1/auth/oidc/callback",
	})
	st := models.AuthSettings{
		ID: 1, LocalLoginEnabled: true, LdapEnabled: false, LdapConfig: datatypes.JSON(ldapEx),
		OidcEnabled: false, OidcConfig: datatypes.JSON(oidcEx),
	}
	_ = db.Create(&st).Error

	highlights, _ := json.Marshal([]string{
		"彻底解决国产化替代合规风险。",
		"资源利用率提升 40% 以上，显著降低 TCO。",
		"原生集成 KubeVirt，容器与虚拟机统一纳管。",
	})
	questions, _ := json.Marshal([]string{
		"面对 VMware 授权涨价，您是否计算过未来三年的总拥有成本（TCO）变化？",
		"现有的传统备份方案在恢复容器化业务时，RTO 能够达到分钟级吗？",
		"您的开发团队是否抱怨过环境交付速度跟不上业务迭代？",
	})
	personas, _ := json.Marshal(map[string]string{
		"CIO":  "关注合规、成本与供应商锁定风险。",
		"架构师": "关注工作负载形态统一与运维复杂度。",
		"业务负责人": "关注交付速度与业务连续性。",
	})

	competitor := "## 竞品对比（埋雷）\n\n| 维度 | 我们（OpenShift + KubeVirt） | 私有云大厂 A/B |\n|------|------------------------------|----------------|\n| 生态 | 开放 Kubernetes 生态，可移植 | 闭源栈深绑定，迁移成本高 |\n| 混合云 | 与主流公有云一致体验 | 难以在公有云复现同等架构 |\n| 风险 | 由于我们拥有 **开放标准与跨云一致性**，所以能避免客户遇到 **单一厂商锁定与合规不可控** 的风险；竞品往往无法同时解决「虚拟机遗留」与「云原生扩展」的 Y 类痛点。|\n"

	var solCat models.SolutionCategory
	var solPtr *uint
	catLabel := "云原生 / 虚拟化迁移"
	if err := db.Where("slug = ?", "cloud-k8s").First(&solCat).Error; err == nil {
		solPtr = &solCat.ID
		catLabel = solCat.Label
	}

	p := models.Product{
		Name:               "VMware 迁移至 OpenShift 方案",
		Category:           catLabel,
		SolutionCategoryID: solPtr,
		VendorMarket:       "foreign",
		ManufacturerName:   "红帽软件（示例）",
		SalesContactName:   "张销售",
		SalesContactPhone:  "13800000000",
		SalesContactEmail:  "sales@example.com",
		PresalesContactName:  "李售前",
		PresalesContactPhone: "13900000000",
		PresalesContactEmail: "presales@example.com",
		Description: "面向 VMware 存量客户，提供可落地的 OpenShift 迁移与并行运行路径，强调合规、成本与统一管控。",
		Highlights:         datatypes.JSON(highlights),
		TargetPersonas:     datatypes.JSON(personas),
		TriggerEvents:      "授权涨价谈判、信创时间表、容器化改造立项、数据中心整合。",
		DiscoveryQuestions: datatypes.JSON(questions),
		CompetitorAnalysis: competitor,
		ROIMetrics:         "典型客户：三年 TCO 下降约 25%–35%（视规模与授权结构）；环境交付从周级缩短到天级。",
	}
	if err := db.Create(&p).Error; err != nil {
		log.Fatal(err)
	}

	scripts := []models.SalesScript{
		{ProductID: p.ID, Scenario: "异议：迁移风险大", Content: "我们建议分阶段并行：先用 KubeVirt 纳管关键 VM，再逐步容器化。每一步都有回滚窗口，业务不中断。"},
		{ProductID: p.ID, Scenario: "开场破冰", Content: "如果未来三年授权成本按当前涨幅外推，这笔预算是否还能支撑贵司的云原生路线图？"},
	}
	for i := range scripts {
		_ = db.Create(&scripts[i]).Error
	}

	cases := []models.Case{
		{
			ProductID: p.ID, ClientName: "某省政务云（示例）",
			PainPoints:     "VMware 授权成本攀升，信创验收窗口紧。",
			Solution:       "OpenShift 为底座，KubeVirt 承载遗留 VM，新应用容器化发布。",
			ValueDelivered: "6 个月内完成核心系统迁移评估与一期上线，形成可复制的标准模板。",
		},
	}
	for i := range cases {
		_ = db.Create(&cases[i]).Error
	}

	log.Println("种子完成：管理员 admin / admin123（生产环境请立即修改密码）")
}
