package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/6547709/sales-pilot/backend/internal/auth"
	"github.com/6547709/sales-pilot/backend/internal/cache"
	"github.com/6547709/sales-pilot/backend/internal/config"
	"github.com/6547709/sales-pilot/backend/internal/middleware"
	"github.com/6547709/sales-pilot/backend/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Server HTTP 依赖
type Server struct {
	DB         *gorm.DB
	Cfg        config.Config
	SearchCache *cache.SearchCache
	OidcStates *auth.OidcStateStore
}

// RegisterRoutes 注册 API
func (s *Server) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")

	api.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
	api.GET("/config", s.getConfig)

	api.POST("/auth/login", s.postLocalLogin)
	api.POST("/auth/login/ldap", s.postLDAPLogin)

	api.GET("/auth/oidc/start", s.getOidcStart)
	api.GET("/auth/oidc/callback", s.getOidcCallback)
	api.GET("/auth/login-options", s.getLoginOptions)

	// 登录后可读：JWT Bearer 或管理端 X-API-Key（与 AdminAuth 一致，普通用户仅 Bearer）
	session := api.Group("")
	session.Use(middleware.AdminAuth(s.Cfg.JWTSecret, s.DB))
	{
		session.GET("/me", s.getMe)
		session.GET("/topology", s.getTopology)
		session.GET("/solution-categories", s.listSolutionCategories)
		session.GET("/products", s.listProducts)
		session.GET("/products/search", s.searchProducts)
		session.GET("/products/:id", s.getProduct)
		session.GET("/products/:id/scripts", s.listScripts)
		session.GET("/products/:id/cases", s.listCases)
	}

	admin := api.Group("/admin")
	admin.Use(middleware.AdminAuth(s.Cfg.JWTSecret, s.DB), middleware.RequireAdmin())
	{
		admin.GET("/api-keys", s.adminListAPIKeys)
		admin.POST("/api-keys", s.adminCreateAPIKey)
		admin.DELETE("/api-keys/:id", s.adminRevokeAPIKey)

		admin.GET("/auth-settings", s.getAuthSettings)
		admin.PATCH("/auth-settings", s.patchAuthSettings)

		admin.GET("/system-settings", s.getSystemSettings)
		admin.PATCH("/system-settings", s.patchSystemSettings)

		admin.GET("/users/local", s.adminListLocalUsers)
		admin.POST("/users/local", s.adminCreateLocalUser)
		admin.PUT("/users/local/:id", s.adminUpdateLocalUser)
		admin.PUT("/users/local/:id/password", s.adminResetLocalUserPassword)
		admin.DELETE("/users/local/:id", s.adminDeleteLocalUser)

		admin.GET("/backup/export", s.adminExportBackup)
		admin.POST("/backup/import", s.adminImportBackup)

		// 供 AI / MCP 拉取产品维护字段与端点说明（机器可读）
		admin.GET("/meta/product-maintain", s.adminProductMaintainMeta)

		admin.GET("/products", s.adminListProducts)
		admin.GET("/products/:id/cases", s.adminListCasesByProduct)
		admin.GET("/products/:id/scripts", s.adminListScriptsByProduct)
		admin.GET("/products/:id", s.adminGetProduct)
		admin.PATCH("/products/:id/draft", s.adminPatchProductDraft)
		admin.POST("/products", s.createProduct)
		admin.PUT("/products/:id", s.updateProduct)
		admin.DELETE("/products/:id", s.deleteProduct)
		admin.POST("/products/:id/scripts", s.createScript)
		admin.PUT("/scripts/:id", s.updateScript)
		admin.DELETE("/scripts/:id", s.deleteScript)
		admin.POST("/products/:id/cases", s.createCase)
		admin.PUT("/cases/:id", s.updateCase)
		admin.DELETE("/cases/:id", s.deleteCase)

		admin.GET("/topology/layers", s.adminListLayers)
		admin.POST("/topology/layers", s.adminCreateLayer)
		admin.PUT("/topology/layers/:id", s.adminUpdateLayer)
		admin.DELETE("/topology/layers/:id", s.adminDeleteLayer)
		admin.POST("/topology/categories", s.adminCreateCategory)
		admin.PUT("/topology/categories/:id", s.adminUpdateCategory)
		admin.DELETE("/topology/categories/:id", s.adminDeleteCategory)
		admin.POST("/topology/vendors", s.adminCreateVendor)
		admin.PUT("/topology/vendors/:id", s.adminUpdateVendor)
		admin.DELETE("/topology/vendors/:id", s.adminDeleteVendor)
	}
}

