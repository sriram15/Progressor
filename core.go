package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/sriram15/progressor-todo-app/internal/connection"
	"github.com/sriram15/progressor-todo-app/internal/database"
	"github.com/sriram15/progressor-todo-app/internal/events"
	"github.com/sriram15/progressor-todo-app/internal/profile"
	"github.com/sriram15/progressor-todo-app/internal/service"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// ProgressorApp is the main application struct. It holds a pointer to the current session,
// allowing for dynamic switching between profiles.
type ProgressorApp struct {
	currentSession *AppSession
	wailsApp       *application.App
	eventBus       *events.EventBus
	profileManager *profile.ProfileManager
	sessionMutex   sync.RWMutex
}

// AppSession holds all services and the database manager for a single, active profile.
type AppSession struct {
	dbManager             *connection.DBManager
	cardService           *service.CardService
	progressService       *service.ProgressService
	settingsService       *service.SettingService
	projectService        *service.ProjectService
	skillService          *service.SkillService
	taskCompletionService *service.TaskCompletionService
	focusTimerService     *service.FocusTimerService
}

// NewProgressorApp creates a new App object and initializes the profile manager.
func NewProgressorApp() *ProgressorApp {
	pm, err := profile.NewManager()
	if err != nil {
		log.Fatalf("Failed to create profile manager: %v", err)
	}
	return &ProgressorApp{
		profileManager: pm,
	}
}

// Startup is called when the app starts. This is where we can initialize things.
func (a *ProgressorApp) Startup(app *application.App) {
	log.Println("Progressor app startup")
	a.wailsApp = app
	a.eventBus = events.NewEventBus()
}

// Shutdown is called when the app is shutting down.
func (a *ProgressorApp) Shutdown() {
	a.sessionMutex.RLock()
	defer a.sessionMutex.RUnlock()
	if a.currentSession != nil {
		a.currentSession.focusTimerService.Shutdown()
		err := a.currentSession.cardService.Cleanup()
		fmt.Println("Shutdown cleanup done")
		if err != nil {
			log.Printf("Error during card service cleanup on shutdown: %v", err)
		}
	}
}

// NewAppSession creates a new session with all services initialized for a given database connection.
func NewAppSession(dbManager *connection.DBManager, eventBus *events.EventBus, wailsApp *application.App) (*AppSession, error) {
	projectService := service.NewProjectService(dbManager)
	taskCompletionService := service.NewTaskCompletionService(dbManager)
	settingsService := service.NewSettingService(dbManager)
	skillService := service.NewSkillService(dbManager, eventBus, projectService)
	progressService := service.NewProgressService(dbManager)
	cardService := service.NewCardService(projectService, taskCompletionService, dbManager, eventBus)
	focusTimerService := service.NewFocusTimerService(cardService, settingsService, eventBus, wailsApp)

	skillService.RegisterEventHandlers()
	focusTimerService.RegisterEventHandlers()

	log.Println("New AppSession created with DBManager")

	return &AppSession{
		dbManager:             dbManager,
		cardService:           cardService,
		progressService:       progressService,
		settingsService:       settingsService,
		projectService:        projectService,
		skillService:          skillService,
		taskCompletionService: taskCompletionService,
		focusTimerService:     focusTimerService,
	}, nil
}

// --- Profile Management Methods ---

func (a *ProgressorApp) GetProfiles() ([]profile.Profile, error) {
	profiles, _ := a.profileManager.GetProfiles()
	// pro, _ := a.profileManager.GetProfile(profiles[0].ID) // Ensure the first profile is loaded
	// log.Printf("Loaded profiles: %s\n", pro.Name)
	return profiles, nil;
}

func (a *ProgressorApp) CreateProfile(p profile.Profile, tursoToken string) (*profile.Profile, error) {
	return a.profileManager.CreateProfile(p, tursoToken)
}

func (a *ProgressorApp) SwitchProfile(profileID string) error {
	log.Printf("Attempting to switch to profile: %s", profileID)

	p, err := a.profileManager.GetProfile(profileID)
	if err != nil {
		return fmt.Errorf("failed to get profile %s: %w", profileID, err)
	}

	log.Printf("Switching to profile: %s (%s)", p.Name, p.ID)

	dbManager, err := connection.NewManagerForProfile(p)
	if err != nil {
		return fmt.Errorf("failed to create db manager for profile %s: %w", p.Name, err)
	}


	newSession, err := NewAppSession(dbManager, a.eventBus, a.wailsApp)
	if err != nil {
		dbManager.Close()
		return fmt.Errorf("failed to create new app session for profile %s: %w", p.Name, err)
	}

	stats, err := newSession.progressService.GetStats()
	if err != nil {
		log.Printf("Failed to get stats for profile %s: %v", p.Name, err)
	}

	log.Printf("Successfully created new session for profile: %s", stats.MonthProgress)
	a.sessionMutex.Lock()
	a.currentSession = newSession
	a.sessionMutex.Unlock()

	a.wailsApp.Event.Emit("profile:switched", p)

	log.Printf("Successfully switched to profile: %s", p.Name)
	return nil
}

