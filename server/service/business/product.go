package business

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/business"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"gorm.io/gorm"
)

type ProductService struct{}

func (s *ProductService) CreateProduct(product *business.Product) error {
	return global.GVA_DB.Create(product).Error
}

func (s *ProductService) GetProduct(id uint) (product *business.Product, err error) {
	err = global.GVA_DB.First(&product, id).Error
	return
}

func (s *ProductService) GetProductList(info *request.PageInfo) (list []business.Product, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&business.Product{})
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Order("updated_at DESC").Limit(limit).Offset(offset).Find(&list).Error
	return
}

func (s *ProductService) UpdateProduct(product *business.Product) error {
	return global.GVA_DB.Save(product).Error
}

func (s *ProductService) DeleteProduct(id uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&business.Product{}, id).Error; err != nil {
			return err
		}
		if err := tx.Where("product_id = ?", id).Delete(&business.SalesScript{}).Error; err != nil {
			return err
		}
		if err := tx.Where("product_id = ?", id).Delete(&business.Case{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *ProductService) GetPublicProductList(info *request.PageInfo, vendorMarket string) (list []business.Product, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&business.Product{}).Where("is_draft = ?", false)

	if vendorMarket == "domestic" || vendorMarket == "foreign" {
		db = db.Where("vendor_market IN ?", []string{"all", vendorMarket})
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Order("updated_at DESC").Limit(limit).Offset(offset).Find(&list).Error
	return
}

// GetAllPublicProducts 获取所有公开产品（供首页全景图使用，不分页）
func (s *ProductService) GetAllPublicProducts() (list []business.Product, err error) {
	err = global.GVA_DB.Where("is_draft = ?", false).Order("updated_at DESC").Find(&list).Error
	return
}

// GetProductsByCategory 获取指定分类的公开产品
func (s *ProductService) GetProductsByCategory(categoryID uint, vendorMarket string) (list []business.Product, err error) {
	db := global.GVA_DB.Model(&business.Product{}).Where("is_draft = ?", false).Where("solution_category_id = ?", categoryID)

	if vendorMarket == "domestic" || vendorMarket == "foreign" {
		db = db.Where("vendor_market IN ?", []string{"all", vendorMarket})
	}

	err = db.Order("updated_at DESC").Find(&list).Error
	return
}