// applyVendorMarketFilter 与首页「国内/国外」分区一致：仅当查询 domestic|foreign 时过滤；all 与空视为不区分
func applyVendorMarketFilter(db *gorm.DB, raw string) *gorm.DB {
	m := strings.TrimSpace(strings.ToLower(raw))
	if m != "domestic" && m != "foreign" {
		return db
	}
	return db.Where("LOWER(COALESCE(NULLIF(TRIM(vendor_market), ''), 'all')) IN ?", []string{"all", m})
}

func normalizeVendorMarket(v string) string {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "domestic", "foreign":
		return strings.ToLower(strings.TrimSpace(v))
	default:
		return "all"
	}
}

type loginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *Server) postLocalLogin(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	var stAuth models.AuthSettings
	if err := s.DB.FirstOrCreate(&stAuth, models.AuthSettings{ID: 1}).Error; err == nil && !stAuth.LocalLoginEnabled {
		c.JSON(http.StatusForbidden, gin.H{"error": "本地账号登录已关闭，请使用企业 SSO 或 LDAP"})
		return
	}
	var u models.User
	if err := s.DB.Where("username = ? AND auth_provider = ?", req.Username, "local").First(&u).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}
	if !u.IsActive || !auth.CheckPassword(u.PasswordHash, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}
	now := time.Now()
	_ = s.DB.Model(&u).Update("last_login", now).Error
	token, err := auth.IssueToken(s.Cfg.JWTSecret, u.ID, u.Username, u.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "签发令牌失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": token, "token_type": "Bearer", "user": u})
}

func (s *Server) postLDAPLogin(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	var st models.AuthSettings
	if err := s.DB.First(&st, 1).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取配置失败"})
		return
	}
	if !st.LdapEnabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "LDAP 未启用"})
		return
	}
	lc, err := auth.ParseLdapConfig([]byte(st.LdapConfig))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "LDAP 配置无效"})
		return
	}
	if err := auth.TryLDAP(lc, req.Username, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "LDAP 认证失败"})
		return
	}
	var u models.User
	err = s.DB.Where("username = ? AND auth_provider = ?", req.Username, "ldap").First(&u).Error
	if err == gorm.ErrRecordNotFound {
		u = models.User{
			Username: req.Username, Email: req.Username + "@ldap.local",
			AuthProvider: "ldap", Role: "user", IsActive: true,
		}
		if err := s.DB.Create(&u).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
			return
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户失败"})
		return
	}
	if !u.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "账号已禁用"})
		return
	}
	now := time.Now()
	_ = s.DB.Model(&u).Update("last_login", now).Error
	token, err := auth.IssueToken(s.Cfg.JWTSecret, u.ID, u.Username, u.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "签发令牌失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": token, "token_type": "Bearer", "user": u})
}

func (s *Server) getOidcStart(c *gin.Context) {
	var st models.AuthSettings
	if err := s.DB.First(&st, 1).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请先配置认证"})
		return
	}
	if !st.OidcEnabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OIDC 未启用"})
		return
	}
	oc, err := auth.ParseOidcConfig([]byte(st.OidcConfig))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OIDC 配置无效"})
		return
	}
	ctx := context.Background()
	oauthCfg, _, err := auth.OidcEndpoints(ctx, oc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	state, err := auth.RandomState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "state 生成失败"})
		return
	}
	s.OidcStates.Put(state)
	url := oauthCfg.AuthCodeURL(state)
	c.Redirect(http.StatusFound, url)
}

