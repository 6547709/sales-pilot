package main

import (
	"log"
	"os"
	"time"

	"github.com/6547709/sales-pilot/backend/internal/auth"
	"github.com/6547709/sales-pilot/backend/internal/cache"
	"github.com/6547709/sales-pilot/backend/internal/config"
	"github.com/6547709/sales-pilot/backend/internal/database"
	"github.com/6547709/sales-pilot/backend/internal/handlers"
	"github.com/6547709/sales-pilot/backend/internal/seed"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	gin.SetMode(cfg.GinMode)

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("数据库连接失败: ", err)
	}
	seed.Run(db)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(corsMiddleware(cfg.FrontendOrigin))

	srv := &handlers.Server{
		DB:          db,
		Cfg:         cfg,
		SearchCache: cache.NewSearchCache(30 * time.Second),
		OidcStates:  auth.NewOidcStateStore(),
	}
	srv.RegisterRoutes(r)

	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	log.Println("Sales-Pilot API 监听", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}

func corsMiddleware(origin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
