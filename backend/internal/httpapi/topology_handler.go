package httpapi

import (
	"github.com/6547709/sales-pilot/backend/internal/apperr"
	"github.com/6547709/sales-pilot/backend/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handlers) getTopology(c *gin.Context) {
	tree, err := h.Topology.Tree(c.Request.Context())
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, tree)
}

func (h *Handlers) listSolutionCategories(c *gin.Context) {
	activeOnly := c.Query("active_only") == "1" || c.Query("active_only") == "true"
	list, err := h.Topology.ListCategories(c.Request.Context(), activeOnly)
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, list)
}

func (h *Handlers) adminListLayers(c *gin.Context) {
	list, err := h.Topology.ListLayers(c.Request.Context())
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, list)
}

func (h *Handlers) adminCreateLayer(c *gin.Context) {
	var body models.TopologyLayer
	if err := c.ShouldBindJSON(&body); err != nil {
		Fail(c, apperr.Invalid("invalid_body", "参数错误"))
		return
	}
	created, err := h.Topology.CreateLayer(c.Request.Context(), &body)
	if err != nil {
		Fail(c, err)
		return
	}
	Created(c, created)
}

func (h *Handlers) adminUpdateLayer(c *gin.Context) {
	id, ok := ParseID(c, "id")
	if !ok {
		return
	}
	var body models.TopologyLayer
	if err := c.ShouldBindJSON(&body); err != nil {
		Fail(c, apperr.Invalid("invalid_body", "参数错误"))
		return
	}
	updated, err := h.Topology.UpdateLayer(c.Request.Context(), id, &body)
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, updated)
}

func (h *Handlers) adminDeleteLayer(c *gin.Context) {
	id, ok := ParseID(c, "id")
	if !ok {
		return
	}
	if err := h.Topology.DeleteLayer(c.Request.Context(), id); err != nil {
		Fail(c, err)
		return
	}
	NoContent(c)
}

func (h *Handlers) adminCreateCategory(c *gin.Context) {
	var body models.SolutionCategory
	if err := c.ShouldBindJSON(&body); err != nil {
		Fail(c, apperr.Invalid("invalid_body", "参数错误"))
		return
	}
	created, err := h.Topology.CreateCategory(c.Request.Context(), &body)
	if err != nil {
		Fail(c, err)
		return
	}
	Created(c, created)
}

func (h *Handlers) adminUpdateCategory(c *gin.Context) {
	id, ok := ParseID(c, "id")
	if !ok {
		return
	}
	var body models.SolutionCategory
	if err := c.ShouldBindJSON(&body); err != nil {
		Fail(c, apperr.Invalid("invalid_body", "参数错误"))
		return
	}
	updated, err := h.Topology.UpdateCategory(c.Request.Context(), id, &body)
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, updated)
}

func (h *Handlers) adminDeleteCategory(c *gin.Context) {
	id, ok := ParseID(c, "id")
	if !ok {
		return
	}
	if err := h.Topology.DeleteCategory(c.Request.Context(), id); err != nil {
		Fail(c, err)
		return
	}
	NoContent(c)
}

type vendorBody struct {
	CategoryID uint   `json:"category_id"`
	Market     string `json:"market"`
	Name       string `json:"name"`
	SortOrder  int    `json:"sort_order"`
}

func (h *Handlers) adminCreateVendor(c *gin.Context) {
	var body vendorBody
	if err := c.ShouldBindJSON(&body); err != nil {
		Fail(c, apperr.Invalid("invalid_body", "参数错误"))
		return
	}
	created, err := h.Topology.CreateVendor(c.Request.Context(), &models.SolutionVendor{
		CategoryID: body.CategoryID, Market: body.Market, Name: body.Name, SortOrder: body.SortOrder,
	})
	if err != nil {
		Fail(c, err)
		return
	}
	Created(c, created)
}

func (h *Handlers) adminUpdateVendor(c *gin.Context) {
	id, ok := ParseID(c, "id")
	if !ok {
		return
	}
	var body models.SolutionVendor
	if err := c.ShouldBindJSON(&body); err != nil {
		Fail(c, apperr.Invalid("invalid_body", "参数错误"))
		return
	}
	updated, err := h.Topology.UpdateVendor(c.Request.Context(), id, &body)
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, updated)
}

func (h *Handlers) adminDeleteVendor(c *gin.Context) {
	id, ok := ParseID(c, "id")
	if !ok {
		return
	}
	if err := h.Topology.DeleteVendor(c.Request.Context(), id); err != nil {
		Fail(c, err)
		return
	}
	NoContent(c)
}
