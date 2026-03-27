package handlers

import (
	"net/http"
	"strconv"

	"github.com/6547709/sales-pilot/backend/internal/auth"
	"github.com/6547709/sales-pilot/backend/internal/middleware"
	"github.com/6547709/sales-pilot/backend/internal/models"
	"github.com/gin-gonic/gin"
)

// 本地用户列表（不含密码）
type localUserOut struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (s *Server) adminListLocalUsers(c *gin.Context) {
	var list []models.User
	if err := s.DB.Where("auth_provider = ?", "local").Order("id asc").Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	out := make([]localUserOut, 0, len(list))
	for _, u := range list {
		out = append(out, localUserOut{
			ID: u.ID, Username: u.Username, Email: u.Email, Role: u.Role, IsActive: u.IsActive,
			CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: u.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	c.JSON(http.StatusOK, out)
}

type createLocalUserReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email"`
	Role     string `json:"role"` // admin | user，默认 user
}

func (s *Server) adminCreateLocalUser(c *gin.Context) {
	var req createLocalUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名、密码必填，密码至少 6 位"})
		return
	}
	role := req.Role
	if role != "admin" && role != "user" {
		role = "user"
	}
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码处理失败"})
		return
	}
	email := req.Email
	if email == "" {
		email = req.Username + "@local"
	}
	u := models.User{
		Username: req.Username, Email: email, PasswordHash: hash,
		AuthProvider: "local", Role: role, IsActive: true,
	}
	if err := s.DB.Create(&u).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "创建失败，用户名可能已存在"})
		return
	}
	c.JSON(http.StatusCreated, localUserOut{
		ID: u.ID, Username: u.Username, Email: u.Email, Role: u.Role, IsActive: u.IsActive,
		CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: u.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

type updateLocalUserReq struct {
	Email    *string `json:"email"`
	Role     *string `json:"role"`
	IsActive *bool   `json:"is_active"`
}

func (s *Server) adminUpdateLocalUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var u models.User
	if err := s.DB.First(&u, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到"})
		return
	}
	if u.AuthProvider != "local" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "仅可维护本地用户"})
		return
	}
	var req updateLocalUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	cl, _ := middleware.Claims(c)

	if req.Role != nil && *req.Role != "admin" && *req.Role != "user" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "角色只能是 admin 或 user"})
		return
	}
	// 降级或停用管理员前，须仍有其他活跃本地管理员
	if req.Role != nil && *req.Role != "admin" && u.Role == "admin" {
		var otherAdmins int64
		s.DB.Model(&models.User{}).Where("auth_provider = ? AND role = ? AND is_active = ? AND id != ?", "local", "admin", true, id).Count(&otherAdmins)
		if otherAdmins == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "至少保留一名可用的本地管理员"})
			return
		}
	}
	if req.IsActive != nil && !*req.IsActive && uint(id) == cl.UserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能停用当前登录账号"})
		return
	}
	if req.IsActive != nil && !*req.IsActive && u.Role == "admin" {
		var otherAdmins int64
		s.DB.Model(&models.User{}).Where("auth_provider = ? AND role = ? AND is_active = ? AND id != ?", "local", "admin", true, id).Count(&otherAdmins)
		if otherAdmins == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "至少保留一名可用的本地管理员"})
			return
		}
	}

	up := map[string]any{}
	if req.Email != nil {
		up["email"] = *req.Email
	}
	if req.Role != nil {
		up["role"] = *req.Role
	}
	if req.IsActive != nil {
		up["is_active"] = *req.IsActive
	}
	if len(up) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无更新字段"})
		return
	}
	if err := s.DB.Model(&u).Updates(up).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_ = s.DB.First(&u, id).Error
	c.JSON(http.StatusOK, localUserOut{
		ID: u.ID, Username: u.Username, Email: u.Email, Role: u.Role, IsActive: u.IsActive,
		CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: u.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

type resetPasswordReq struct {
	Password string `json:"password" binding:"required,min=6"`
}

func (s *Server) adminResetLocalUserPassword(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var u models.User
	if err := s.DB.First(&u, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到"})
		return
	}
	if u.AuthProvider != "local" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "仅可重置本地用户密码"})
		return
	}
	var req resetPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "新密码至少 6 位"})
		return
	}
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码处理失败"})
		return
	}
	if err := s.DB.Model(&u).Update("password_hash", hash).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (s *Server) adminDeleteLocalUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	cl, _ := middleware.Claims(c)
	if uint(id) == cl.UserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能删除当前登录账号"})
		return
	}
	var u models.User
	if err := s.DB.First(&u, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到"})
		return
	}
	if u.AuthProvider != "local" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "仅可删除本地用户"})
		return
	}
	if u.Role == "admin" {
		var cnt int64
		s.DB.Model(&models.User{}).Where("auth_provider = ? AND role = ? AND is_active = ?", "local", "admin", true).Count(&cnt)
		if cnt <= 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "至少保留一名本地管理员"})
			return
		}
	}
	if err := s.DB.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
