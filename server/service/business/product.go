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
