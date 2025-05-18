package service

import (
	"context"
	"database/sql"

	"github.com/sriram15/progressor-todo-app/internal/database"
)

type ITaskCompletionService interface {
	CreateTaskCompletion(cardId int64, userId int64, baseExp int64, timeBonusExp int64, streakBonusExp int64) (database.TaskCompletion, error)
	GetTaskCompletion(cardId int64, userId int64) (database.TaskCompletion, error)
	ListTaskCompletionsByUser(userId int64) ([]database.TaskCompletion, error)
	TotalUserExp(userId int64) (float64, error)
}

type TaskCompletionService struct {
	ctx     context.Context
	db      *sql.DB
	queries *database.Queries
}

func NewTaskCompletionService(db *sql.DB, queries *database.Queries) *TaskCompletionService {
	return &TaskCompletionService{
		ctx:     context.Background(),
		db:      db,
		queries: queries,
	}
}

// CreateTaskCompletion creates a new TaskCompletion record and returns it
func (t *TaskCompletionService) CreateTaskCompletion(cardId int64, userId int64, baseExp int64, timeBonusExp int64, streakBonusExp int64) (database.TaskCompletion, error) {
	totalExp := baseExp + timeBonusExp + streakBonusExp

	taskValue, err := t.queries.CreateTaskCompletion(t.ctx, database.CreateTaskCompletionParams{
		Cardid:         cardId,
		Userid:         userId,
		Baseexp:        baseExp,
		Timebonusexp:   timeBonusExp,
		Streakbonusexp: streakBonusExp,
		Totalexp:       totalExp,
	})

	if err != nil {
		return database.TaskCompletion{}, err
	}

	return taskValue, nil
}

// GetTaskCompletion retrieves a TaskCompletion record using cardId and userId
func (t *TaskCompletionService) GetTaskCompletion(cardId int64, userId int64) (database.TaskCompletion, error) {
	taskCompletion, err := t.queries.GetTaskCompletion(t.ctx, database.GetTaskCompletionParams{
		Cardid: cardId,
		Userid: userId,
	})
	if err != nil {
		return database.TaskCompletion{}, err
	}

	return taskCompletion, nil
}

// ListTaskCompletionsByUser lists all task completions for a user
func (t *TaskCompletionService) ListTaskCompletionsByUser(userId int64) ([]database.TaskCompletion, error) {
	taskCompletions, err := t.queries.ListTaskCompletionsByUser(t.ctx, userId)
	if err != nil {
		return nil, err
	}

	return taskCompletions, nil
}

// TotalUserExp calculates total user exp
func (t *TaskCompletionService) TotalUserExp(userId int64) (float64, error) {
	totalExp, err := t.queries.TotalUserExp(t.ctx, userId)
	if err != nil {
		return 0, err
	}
	return totalExp, nil
}
