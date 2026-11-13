package httpapi

import (
	"strings"

	"github.com/6547709/sales-pilot/backend/internal/apperr"
	"github.com/6547709/sales-pilot/backend/internal/models"
	"github.com/6547709/sales-pilot/backend/internal/paging"
	"github.com/6547709/sales-pilot/backend/internal/repository"
	"github.com/gin-gonic/gin"
)

func adminProductFilter(c *gin.Context) repository.AdminProductFilter {
	market := c.Query("market")
	if market == "" {
		market = c.Query("vendor_market")
	}
	return repository.AdminProductFilter{
		Q:      c.Query("q"),
		Status: c.Query("status"),
		Market: market,
	}
}

func (h *Handlers) adminListProducts(c *gin.Context) {
	ctx := c.Request.Context()
	if strings.TrimSpace(c.Query("page")) == "" {
		list, err := h.Products.ListAdminAll(ctx, adminProductFilter(c))
		if err != nil {
			Fail(c, err)
			return
		}
		OK(c, list)
		return
	}
	p := paging.New(QueryInt(c, "page", 1), QueryInt(c, "page_size", 20))
	p.Sort = c.Query("sort")
	p.Order = c.Query("order")
	list, total, err := h.Products.ListAdmin(ctx, adminProductFilter(c), p)
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, paging.Page[models.Product]{Items: list, Total: total, Page: p.Page, PageSize: p.PageSize})
}

func (h *Handlers) adminGetProduct(c *gin.Context) {
	id, ok := ParseID(c, "id")
	if !ok {
		return
	}
	p, err := h.Products.GetAdmin(c.Request.Context(), id)
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, p)
}

func (h *Handlers) adminListCasesByProduct(c *gin.Context) {
	id, ok := ParseID(c, "id")
	if !ok {
		return
	}
	list, err := h.Products.ListCasesAdmin(c.Request.Context(), id)
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, list)
}

func (h *Handlers) adminListScriptsByProduct(c *gin.Context) {
	id, ok := ParseID(c, "id")
	if !ok {
		return
	}
	list, err := h.Products.ListScriptsAdmin(c.Request.Context(), id)
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, list)
}

func (h *Handlers) adminPatchProductDraft(c *gin.Context) {
	id, ok := ParseID(c, "id")
	if !ok {
		return
	}
	var body struct {
		IsDraft *bool `json:"is_draft"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.IsDraft == nil {
		Fail(c, apperr.Invalid("invalid_body", "请提供 is_draft 布尔字段"))
		return
	}
	p, err := h.Products.SetDraft(c.Request.Context(), id, *body.IsDraft)
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, p)
}

func (h *Handlers) createProduct(c *gin.Context) {
	var p models.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		Fail(c, apperr.Invalid("invalid_body", "参数错误"))
		return
	}
	created, err := h.Products.Create(c.Request.Context(), &p)
	if err != nil {
		Fail(c, err)
		return
	}
	Created(c, created)
}

func (h *Handlers) updateProduct(c *gin.Context) {
	id, ok := ParseID(c, "id")
	if !ok {
		return
	}
	var p models.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		Fail(c, apperr.Invalid("invalid_body", "参数错误"))
		return
	}
	updated, err := h.Products.Update(c.Request.Context(), id, &p)
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, updated)
}

func (h *Handlers) deleteProduct(c *gin.Context) {
	id, ok := ParseID(c, "id")
	if !ok {
		return
	}
	if err := h.Products.Delete(c.Request.Context(), id); err != nil {
		Fail(c, err)
		return
	}
	NoContent(c)
}

func (h *Handlers) createScript(c *gin.Context) {
	id, ok := ParseID(c, "id")
	if !ok {
		return
	}
	var sc models.SalesScript
	if err := c.ShouldBindJSON(&sc); err != nil {
		Fail(c, apperr.Invalid("invalid_body", "参数错误"))
		return
	}
	created, err := h.Products.CreateScript(c.Request.Context(), id, &sc)
	if err != nil {
		Fail(c, err)
		return
	}
	Created(c, created)
}

func (h *Handlers) updateScript(c *gin.Context) {
	id, ok := ParseID(c, "id")
	if !ok {
		return
	}
	var sc models.SalesScript
	if err := c.ShouldBindJSON(&sc); err != nil {
		Fail(c, apperr.Invalid("invalid_body", "参数错误"))
		return
	}
	updated, err := h.Products.UpdateScript(c.Request.Context(), id, &sc)
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, updated)
}

func (h *Handlers) deleteScript(c *gin.Context) {
	id, ok := ParseID(c, "id")
	if !ok {
		return
	}
	if err := h.Products.DeleteScript(c.Request.Context(), id); err != nil {
		Fail(c, err)
		return
	}
	NoContent(c)
}

func (h *Handlers) createCase(c *gin.Context) {
	id, ok := ParseID(c, "id")
	if !ok {
		return
	}
	var ca models.Case
	if err := c.ShouldBindJSON(&ca); err != nil {
		Fail(c, apperr.Invalid("invalid_body", "参数错误"))
		return
	}
	created, err := h.Products.CreateCase(c.Request.Context(), id, &ca)
	if err != nil {
		Fail(c, err)
		return
	}
	Created(c, created)
}

func (h *Handlers) updateCase(c *gin.Context) {
	id, ok := ParseID(c, "id")
	if !ok {
		return
	}
	var ca models.Case
	if err := c.ShouldBindJSON(&ca); err != nil {
		Fail(c, apperr.Invalid("invalid_body", "参数错误"))
		return
	}
	updated, err := h.Products.UpdateCase(c.Request.Context(), id, &ca)
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, updated)
}

func (h *Handlers) deleteCase(c *gin.Context) {
	id, ok := ParseID(c, "id")
	if !ok {
		return
	}
	if err := h.Products.DeleteCase(c.Request.Context(), id); err != nil {
		Fail(c, err)
		return
	}
	NoContent(c)
}
