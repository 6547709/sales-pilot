package httpapi

import (
	"strconv"
	"strings"

	"github.com/6547709/sales-pilot/backend/internal/service"
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
	list, err := h.Products.ListPublic(c.Request.Context(), publicFilter(c))
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, list)
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
	p, err := h.Products.GetPublic(c.Request.Context(), id)
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, p)
}

func (h *Handlers) listScripts(c *gin.Context) {
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
	list, err := h.Products.ListCasesPublic(c.Request.Context(), id)
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, list)
}
