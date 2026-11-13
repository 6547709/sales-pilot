package service

import (
	"context"
	"errors"
	"strings"

	"github.com/6547709/sales-pilot/backend/internal/apperr"
	"github.com/6547709/sales-pilot/backend/internal/models"
	"github.com/6547709/sales-pilot/backend/internal/view"
	"gorm.io/gorm"
)

type ScenarioSummary struct {
	Scenario     models.Scenario `json:"scenario"`
	ProductCount int64           `json:"product_count"`
	ServiceCount int64           `json:"service_count"`
}

type ScenarioGroupView struct {
	Group     models.ScenarioGroup `json:"group"`
	Scenarios []ScenarioSummary    `json:"scenarios"`
}

type ScenarioDetail struct {
	Scenario models.Scenario       `json:"scenario"`
	Group    *models.ScenarioGroup `json:"group,omitempty"`
	Products []models.Product      `json:"products"`
	Services []models.Service      `json:"services"`
	Cases    []models.Case         `json:"cases"`
}

type ScenarioService struct {
	repo ScenarioRepo
}

func NewScenarioService(repo ScenarioRepo) *ScenarioService {
	return &ScenarioService{repo: repo}
}

const MaxScenarioGroups = 6

func (s *ScenarioService) Groups(ctx context.Context) ([]ScenarioGroupView, error) {
	groups, err := s.repo.Groups(ctx)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	scenarios, err := s.repo.Scenarios(ctx)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	pc, err := s.repo.CountProducts(ctx)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	sc, err := s.repo.CountServices(ctx)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	byGroup := make(map[uint][]ScenarioSummary, len(groups))
	for _, item := range scenarios {
		byGroup[item.GroupID] = append(byGroup[item.GroupID], ScenarioSummary{
			Scenario: item, ProductCount: pc[item.ID], ServiceCount: sc[item.ID],
		})
	}
	out := make([]ScenarioGroupView, 0, len(groups))
	for _, g := range groups {
		out = append(out, ScenarioGroupView{Group: g, Scenarios: byGroup[g.ID]})
	}
	return out, nil
}

func (s *ScenarioService) GetBySlug(ctx context.Context, slug string, mode view.Mode) (*ScenarioDetail, error) {
	scenario, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound("scenario_not_found", "未找到")
		}
		return nil, apperr.Internal(err)
	}
	if !scenario.IsActive {
		return nil, apperr.NotFound("scenario_not_found", "未找到")
	}
	group, _ := s.repo.FindGroup(ctx, scenario.GroupID)
	productIDs, err := s.repo.ProductIDs(ctx, scenario.ID)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	serviceIDs, err := s.repo.ServiceIDs(ctx, scenario.ID)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	products, err := s.repo.ListProductsByIDs(ctx, productIDs)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	services, err := s.repo.ListServicesByIDs(ctx, serviceIDs)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	var cases []models.Case
	if mode.IsCustomer() {
		cases, err = s.repo.ListPublicCasesByProductIDs(ctx, productIDs)
	} else {
		cases, err = s.repo.ListCasesByProductIDs(ctx, productIDs)
	}
	if err != nil {
		return nil, apperr.Internal(err)
	}
	return &ScenarioDetail{
		Scenario: *scenario, Group: group,
		Products: products, Services: services, Cases: cases,
	}, nil
}

// GetAdminDetail 后台关联编辑：不过滤停用场景、草稿产品与停用服务。
func (s *ScenarioService) GetAdminDetail(ctx context.Context, id uint) (*ScenarioDetail, error) {
	scenario, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound("scenario_not_found", "未找到")
		}
		return nil, apperr.Internal(err)
	}
	group, _ := s.repo.FindGroup(ctx, scenario.GroupID)
	productIDs, err := s.repo.ProductIDs(ctx, id)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	serviceIDs, err := s.repo.ServiceIDs(ctx, id)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	products, err := s.repo.ListProductsByIDsAll(ctx, productIDs)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	services, err := s.repo.ListServicesByIDsAll(ctx, serviceIDs)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	return &ScenarioDetail{Scenario: *scenario, Group: group, Products: products, Services: services}, nil
}

