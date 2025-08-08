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
	GetTotalExpForUser(userID int64) (float64, error)
}

type ProgressService struct {
	ctx       context.Context
	dbManager *connection.DBManager
}

func NewProgressService(dbManager *connection.DBManager) *ProgressService {
	return &ProgressService{
		ctx:       context.Background(),
		dbManager: dbManager,
	}
}

func (p *ProgressService) GetStats() (GetStatsResult, error) {
	var result GetStatsResult
	var err error

	readQueries := p.dbManager.Queries(p.ctx)

	weekMins, err := readQueries.AggregateWeekHours(p.ctx, int64(1))
	if err != nil {
		return result, err
	}

	prevWeekMins, err := readQueries.AggregatePreviousWeekHours(p.ctx, int64(1))
	if err != nil {
		return result, err
	}

	monthMins, err := readQueries.AggregateMonthHours(p.ctx, int64(1))
	if err != nil {
		return result, err
	}

	prevMonthMins, err := readQueries.AggregatePreviousMonthHours(p.ctx, int64(1))
	if err != nil {
		return result, err
	}

	weekProgressDays, err := readQueries.GetWeeklyProgress(p.ctx)
	if err != nil {
		return result, err
	}

	previousWeekProgressDays, err := readQueries.GetPreviousWeeklyProgress(p.ctx)
	if err != nil {
		return result, err
	}

	monthProgressDays, err := readQueries.GetMonthlyProgress(p.ctx)
	if err != nil {
		return result, err
	}

	previousMonthProgressDays, err := readQueries.GetPreviousMonthlyProgress(p.ctx)
	if err != nil {
		return result, err
	}

	// Convert to hours from mins
	weekHours := math.Ceil(weekMins / 60.0)
	monthHours := math.Ceil(monthMins / 60.0)
	prevWeekHours := math.Ceil(prevWeekMins / 60.0)
	prevMonthHours := math.Ceil(prevMonthMins / 60.0)

	result = GetStatsResult{
		WeekHrs:       StatCardData{Value: int(weekHours), PrevValue: int(prevWeekHours)},
		MonthHrs:      StatCardData{Value: int(monthHours), PrevValue: int(prevMonthHours)},
		WeekProgress:  StatCardData{Value: int(weekProgressDays), PrevValue: int(previousWeekProgressDays)},
		MonthProgress: StatCardData{Value: int(monthProgressDays), PrevValue: int(previousMonthProgressDays)},
	}
	return result, nil
}

func (p *ProgressService) GetDailyTotalMinutes() ([]database.GetDailyTotalMinutesRow, error) {
	var result []database.GetDailyTotalMinutesRow
	var err error

	readQueries := p.dbManager.Queries(p.ctx)
	result, err = readQueries.GetDailyTotalMinutes(p.ctx)
	return result, err
}

func (p *ProgressService) GetTotalExpForUser(userID int64) (float64, error) {
	readQueries := p.dbManager.Queries(p.ctx)
	totalExp, err := readQueries.TotalUserExp(p.ctx, userID)
	if err != nil {
		return 0, err
	}
	return totalExp, nil
}
