package httpapi

import (
	"strconv"
	"strings"

	"github.com/6547709/sales-pilot/backend/internal/httpapi/dto"
	"github.com/6547709/sales-pilot/backend/internal/models"
	"github.com/6547709/sales-pilot/backend/internal/service"
	"github.com/6547709/sales-pilot/backend/internal/view"
	"github.com/gin-gonic/gin"
)

func publicFilter(c *gin.Context) service.PublicProductFilter {
	f := service.PublicProductFilter{
		Market:       c.Query("vendor_market"),
		Manufacturer: c.Query("manufacturer"),
		Keyword:      c.Query("q"),
	}
	if sid := strings.TrimSpace(c.Query("solution_category_id")); sid != "" {
		if id, err := strconv.ParseUint(sid, 10, 64); err == nil && id > 0 {
			v := uint(id)
			f.SolutionCategoryID = &v
		}
	}
	return f
}

func (h *Handlers) listProducts(c *gin.Context) {
	mode := view.Parse(c.Query("view"))
	list, err := h.Products.ListPublic(c.Request.Context(), publicFilter(c))
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, dto.ProjectProducts(list, mode))
}

func (h *Handlers) searchProducts(c *gin.Context) {
	q := strings.TrimSpace(c.Query("q"))
	sid := strings.TrimSpace(c.Query("solution_category_id"))
	if q == "" && sid == "" {
		h.listProducts(c)
		return
	}
	h.listProducts(c)
}

func (h *Handlers) getProduct(c *gin.Context) {
	id, ok := ParseID(c, "id")
	if !ok {
		return
	}
	mode := view.Parse(c.Query("view"))
	p, err := h.Products.GetPublic(c.Request.Context(), id)
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, dto.ProjectProduct(*p, mode))
}

func (h *Handlers) listScripts(c *gin.Context) {
	mode := view.Parse(c.Query("view"))
	if mode.IsCustomer() {
		OK(c, []models.SalesScript{})
		return
	}
	id, ok := ParseID(c, "id")
	if !ok {
		return
	}
	list, err := h.Products.ListScriptsPublic(c.Request.Context(), id)
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, list)
}

func (h *Handlers) listCases(c *gin.Context) {
	id, ok := ParseID(c, "id")
	if !ok {
		return
	}
	mode := view.Parse(c.Query("view"))
	list, err := h.Products.ListCasesPublic(c.Request.Context(), id)
	if err != nil {
		Fail(c, err)
		return
	}
	if mode.IsCustomer() {
		public := make([]models.Case, 0, len(list))
		for _, item := range list {
			if item.IsPublic {
				public = append(public, item)
			}
		}
		list = public
	}
	OK(c, dto.ProjectCases(list, mode))
}

func (h *Handlers) getProductScenarios(c *gin.Context) {
	id, ok := ParseID(c, "id")
	if !ok {
		return
	}
	mode := view.Parse(c.Query("view"))
	list, err := h.Scenarios.ListForProduct(c.Request.Context(), id, mode)
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, dto.ProjectScenarios(list, mode))
}
