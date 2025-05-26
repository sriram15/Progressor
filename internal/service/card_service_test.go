package service

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"github.com/sriram15/progressor-todo-app/internal/database"
	"github.com/stretchr/testify/assert"
)

var db *sql.DB

func SetupDB(t *testing.T) *database.Queries {

	t.Helper()
	var err error
	db, err = goose.OpenDBWithDriver("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	goose.SetDialect("sqlite3")

	// Run migrations
	if err := goose.Up(db, "../database/migrations"); err != nil {
		panic(err)
	}
	assert.NoError(t, err)

	queries := database.New(db)
	return queries
}

func teardown(t *testing.T) {
	t.Helper()
	err := db.Close()
	assert.NoError(t, err)
}

func TestCardService(t *testing.T) {

	queries := SetupDB(t)
	defer teardown(t)
	t.Run("Get All", func(t *testing.T) {

		cardService := NewCardService(queries)
		cards, err := cardService.GetAll()

		assert.NoError(t, err)
		assert.Equal(t, []database.Card(nil), cards)
	})
}

func TestUpdateCardStatus_ActiveToDone(t *testing.T) {
	queries := SetupDB(t)
	defer teardown(t)

	// Create mock services
	projectService := NewProjectService()
	taskCompletionService := NewTaskCompletionService()
	cardService := NewCardService(projectService, taskCompletionService)

	// 1. Create a project and a card
	projectName := "Test Project"
	project, err := projectService.AddProject(projectName)
	assert.NoError(t, err)

	cardTitle := "Test Card"
	estimatedMins := uint(30)
	err = cardService.AddCard(uint(project.ID), cardTitle, estimatedMins)
	assert.NoError(t, err)

	// Retrieve the card to get its ID
	cards, err := cardService.GetAll(uint(project.ID), Todo)
	assert.NoError(t, err)
	assert.Len(t, cards, 1)
	cardId := cards[0].ID

	// 2. Starts the card (makes it active).
	err = cardService.StartCard(uint(project.ID), uint(cardId))
	assert.NoError(t, err)

	// 3. Marks the card as complete using `UpdateCardStatus`.
	err = cardService.UpdateCardStatus(uint(project.ID), uint(cardId), Done)
	assert.NoError(t, err)

	// 4. Verifies that the card's status is updated to "Done".
	updatedCard, err := cardService.GetCardById(uint(project.ID), uint(cardId))
	assert.NoError(t, err)
	assert.Equal(t, int64(Done), updatedCard.Status)
	assert.True(t, updatedCard.Completedat.Valid) // Ensure Completedat is set

	// 5. Verifies that a corresponding task completion record is created.
	// Assuming userId is 1 as per card_service.go
	taskCompletion, err := queries.GetTaskCompletion(cardService.ctx, database.GetTaskCompletionParams{
		Cardid: cardId,
		Userid: 1,
	})
	assert.NoError(t, err)
	assert.NotNil(t, taskCompletion)
	assert.Equal(t, cardId, taskCompletion.Cardid)
	assert.Equal(t, int64(1), taskCompletion.Userid) // Ensure correct UserId
	assert.Greater(t, taskCompletion.Totalexp, int64(0)) // Ensure some EXP is awarded
}

// func TestCardService_AddCard(t *testing.T) {

// 	mockRewardRepo := new(repository.MockRewardRepository)
// 	mockCardRepo := new(repository.MockCardRepository)

// 	t.Run("Add Card", func(t *testing.T) {

// 		mockCardRepo.On("AddCard", mock.Anything).Return(nil)

// 		cardService := NewCardService(mockCardRepo, mockRewardRepo)
// 		err := cardService.AddCard(model.Card{Title: "mock card", Description: "test"})

// 		assert.NoError(t, err)
// 		mockCardRepo.AssertExpectations(t)

// 	})

// 	t.Run("Add Card with empty title should return error", func(t *testing.T) {

// 		mockRewardRepo := new(repository.MockRewardRepository)
// 		mockCardRepo := new(repository.MockCardRepository)

// 		cardService := NewCardService(mockCardRepo, mockRewardRepo)
// 		err := cardService.AddCard(model.Card{Title: "", Description: "test"})

// 		assert.Error(t, err)
// 		mockCardRepo.AssertNotCalled(t, "AddCard", mock.Anything)
// 	})
// }

// func TestCardService_GetCardById(t *testing.T) {

// 	mockRewardRepo := new(repository.MockRewardRepository)
// 	mockCardRepo := new(repository.MockCardRepository)

// 	t.Run("Get Card By Id - Success", func(t *testing.T) {
// 		expectedCard := model.Card{Title: "mock card", Description: "test"}
// 		mockCardRepo.On("GetById", uint(1)).Return(expectedCard, nil)

// 		cardService := NewCardService(mockCardRepo, mockRewardRepo)
// 		card, err := cardService.GetCardById(1)

// 		assert.NoError(t, err)
// 		assert.Equal(t, &expectedCard, card)
// 		mockCardRepo.AssertExpectations(t)
// 	})

// 	t.Run("Get Card By Id - Error", func(t *testing.T) {
// 		mockCardRepo.On("GetById", uint(2)).Return(model.Card{}, errors.New("not found"))

// 		cardService := NewCardService(mockCardRepo, mockRewardRepo)
// 		_, err := cardService.GetCardById(2)

// 		assert.Error(t, err)
// 		assert.EqualError(t, err, ErrNotFound.Error())
// 		mockCardRepo.AssertExpectations(t)
// 	})
// }

// func TestCardService_DeleteCard(t *testing.T) {

// 	mockCardRepo := new(repository.MockCardRepository)

// 	t.Run("Delete Card - Success", func(t *testing.T) {
// 		mockCardRepo.On("DeleteCard", uint(1)).Return(nil)

// 		cardService := NewCardService(mockCardRepo, nil)
// 		err := cardService.DeleteCard(1)

// 		assert.NoError(t, err)
// 		mockCardRepo.AssertExpectations(t)
// 	})

// 	t.Run("Delete Card - Error", func(t *testing.T) {
// 		mockCardRepo.On("DeleteCard", uint(2)).Return(errors.New("not found"))

// 		cardService := NewCardService(mockCardRepo, nil)
// 		err := cardService.DeleteCard(2)

// 		assert.Error(t, err)
// 		assert.EqualError(t, err, ErrNotFound.Error())
// 		mockCardRepo.AssertExpectations(t)
// 	})
// }

// func TestCardService_UpdateCardStatus(t *testing.T) {

// 	t.Run("Update Card Status - Success", func(t *testing.T) {
// 		mockRewardRepo := new(repository.MockRewardRepository)
// 		mockCardRepo := new(repository.MockCardRepository)

// 		card := model.Card{Title: "mock card", Description: "test", Status: model.Todo}
// 		updatedStatus := model.Done

// 		mockCardRepo.On("GetById", uint(1)).Return(card, nil)
// 		mockCardRepo.On("UpdateCard", mock.Anything).Return(nil)
// 		mockRewardRepo.On("CreateReward", mock.Anything).Return(nil)

// 		cardService := NewCardService(mockCardRepo, mockRewardRepo)
// 		err := cardService.UpdateCardStatus(1, updatedStatus)

// 		assert.NoError(t, err)
// 		mockCardRepo.AssertExpectations(t)
// 		mockRewardRepo.AssertExpectations(t)
// 	})

// 	t.Run("Update Card Status - Invalid Status", func(t *testing.T) {
// 		mockRewardRepo := new(repository.MockRewardRepository)
// 		mockCardRepo := new(repository.MockCardRepository)

// 		card := model.Card{Title: "mock card", Description: "test", Status: model.Todo}
// 		invalidStatus := model.CardStatus("in-progress")

// 		mockCardRepo.On("GetById", uint(1)).Return(card, nil)

// 		cardService := NewCardService(mockCardRepo, mockRewardRepo)
// 		err := cardService.UpdateCardStatus(1, invalidStatus)

// 		assert.Error(t, err)
// 		assert.EqualError(t, err, ErrInvalidStatus.Error())
// 		mockCardRepo.AssertNotCalled(t, "UpdateCard", mock.Anything)
// 		mockRewardRepo.AssertNotCalled(t, "CreateReward", mock.Anything)
// 	})

// 	t.Run("Update Card Status - Card Not Found", func(t *testing.T) {
// 		mockRewardRepo := new(repository.MockRewardRepository)
// 		mockCardRepo := new(repository.MockCardRepository)

// 		mockCardRepo.On("GetById", uint(2)).Return(model.Card{}, errors.New("not found"))

// 		cardService := NewCardService(mockCardRepo, mockRewardRepo)
// 		err := cardService.UpdateCardStatus(2, model.Done)

// 		assert.Error(t, err)
// 		assert.EqualError(t, err, ErrNotFound.Error())
// 		mockCardRepo.AssertExpectations(t)
// 		mockRewardRepo.AssertNotCalled(t, "CreateReward", mock.Anything)
// 	})
// }

// func TestCardService_StartCard(t *testing.T) {
// 	t.Run("Start Card should start the timer", func(t *testing.T) {

// 		mockRewardRepo := new(repository.MockRewardRepository)
// 		mockCardRepo := new(repository.MockCardRepository)

// 		mockCard := model.Card{Title: "mock card", Description: "test", Status: model.Todo}
// 		mockCardRepo.On("GetById", mock.Anything).Return(mockCard, nil)
// 		mockCardRepo.On("UpdateCard", mock.Anything).Return(nil)

// 		cardService := NewCardService(mockCardRepo, mockRewardRepo)
// 		err := cardService.StartCard(1)

// 		assert.NoError(t, err)

// 	})

// 	t.Run("Start Card should start the timer when the previous timer is complete", func(t *testing.T) {

// 		mockRewardRepo := new(repository.MockRewardRepository)
// 		mockCardRepo := new(repository.MockCardRepository)

// 		mockCard := model.Card{
// 			Title: "mock card", Description: "test", Status: model.Todo,
// 			TimeEntries: []model.TimeEntry{
// 				{
// 					StartTime: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
// 					EndTime:   sql.NullTime{Time: time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC), Valid: true},
// 				},
// 			},
// 		}
// 		mockCardRepo.On("GetById", mock.Anything).Return(mockCard, nil)
// 		mockCardRepo.On("UpdateCard", mock.Anything).Return(nil)

// 		cardService := NewCardService(mockCardRepo, mockRewardRepo)
// 		err := cardService.StartCard(1)

// 		assert.NoError(t, err)

// 	})

// 	t.Run("Start Card should return error when the card is not found", func(t *testing.T) {

// 		mockRewardRepo := new(repository.MockRewardRepository)
// 		mockCardRepo := new(repository.MockCardRepository)

// 		mockCardRepo.On("GetById", mock.Anything).Return(model.Card{}, errors.New("not found"))
// 		mockCardRepo.On("UpdateCard", mock.Anything).Return(nil)

// 		cardService := NewCardService(mockCardRepo, mockRewardRepo)
// 		err := cardService.StartCard(1)

// 		assert.Error(t, err)
// 		assert.EqualError(t, err, ErrNotFound.Error())
// 	})

// 	t.Run("Start Card should return error when the card is already started", func(t *testing.T) {

// 		mockRewardRepo := new(repository.MockRewardRepository)
// 		mockCardRepo := new(repository.MockCardRepository)

// 		now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
// 		mockCard := model.Card{Title: "mock card", Description: "test", Status: model.Todo, TimeEntries: []model.TimeEntry{
// 			{StartTime: now},
// 			{StartTime: now, EndTime: sql.NullTime{Time: now, Valid: true}},
// 		}}
// 		mockCardRepo.On("GetById", mock.Anything).Return(mockCard, nil)
// 		mockCardRepo.On("UpdateCard", mock.Anything).Return(nil)

// 		cardService := NewCardService(mockCardRepo, mockRewardRepo)
// 		err := cardService.StartCard(1)

// 		assert.Error(t, err)
// 		assert.EqualError(t, err, ErrCardTrackingStarted.Error())
// 	})
// }

// func TestCardService_StopCard(t *testing.T) {

// 	t.Run("Stop Card should return error when the card is not found", func(t *testing.T) {

// 		mockRewardRepo := new(repository.MockRewardRepository)
// 		mockCardRepo := new(repository.MockCardRepository)

// 		mockCardRepo.On("GetById", mock.Anything).Return(model.Card{}, errors.New("not found"))

// 		cardService := NewCardService(mockCardRepo, mockRewardRepo)
// 		err := cardService.StopCard(1)

// 		assert.Error(t, err)
// 		assert.EqualError(t, err, ErrNotFound.Error())
// 	})

// 	t.Run("Stop Card should return error when the card is already stopped", func(t *testing.T) {

// 		mockRewardRepo := new(repository.MockRewardRepository)
// 		mockCardRepo := new(repository.MockCardRepository)

// 		mockCard := model.Card{
// 			Title: "not available",
// 			TimeEntries: []model.TimeEntry{
// 				{StartTime: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), EndTime: sql.NullTime{Time: time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC), Valid: true}},
// 			},
// 		}
// 		mockCardRepo.On("GetById", mock.Anything).Return(mockCard, nil)

// 		cardService := NewCardService(mockCardRepo, mockRewardRepo)
// 		err := cardService.StopCard(1)

// 		assert.Error(t, err)
// 		assert.EqualError(t, err, ErrCardTrackingStopped.Error())
// 	})

// 	t.Run("Stop card should calculate total time", func(t *testing.T) {

// 		mockRewardRepo := new(repository.MockRewardRepository)
// 		mockCardRepo := new(repository.MockCardRepository)

// 		mockCard := model.Card{
// 			Title:       "mock card",
// 			Description: "test",
// 			Status:      model.Todo,
// 			TimeEntries: []model.TimeEntry{
// 				{StartTime: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
// 			},
// 		}
// 		mockCardRepo.On("GetById", mock.Anything).Return(mockCard, nil)
// 		mockCardRepo.On("UpdateCard", mock.Anything).Return(nil)

// 		cardService := NewCardService(mockCardRepo, mockRewardRepo)
// 		err := cardService.StopCard(1)

// 		assert.NoError(t, err)
// 	})

// }
