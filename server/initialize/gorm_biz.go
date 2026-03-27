package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/business"
)

func bizModel() error {
	db := global.GVA_DB
	err := db.AutoMigrate(
		&business.TopologyLayer{},
		&business.TopologyCategory{},
	)
	if err != nil {
		return err
	}
	return nil
}
