package repository

import (
	"context"
	"zzz/internal/dto"

	"gorm.io/gorm"
)

type StatsRepo struct {
	db *gorm.DB
}

func NewStatsRepo(db *gorm.DB) *StatsRepo {
	return &StatsRepo{db: db}
}

func (s *StatsRepo) GetStats(ctx context.Context) (dto.StatsResponse, error) {
	var stats dto.StatsResponse
	query := `
		SELECT 
			(SELECT COUNT(*) FROM products) AS products,
			(SELECT COUNT(*) FROM clients) AS clients,
			(SELECT COUNT(*) FROM brands) AS brands,
			(SELECT COUNT(*) FROM categories) AS categories;
	`

	if err := s.db.WithContext(ctx).Raw(query).Scan(&stats).Error; err != nil {
		return dto.StatsResponse{}, err
	}

	return stats, nil
}
