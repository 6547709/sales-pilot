package httpapi

import (
	"net/http"

	"github.com/6547709/sales-pilot/backend/internal/httpapi/middleware"
	"github.com/gin-gonic/gin"
)

// NewRouter 注册全部路由
func NewRouter(d Deps) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(corsMiddleware(d.Cfg.FrontendOrigin))
	h := NewHandlers(d)

	api := r.Group("/api/v1")
	api.GET("/health", func(c *gin.Context) { OK(c, gin.H{"status": "ok"}) })
	api.GET("/config", h.getConfig)

	api.POST("/auth/login", h.postLocalLogin)
	api.POST("/auth/login/ldap", h.postLDAPLogin)
	api.GET("/auth/oidc/start", h.getOidcStart)
	api.GET("/auth/oidc/callback", h.getOidcCallback)
	api.GET("/auth/login-options", h.getLoginOptions)

	authMW := middleware.AdminAuth(d.Cfg.JWTSecret, d.APIKeys)
	session := api.Group("")
	session.Use(authMW)
	{
		session.GET("/me", h.getMe)
		session.GET("/topology", h.getTopology)
		session.GET("/solution-categories", h.listSolutionCategories)
		session.GET("/products", h.listProducts)
		session.GET("/products/search", h.searchProducts)
		session.GET("/products/:id", h.getProduct)
		session.GET("/products/:id/scripts", h.listScripts)
		session.GET("/products/:id/cases", h.listCases)
		session.GET("/products/:id/scenarios", h.getProductScenarios)
		session.GET("/services", h.listServices)
		session.GET("/services/:slug", h.getService)
		session.GET("/services/:slug/scenarios", h.getServiceScenarios)
		session.GET("/scenarios", h.listScenarios)
		session.GET("/scenarios/:slug", h.getScenario)
	}

	admin := api.Group("/admin")
	admin.Use(authMW, middleware.RequireAdmin())
	{
		admin.GET("/api-keys", h.adminListAPIKeys)
		admin.POST("/api-keys", h.adminCreateAPIKey)
		admin.DELETE("/api-keys/:id", h.adminRevokeAPIKey)

		admin.GET("/auth-settings", h.getAuthSettings)
		admin.PATCH("/auth-settings", h.patchAuthSettings)
		admin.GET("/system-settings", h.getSystemSettings)
		admin.PATCH("/system-settings", h.patchSystemSettings)

		admin.GET("/users/local", h.adminListLocalUsers)
		admin.POST("/users/local", h.adminCreateLocalUser)
		admin.PUT("/users/local/:id", h.adminUpdateLocalUser)
		admin.PUT("/users/local/:id/password", h.adminResetLocalUserPassword)
		admin.DELETE("/users/local/:id", h.adminDeleteLocalUser)

		admin.GET("/backup/export", h.adminExportBackup)
		admin.POST("/backup/import", h.adminImportBackup)
		admin.GET("/backup/export.zip", h.adminExportBackupZip)
		admin.POST("/backup/import-zip", h.adminImportBackupZip)

		admin.GET("/meta/product-maintain", h.adminProductMaintainMeta)

		admin.GET("/products", h.adminListProducts)
		admin.GET("/products/:id/cases", h.adminListCasesByProduct)
		admin.GET("/products/:id/scripts", h.adminListScriptsByProduct)
		admin.GET("/products/:id", h.adminGetProduct)
		admin.PATCH("/products/:id/draft", h.adminPatchProductDraft)
		admin.POST("/products", h.createProduct)
		admin.PUT("/products/:id", h.updateProduct)
		admin.GET("/products/:id/links", h.adminGetProductLinks)
		admin.PUT("/products/:id/links", h.adminSetProductLinks)
		admin.DELETE("/products/:id", h.deleteProduct)
		admin.POST("/products/:id/scripts", h.createScript)
		admin.PUT("/scripts/:id", h.updateScript)
		admin.DELETE("/scripts/:id", h.deleteScript)
		admin.POST("/products/:id/cases", h.createCase)
		admin.PUT("/cases/:id", h.updateCase)
		admin.DELETE("/cases/:id", h.deleteCase)
		admin.PATCH("/cases/:id/publish", h.adminPatchCasePublish)

		admin.GET("/services", h.adminListServices)
		admin.POST("/services", h.adminCreateService)
		admin.PUT("/services/:id", h.adminUpdateService)
		admin.DELETE("/services/:id", h.adminDeleteService)
		admin.PUT("/services/:id/categories", h.adminSetServiceCategories)

		admin.GET("/scenario-groups", h.adminListScenarioGroups)
		admin.POST("/scenario-groups", h.adminCreateScenarioGroup)
		admin.PUT("/scenario-groups/:id", h.adminUpdateScenarioGroup)
		admin.DELETE("/scenario-groups/:id", h.adminDeleteScenarioGroup)
		admin.GET("/scenarios", h.adminListScenarios)
		admin.POST("/scenarios", h.adminCreateScenario)
		admin.GET("/scenarios/:id", h.adminGetScenario)
		admin.PUT("/scenarios/:id", h.adminUpdateScenario)
		admin.DELETE("/scenarios/:id", h.adminDeleteScenario)
		admin.PUT("/scenarios/:id/links", h.adminSetScenarioLinks)
		admin.PUT("/solution-categories/:id/customer", h.adminUpdateCategoryCustomer)

		admin.GET("/topology/layers", h.adminListLayers)
		admin.POST("/topology/layers", h.adminCreateLayer)
		admin.PUT("/topology/layers/:id", h.adminUpdateLayer)
		admin.DELETE("/topology/layers/:id", h.adminDeleteLayer)
		admin.POST("/topology/categories", h.adminCreateCategory)
		admin.PUT("/topology/categories/:id", h.adminUpdateCategory)
		admin.DELETE("/topology/categories/:id", h.adminDeleteCategory)
		admin.POST("/topology/vendors", h.adminCreateVendor)
		admin.PUT("/topology/vendors/:id", h.adminUpdateVendor)
		admin.DELETE("/topology/vendors/:id", h.adminDeleteVendor)
	}
	return r
}

func corsMiddleware(origin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-API-Key")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
