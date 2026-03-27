package business

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/business"
	"gorm.io/gorm"
)

type TopologyService struct{}

// Layer CRUD
func (s *TopologyService) CreateLayer(layer *business.TopologyLayer) error {
	return global.GVA_DB.Create(layer).Error
}

func (s *TopologyService) GetLayer(id uint) (layer *business.TopologyLayer, err error) {
	err = global.GVA_DB.First(&layer, id).Error
	return
}

func (s *TopologyService) GetAllLayers() (layers []business.TopologyLayer, err error) {
	err = global.GVA_DB.Order("sort_order ASC").Find(&layers).Error
	return
}

func (s *TopologyService) UpdateLayer(layer *business.TopologyLayer) error {
	return global.GVA_DB.Save(layer).Error
}

func (s *TopologyService) DeleteLayer(id uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("layer_id = ?", id).Delete(&business.TopologyCategory{}).Error; err != nil {
			return err
		}
		return tx.Delete(&business.TopologyLayer{}, id).Error
	})
}

// Category CRUD
func (s *TopologyService) CreateCategory(cat *business.TopologyCategory) error {
	return global.GVA_DB.Create(cat).Error
}

func (s *TopologyService) GetCategory(id uint) (cat *business.TopologyCategory, err error) {
	err = global.GVA_DB.First(&cat, id).Error
	return
}

func (s *TopologyService) GetCategoriesByLayerID(layerID uint) (cats []business.TopologyCategory, err error) {
	err = global.GVA_DB.Where("layer_id = ?", layerID).Find(&cats).Error
	return
}

func (s *TopologyService) UpdateCategory(cat *business.TopologyCategory) error {
	return global.GVA_DB.Save(cat).Error
}

func (s *TopologyService) DeleteCategory(id uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("category_id = ?", id).Delete(&business.TopologyVendor{}).Error; err != nil {
			return err
		}
		return tx.Delete(&business.TopologyCategory{}, id).Error
	})
}

// Vendor CRUD
func (s *TopologyService) CreateVendor(vendor *business.TopologyVendor) error {
	return global.GVA_DB.Create(vendor).Error
}

func (s *TopologyService) GetVendor(id uint) (vendor *business.TopologyVendor, err error) {
	err = global.GVA_DB.First(&vendor, id).Error
	return
}

func (s *TopologyService) GetVendorsByCategoryID(categoryID uint) (vendors []business.TopologyVendor, err error) {
	err = global.GVA_DB.Where("category_id = ?", categoryID).Find(&vendors).Error
	return
}

func (s *TopologyService) UpdateVendor(vendor *business.TopologyVendor) error {
	return global.GVA_DB.Save(vendor).Error
}

func (s *TopologyService) DeleteVendor(id uint) error {
	return global.GVA_DB.Delete(&business.TopologyVendor{}, id).Error
}

// GetFullTopology 返回完整的拓扑结构（供首页全景图使用）
func (s *TopologyService) GetFullTopology() (*business.TopologyFullResponse, error) {
	var layers []business.TopologyLayer
	if err := global.GVA_DB.Order("sort_order ASC").Find(&layers).Error; err != nil {
		return nil, err
	}

	var categories []business.TopologyCategory
	if err := global.GVA_DB.Find(&categories).Error; err != nil {
		return nil, err
	}

	// 按 layer_id 分组
	catsByLayer := make(map[uint][]business.TopologyCategory)
	for _, cat := range categories {
		catsByLayer[cat.LayerID] = append(catsByLayer[cat.LayerID], cat)
	}

	// 构建响应
	resp := &business.TopologyFullResponse{
		CentralLayers: []business.TopologyLayerBlock{},
	}

	for _, layer := range layers {
		block := business.TopologyLayerBlock{
			Layer:      layer,
			Categories: catsByLayer[layer.ID],
		}
		resp.CentralLayers = append(resp.CentralLayers, block)
	}

	return resp, nil
}
