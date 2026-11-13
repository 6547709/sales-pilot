package main

import (
	"log"
	"os"
	"time"

	"github.com/6547709/sales-pilot/backend/internal/auth"
	"github.com/6547709/sales-pilot/backend/internal/cache"
	"github.com/6547709/sales-pilot/backend/internal/config"
	"github.com/6547709/sales-pilot/backend/internal/database"
	"github.com/6547709/sales-pilot/backend/internal/httpapi"
	"github.com/6547709/sales-pilot/backend/internal/repository"
	"github.com/6547709/sales-pilot/backend/internal/seed"
	"github.com/6547709/sales-pilot/backend/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	gin.SetMode(cfg.GinMode)

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("数据库连接失败: ", err)
	}
	seed.EnsureTopology(db)
	seed.FixLegacyTopologyHints(db)
	seed.Run(db)
	seed.EnsureServices(db)
	seed.EnsureScenarios(db)

	userRepo := repository.NewUserRepository(db)
	apiKeyRepo := repository.NewAPIKeyRepository(db)
	productRepo := repository.NewProductRepository(db)
	topologyRepo := repository.NewTopologyRepository(db)
	settingsRepo := repository.NewSettingsRepository(db)
	backupRepo := repository.NewBackupRepository(db)
	serviceRepo := repository.NewServiceRepository(db)
	scenarioRepo := repository.NewScenarioRepository(db)

	searchCache := cache.NewSearchCache(30 * time.Second)
	authSvc := service.NewAuthService(cfg, userRepo, settingsRepo)
	productSvc := service.NewProductService(productRepo)
	topologySvc := service.NewTopologyService(topologyRepo, serviceRepo)
	serviceSvc := service.NewServiceService(serviceRepo)
	scenarioSvc := service.NewScenarioService(scenarioRepo)
	userSvc := service.NewUserService(userRepo)
	apiKeySvc := service.NewAPIKeyService(apiKeyRepo, userRepo)
	settingsSvc := service.NewSettingsService(settingsRepo)
	backupSvc := service.NewBackupService(backupRepo, searchCache)

	r := httpapi.NewRouter(httpapi.Deps{
		Cfg:        cfg,
		Auth:       authSvc,
		Products:   productSvc,
		Topology:   topologySvc,
		Users:      userSvc,
		APIKeys:    apiKeySvc,
		Settings:   settingsSvc,
		Backup:     backupSvc,
		Services:   serviceSvc,
		Scenarios:  scenarioSvc,
		OidcStates: auth.NewOidcStateStore(),
	})

	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	log.Println("Sales-Pilot API 监听", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
