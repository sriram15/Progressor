package service

import (
	"context"
	"fmt"
	"log"

	"github.com/sriram15/progressor-todo-app/internal/connection"
	"github.com/sriram15/progressor-todo-app/internal/database"
)

type IProjectService interface {
	IsValidProject(projectId uint) (bool, error)
	AddProjectSkill(ctx context.Context, projectID, skillID int64) error
	RemoveProjectSkill(ctx context.Context, projectID, skillID int64) error
	GetSkillsForProject(ctx context.Context, projectID int64) ([]database.UserSkill, error)
	GetProjects() ([]database.Project, error)
}

type ProjectService struct {
	dbManager *connection.DBManager
}

func NewProjectService(dbManager *connection.DBManager) *ProjectService {
	return &ProjectService{
		dbManager: dbManager,
	}
}

func (p *ProjectService) IsValidProject(projectId uint) (bool, error) {

	// TODO: Access the DB adn validate the projectId
	// Returning based ont he default project for now.
	if projectId == 1 {
		return true, nil
	}
	return false, ErrInvalidProject
}

func (p *ProjectService) AddProjectSkill(ctx context.Context, projectID, skillID int64) error {
	return p.dbManager.Execute(ctx, func(q *database.Queries) error {
		err := q.AddProjectSkill(ctx, database.AddProjectSkillParams{
			ProjectID: projectID,
			SkillID:   skillID,
		})
		if err != nil {
			log.Printf("Error adding project skill: %v", err)
			return fmt.Errorf("failed to add project skill: %w", err)
		}
		return nil
	})
}

func (p *ProjectService) RemoveProjectSkill(ctx context.Context, projectID, skillID int64) error {
	return p.dbManager.Execute(ctx, func(q *database.Queries) error {
		err := q.RemoveProjectSkill(ctx, database.RemoveProjectSkillParams{
			ProjectID: projectID,
			SkillID:   skillID,
		})
		if err != nil {
			log.Printf("Error removing project skill: %v", err)
			return fmt.Errorf("failed to remove project skill: %w", err)
		}
		return nil
	})
}

func (p *ProjectService) GetSkillsForProject(ctx context.Context, projectID int64) ([]database.UserSkill, error) {
	queries := p.dbManager.Queries(ctx)
	skills, err := queries.GetSkillsForProject(ctx, projectID)
	if err != nil {
		log.Printf("Error getting skills for project: %v", err)
		return nil, fmt.Errorf("failed to get skills for project: %w", err)
	}
	return skills, nil
}

func (p *ProjectService) GetProjects() ([]database.Project, error) {
	return []database.Project{
		{ID: 1, Name: "Inbox"},
	}, nil
}
