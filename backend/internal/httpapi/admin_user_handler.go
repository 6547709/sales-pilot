package httpapi

import (
	"strings"

	"github.com/6547709/sales-pilot/backend/internal/apperr"
	"github.com/6547709/sales-pilot/backend/internal/models"
	"github.com/6547709/sales-pilot/backend/internal/paging"
	"github.com/6547709/sales-pilot/backend/internal/repository"
	"github.com/6547709/sales-pilot/backend/internal/service"
	"github.com/gin-gonic/gin"
)

func (h *Handlers) adminListLocalUsers(c *gin.Context) {
	ctx := c.Request.Context()
	if strings.TrimSpace(c.Query("page")) == "" {
		list, err := h.Users.ListLocal(ctx)
		if err != nil {
			Fail(c, err)
			return
		}
		OK(c, list)
		return
	}
	p := paging.New(QueryInt(c, "page", 1), QueryInt(c, "page_size", 20))
	f := repository.LocalUserFilter{Q: c.Query("q"), Role: c.Query("role")}
	if raw := c.Query("active"); raw == "true" || raw == "false" {
		v := raw == "true"
		f.Active = &v
	}
	list, total, err := h.Users.ListLocalPage(ctx, f, p)
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, paging.Page[models.User]{Items: list, Total: total, Page: p.Page, PageSize: p.PageSize})
}

func (h *Handlers) adminCreateLocalUser(c *gin.Context) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Role     string `json:"role"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		Fail(c, apperr.Invalid("invalid_body", "用户名、密码必填，密码至少 6 位"))
		return
	}
	u, err := h.Users.CreateLocal(c.Request.Context(), service.CreateLocalUserInput{
		Username: body.Username, Password: body.Password, Email: body.Email, Role: body.Role,
	})
	if err != nil {
		Fail(c, err)
		return
	}
	Created(c, u)
}

func (h *Handlers) adminUpdateLocalUser(c *gin.Context) {
	id, ok := ParseID(c, "id")
	if !ok {
		return
	}
	var body struct {
		Email    *string `json:"email"`
		Role     *string `json:"role"`
		IsActive *bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		Fail(c, apperr.Invalid("invalid_body", "参数错误"))
		return
	}
	u, err := h.Users.UpdateLocal(c.Request.Context(), id, service.UpdateLocalUserInput{
		Email: body.Email, Role: body.Role, IsActive: body.IsActive,
	})
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, u)
}

func (h *Handlers) adminResetLocalUserPassword(c *gin.Context) {
	id, ok := ParseID(c, "id")
	if !ok {
		return
	}
	var body struct {
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		Fail(c, apperr.Invalid("invalid_body", "参数错误"))
		return
	}
	if err := h.Users.ResetPassword(c.Request.Context(), id, body.Password); err != nil {
		Fail(c, err)
		return
	}
	OK(c, gin.H{"ok": true})
}

func (h *Handlers) adminDeleteLocalUser(c *gin.Context) {
	id, ok := ParseID(c, "id")
	if !ok {
		return
	}
	if err := h.Users.DeleteLocal(c.Request.Context(), id); err != nil {
		Fail(c, err)
		return
	}
	NoContent(c)
}
