package service

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/sriram15/progressor-todo-app/internal/connection"
	"github.com/sriram15/progressor-todo-app/internal/database"
	"github.com/sriram15/progressor-todo-app/internal/events"
)

var (
	ErrorUnknown           = errors.New("unknown error")
	ErrNotFound            = errors.New("not found")
	ErrInvalidProject      = errors.New("invalid project")
	ErrInvalidStatus       = errors.New("invalid status")
	ErrInvalidUpdate       = errors.New("invalid update")
	ErrCardTitleRequired   = errors.New("card title is required")
	ErrCardTrackingStarted = errors.New("card tracking already in progress")
	ErrCardTrackingStopped = errors.New("card tracking already stopped")
)

type CardStatus int

const (
	Todo CardStatus = iota
	Done
	Active
)

type UpdateCardParams struct {
	Title         string `json:"title"`
	EstimatedMins int    `json:"estimatedMins"`
	Description   string `json:"description"`
}

const userId = 1

type ICardService interface {
	GetAll(projectId uint, status CardStatus) ([]database.ListCardsRow, error)
	GetCardById(projectId uint, id uint) (*database.GetCardRow, error)
	GetActiveTimeEntry(projectId uint, id uint) (*database.TimeEntry, error)
	DeleteCard(projectId uint, id uint) error
	UpdateCard(projectId uint, id uint, updateCardParam UpdateCardParams) error
	UpdateCardStatus(projectId uint, id uint, status CardStatus) error
	AddCard(projectId uint, cardTitle string, estimatedMins uint) error
	StartCard(projectId uint, id uint) error
	StopCard(projectId uint, id uint) error
	Cleanup() error
}

type CardService struct {
	ctx                   context.Context
	projectService        IProjectService
	taskCompletionService ITaskCompletionService
	dbManager             *connection.DBManager
	eventBus              *events.EventBus
}

func NewCardService(projectService IProjectService, taskCompletionService ITaskCompletionService, dbManager *connection.DBManager, eventBus *events.EventBus) *CardService {
	return &CardService{
		ctx:                   context.Background(),
		projectService:        projectService,
		taskCompletionService: taskCompletionService,
		dbManager:             dbManager,
		eventBus:              eventBus,
	}
}

func (c *CardService) GetAll(projectId uint, status CardStatus) ([]database.ListCardsRow, error) {
	if _, err := c.projectService.IsValidProject(projectId); err != nil {
		return nil, err
	}

	queries := c.dbManager.Queries(c.ctx)
	return queries.ListCards(c.ctx, database.ListCardsParams{Projectid: int64(projectId), Status: int64(status)})
}

func (c *CardService) GetCardById(projectId uint, id uint) (*database.GetCardRow, error) {
	if _, err := c.projectService.IsValidProject(projectId); err != nil {
		return nil, err
	}

	queries := c.dbManager.Queries(c.ctx)
	card, err := queries.GetCard(c.ctx, database.GetCardParams{
		ID:        int64(id),
		Projectid: int64(projectId),
	})
	if err != nil {
		return nil, ErrNotFound
	}
	return &card, nil
}

func (c *CardService) GetActiveTimeEntry(projectId uint, id uint) (*database.TimeEntry, error) {
	if _, err := c.projectService.IsValidProject(projectId); err != nil {
		return nil, err
	}

	queries := c.dbManager.Queries(c.ctx)

	res, err := queries.GetActiveTimeEntry(c.ctx, int64(id))
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *CardService) DeleteCard(projectId uint, id uint) error {
	if _, err := c.projectService.IsValidProject(projectId); err != nil {
		return err
	}

	return c.dbManager.Execute(c.ctx, func(q *database.Queries) error {
		return q.DeleteCard(c.ctx, database.DeleteCardParams{
			ID:        int64(id),
			Projectid: int64(projectId),
		})
	})
}

