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

	api.POST("/auth/login", s.postLocalLogin)
	api.POST("/auth/login/ldap", s.postLDAPLogin)

	api.GET("/auth/oidc/start", s.getOidcStart)
	api.GET("/auth/oidc/callback", s.getOidcCallback)

	authorized := api.Group("")
	authorized.Use(middleware.JWTAuth(s.Cfg.JWTSecret))
	{
		authorized.GET("/me", s.getMe)
		authorized.GET("/products", s.listProducts)
		authorized.GET("/products/search", s.searchProducts)
		authorized.GET("/products/:id", s.getProduct)
		authorized.GET("/products/:id/scripts", s.listScripts)
		authorized.GET("/products/:id/cases", s.listCases)
	}

	admin := api.Group("/admin")
	admin.Use(middleware.JWTAuth(s.Cfg.JWTSecret), middleware.RequireAdmin())
	{
		admin.GET("/auth-settings", s.getAuthSettings)
		admin.PATCH("/auth-settings", s.patchAuthSettings)
		admin.POST("/products", s.createProduct)
		admin.PUT("/products/:id", s.updateProduct)
		admin.DELETE("/products/:id", s.deleteProduct)
		admin.POST("/products/:id/scripts", s.createScript)
		admin.PUT("/scripts/:id", s.updateScript)
		admin.DELETE("/scripts/:id", s.deleteScript)
		admin.POST("/products/:id/cases", s.createCase)
		admin.PUT("/cases/:id", s.updateCase)
		admin.DELETE("/cases/:id", s.deleteCase)
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
		LdapEnabled *bool            `json:"ldap_enabled"`
		LdapConfig  *datatypes.JSON `json:"ldap_config"`
		OidcEnabled *bool            `json:"oidc_enabled"`
		OidcConfig  *datatypes.JSON `json:"oidc_config"`
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
	var list []models.Product
	if err := s.DB.Order("updated_at desc").Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (s *Server) searchProducts(c *gin.Context) {
	q := strings.TrimSpace(strings.ToLower(c.Query("q")))
	if q == "" {
		s.listProducts(c)
		return
	}
	cacheKey := "search:" + q
	if b, hit := s.SearchCache.Get(cacheKey); hit {
		c.Data(http.StatusOK, "application/json; charset=utf-8", b)
		return
	}
	var list []models.Product
	pat := "%" + q + "%"
	if err := s.DB.Where("LOWER(name) LIKE ? OR LOWER(category) LIKE ? OR LOWER(description) LIKE ?", pat, pat, pat).
		Order("updated_at desc").Find(&list).Error; err != nil {
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
	if err := s.DB.First(&p, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到"})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (s *Server) createProduct(c *gin.Context) {
	var p models.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.DB.Create(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
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
	if err := s.DB.Model(&models.Product{}).Where("id = ?", id).Updates(map[string]any{
		"name": p.Name, "category": p.Category, "description": p.Description,
		"highlights": p.Highlights, "target_personas": p.TargetPersonas, "trigger_events": p.TriggerEvents,
		"discovery_questions": p.DiscoveryQuestions, "competitor_analysis": p.CompetitorAnalysis, "roi_metrics": p.ROIMetrics,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_ = s.DB.First(&p, id).Error
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
