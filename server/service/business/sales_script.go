package business

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/business"
)

type SalesScriptService struct{}

func (s *SalesScriptService) CreateSalesScript(script *business.SalesScript) error {
	return global.GVA_DB.Create(script).Error
}

func (s *SalesScriptService) GetSalesScript(id uint) (script *business.SalesScript, err error) {
	err = global.GVA_DB.First(&script, id).Error
	return
}

func (s *SalesScriptService) GetSalesScriptsByProductID(productID uint) (scripts []business.SalesScript, err error) {
	err = global.GVA_DB.Where("product_id = ?", productID).Find(&scripts).Error
	return
}

func (s *SalesScriptService) UpdateSalesScript(script *business.SalesScript) error {
	return global.GVA_DB.Save(script).Error
}

func (s *SalesScriptService) DeleteSalesScript(id uint) error {
	return global.GVA_DB.Delete(&business.SalesScript{}, id).Error
}