func (c *CardService) UpdateCard(projectId uint, id uint, updateCardParam UpdateCardParams) error {
	if _, err := c.projectService.IsValidProject(projectId); err != nil {
		return err
	}

	return c.dbManager.Execute(c.ctx, func(q *database.Queries) error {
		card, err := q.GetCard(c.ctx, database.GetCardParams{
			ID:        int64(id),
			Projectid: int64(projectId),
		})
		if err != nil {
			return err
		}

		if updateCardParam.Title == "" {
			return ErrInvalidUpdate
		}

		var description sql.NullString
		if updateCardParam.Description != "" {
			description = sql.NullString{Valid: true, String: updateCardParam.Description}
		}

		return q.UpdateCard(c.ctx, database.UpdateCardParams{
			Title:         updateCardParam.Title,
			Description:   description,
			ID:            card.CardID,
			Status:        card.Status,
			Trackedmins:   card.Trackedmins,
			Estimatedmins: int64(updateCardParam.EstimatedMins),
			Completedat:   card.Completedat,
		})
	})
}

func (c *CardService) UpdateCardStatus(projectId uint, id uint, status CardStatus) error {
	if _, err := c.projectService.IsValidProject(projectId); err != nil {
		return err
	}

	return c.dbManager.Execute(c.ctx, func(q *database.Queries) error {
		card, err := q.GetCard(c.ctx, database.GetCardParams{
			ID:        int64(id),
			Projectid: int64(projectId),
		})
		if err != nil {
			return err
		}

		if status != Todo && status != Done {
			return ErrInvalidStatus
		}

		completedAt := sql.NullTime{}
		if status == Done {
			completedAt = sql.NullTime{Valid: true, Time: time.Now().UTC()}
		}

		err = q.UpdateCard(c.ctx, database.UpdateCardParams{
			Title:         card.Title,
			Description:   card.Description,
			ID:            card.CardID,
			Status:        int64(status),
			Trackedmins:   card.Trackedmins,
			Estimatedmins: card.Estimatedmins,
			Completedat:   completedAt,
		})
		if err != nil {
			return err
		}

		if status == Done {
			baseExp := int64(10)
			timeBonusExp := int64(card.Trackedmins / 5)
			streakBonusExp := int64(0)

			_, err := q.GetTaskCompletion(c.ctx, database.GetTaskCompletionParams{
				Cardid: card.CardID,
				Userid: userId,
			})

			if errors.Is(err, sql.ErrNoRows) {
				_, err = q.CreateTaskCompletion(c.ctx, database.CreateTaskCompletionParams{
					Cardid:         card.CardID,
					Userid:         userId,
					Baseexp:        baseExp,
					Timebonusexp:   timeBonusExp,
					Streakbonusexp: streakBonusExp,
					Totalexp:       baseExp + timeBonusExp + streakBonusExp,
				})
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (c *CardService) AddCard(projectId uint, cardTitle string, estimatedMins uint) error {
	if cardTitle == "" {
		return ErrCardTitleRequired
	}

	card := database.CreateCardParams{
		Title:         cardTitle,
		Status:        int64(Todo),
		Projectid:     int64(projectId),
		Estimatedmins: int64(estimatedMins),
	}

	return c.dbManager.Execute(c.ctx, func(q *database.Queries) error {
		return q.CreateCard(c.ctx, card)
	})
}

func (c *CardService) StartCard(projectId uint, id uint) error {
	if _, err := c.projectService.IsValidProject(projectId); err != nil {
		return err
	}

	var startedEvent events.CardStartedEvent
	var stoppedEvent *events.CardStoppedEvent

	err := c.dbManager.Execute(c.ctx, func(q *database.Queries) error {
		card, err := q.GetCard(c.ctx, database.GetCardParams{ID: int64(id), Projectid: int64(projectId)})
		if err != nil {
			return err
		}
		if card.Isactive {
			return ErrCardTrackingStarted
		}

		activeCard, err := q.GetActiveCard(c.ctx)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}

		if !errors.Is(err, sql.ErrNoRows) {
			event, err := c.stopCardLogic(q, uint(activeCard.Projectid), uint(activeCard.ID))
			if err != nil {
				return err
			}
			stoppedEvent = &event
		}

		event, err := c.startCardLogic(q, card)
		if err != nil {
			return err
		}
		startedEvent = event
		return nil
	})

	if err != nil {
		return err
	}

	if stoppedEvent != nil {
		c.eventBus.Publish(events.CardStoppedTopic, *stoppedEvent)
		log.Printf("Published CardStoppedEvent: %+v", *stoppedEvent)
	}
	c.eventBus.Publish(events.CardStartedTopic, startedEvent)
	log.Printf("Published CardStartedEvent: %+v", startedEvent)

	return nil
}

func (c *CardService) StopCard(projectId uint, id uint) error {
	if _, err := c.projectService.IsValidProject(projectId); err != nil {
		return err
	}

	var stoppedEvent events.CardStoppedEvent
	err := c.dbManager.Execute(c.ctx, func(q *database.Queries) error {
		event, err := c.stopCardLogic(q, projectId, id)
		if err != nil {
			return err
		}
		stoppedEvent = event
		return nil
	})

	if err != nil {
		return err
	}

	c.eventBus.Publish(events.CardStoppedTopic, stoppedEvent)
	log.Printf("Published CardStoppedEvent: %+v", stoppedEvent)
	return nil
}

func (c *CardService) Cleanup() error {
	log.Println("Cleaning up active card if any...")
	queries := c.dbManager.Queries(c.ctx)
	activeCard, err := queries.GetActiveCard(c.ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("No active card found. Cleanup completed successfully.")
			return nil
		}
		return err
	}

	log.Println("Active card found:", activeCard)
	err = c.StopCard(uint(activeCard.Projectid), uint(activeCard.ID))
	if err != nil {
		return err
	}

	log.Printf("Cleanup completed successfully. Active card with %d id have been stopped.", activeCard.ID)
	return nil
}

// stopCardLogic contains the core logic for stopping a card, designed to be used within a transaction.
func (c *CardService) stopCardLogic(q *database.Queries, projectId uint, id uint) (events.CardStoppedEvent, error) {
	card, err := q.GetCard(c.ctx, database.GetCardParams{ID: int64(id), Projectid: int64(projectId)})
	if err != nil {
		return events.CardStoppedEvent{}, err
	}
	if !card.Isactive {
		return events.CardStoppedEvent{}, ErrCardTrackingStopped
	}

	activeTimeEntry, err := q.GetActiveTimeEntry(c.ctx, int64(id))
	if err != nil {
		return events.CardStoppedEvent{}, err
	}

	currentEndTime := time.Now().UTC()
	duration := currentEndTime.Sub(activeTimeEntry.Starttime).Minutes()

	err = q.UpdateActiveTimeEntry(c.ctx, database.UpdateActiveTimeEntryParams{
		ID:       activeTimeEntry.ID,
		Endtime:  currentEndTime,
		Duration: int64(duration),
	})
	if err != nil {
		return events.CardStoppedEvent{}, err
	}

	newTrackedMins := card.Trackedmins + int64(duration)
	err = q.UpdateCardActive(c.ctx, database.UpdateCardActiveParams{
		ID:          int64(id),
		Isactive:    false,
		Trackedmins: newTrackedMins,
	})
	if err != nil {
		return events.CardStoppedEvent{}, err
	}

	log.Println("Card updated to inactive:", id, "with tracked mins:", newTrackedMins)
	return events.CardStoppedEvent{
		CardID:    card.CardID,
		ProjectID: card.Projectid,
		UserID:    userId,
		TimeSpent: time.Duration(duration) * time.Minute,
		StoppedAt: currentEndTime,
	}, nil
}

// startCardLogic contains the core logic for starting a card, designed to be used within a transaction.
func (c *CardService) startCardLogic(q *database.Queries, card database.GetCardRow) (events.CardStartedEvent, error) {
	err := q.UpdateCardActive(c.ctx, database.UpdateCardActiveParams{
		ID:          card.CardID,
		Isactive:    true,
		Trackedmins: card.Trackedmins,
	})
	if err != nil {
		return events.CardStartedEvent{}, err
	}

	currentStartTime := time.Now().UTC()
	_, err = q.CreateTimeEntry(c.ctx, database.CreateTimeEntryParams{
		Cardid:    card.CardID,
		Starttime: currentStartTime,
		Endtime:   currentStartTime,
	})
	if err != nil {
		return events.CardStartedEvent{}, err
	}

	return events.CardStartedEvent{
		CardID:    card.CardID,
		ProjectID: card.Projectid,
		UserID:    userId,
		StartedAt: currentStartTime,
	}, nil
}
