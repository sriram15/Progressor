package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/sriram15/progressor-todo-app/internal/database"
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

type CardService interface {
	GetAll(projectId uint) ([]database.ListCardsRow, error)
	GetOpenCards(projectId uint) ([]database.ListOpenOrCTCardsRow, error)
	GetCardById(projectId uint, id uint) (database.GetCardRow, error)
	DeleteCard(projectId uint, id uint) error
	UpdateCard(projectId uint, id uint, updateCardParam UpdateCardParams) error
	UpdateCardStatus(projectId uint, id uint, status CardStatus) error
	AddCard(projectId uint, cardTitle string, estimatedMins uint) error
	StartCard(projectId uint, id uint) error
	StopCard(projectId uint, id uint) error
	Cleanup() error
}

type cardService struct {
	ctx            context.Context
	db             *sql.DB
	queries        *database.Queries
	projectService ProjectService
}

func NewCardService(db *sql.DB, queries *database.Queries, projectService ProjectService) CardService {
	return &cardService{
		ctx:            context.Background(),
		db:             db,
		queries:        queries,
		projectService: projectService,
	}
}

func (c *cardService) GetAll(projectId uint) ([]database.ListCardsRow, error) {

	_, err := c.projectService.IsValidProject(projectId)
	if err != nil {
		return []database.ListCardsRow{}, err
	}
	cards, err := c.queries.ListCards(c.ctx, int64(projectId))
	return cards, err
}
func (c *cardService) GetOpenCards(projectId uint) ([]database.ListOpenOrCTCardsRow, error) {

	_, err := c.projectService.IsValidProject(projectId)
	if err != nil {
		return []database.ListOpenOrCTCardsRow{}, err
	}

	cards, err := c.queries.ListOpenOrCTCards(c.ctx, int64(projectId))
	return cards, err
}

func (c *cardService) GetCardById(projectId uint, id uint) (database.GetCardRow, error) {
	_, err := c.projectService.IsValidProject(projectId)
	if err != nil {
		return database.GetCardRow{}, err
	}

	card, err := c.queries.GetCard(c.ctx, database.GetCardParams{
		ID:        int64(id),
		Projectid: int64(projectId),
	})
	if err != nil {
		return database.GetCardRow{}, ErrNotFound
	}
	return card, nil
}

func (c *cardService) DeleteCard(projectId uint, id uint) error {

	_, err := c.projectService.IsValidProject(projectId)
	if err != nil {
		return err
	}

	return c.queries.DeleteCard(c.ctx, database.DeleteCardParams{
		ID:        int64(id),
		Projectid: int64(projectId),
	})
}