func (s *Server) getOidcCallback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")
	if !s.OidcStates.Take(state) {
		c.String(http.StatusBadRequest, "state 无效或过期")
		return
	}
	var st models.AuthSettings
	if err := s.DB.First(&st, 1).Error; err != nil {
		c.String(http.StatusInternalServerError, "配置读取失败")
		return
	}
	oc, err := auth.ParseOidcConfig([]byte(st.OidcConfig))
	if err != nil {
		c.String(http.StatusInternalServerError, "OIDC 配置无效")
		return
	}
	ctx := context.Background()
	oauthCfg, provider, err := auth.OidcEndpoints(ctx, oc)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	ui, err := auth.ExchangeAndVerify(ctx, oauthCfg, provider, code)
	if err != nil {
		c.String(http.StatusUnauthorized, "OIDC 校验失败: "+err.Error())
		return
	}
	username := ui.PreferredName
	if username == "" {
		username = ui.Sub
	}
	var u models.User
	err = s.DB.Where("username = ? AND auth_provider = ?", username, "oidc").First(&u).Error
	if err == gorm.ErrRecordNotFound {
		u = models.User{
			Username: username, Email: ui.Email,
			AuthProvider: "oidc", Role: "user", IsActive: true,
		}
		if err := s.DB.Create(&u).Error; err != nil {
			c.String(http.StatusInternalServerError, "创建用户失败")
			return
		}
	} else if err != nil {
		c.String(http.StatusInternalServerError, "查询用户失败")
		return
	}
	if !u.IsActive {
		c.String(http.StatusForbidden, "账号已禁用")
		return
	}
	now := time.Now()
	_ = s.DB.Model(&u).Update("last_login", now).Error
	token, err := auth.IssueToken(s.Cfg.JWTSecret, u.ID, u.Username, u.Role)
	if err != nil {
		c.String(http.StatusInternalServerError, "签发令牌失败")
		return
	}
	redir := s.Cfg.FrontendOrigin + "/auth/callback?token=" + token
	c.Redirect(http.StatusFound, redir)
}

// getLoginOptions 公开：供登录页决定是否展示本地 / LDAP / SSO 入口
func (s *Server) getLoginOptions(c *gin.Context) {
	var st models.AuthSettings
	if err := s.DB.FirstOrCreate(&st, models.AuthSettings{ID: 1}).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"local_enabled": true,
			"ldap_enabled":  false,
			"oidc_enabled":  false,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"local_enabled": st.LocalLoginEnabled,
		"ldap_enabled":  st.LdapEnabled,
		"oidc_enabled":  st.OidcEnabled,
	})
}

func (s *Server) getMe(c *gin.Context) {
	cl, ok := middleware.Claims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}
	var u models.User
	if err := s.DB.First(&u, cl.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}
	c.JSON(http.StatusOK, u)
}

func (s *Server) getAuthSettings(c *gin.Context) {
	var st models.AuthSettings
	if err := s.DB.FirstOrCreate(&st, models.AuthSettings{ID: 1}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, st)
}

