package httpapi

import (
	"github.com/6547709/sales-pilot/backend/internal/apperr"
	"github.com/6547709/sales-pilot/backend/internal/httpapi/dto"
	"github.com/6547709/sales-pilot/backend/internal/models"
	"github.com/6547709/sales-pilot/backend/internal/view"
	"github.com/gin-gonic/gin"
)

type scenarioOut struct {
	dto.Scenario
	ProductCount int64 `json:"product_count"`
	ServiceCount int64 `json:"service_count"`
}

type scenarioGroupOut struct {
	ID          uint          `json:"id"`
	Slug        string        `json:"slug"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	SortOrder   int           `json:"sort_order"`
	Scenarios   []scenarioOut `json:"scenarios"`
}

func (h *Handlers) listScenarios(c *gin.Context) {
	mode := view.Parse(c.Query("view"))
	groups, err := h.Scenarios.Groups(c.Request.Context())
	if err != nil {
		Fail(c, err)
		return
	}
	out := make([]scenarioGroupOut, 0, len(groups))
	for _, g := range groups {
		items := make([]scenarioOut, 0, len(g.Scenarios))
		for _, item := range g.Scenarios {
			if mode.IsCustomer() && !item.Scenario.IsActive {
				continue
			}
			items = append(items, scenarioOut{
				Scenario:     dto.ProjectScenario(item.Scenario, mode),
				ProductCount: item.ProductCount,
				ServiceCount: item.ServiceCount,
			})
		}
		out = append(out, scenarioGroupOut{
			ID: g.Group.ID, Slug: g.Group.Slug, Name: g.Group.Name,
			Description: g.Group.Description, SortOrder: g.Group.SortOrder,
			Scenarios: items,
		})
	}
	OK(c, gin.H{"groups": out})
}

func (h *Handlers) getScenario(c *gin.Context) {
	mode := view.Parse(c.Query("view"))
	detail, err := h.Scenarios.GetBySlug(c.Request.Context(), c.Param("slug"), mode)
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, gin.H{
		"scenario": dto.ProjectScenario(detail.Scenario, mode),
		"group":    detail.Group,
		"services": dto.ProjectServices(detail.Services, mode),
		"products": dto.ProjectProducts(detail.Products, mode),
		"cases":    dto.ProjectCases(detail.Cases, mode),
	})
}

func (h *Handlers) adminGetScenario(c *gin.Context) {
	id, ok := ParseID(c, "id")
	if !ok {
		return
	}
	detail, err := h.Scenarios.GetAdminDetail(c.Request.Context(), id)
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, gin.H{
		"scenario": detail.Scenario,
		"group":    detail.Group,
		"services": detail.Services,
		"products": detail.Products,
	})
}

func (h *Handlers) adminListScenarios(c *gin.Context) {
	h.listScenarios(c)
}

func (h *Handlers) adminCreateScenarioGroup(c *gin.Context) {
	var body struct {
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		Fail(c, apperr.Invalid("invalid_body", "参数错误"))
		return
	}
	created, err := h.Scenarios.CreateGroup(c.Request.Context(), body.Name, body.Slug, body.Description)
	if err != nil {
		Fail(c, err)
		return
	}
	Created(c, created)
}

func (h *Handlers) adminUpdateScenarioGroup(c *gin.Context) {
	id, ok := ParseID(c, "id")
	if !ok {
		return
	}
	var body models.ScenarioGroup
	if err := c.ShouldBindJSON(&body); err != nil {
		Fail(c, apperr.Invalid("invalid_body", "参数错误"))
		return
	}
	updated, err := h.Scenarios.UpdateGroup(c.Request.Context(), id, &body)
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, updated)
}

func (h *Handlers) adminDeleteScenarioGroup(c *gin.Context) {
	id, ok := ParseID(c, "id")
	if !ok {
		return
	}
	if err := h.Scenarios.DeleteGroup(c.Request.Context(), id); err != nil {
		Fail(c, err)
		return
	}
	NoContent(c)
}

func (h *Handlers) adminCreateScenario(c *gin.Context) {
	var body models.Scenario
	if err := c.ShouldBindJSON(&body); err != nil {
		Fail(c, apperr.Invalid("invalid_body", "参数错误"))
		return
	}
	created, err := h.Scenarios.Create(c.Request.Context(), &body)
	if err != nil {
		Fail(c, err)
		return
	}
	Created(c, created)
}

func (h *Handlers) adminUpdateScenario(c *gin.Context) {
	id, ok := ParseID(c, "id")
	if !ok {
		return
	}
	var body models.Scenario
	if err := c.ShouldBindJSON(&body); err != nil {
		Fail(c, apperr.Invalid("invalid_body", "参数错误"))
		return
	}
	updated, err := h.Scenarios.Update(c.Request.Context(), id, &body)
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, updated)
}

func (h *Handlers) adminDeleteScenario(c *gin.Context) {
	id, ok := ParseID(c, "id")
	if !ok {
		return
	}
	if err := h.Scenarios.Delete(c.Request.Context(), id); err != nil {
		Fail(c, err)
		return
	}
	NoContent(c)
}

func (h *Handlers) adminSetScenarioLinks(c *gin.Context) {
	id, ok := ParseID(c, "id")
	if !ok {
		return
	}
	var body struct {
		ProductIDs []uint `json:"product_ids"`
		ServiceIDs []uint `json:"service_ids"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		Fail(c, apperr.Invalid("invalid_body", "参数错误"))
		return
	}
	if err := h.Scenarios.SetLinks(c.Request.Context(), id, body.ProductIDs, body.ServiceIDs); err != nil {
		Fail(c, err)
		return
	}
	OK(c, gin.H{"ok": true})
}

func (h *Handlers) adminListScenarioGroups(c *gin.Context) {
	groups, err := h.Scenarios.Groups(c.Request.Context())
	if err != nil {
		Fail(c, err)
		return
	}
	list := make([]gin.H, 0, len(groups))
	for _, g := range groups {
		list = append(list, gin.H{
			"id": g.Group.ID, "slug": g.Group.Slug, "name": g.Group.Name,
			"description": g.Group.Description, "sort_order": g.Group.SortOrder,
			"scenario_count": len(g.Scenarios),
		})
	}
	OK(c, list)
}

func (h *Handlers) adminGetProductLinks(c *gin.Context) {
	id, ok := ParseID(c, "id")
	if !ok {
		return
	}
	ids, err := h.Scenarios.ScenarioIDsForProduct(c.Request.Context(), id)
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, gin.H{"scenario_ids": ids})
}

func (h *Handlers) adminSetProductLinks(c *gin.Context) {
	id, ok := ParseID(c, "id")
	if !ok {
		return
	}
	var body struct {
		ScenarioIDs []uint `json:"scenario_ids"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		Fail(c, apperr.Invalid("invalid_body", "参数错误"))
		return
	}
	if err := h.Scenarios.SetProductScenarios(c.Request.Context(), id, body.ScenarioIDs); err != nil {
		Fail(c, err)
		return
	}
	OK(c, gin.H{"ok": true})
}
