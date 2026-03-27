package business

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/business"
	"gorm.io/gorm"
)

// BackupData 备份数据结构
type BackupData struct {
	SolutionCategories []business.SolutionCategory    `json:"solution_categories"`
	Products           []business.Product              `json:"products"`
	SalesScripts       []business.SalesScript          `json:"sales_scripts"`
	Cases              []business.Case                 `json:"cases"`
	TopologyLayers     []business.TopologyLayer        `json:"topology_layers"`
	TopologyCategories []business.TopologyCategory     `json:"topology_categories"`
	TopologyVendors    []business.TopologyVendor       `json:"topology_vendors"`
}

type BackupService struct{}

// ExportAll 导出所有业务数据
func (s *BackupService) ExportAll() (*BackupData, error) {
	data := &BackupData{}

	if err := global.GVA_DB.Find(&data.SolutionCategories).Error; err != nil {
		return nil, err
	}
	if err := global.GVA_DB.Find(&data.Products).Error; err != nil {
		return nil, err
	}
	if err := global.GVA_DB.Find(&data.SalesScripts).Error; err != nil {
		return nil, err
	}
	if err := global.GVA_DB.Find(&data.Cases).Error; err != nil {
		return nil, err
	}
	if err := global.GVA_DB.Find(&data.TopologyLayers).Error; err != nil {
		return nil, err
	}
	if err := global.GVA_DB.Find(&data.TopologyCategories).Error; err != nil {
		return nil, err
	}
	if err := global.GVA_DB.Find(&data.TopologyVendors).Error; err != nil {
		return nil, err
	}

	return data, nil
}

// ImportAll 导入所有业务数据（覆盖式）
func (s *BackupService) ImportAll(data *BackupData) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 清空并重建表
		if err := tx.Where("1=1").Delete(&business.TopologyVendor{}).Error; err != nil {
			return err
		}
		if err := tx.Where("1=1").Delete(&business.TopologyCategory{}).Error; err != nil {
			return err
		}
		if err := tx.Where("1=1").Delete(&business.TopologyLayer{}).Error; err != nil {
			return err
		}
		if err := tx.Where("1=1").Delete(&business.Case{}).Error; err != nil {
			return err
		}
		if err := tx.Where("1=1").Delete(&business.SalesScript{}).Error; err != nil {
			return err
		}
		if err := tx.Where("1=1").Delete(&business.Product{}).Error; err != nil {
			return err
		}
		if err := tx.Where("1=1").Delete(&business.SolutionCategory{}).Error; err != nil {
			return err
		}

		// 重新插入数据
		if len(data.SolutionCategories) > 0 {
			if err := tx.Create(&data.SolutionCategories).Error; err != nil {
				return err
			}
		}
		if len(data.TopologyLayers) > 0 {
			if err := tx.Create(&data.TopologyLayers).Error; err != nil {
				return err
			}
		}
		if len(data.TopologyCategories) > 0 {
			if err := tx.Create(&data.TopologyCategories).Error; err != nil {
				return err
			}
		}
		if len(data.TopologyVendors) > 0 {
			if err := tx.Create(&data.TopologyVendors).Error; err != nil {
				return err
			}
		}
		if len(data.Products) > 0 {
			if err := tx.Create(&data.Products).Error; err != nil {
				return err
			}
		}
		if len(data.SalesScripts) > 0 {
			if err := tx.Create(&data.SalesScripts).Error; err != nil {
				return err
			}
		}
		if len(data.Cases) > 0 {
			if err := tx.Create(&data.Cases).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
