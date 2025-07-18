package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/sriram15/progressor-todo-app/internal/connection"
	"github.com/sriram15/progressor-todo-app/internal/database"
	"github.com/sriram15/progressor-todo-app/internal/events"
)

const skillServiceUserID = 1 // Assuming a single user system for now

type ISkillService interface {
	CreateSkill(ctx context.Context, userID int64, name string, description string) (*database.UserSkill, error)
	GetSkillByID(ctx context.Context, id int64) (*database.UserSkill, error)
	GetSkillsByUserID(ctx context.Context, userID int64) ([]database.UserSkill, error)
	UpdateSkill(ctx context.Context, id int64, name string, description string) (*database.UserSkill, error)
	DeleteSkill(ctx context.Context, id int64) error
	GetUserSkillProgress(ctx context.Context, userID, skillID int64) (*database.UserSkillProgress, error)
}

type SkillService struct {
	queries        *database.Queries
	eventBus       *events.EventBus
	projectService IProjectService
}


func NewSkillService(queries *database.Queries, eventBus *events.EventBus, projectService IProjectService) *SkillService {
	return &SkillService{
		queries:        queries,
		eventBus:       eventBus,
		projectService: projectService,
	}
}

func (s *SkillService) CreateSkill(ctx context.Context, userID int64, name string, description string) (*database.UserSkill, error) {
	db, unlock := connection.GetDB()
	defer unlock()
	skill, err := s.queries.CreateSkill(ctx, db, database.CreateSkillParams{
		UserID:      userID,
		Name:        name,
		Description: sql.NullString{String: description, Valid: description != ""},
	})
	if err != nil {
		log.Printf("Error creating skill: %v", err)
		return nil, fmt.Errorf("failed to create skill: %w", err)
	}
	return &skill, nil
}

func (s *SkillService) GetSkillByID(ctx context.Context, id int64) (*database.UserSkill, error) {
	db, unlock := connection.GetDB()
	defer unlock()
	skill, err := s.queries.GetSkillByID(ctx, db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("skill with ID %d not found", id)
		}
		log.Printf("Error getting skill by ID: %v", err)
		return nil, fmt.Errorf("failed to get skill: %w", err)
	}
	return &skill, nil
}

func (s *SkillService) GetSkillsByUserID(ctx context.Context, userID int64) ([]database.UserSkill, error) {
	db, unlock := connection.GetDB()
	defer unlock()
	skills, err := s.queries.GetSkillsByUserID(ctx, db, userID)
	if err != nil {
		log.Printf("Error getting skills by user ID: %v", err)
		return nil, fmt.Errorf("failed to get skills for user: %w", err)
	}
	return skills, nil
}

func (s *SkillService) UpdateSkill(ctx context.Context, id int64, name string, description string) (*database.UserSkill, error) {
	db, unlock := connection.GetDB()
	defer unlock()
	skill, err := s.queries.UpdateSkill(ctx, db, database.UpdateSkillParams{
		ID:          id,
		Name:        name,
		Description: sql.NullString{String: description, Valid: description != ""},
	})
	if err != nil {
		log.Printf("Error updating skill: %v", err)
		return nil, fmt.Errorf("failed to update skill: %w", err)
	}
	return &skill, nil
}

func (s *SkillService) DeleteSkill(ctx context.Context, id int64) error {
	db, unlock := connection.GetDB()
	defer unlock()
	err := s.queries.DeleteSkill(ctx, db, id)
	if err != nil {
		log.Printf("Error deleting skill: %v", err)
		return fmt.Errorf("failed to delete skill: %w", err)
	}
	return nil
}

func (s *SkillService) GetUserSkillProgress(ctx context.Context, userID, skillID int64) (*database.UserSkillProgress, error) {
	db, unlock := connection.GetDB()
	defer unlock()
	progress, err := s.queries.GetUserSkillProgress(ctx, db, database.GetUserSkillProgressParams{
		UserID:  userID,
		SkillID: skillID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No progress yet, not an error
		}
		log.Printf("Error getting user skill progress: %v", err)
		return nil, fmt.Errorf("failed to get user skill progress: %w", err)
	}
	return &progress, nil
}

func (s *SkillService) RegisterEventHandlers() {
	s.eventBus.Subscribe(events.CardStoppedTopic, s.handleCardStopped)
}

func (s *SkillService) handleCardStopped(eventData interface{}) {
	event, ok := eventData.(events.CardStoppedEvent)
	if !ok {
		log.Printf("Error: received non-CardStoppedEvent for topic %s", events.CardStoppedTopic)
		return
	}
	log.Printf("Received CardStoppedEvent: %+v", event)

	ctx := context.Background()

	projectSkills, err := s.projectService.GetSkillsForProject(ctx, event.ProjectID)
	if err != nil {
		log.Printf("Error getting skills for project %d: %v", event.ProjectID, err)
		return
	}

	durationMins := int64(event.TimeSpent.Minutes())

	// Upsert the user's skill progress for each skill associated with the project.
	db, unlock := connection.GetDB()
	defer unlock()

	// Upsert the user's skill progress for each skill associated with the project.
	for _, skill := range projectSkills {
		_, err := s.queries.UpsertUserSkillProgress(ctx, db, database.UpsertUserSkillProgressParams{
			UserID:              event.UserID,
			SkillID:             skill.ID,
			TotalMinutesTracked: sql.NullInt64{Int64: durationMins, Valid: true},
		})
		if err != nil {
			log.Printf("Error upserting skill progress for skill %d: %v", skill.ID, err)
			continue
		}
		log.Printf("Successfully updated skill progress for skill %d by %d minutes.", skill.ID, durationMins)
	}
}