func (s *ScenarioService) CreateGroup(ctx context.Context, name, slug, description string) (*models.ScenarioGroup, error) {
	if strings.TrimSpace(name) == "" || strings.TrimSpace(slug) == "" {
		return nil, apperr.Invalid("invalid_body", "name 与 slug 必填")
	}
	n, err := s.repo.CountGroups(ctx)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	if n >= MaxScenarioGroups {
		return nil, apperr.Conflict("scenario_group_limit", "场景组最多 6 个")
	}
	g := &models.ScenarioGroup{Name: name, Slug: slug, Description: description}
	if err := s.repo.CreateGroup(ctx, g); err != nil {
		return nil, apperr.Conflict("scenario_group_slug_taken", "创建失败，slug 可能已存在")
	}
	return g, nil
}

func (s *ScenarioService) UpdateGroup(ctx context.Context, id uint, g *models.ScenarioGroup) (*models.ScenarioGroup, error) {
	if _, err := s.repo.FindGroup(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound("scenario_group_not_found", "未找到")
		}
		return nil, apperr.Internal(err)
	}
	g.ID = id
	if err := s.repo.UpdateGroupFields(ctx, id, map[string]any{
		"name": g.Name, "slug": g.Slug, "description": g.Description,
	}); err != nil {
		return nil, apperr.Internal(err)
	}
	return g, nil
}

func (s *ScenarioService) DeleteGroup(ctx context.Context, id uint) error {
	n, err := s.repo.CountScenariosInGroup(ctx, id)
	if err != nil {
		return apperr.Internal(err)
	}
	if n > 0 {
		return apperr.Conflict("scenario_group_not_empty", "该组下仍有场景，请先移动或删除")
	}
	if err := s.repo.DeleteGroup(ctx, id); err != nil {
		return apperr.Internal(err)
	}
	return nil
}

func (s *ScenarioService) Create(ctx context.Context, sc *models.Scenario) (*models.Scenario, error) {
	if strings.TrimSpace(sc.Name) == "" || strings.TrimSpace(sc.Slug) == "" {
		return nil, apperr.Invalid("invalid_body", "name 与 slug 必填")
	}
	if _, err := s.repo.FindGroup(ctx, sc.GroupID); err != nil {
		return nil, apperr.Invalid("invalid_scenario_group", "场景组不存在")
	}
	if err := s.repo.Create(ctx, sc); err != nil {
		return nil, apperr.Conflict("scenario_slug_taken", "创建失败，slug 可能已存在")
	}
	return sc, nil
}

func (s *ScenarioService) Update(ctx context.Context, id uint, sc *models.Scenario) (*models.Scenario, error) {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound("scenario_not_found", "未找到")
		}
		return nil, apperr.Internal(err)
	}
	if _, err := s.repo.FindGroup(ctx, sc.GroupID); err != nil {
		return nil, apperr.Invalid("invalid_scenario_group", "场景组不存在")
	}
	sc.ID = id
	if err := s.repo.UpdateFields(ctx, id, map[string]any{
		"group_id": sc.GroupID, "slug": sc.Slug, "name": sc.Name, "summary": sc.Summary,
		"pain_point": sc.PainPoint, "value_prop": sc.ValueProp, "icon_key": sc.IconKey,
		"sort_order": sc.SortOrder, "is_active": sc.IsActive,
	}); err != nil {
		return nil, apperr.Internal(err)
	}
	return sc, nil
}

func (s *ScenarioService) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return apperr.Internal(err)
	}
	return nil
}

func (s *ScenarioService) SetLinks(ctx context.Context, id uint, productIDs, serviceIDs []uint) error {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperr.NotFound("scenario_not_found", "未找到")
		}
		return apperr.Internal(err)
	}
	if err := s.repo.SetLinks(ctx, id, productIDs, serviceIDs); err != nil {
		return apperr.Internal(err)
	}
	return nil
}
