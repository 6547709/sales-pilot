package httpapi

import (
	"github.com/6547709/sales-pilot/backend/internal/apperr"
	"github.com/6547709/sales-pilot/backend/internal/httpapi/dto"
	"github.com/6547709/sales-pilot/backend/internal/models"
	"github.com/6547709/sales-pilot/backend/internal/view"
	"github.com/gin-gonic/gin"
)

func (h *Handlers) listServices(c *gin.Context) {
	mode := view.Parse(c.Query("view"))
	list, err := h.Services.List(c.Request.Context(), c.Query("kind"), true)
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, dto.ProjectServices(list, mode))
}

func (h *Handlers) getService(c *gin.Context) {
	mode := view.Parse(c.Query("view"))
	svc, err := h.Services.GetBySlug(c.Request.Context(), c.Param("slug"))
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, dto.ProjectService(*svc, mode))
}

func (h *Handlers) adminListServices(c *gin.Context) {
	list, err := h.Services.List(c.Request.Context(), c.Query("kind"), false)
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, list)
}

type serviceBody struct {
	models.Service
	CategoryIDs []uint `json:"category_ids"`
}

func (h *Handlers) adminCreateService(c *gin.Context) {
	var body serviceBody
	if err := c.ShouldBindJSON(&body); err != nil {
		Fail(c, apperr.Invalid("invalid_body", "参数错误"))
		return
	}
	created, err := h.Services.Create(c.Request.Context(), &body.Service, body.CategoryIDs)
	if err != nil {
		Fail(c, err)
		return
	}
	Created(c, created)
}

func (h *Handlers) adminUpdateService(c *gin.Context) {
	id, ok := ParseID(c, "id")
	if !ok {
		return
	}
	var body serviceBody
	if err := c.ShouldBindJSON(&body); err != nil {
		Fail(c, apperr.Invalid("invalid_body", "参数错误"))
		return
	}
	updated, err := h.Services.Update(c.Request.Context(), id, &body.Service, body.CategoryIDs)
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, updated)
}

func (h *Handlers) adminDeleteService(c *gin.Context) {
	id, ok := ParseID(c, "id")
	if !ok {
		return
	}
	if err := h.Services.Delete(c.Request.Context(), id); err != nil {
		Fail(c, err)
		return
	}
	NoContent(c)
}

func (h *Handlers) adminSetServiceCategories(c *gin.Context) {
	id, ok := ParseID(c, "id")
	if !ok {
		return
	}
	var body struct {
		CategoryIDs []uint `json:"category_ids"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		Fail(c, apperr.Invalid("invalid_body", "参数错误"))
		return
	}
	if err := h.Services.SetCategories(c.Request.Context(), id, body.CategoryIDs); err != nil {
		Fail(c, err)
		return
	}
	OK(c, gin.H{"ok": true})
}

func (h *Handlers) getServiceScenarios(c *gin.Context) {
	mode := view.Parse(c.Query("view"))
	svc, err := h.Services.GetBySlug(c.Request.Context(), c.Param("slug"))
	if err != nil {
		Fail(c, err)
		return
	}
	list, err := h.Scenarios.ListForService(c.Request.Context(), svc.ID, mode)
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, dto.ProjectScenarios(list, mode))
}
