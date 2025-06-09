package service

import (
	"context"

	"github.com/sriram15/progressor-todo-app/internal/connection"
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
	queries *database.Queries
}

func NewTaskCompletionService(queries *database.Queries) *TaskCompletionService {
	return &TaskCompletionService{
		ctx:     context.Background(),
		queries: queries,
	}
}

// CreateTaskCompletion creates a new TaskCompletion record and returns it
func (t *TaskCompletionService) CreateTaskCompletion(cardId int64, userId int64, baseExp int64, timeBonusExp int64, streakBonusExp int64) (database.TaskCompletion, error) {
	totalExp := baseExp + timeBonusExp + streakBonusExp

	db, unlock := connection.GetDB()
	defer unlock()

	taskValue, err := t.queries.CreateTaskCompletion(t.ctx, db, database.CreateTaskCompletionParams{
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
	db, unlock := connection.GetDB()
	defer unlock()

	taskCompletion, err := t.queries.GetTaskCompletion(t.ctx, db, database.GetTaskCompletionParams{
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
	db, unlock := connection.GetDB()
	defer unlock()

	taskCompletions, err := t.queries.ListTaskCompletionsByUser(t.ctx, db, userId)
	if err != nil {
		return nil, err
	}

	return taskCompletions, nil
}

// TotalUserExp calculates total user exp
func (t *TaskCompletionService) TotalUserExp(userId int64) (float64, error) {
	db, unlock := connection.GetDB()
	defer unlock()

	totalExp, err := t.queries.TotalUserExp(t.ctx, db, userId)
	if err != nil {
		return 0, err
	}
	return totalExp, nil
}
