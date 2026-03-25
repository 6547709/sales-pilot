package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/6547709/sales-pilot/backend/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// --- 公开 DTO ---

type vendorMini struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Order int    `json:"sort_order"`
}

type vendorsSplit struct {
	Domestic []vendorMini `json:"domestic"`
	Foreign  []vendorMini `json:"foreign"`
}

type categoryNode struct {
	ID       uint         `json:"id"`
	Slug     string       `json:"slug"`
	Label    string       `json:"label"`
	IconKey  string       `json:"icon_key"`
	Keywords []string     `json:"keywords"`
	Hint     string       `json:"hint"`
	Vendors  vendorsSplit `json:"vendors"`
}

type layerBlock struct {
	Layer      models.TopologyLayer `json:"layer"`
	Categories []categoryNode       `json:"categories"`
}

func keywordsSlice(j []byte) []string {
	var k []string
	if len(j) == 0 {
		return k
	}
	_ = json.Unmarshal(j, &k)
	return k
}

func buildCategoryNode(c models.SolutionCategory) categoryNode {
	dom := make([]vendorMini, 0)
	frn := make([]vendorMini, 0)
	for _, v := range c.Vendors {
		m := vendorMini{ID: v.ID, Name: v.Name, Order: v.SortOrder}
		if v.Market == "foreign" {
			frn = append(frn, m)
		} else {
			dom = append(dom, m)
		}
	}
	return categoryNode{
		ID:       c.ID,
		Slug:     c.Slug,
		Label:    c.Label,
		IconKey:  c.IconKey,
		Keywords: keywordsSlice(c.Keywords),
		Hint:     c.Hint,
		Vendors:  vendorsSplit{Domestic: dom, Foreign: frn},
	}
}

