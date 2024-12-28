package service

import (
	"context"

	"github.com/sriram15/progressor-todo-app/internal/database"
)

type GetStatsResult struct {
	WeekHrs  float64 `json:"weekHrs"`
	MonthHrs float64 `json:"monthHrs"`
	YearHrs  float64 `json:"yearHrs"`
}

type ProgressService interface {
	GetStats() (GetStatsResult, error)
}

type progressService struct {
	ctx     context.Context
	queries *database.Queries
}

func NewProgressService(queries *database.Queries) ProgressService {
	return &progressService{
		ctx:     context.Background(),
		queries: queries,
	}

}

func (p *progressService) GetStats() (GetStatsResult, error) {

	weekHrs, err := p.queries.AggregateWeekHours(p.ctx, int64(1))
	if err != nil || !weekHrs.Valid {
		return GetStatsResult{}, err
	}
	monthHrs, err := p.queries.AggregateMonthHours(p.ctx, int64(1))
	if err != nil || !monthHrs.Valid {
		return GetStatsResult{}, err
	}

	yearHrs, err := p.queries.AggregateYearHours(p.ctx, int64(1))
	if err != nil || !yearHrs.Valid {
		return GetStatsResult{}, err
	}

	return GetStatsResult{
		WeekHrs:  weekHrs.Float64,
		MonthHrs: monthHrs.Float64,
		YearHrs:  yearHrs.Float64,
	}, nil
}
