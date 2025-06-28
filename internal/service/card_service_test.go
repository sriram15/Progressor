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