// getTopology 首页拓扑（仅启用项）
func (s *Server) getTopology(c *gin.Context) {
	var layers []models.TopologyLayer
	if err := s.DB.Order("sort_order ASC, level DESC").Find(&layers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var sec []models.SolutionCategory
	if err := s.DB.Where("column_type = ? AND is_active = ?", "security", true).
		Preload("Vendors", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order ASC, id ASC") }).
		Order("sort_order ASC, id ASC").Find(&sec).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	secNodes := make([]categoryNode, 0, len(sec))
	for _, x := range sec {
		secNodes = append(secNodes, buildCategoryNode(x))
	}

	var ops []models.SolutionCategory
	if err := s.DB.Where("column_type = ? AND is_active = ?", "ops", true).
		Preload("Vendors", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order ASC, id ASC") }).
		Order("sort_order ASC, id ASC").Find(&ops).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	opsNodes := make([]categoryNode, 0, len(ops))
	for _, x := range ops {
		opsNodes = append(opsNodes, buildCategoryNode(x))
	}

	central := make([]layerBlock, 0, len(layers))
	for _, L := range layers {
		var cats []models.SolutionCategory
		if err := s.DB.Where("column_type = ? AND layer_id = ? AND is_active = ?", "central", L.ID, true).
			Preload("Vendors", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order ASC, id ASC") }).
			Order("sort_order ASC, id ASC").Find(&cats).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		nodes := make([]categoryNode, 0, len(cats))
		for _, x := range cats {
			nodes = append(nodes, buildCategoryNode(x))
		}
		central = append(central, layerBlock{Layer: L, Categories: nodes})
	}

	c.JSON(http.StatusOK, gin.H{
		"security":        secNodes,
		"ops":             opsNodes,
		"central_layers":  central,
	})
}

// listSolutionCategories 扁平列表（方案下拉，含停用项供管理端也可复用带 query）
func (s *Server) listSolutionCategories(c *gin.Context) {
	q := s.DB.Model(&models.SolutionCategory{}).Order("column_type ASC, sort_order ASC, id ASC")
	if c.Query("active_only") == "1" || c.Query("active_only") == "true" {
		q = q.Where("is_active = ?", true)
	}
	var list []models.SolutionCategory
	if err := q.Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// --- 管理：分层 ---

func (s *Server) adminListLayers(c *gin.Context) {
	var layers []models.TopologyLayer
	if err := s.DB.Order("sort_order ASC, level DESC").Find(&layers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, layers)
}

func (s *Server) adminCreateLayer(c *gin.Context) {
	var body models.TopologyLayer
	if err := c.ShouldBindJSON(&body); err != nil || body.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if err := s.DB.Create(&body).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, body)
}

func (s *Server) adminUpdateLayer(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var L models.TopologyLayer
	if err := s.DB.First(&L, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到"})
		return
	}
	if err := c.ShouldBindJSON(&L); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	L.ID = uint(id)
	if err := s.DB.Model(&models.TopologyLayer{}).Where("id = ?", id).Updates(map[string]any{
		"level": L.Level, "title": L.Title, "subtitle": L.Subtitle, "sort_order": L.SortOrder,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_ = s.DB.First(&L, id).Error
	c.JSON(http.StatusOK, L)
}

func (s *Server) adminDeleteLayer(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var cnt int64
	s.DB.Model(&models.SolutionCategory{}).Where("layer_id = ?", id).Count(&cnt)
	if cnt > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该层下仍有解决方案，请先移动或删除分类"})
		return
	}
	if err := s.DB.Delete(&models.TopologyLayer{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// --- 管理：解决方案分类 ---

func (s *Server) adminCreateCategory(c *gin.Context) {
	var body models.SolutionCategory
	if err := c.ShouldBindJSON(&body); err != nil || body.Slug == "" || body.Label == "" || body.ColumnType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "slug、label、column_type 必填"})
		return
	}
	if err := s.DB.Create(&body).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, body)
}

func (s *Server) adminUpdateCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var cat models.SolutionCategory
	if err := s.DB.First(&cat, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到"})
		return
	}
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cat.ID = uint(id)
	up := map[string]any{
		"slug": cat.Slug, "label": cat.Label, "icon_key": cat.IconKey,
		"column_type": cat.ColumnType, "layer_id": cat.LayerID, "sort_order": cat.SortOrder,
		"keywords": cat.Keywords, "hint": cat.Hint, "is_active": cat.IsActive,
	}
	if err := s.DB.Model(&models.SolutionCategory{}).Where("id = ?", id).Updates(up).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_ = s.DB.Preload("Vendors").First(&cat, id).Error
	c.JSON(http.StatusOK, cat)
}

func (s *Server) adminDeleteCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var cnt int64
	s.DB.Model(&models.Product{}).Where("solution_category_id = ?", id).Count(&cnt)
	if cnt > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "仍有产品关联此解决方案，无法删除"})
		return
	}
	_ = s.DB.Where("category_id = ?", id).Delete(&models.SolutionVendor{}).Error
	if err := s.DB.Delete(&models.SolutionCategory{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// --- 管理：厂商 ---

type vendorBody struct {
	CategoryID uint   `json:"category_id" binding:"required"`
	Market     string `json:"market" binding:"required"` // domestic | foreign
	Name       string `json:"name" binding:"required"`
	SortOrder  int    `json:"sort_order"`
}

func (s *Server) adminCreateVendor(c *gin.Context) {
	var body vendorBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	v := models.SolutionVendor{
		CategoryID: body.CategoryID, Market: body.Market, Name: body.Name, SortOrder: body.SortOrder,
	}
	if err := s.DB.Create(&v).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, v)
}

func (s *Server) adminUpdateVendor(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var v models.SolutionVendor
	if err := s.DB.First(&v, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到"})
		return
	}
	if err := c.ShouldBindJSON(&v); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	v.ID = uint(id)
	if err := s.DB.Model(&v).Updates(map[string]any{
		"category_id": v.CategoryID, "market": v.Market, "name": v.Name, "sort_order": v.SortOrder,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_ = s.DB.First(&v, id).Error
	c.JSON(http.StatusOK, v)
}

func (s *Server) adminDeleteVendor(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	_ = s.DB.Delete(&models.SolutionVendor{}, id).Error
	c.Status(http.StatusNoContent)
}