func (s *Server) patchAuthSettings(c *gin.Context) {
	var body struct {
		LocalLoginEnabled *bool            `json:"local_login_enabled"`
		LdapEnabled       *bool            `json:"ldap_enabled"`
		LdapConfig        *datatypes.JSON `json:"ldap_config"`
		OidcEnabled       *bool            `json:"oidc_enabled"`
		OidcConfig        *datatypes.JSON `json:"oidc_config"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	var st models.AuthSettings
	if err := s.DB.FirstOrCreate(&st, models.AuthSettings{ID: 1}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if body.LocalLoginEnabled != nil {
		st.LocalLoginEnabled = *body.LocalLoginEnabled
	}
	if body.LdapEnabled != nil {
		st.LdapEnabled = *body.LdapEnabled
	}
	if body.LdapConfig != nil {
		st.LdapConfig = *body.LdapConfig
	}
	if body.OidcEnabled != nil {
		st.OidcEnabled = *body.OidcEnabled
	}
	if body.OidcConfig != nil {
		st.OidcConfig = *body.OidcConfig
	}
	if err := s.DB.Save(&st).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, st)
}

func (s *Server) listProducts(c *gin.Context) {
	qb := s.DB.Model(&models.Product{}).Where("is_draft = ?", false)
	qb = applyVendorMarketFilter(qb, c.Query("vendor_market"))
	if sid := strings.TrimSpace(c.Query("solution_category_id")); sid != "" {
		qb = qb.Where("solution_category_id = ?", sid)
	}
	if mf := strings.TrimSpace(c.Query("manufacturer")); mf != "" {
		qb = qb.Where("LOWER(TRIM(manufacturer_name)) = LOWER(?)", mf)
	}
	qb = qb.Order("updated_at desc")
	var list []models.Product
	if err := qb.Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (s *Server) searchProducts(c *gin.Context) {
	q := strings.TrimSpace(strings.ToLower(c.Query("q")))
	sid := strings.TrimSpace(c.Query("solution_category_id"))
	if q == "" && sid == "" {
		s.listProducts(c)
		return
	}
	vm := strings.TrimSpace(c.Query("vendor_market"))
	mfg := strings.TrimSpace(c.Query("manufacturer"))
	cacheKey := "search:" + q + ":sc:" + sid + ":vm:" + vm + ":mfg:" + mfg
	if b, hit := s.SearchCache.Get(cacheKey); hit {
		c.Data(http.StatusOK, "application/json; charset=utf-8", b)
		return
	}
	db := s.DB.Model(&models.Product{}).Where("is_draft = ?", false)
	db = applyVendorMarketFilter(db, vm)
	if mfg != "" {
		db = db.Where("LOWER(TRIM(manufacturer_name)) = LOWER(?)", mfg)
	}
	if q != "" {
		pat := "%" + q + "%"
		db = db.Where("LOWER(name) LIKE ? OR LOWER(category) LIKE ? OR LOWER(description) LIKE ?", pat, pat, pat)
	}
	if sid != "" {
		db = db.Where("solution_category_id = ?", sid)
	}
	var list []models.Product
	if err := db.Order("updated_at desc").Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	raw, _ := json.Marshal(list)
	s.SearchCache.Set(cacheKey, raw)
	c.Data(http.StatusOK, "application/json; charset=utf-8", raw)
}

func (s *Server) getProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var p models.Product
	if err := s.DB.Preload("SolutionCategory").First(&p, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到"})
		return
	}
	if p.IsDraft {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到"})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (s *Server) adminListProducts(c *gin.Context) {
	var list []models.Product
	if err := s.DB.Preload("SolutionCategory").Order("updated_at desc").Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (s *Server) adminGetProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var p models.Product
	if err := s.DB.Preload("SolutionCategory").First(&p, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到"})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (s *Server) adminListCasesByProduct(c *gin.Context) {
	pid, _ := strconv.Atoi(c.Param("id"))
	var list []models.Case
	if err := s.DB.Where("product_id = ?", pid).Order("id asc").Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (s *Server) adminListScriptsByProduct(c *gin.Context) {
	pid, _ := strconv.Atoi(c.Param("id"))
	var list []models.SalesScript
	if err := s.DB.Where("product_id = ?", pid).Order("id asc").Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (s *Server) adminPatchProductDraft(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var body struct {
		IsDraft *bool `json:"is_draft"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.IsDraft == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供 is_draft 布尔字段"})
		return
	}
	var p models.Product
	if err := s.DB.First(&p, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到"})
		return
	}
	if err := s.DB.Model(&models.Product{}).Where("id = ?", id).Update("is_draft", *body.IsDraft).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_ = s.DB.Preload("SolutionCategory").First(&p, id).Error
	c.JSON(http.StatusOK, p)
}

func (s *Server) syncProductCategoryFromSolution(p *models.Product) {
	if p.SolutionCategoryID == nil || *p.SolutionCategoryID == 0 {
		return
	}
	var cat models.SolutionCategory
	if err := s.DB.First(&cat, *p.SolutionCategoryID).Error; err == nil {
		p.Category = cat.Label
	}
}

