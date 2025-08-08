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
	ctx       context.Context
	dbManager *connection.DBManager
}

func NewTaskCompletionService(dbManager *connection.DBManager) *TaskCompletionService {
	return &TaskCompletionService{
		ctx:       context.Background(),
		dbManager: dbManager,
	}
}

// CreateTaskCompletion creates a new TaskCompletion record and returns it
func (t *TaskCompletionService) CreateTaskCompletion(cardId int64, userId int64, baseExp int64, timeBonusExp int64, streakBonusExp int64) (database.TaskCompletion, error) {
	totalExp := baseExp + timeBonusExp + streakBonusExp
	var taskValue database.TaskCompletion

	err := t.dbManager.Execute(t.ctx, func(q *database.Queries) error {
		var err error
		taskValue, err = q.CreateTaskCompletion(t.ctx, database.CreateTaskCompletionParams{
			Cardid:         cardId,
			Userid:         userId,
			Baseexp:        baseExp,
			Timebonusexp:   timeBonusExp,
			Streakbonusexp: streakBonusExp,
			Totalexp:       totalExp,
		})
		return err
	})

	if err != nil {
		return database.TaskCompletion{}, err
	}

	return taskValue, nil
}

// GetTaskCompletion retrieves a TaskCompletion record using cardId and userId
func (t *TaskCompletionService) GetTaskCompletion(cardId int64, userId int64) (database.TaskCompletion, error) {
	queries := t.dbManager.Queries(t.ctx)
	return queries.GetTaskCompletion(t.ctx, database.GetTaskCompletionParams{
		Cardid: cardId,
		Userid: userId,
	})
}

// ListTaskCompletionsByUser lists all task completions for a user
func (t *TaskCompletionService) ListTaskCompletionsByUser(userId int64) ([]database.TaskCompletion, error) {
	queries := t.dbManager.Queries(t.ctx)
	return queries.ListTaskCompletionsByUser(t.ctx, userId)
}

// TotalUserExp calculates total user exp
func (t *TaskCompletionService) TotalUserExp(userId int64) (float64, error) {
	queries := t.dbManager.Queries(t.ctx)
	return queries.TotalUserExp(t.ctx, userId)
}
