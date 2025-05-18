package service

import (
	"context"
	"math"

	"github.com/sriram15/progressor-todo-app/internal/database"
)

type GetStatsResult struct {
	WeekHrs  float64 `json:"weekHrs"`
	MonthHrs float64 `json:"monthHrs"`
	YearHrs  float64 `json:"yearHrs"`
}

type IProgressService interface {
	GetStats() (GetStatsResult, error)
	GetDailyTotalMinutes() ([]database.GetDailyTotalMinutesRow, error)
}

type ProgressService struct {
	ctx                   context.Context
	queries               *database.Queries
	taskCompletionService ITaskCompletionService
}

func NewProgressService(queries *database.Queries, taskCompletionService ITaskCompletionService) *ProgressService {
	return &ProgressService{
		ctx:                   context.Background(),
		queries:               queries,
		taskCompletionService: taskCompletionService,
	}

}

func (p *ProgressService) GetStats() (GetStatsResult, error) {

	weekMins, err := p.queries.AggregateWeekHours(p.ctx, int64(1))
	if err != nil {
		return GetStatsResult{}, err
	}
	monthMins, err := p.queries.AggregateMonthHours(p.ctx, int64(1))
	if err != nil {
		return GetStatsResult{}, err
	}

	yearMins, err := p.queries.AggregateYearHours(p.ctx, int64(1))
	if err != nil {
		return GetStatsResult{}, err
	}

	// Convert to hours from mins

	weekHours := math.Ceil(weekMins / 60.0)
	monthHours := math.Ceil(monthMins / 60.0)
	yearHours := math.Ceil(yearMins / 60.0)

	return GetStatsResult{
		WeekHrs:  weekHours,
		MonthHrs: monthHours,
		YearHrs:  yearHours,
	}, nil
}

func (p *ProgressService) GetDailyTotalMinutes() ([]database.GetDailyTotalMinutesRow, error) {
	return p.queries.GetDailyTotalMinutes(p.ctx)
}

func (p *ProgressService) GetTotalExpForUser() (float64, error) {
	return p.taskCompletionService.TotalUserExp(userId)
}
