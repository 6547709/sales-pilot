package business

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/business"
)

type CaseService struct{}

func (s *CaseService) CreateCase(c *business.Case) error {
	return global.GVA_DB.Create(c).Error
}

func (s *CaseService) GetCase(id uint) (c *business.Case, err error) {
	err = global.GVA_DB.First(&c, id).Error
	return
}

func (s *CaseService) GetCasesByProductID(productID uint) (cases []business.Case, err error) {
	err = global.GVA_DB.Where("product_id = ?", productID).Find(&cases).Error
	return
}

func (s *CaseService) UpdateCase(c *business.Case) error {
	return global.GVA_DB.Save(c).Error
}

func (s *CaseService) DeleteCase(id uint) error {
	return global.GVA_DB.Delete(&business.Case{}, id).Error
}