// --- Service Method Delegates ---

func (a *ProgressorApp) withSession(fn func(s *AppSession) (interface{}, error)) (interface{}, error) {
	a.sessionMutex.RLock()
	defer a.sessionMutex.RUnlock()
	if a.currentSession == nil {
		return nil, fmt.Errorf("no active profile session")
	}
	return fn(a.currentSession)
}

// CardService delegates
func (a *ProgressorApp) AddCard(projectID uint, cardTitle string, estimatedMins uint) error {
	_, err := a.withSession(func(s *AppSession) (interface{}, error) {
		return nil, s.cardService.AddCard(projectID, cardTitle, estimatedMins)
	})
	return err
}

func (a *ProgressorApp) Cleanup() error {
	_, err := a.withSession(func(s *AppSession) (interface{}, error) {
		return nil, s.cardService.Cleanup()
	})
	return err
}

func (a *ProgressorApp) DeleteCard(projectID uint, id uint) error {
	_, err := a.withSession(func(s *AppSession) (interface{}, error) {
		return nil, s.cardService.DeleteCard(projectID, id)
	})
	return err
}

func (a *ProgressorApp) GetActiveTimeEntry(projectID uint, id uint) (*database.TimeEntry, error) {
	res, err := a.withSession(func(s *AppSession) (interface{}, error) {
		return s.cardService.GetActiveTimeEntry(projectID, id)
	})
	if err != nil {
		return nil, err
	}
	return res.(*database.TimeEntry), nil
}

func (a *ProgressorApp) GetAll(projectID uint, status service.CardStatus) ([]database.ListCardsRow, error) {
	res, err := a.withSession(func(s *AppSession) (interface{}, error) {
		return s.cardService.GetAll(projectID, status)
	})
	if err != nil {
		return nil, err
	}
	return res.([]database.ListCardsRow), nil
}

func (a *ProgressorApp) GetCardById(projectID uint, id uint) (*database.GetCardRow, error) {
	res, err := a.withSession(func(s *AppSession) (interface{}, error) {
		return s.cardService.GetCardById(projectID, id)
	})
	if err != nil {
		return nil, err
	}
	return res.(*database.GetCardRow), nil
}

func (a *ProgressorApp) StartCard(projectID uint, id uint) error {
	_, err := a.withSession(func(s *AppSession) (interface{}, error) {
		return nil, s.cardService.StartCard(projectID, id)
	})
	return err
}

func (a *ProgressorApp) StopCard(projectID uint, id uint) error {
	_, err := a.withSession(func(s *AppSession) (interface{}, error) {
		return nil, s.cardService.StopCard(projectID, id)
	})
	return err
}

func (a *ProgressorApp) UpdateCard(projectID uint, id uint, updateCardParam service.UpdateCardParams) error {
	_, err := a.withSession(func(s *AppSession) (interface{}, error) {
		return nil, s.cardService.UpdateCard(projectID, id, updateCardParam)
	})
	return err
}

func (a *ProgressorApp) UpdateCardStatus(projectID uint, id uint, status service.CardStatus) error {
	_, err := a.withSession(func(s *AppSession) (interface{}, error) {
		return nil, s.cardService.UpdateCardStatus(projectID, id, status)
	})
	return err
}

// ProgressService delegates
func (a *ProgressorApp) GetDailyTotalMinutes() ([]database.GetDailyTotalMinutesRow, error) {
	res, err := a.withSession(func(s *AppSession) (interface{}, error) {
		return s.progressService.GetDailyTotalMinutes()
	})
	if err != nil {
		return nil, err
	}
	return res.([]database.GetDailyTotalMinutesRow), nil
}

func (a *ProgressorApp) GetStats() (service.GetStatsResult, error) {
	res, err := a.withSession(func(s *AppSession) (interface{}, error) {
		return s.progressService.GetStats()
	})
	if err != nil {
		return service.GetStatsResult{}, err
	}
	return res.(service.GetStatsResult), nil
}

