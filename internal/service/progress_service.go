package service

import (
	"context"
	"math"

	"github.com/sriram15/progressor-todo-app/internal/connection"
	"github.com/sriram15/progressor-todo-app/internal/database"
)

type StatCardData struct {
	Value     int `json:"value"`
	PrevValue int `json:"prevValue"`
}

type GetStatsResult struct {
	WeekHrs       StatCardData `json:"weekHrs"`
	MonthHrs      StatCardData `json:"monthHrs"`
	WeekProgress  StatCardData `json:"weekProgress"`
	MonthProgress StatCardData `json:"monthProgress"`
}

type IProgressService interface {
	GetStats() (GetStatsResult, error)
	GetDailyTotalMinutes() ([]database.GetDailyTotalMinutesRow, error)
}

type ProgressService struct {
	ctx                   context.Context
	taskCompletionService ITaskCompletionService
	queries               *database.Queries
}

func NewProgressService(taskCompletionService ITaskCompletionService, queries *database.Queries) *ProgressService {
	return &ProgressService{
		ctx:                   context.Background(),
		taskCompletionService: taskCompletionService,
		queries:               queries,
	}
}

func (p *ProgressService) GetStats() (GetStatsResult, error) {
	db, unlock := connection.GetDB()
	defer unlock()

	weekMins, err := p.queries.AggregateWeekHours(p.ctx, db, int64(1))
	if err != nil {
		return GetStatsResult{}, err
	}

	prevWeekMins, err := p.queries.AggregatePreviousWeekHours(p.ctx, db, int64(1))
	if err != nil {
		return GetStatsResult{}, err
	}

	monthMins, err := p.queries.AggregateMonthHours(p.ctx, db, int64(1))
	if err != nil {
		return GetStatsResult{}, err
	}

	prevMonthMins, err := p.queries.AggregatePreviousMonthHours(p.ctx, db, int64(1))
	if err != nil {
		return GetStatsResult{}, err
	}

	weekProgressDays, err := p.queries.GetWeeklyProgress(p.ctx, db)
	if err != nil {
		return GetStatsResult{}, err
	}

	previousWeekProgressDays, err := p.queries.GetPreviousWeeklyProgress(p.ctx, db)
	if err != nil {
		return GetStatsResult{}, err
	}

	monthProgressDays, err := p.queries.GetMonthlyProgress(p.ctx, db)
	if err != nil {
		return GetStatsResult{}, err
	}

	previousMonthProgressDays, err := p.queries.GetPreviousMonthlyProgress(p.ctx, db)
	if err != nil {
		return GetStatsResult{}, err
	}

	// Convert to hours from mins
	weekHours := math.Ceil(weekMins / 60.0)
	monthHours := math.Ceil(monthMins / 60.0)
	prevWeekHours := math.Ceil(prevWeekMins / 60.0)
	prevMonthHours := math.Ceil(prevMonthMins / 60.0)

	return GetStatsResult{
		WeekHrs:       StatCardData{Value: int(weekHours), PrevValue: int(prevWeekHours)},
		MonthHrs:      StatCardData{Value: int(monthHours), PrevValue: int(prevMonthHours)},
		WeekProgress:  StatCardData{Value: int(weekProgressDays), PrevValue: int(previousWeekProgressDays)},
		MonthProgress: StatCardData{Value: int(monthProgressDays), PrevValue: int(previousMonthProgressDays)},
	}, nil
}

func (p *ProgressService) GetDailyTotalMinutes() ([]database.GetDailyTotalMinutesRow, error) {
	db, unlock := connection.GetDB()
	defer unlock()
	return p.queries.GetDailyTotalMinutes(p.ctx, db)
}

func (p *ProgressService) GetTotalExpForUser() (float64, error) {
	return p.taskCompletionService.TotalUserExp(userId)
}