func (c *cardService) UpdateCard(projectId uint, id uint, updateCardParam UpdateCardParams) error {

	_, err := c.projectService.IsValidProject(projectId)
	if err != nil {
		return err
	}

	card, err := c.queries.GetCard(c.ctx, database.GetCardParams{
		ID:        int64(id),
		Projectid: int64(projectId),
	})
	if err != nil {
		return err
	}

	// TODO: Add update card validation
	if updateCardParam.Title == "" {
		return ErrInvalidUpdate
	}

	var description sql.NullString
	if updateCardParam.Description == "" {
		description = sql.NullString{Valid: true, String: ""}
	} else {
		description = sql.NullString{Valid: true, String: updateCardParam.Description}
	}

	return c.queries.UpdateCard(c.ctx, database.UpdateCardParams{
		Title:         updateCardParam.Title,
		Description:   description,
		ID:            card.CardID,
		Status:        card.Status,
		Trackedmins:   card.Trackedmins,
		Estimatedmins: int64(updateCardParam.EstimatedMins),
		Completedat:   card.Completedat,
	})

}
func (c *cardService) UpdateCardStatus(projectId uint, id uint, status CardStatus) error {
	_, err := c.projectService.IsValidProject(projectId)
	if err != nil {
		return err
	}

	card, err := c.queries.GetCard(c.ctx, database.GetCardParams{
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

	return c.queries.UpdateCard(c.ctx, database.UpdateCardParams{
		Title:         card.Title,
		Description:   card.Description,
		ID:            card.CardID,
		Status:        int64(status),
		Trackedmins:   card.Trackedmins,
		Estimatedmins: card.Estimatedmins,
		Completedat:   completedAt,
	})
}

func (c *cardService) AddCard(projectId uint, cardTitle string, estimatedMins uint) error {

	if cardTitle == "" {
		return ErrCardTitleRequired
	}

	card := database.CreateCardParams{
		Title:         cardTitle,
		Status:        int64(Todo),
		Projectid:     int64(projectId),
		Estimatedmins: int64(estimatedMins),
	}

	return c.queries.CreateCard(c.ctx, card)
}

func (c *cardService) StartCard(projectId uint, id uint) error {
	_, err := c.projectService.IsValidProject(projectId)
	if err != nil {
		return err
	}

	card, err := c.queries.GetCard(c.ctx, database.GetCardParams{
		ID:        int64(id),
		Projectid: int64(projectId),
	})
	if err != nil {
		return err
	}

	// Check if card is already active
	if card.Isactive {
		return ErrCardTrackingStarted
	}

	// Check for other open cards which is currently in progress and stop the timer there
	activeCard, err := c.queries.GetActiveCard(c.ctx)

	// When the active card is empty. It will throw sql.ErrNoRows. If the err is not that, then return err
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	} else { // We have an active card already. Try to stop that first and return err if that fails
		err := c.StopCard(projectId, uint(activeCard.ID))
		if err != nil {
			return err
		}
	}

	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := c.queries.WithTx(tx)

	err = qtx.UpdateCardActive(c.ctx, database.UpdateCardActiveParams{
		ID:          int64(id),
		Isactive:    true,
		Trackedmins: card.Trackedmins,
	})
	if err != nil {
		return err
	}

	// Create a new Timeentry object and add it to the card
	currentStartTime := time.Now().UTC()
	_, err = qtx.CreateTimeEntry(c.ctx, database.CreateTimeEntryParams{
		Cardid:    int64(id),
		Starttime: currentStartTime,
		Endtime:   currentStartTime,
	})
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (c *cardService) StopCard(projectId uint, id uint) error {
	_, err := c.projectService.IsValidProject(projectId)
	if err != nil {
		return err
	}

	card, err := c.queries.GetCard(c.ctx, database.GetCardParams{
		ID:        int64(id),
		Projectid: int64(projectId),
	})
	if err != nil {
		return err
	}

	// Check if card is still active
	if !card.Isactive {
		return ErrCardTrackingStopped
	}

	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := c.queries.WithTx(tx)

	// Get the active time entry
	activeTimeentry, err := qtx.GetActiveTimeEntry(c.ctx, int64(id))
	if err != nil {
		return err
	}

	currentEndTime := time.Now().UTC()
	duration := currentEndTime.Sub(activeTimeentry.Starttime).Minutes()
	err = qtx.UpdateActiveTimeEntry(c.ctx, database.UpdateActiveTimeEntryParams{
		ID:       activeTimeentry.ID,
		Endtime:  currentEndTime,
		Duration: int64(duration),
	})
	if err != nil {
		return err
	}

	err = qtx.UpdateCardActive(c.ctx, database.UpdateCardActiveParams{
		ID:          int64(id),
		Isactive:    false,
		Trackedmins: card.Trackedmins + int64(duration),
	})
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (c *cardService) Cleanup() error {

	// Check for other open cards which is currently in progress and stop the timer there
	activeCard, err := c.queries.GetActiveCard(c.ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		} else {
			return err
		}
	}

	// Stop the active Card Now
	err = c.StopCard(uint(activeCard.Projectid), uint(activeCard.ID))
	if err != nil {
		return err
	}

	return nil
}