func (a *ProgressorApp) GetTotalExpForUser(userID int64) (float64, error) {
	res, err := a.withSession(func(s *AppSession) (interface{}, error) {
		return s.progressService.GetTotalExpForUser(userID)
	})
	if err != nil {
		return 0, err
	}
	return res.(float64), nil
}

// SettingService delegates
func (a *ProgressorApp) GetAllSettings() ([]service.SettingsItem, error) {
	res, err := a.withSession(func(s *AppSession) (interface{}, error) {
		return s.settingsService.GetAllSettings()
	})
	if err != nil {
		return nil, err
	}
	return res.([]service.SettingsItem), nil
}

func (a *ProgressorApp) GetSetting(key string) (string, error) {
	res, err := a.withSession(func(s *AppSession) (interface{}, error) {
		return s.settingsService.GetSetting(key)
	})
	if err != nil {
		return "", err
	}
	return res.(string), nil
}

func (a *ProgressorApp) SetSetting(key string, value string) error {
	_, err := a.withSession(func(s *AppSession) (interface{}, error) {
		return nil, s.settingsService.SetSetting(key, value)
	})
	return err
}

// SkillService delegates
func (a *ProgressorApp) CreateSkill(userID int64, name string, description string) (*database.UserSkill, error) {
	res, err := a.withSession(func(s *AppSession) (interface{}, error) {
		return s.skillService.CreateSkill(context.Background(), userID, name, description)
	})
	if err != nil {
		return nil, err
	}
	return res.(*database.UserSkill), nil
}

func (a *ProgressorApp) DeleteSkill(id int64) error {
	_, err := a.withSession(func(s *AppSession) (interface{}, error) {
		return nil, s.skillService.DeleteSkill(context.Background(), id)
	})
	return err
}

func (a *ProgressorApp) GetSkillByID(id int64) (*database.UserSkill, error) {
	res, err := a.withSession(func(s *AppSession) (interface{}, error) {
		return s.skillService.GetSkillByID(context.Background(), id)
	})
	if err != nil {
		return nil, err
	}
	return res.(*database.UserSkill), nil
}

func (a *ProgressorApp) GetSkillsByUserID(userID int64) ([]database.UserSkill, error) {
	res, err := a.withSession(func(s *AppSession) (interface{}, error) {
		return s.skillService.GetSkillsByUserID(context.Background(), userID)
	})
	if err != nil {
		return nil, err
	}
	return res.([]database.UserSkill), nil
}

func (a *ProgressorApp) GetUserSkillProgress(userID int64, skillID int64) (*database.UserSkillProgress, error) {
	res, err := a.withSession(func(s *AppSession) (interface{}, error) {
		return s.skillService.GetUserSkillProgress(context.Background(), userID, skillID)
	})
	if err != nil {
		return nil, err
	}
	return res.(*database.UserSkillProgress), nil
}

func (a *ProgressorApp) UpdateSkill(id int64, name string, description string) (*database.UserSkill, error) {
	res, err := a.withSession(func(s *AppSession) (interface{}, error) {
		return s.skillService.UpdateSkill(context.Background(), id, name, description)
	})
	if err != nil {
		return nil, err
	}
	return res.(*database.UserSkill), nil
}

// ProjectService delegates
func (a *ProgressorApp) AddProjectSkill(projectID int64, skillID int64) error {
	_, err := a.withSession(func(s *AppSession) (interface{}, error) {
		return nil, s.projectService.AddProjectSkill(context.Background(), projectID, skillID)
	})
	return err
}

func (a *ProgressorApp) GetProjects() ([]database.Project, error) {
	res, err := a.withSession(func(s *AppSession) (interface{}, error) {
		return s.projectService.GetProjects()
	})
	if err != nil {
		return nil, err
	}
	return res.([]database.Project), nil
}

func (a *ProgressorApp) GetSkillsForProject(projectID int64) ([]database.UserSkill, error) {
	res, err := a.withSession(func(s *AppSession) (interface{}, error) {
		return s.projectService.GetSkillsForProject(context.Background(), projectID)
	})
	if err != nil {
		return nil, err
	}
	return res.([]database.UserSkill), nil
}

func (a *ProgressorApp) IsValidProject(projectID uint) (bool, error) {
	res, err := a.withSession(func(s *AppSession) (interface{}, error) {
		return s.projectService.IsValidProject(projectID)
	})
	if err != nil {
		return false, err
	}
	return res.(bool), nil
}

func (a *ProgressorApp) RemoveProjectSkill(projectID int64, skillID int64) error {
	_, err := a.withSession(func(s *AppSession) (interface{}, error) {
		return nil, s.projectService.RemoveProjectSkill(context.Background(), projectID, skillID)
	})
	return err
}
