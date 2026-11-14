package service

import (
	"context"

	"github.com/6547709/sales-pilot/backend/internal/apperr"
	"github.com/6547709/sales-pilot/backend/internal/repository"
)

type StatsRepo interface {
	ProductsByCategory(ctx context.Context) ([]repository.LabelCount, error)
	ProductsByMarket(ctx context.Context) ([]repository.LabelCount, error)
	UsersSummary(ctx context.Context) (total, admins, active int64, err error)
	ScenariosByGroup(ctx context.Context) ([]repository.LabelCount, error)
	ServicesSummary(ctx context.Context) (total, active int64, err error)
}

type ProductStats struct {
	Total      int64                   `json:"total"`
	ByCategory []repository.LabelCount `json:"by_category"`
	ByMarket   []repository.LabelCount `json:"by_market"`
}

type UserStats struct {
	Total  int64 `json:"total"`
	Admins int64 `json:"admins"`
	Active int64 `json:"active"`
}

type ScenarioStats struct {
	Total  int64                   `json:"total"`
	Groups []repository.LabelCount `json:"groups"`
}

type ServiceStats struct {
	Total  int64 `json:"total"`
	Active int64 `json:"active"`
}

type StatsOverview struct {
	Products  ProductStats  `json:"products"`
	Users     UserStats     `json:"users"`
	Scenarios ScenarioStats `json:"scenarios"`
	Services  ServiceStats  `json:"services"`
}

type StatsService struct {
	repo StatsRepo
}

func NewStatsService(repo StatsRepo) *StatsService {
	return &StatsService{repo: repo}
}

func (s *StatsService) Overview(ctx context.Context) (*StatsOverview, error) {
	byCategory, err := s.repo.ProductsByCategory(ctx)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	byMarket, err := s.repo.ProductsByMarket(ctx)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	userTotal, admins, activeUsers, err := s.repo.UsersSummary(ctx)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	groups, err := s.repo.ScenariosByGroup(ctx)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	svcTotal, svcActive, err := s.repo.ServicesSummary(ctx)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	var scenarioTotal int64
	for _, g := range groups {
		scenarioTotal += g.Count
	}
	var productTotal int64
	for _, m := range byMarket {
		productTotal += m.Count
	}
	return &StatsOverview{
		Products:  ProductStats{Total: productTotal, ByCategory: byCategory, ByMarket: byMarket},
		Users:     UserStats{Total: userTotal, Admins: admins, Active: activeUsers},
		Scenarios: ScenarioStats{Total: scenarioTotal, Groups: groups},
		Services:  ServiceStats{Total: svcTotal, Active: svcActive},
	}, nil
}