func (s *Server) createProduct(c *gin.Context) {
	var p models.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p.VendorMarket = normalizeVendorMarket(p.VendorMarket)
	s.syncProductCategoryFromSolution(&p)
	if err := s.DB.Create(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_ = s.DB.Preload("SolutionCategory").First(&p, p.ID).Error
	c.JSON(http.StatusCreated, p)
}

func (s *Server) updateProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var p models.Product
	if err := s.DB.First(&p, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到"})
		return
	}
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p.ID = uint(id)
	p.VendorMarket = normalizeVendorMarket(p.VendorMarket)
	s.syncProductCategoryFromSolution(&p)
	up := map[string]any{
		"name":                    p.Name,
		"category":                p.Category,
		"solution_category_id":      p.SolutionCategoryID,
		"vendor_market":             p.VendorMarket,
		"manufacturer_name":         p.ManufacturerName,
		"manufacturer_logo":         p.ManufacturerLogo,
		"sales_contact_name":        p.SalesContactName,
		"sales_contact_phone":       p.SalesContactPhone,
		"sales_contact_email":       p.SalesContactEmail,
		"presales_contact_name":     p.PresalesContactName,
		"presales_contact_phone":    p.PresalesContactPhone,
		"presales_contact_email":    p.PresalesContactEmail,
		"description":               p.Description,
		"highlights":                p.Highlights,
		"target_personas":           p.TargetPersonas,
		"trigger_events":            p.TriggerEvents,
		"discovery_questions":       p.DiscoveryQuestions,
		"competitor_analysis":       p.CompetitorAnalysis,
		"roi_metrics":               p.ROIMetrics,
		"is_draft":                  p.IsDraft,
	}
	if err := s.DB.Model(&models.Product{}).Where("id = ?", id).Updates(up).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_ = s.DB.Preload("SolutionCategory").First(&p, id).Error
	c.JSON(http.StatusOK, p)
}

func (s *Server) deleteProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := s.DB.Delete(&models.Product{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_ = s.DB.Where("product_id = ?", id).Delete(&models.SalesScript{}).Error
	_ = s.DB.Where("product_id = ?", id).Delete(&models.Case{}).Error
	c.Status(http.StatusNoContent)
}

func (s *Server) listScripts(c *gin.Context) {
	pid, _ := strconv.Atoi(c.Param("id"))
	var prod models.Product
	if err := s.DB.Select("id", "is_draft").First(&prod, pid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到"})
		return
	}
	if prod.IsDraft {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到"})
		return
	}
	var list []models.SalesScript
	if err := s.DB.Where("product_id = ?", pid).Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (s *Server) createScript(c *gin.Context) {
	pid, _ := strconv.Atoi(c.Param("id"))
	var sc models.SalesScript
	if err := c.ShouldBindJSON(&sc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sc.ProductID = uint(pid)
	if err := s.DB.Create(&sc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, sc)
}

func (s *Server) updateScript(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var sc models.SalesScript
	if err := s.DB.First(&sc, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到"})
		return
	}
	if err := c.ShouldBindJSON(&sc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sc.ID = uint(id)
	if err := s.DB.Model(&sc).Updates(map[string]any{"scenario": sc.Scenario, "content": sc.Content}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_ = s.DB.First(&sc, id).Error
	c.JSON(http.StatusOK, sc)
}

func (s *Server) deleteScript(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	_ = s.DB.Delete(&models.SalesScript{}, id).Error
	c.Status(http.StatusNoContent)
}

func (s *Server) listCases(c *gin.Context) {
	pid, _ := strconv.Atoi(c.Param("id"))
	var prod models.Product
	if err := s.DB.Select("id", "is_draft").First(&prod, pid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到"})
		return
	}
	if prod.IsDraft {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到"})
		return
	}
	var list []models.Case
	if err := s.DB.Where("product_id = ?", pid).Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (s *Server) createCase(c *gin.Context) {
	pid, _ := strconv.Atoi(c.Param("id"))
	var ca models.Case
	if err := c.ShouldBindJSON(&ca); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ca.ProductID = uint(pid)
	if err := s.DB.Create(&ca).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, ca)
}

func (s *Server) updateCase(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var ca models.Case
	if err := s.DB.First(&ca, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到"})
		return
	}
	if err := c.ShouldBindJSON(&ca); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ca.ID = uint(id)
	if err := s.DB.Model(&ca).Updates(map[string]any{
		"client_name": ca.ClientName, "pain_points": ca.PainPoints, "solution": ca.Solution, "value_delivered": ca.ValueDelivered,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_ = s.DB.First(&ca, id).Error
	c.JSON(http.StatusOK, ca)
}

func (s *Server) deleteCase(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	_ = s.DB.Delete(&models.Case{}, id).Error
	c.Status(http.StatusNoContent)
}
