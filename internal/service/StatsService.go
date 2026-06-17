package service

import (
	"context"
	"zzz/internal/dto"
	"zzz/internal/repository"
)

type StatsService struct {
	repo *repository.StatsRepo
}

func NewStatsService(repo *repository.StatsRepo) *StatsService {
	return &StatsService{repo: repo}
}

func (h *StatsService) GetStats(ctx context.Context) (*dto.StatsResponse, error) {
	stats, err := h.repo.GetStats(ctx)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}
