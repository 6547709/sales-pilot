package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/6547709/sales-pilot/backend/internal/apperr"
	"github.com/6547709/sales-pilot/backend/internal/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type VendorMini struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Order int    `json:"sort_order"`
}

type VendorsSplit struct {
	Domestic []VendorMini `json:"domestic"`
	Foreign  []VendorMini `json:"foreign"`
}

type ServiceMini struct {
	ID      uint   `json:"id"`
	Slug    string `json:"slug"`
	Name    string `json:"name"`
	Summary string `json:"summary"`
}

type CategoryNode struct {
	ID       uint          `json:"id"`
	Slug     string        `json:"slug"`
	Label    string        `json:"label"`
	IconKey  string        `json:"icon_key"`
	Keywords []string      `json:"keywords"`
	Hint     string        `json:"hint"`
	Vendors  VendorsSplit  `json:"vendors"`
	Services []ServiceMini `json:"services"`
}

type LayerBlock struct {
	Layer      models.TopologyLayer `json:"layer"`
	Categories []CategoryNode       `json:"categories"`
}

type TreeResponse struct {
	Security      []CategoryNode `json:"security"`
	Ops           []CategoryNode `json:"ops"`
	CentralLayers []LayerBlock   `json:"central_layers"`
}

type TopologyService struct {
	repo     TopologyRepo
	services ServiceRepo
}

func NewTopologyService(repo TopologyRepo, services ServiceRepo) *TopologyService {
	return &TopologyService{repo: repo, services: services}
}

func keywordsSlice(j []byte) []string {
	var k []string
	if len(j) == 0 {
		return k
	}
	_ = json.Unmarshal(j, &k)
	return k
}

func buildNode(c models.SolutionCategory) CategoryNode {
	dom := make([]VendorMini, 0)
	frn := make([]VendorMini, 0)
	for _, v := range c.Vendors {
		m := VendorMini{ID: v.ID, Name: v.Name, Order: v.SortOrder}
		if v.Market == "foreign" {
			frn = append(frn, m)
		} else {
			dom = append(dom, m)
		}
	}
	return CategoryNode{
		ID: c.ID, Slug: c.Slug, Label: c.Label, IconKey: c.IconKey,
		Keywords: keywordsSlice(c.Keywords), Hint: c.Hint,
		Vendors: VendorsSplit{Domestic: dom, Foreign: frn},
	}
}

func buildNodes(cats []models.SolutionCategory) []CategoryNode {
	nodes := make([]CategoryNode, 0, len(cats))
	for _, c := range cats {
		nodes = append(nodes, buildNode(c))
	}
	return nodes
}

func (s *TopologyService) Tree(ctx context.Context) (*TreeResponse, error) {
	snap, err := s.repo.Snapshot(ctx)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	var all []models.SolutionCategory
	all = append(all, snap.Security...)
	all = append(all, snap.Ops...)
	for _, cats := range snap.Central {
		all = append(all, cats...)
	}
	ids := make([]uint, 0, len(all))
	for _, c := range all {
		ids = append(ids, c.ID)
	}
	svcMap := map[uint][]models.Service{}
	if s.services != nil && len(ids) > 0 {
		if svcMap, err = s.services.ListByCategoryIDs(ctx, ids); err != nil {
			return nil, apperr.Internal(err)
		}
	}
	withServices := func(cats []models.SolutionCategory) []CategoryNode {
		nodes := buildNodes(cats)
		for i, c := range cats {
			for _, svc := range svcMap[c.ID] {
				nodes[i].Services = append(nodes[i].Services, ServiceMini{
					ID: svc.ID, Slug: svc.Slug, Name: svc.Name, Summary: svc.Summary,
				})
			}
		}
		return nodes
	}
	central := make([]LayerBlock, 0, len(snap.Layers))
	for _, l := range snap.Layers {
		central = append(central, LayerBlock{Layer: l, Categories: withServices(snap.Central[l.ID])})
	}
	return &TreeResponse{
		Security:      withServices(snap.Security),
		Ops:           withServices(snap.Ops),
		CentralLayers: central,
	}, nil
}

func (s *TopologyService) ListCategories(ctx context.Context, activeOnly bool) ([]models.SolutionCategory, error) {
	list, err := s.repo.Categories(ctx, activeOnly)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	return list, nil
}

func (s *TopologyService) ListLayers(ctx context.Context) ([]models.TopologyLayer, error) {
	list, err := s.repo.Layers(ctx)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	return list, nil
}

func (s *TopologyService) CreateLayer(ctx context.Context, l *models.TopologyLayer) (*models.TopologyLayer, error) {
	if l.Title == "" {
		return nil, apperr.Invalid("invalid_body", "参数错误")
	}
	if err := s.repo.CreateLayer(ctx, l); err != nil {
		return nil, apperr.Internal(err)
	}
	return l, nil
}

