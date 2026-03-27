package business

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/business"
)

type SolutionCategoryService struct{}

func (s *SolutionCategoryService) Create(cat *business.SolutionCategory) error {
	return global.GVA_DB.Create(cat).Error
}

func (s *SolutionCategoryService) Get(id uint) (cat *business.SolutionCategory, err error) {
	err = global.GVA_DB.First(&cat, id).Error
	return
}

func (s *SolutionCategoryService) GetAll() (cats []business.SolutionCategory, err error) {
	err = global.GVA_DB.Find(&cats).Error
	return
}

func (s *SolutionCategoryService) Update(cat *business.SolutionCategory) error {
	return global.GVA_DB.Save(cat).Error
}

func (s *SolutionCategoryService) Delete(id uint) error {
	return global.GVA_DB.Delete(&business.SolutionCategory{}, id).Error
}
