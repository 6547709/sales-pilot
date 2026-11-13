package httpapi

import (
	"github.com/6547709/sales-pilot/backend/internal/auth"
	"github.com/6547709/sales-pilot/backend/internal/config"
	"github.com/6547709/sales-pilot/backend/internal/service"
)

// Deps 路由所需的全部依赖
type Deps struct {
	Cfg        config.Config
	Auth       *service.AuthService
	Products   *service.ProductService
	Topology   *service.TopologyService
	Users      *service.UserService
	APIKeys    *service.APIKeyService
	Settings   *service.SettingsService
	Backup     *service.BackupService
	OidcStates *auth.OidcStateStore
}

type Handlers struct {
	Deps
}

func NewHandlers(d Deps) *Handlers {
	return &Handlers{Deps: d}
}