func (s *TopologyService) UpdateLayer(ctx context.Context, id uint, l *models.TopologyLayer) (*models.TopologyLayer, error) {
	if _, err := s.repo.FindLayer(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound("layer_not_found", "未找到")
		}
		return nil, apperr.Internal(err)
	}
	l.ID = id
	fields := map[string]any{
		"level": l.Level, "title": l.Title, "subtitle": l.Subtitle, "sort_order": l.SortOrder,
	}
	if err := s.repo.UpdateLayerFields(ctx, id, fields); err != nil {
		return nil, apperr.Internal(err)
	}
	return l, nil
}

func (s *TopologyService) DeleteLayer(ctx context.Context, id uint) error {
	n, err := s.repo.CountCategoriesInLayer(ctx, id)
	if err != nil {
		return apperr.Internal(err)
	}
	if n > 0 {
		return apperr.Conflict("layer_not_empty", "该层下仍有解决方案，请先移动或删除分类")
	}
	if err := s.repo.DeleteLayer(ctx, id); err != nil {
		return apperr.Internal(err)
	}
	return nil
}

func (s *TopologyService) CreateCategory(ctx context.Context, c *models.SolutionCategory) (*models.SolutionCategory, error) {
	if c.Slug == "" || c.Label == "" || c.ColumnType == "" {
		return nil, apperr.Invalid("invalid_body", "slug、label、column_type 必填")
	}
	if err := s.repo.CreateCategory(ctx, c); err != nil {
		return nil, apperr.Internal(err)
	}
	return c, nil
}

func (s *TopologyService) UpdateCategory(ctx context.Context, id uint, c *models.SolutionCategory) (*models.SolutionCategory, error) {
	if _, err := s.repo.FindCategory(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound("category_not_found", "未找到")
		}
		return nil, apperr.Internal(err)
	}
	c.ID = id
	fields := map[string]any{
		"slug": c.Slug, "label": c.Label, "icon_key": c.IconKey,
		"column_type": c.ColumnType, "layer_id": c.LayerID, "sort_order": c.SortOrder,
		"keywords": c.Keywords, "hint": c.Hint, "is_active": c.IsActive,
		"customer_value": c.CustomerValue, "customer_features": c.CustomerFeatures,
	}
	if err := s.repo.UpdateCategoryFields(ctx, id, fields); err != nil {
		return nil, apperr.Internal(err)
	}
	return s.repo.FindCategory(ctx, id)
}

func (s *TopologyService) UpdateCategoryCustomer(ctx context.Context, id uint, value string, features datatypes.JSON) (*models.SolutionCategory, error) {
	if _, err := s.repo.FindCategory(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound("category_not_found", "未找到")
		}
		return nil, apperr.Internal(err)
	}
	if err := s.repo.UpdateCategoryFields(ctx, id, map[string]any{
		"customer_value": value, "customer_features": features,
	}); err != nil {
		return nil, apperr.Internal(err)
	}
	return s.repo.FindCategory(ctx, id)
}

func (s *TopologyService) DeleteCategory(ctx context.Context, id uint) error {
	n, err := s.repo.CountProductsInCategory(ctx, id)
	if err != nil {
		return apperr.Internal(err)
	}
	if n > 0 {
		return apperr.Conflict("category_not_empty", "仍有产品关联此解决方案，无法删除")
	}
	if err := s.repo.DeleteCategory(ctx, id); err != nil {
		return apperr.Internal(err)
	}
	return nil
}

func (s *TopologyService) CreateVendor(ctx context.Context, v *models.SolutionVendor) (*models.SolutionVendor, error) {
	if v.CategoryID == 0 || v.Market == "" || v.Name == "" {
		return nil, apperr.Invalid("invalid_body", "参数错误")
	}
	if err := s.repo.CreateVendor(ctx, v); err != nil {
		return nil, apperr.Internal(err)
	}
	return v, nil
}

func (s *TopologyService) UpdateVendor(ctx context.Context, id uint, v *models.SolutionVendor) (*models.SolutionVendor, error) {
	if _, err := s.repo.FindVendor(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound("vendor_not_found", "未找到")
		}
		return nil, apperr.Internal(err)
	}
	v.ID = id
	fields := map[string]any{
		"category_id": v.CategoryID, "market": v.Market, "name": v.Name, "sort_order": v.SortOrder,
	}
	if err := s.repo.UpdateVendorFields(ctx, id, fields); err != nil {
		return nil, apperr.Internal(err)
	}
	return v, nil
}

func (s *TopologyService) DeleteVendor(ctx context.Context, id uint) error {
	if err := s.repo.DeleteVendor(ctx, id); err != nil {
		return apperr.Internal(err)
	}
	return nil
}
