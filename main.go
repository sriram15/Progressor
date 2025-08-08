package main

import (
	"embed"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/sriram15/progressor-todo-app/internal/connection"
	"github.com/sriram15/progressor-todo-app/internal/events"
	"github.com/sriram15/progressor-todo-app/internal/service"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/services/notifications"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env file (this may be normal if env vars are set externally):", err)
	}

	err := connection.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	dbManager := connection.GetManager()
	eventBus := events.NewEventBus()

	projectService := service.NewProjectService(dbManager)
	taskCompletionService := service.NewTaskCompletionService(dbManager)
	cardService := service.NewCardService(projectService, taskCompletionService, dbManager, eventBus)
	progressService := service.NewProgressService(dbManager)
	settingsService := service.NewSettingService(dbManager)
	skillService := service.NewSkillService(dbManager, eventBus, projectService)
	skillService.RegisterEventHandlers()
	// shortcuts := internal.NewShortcut()

	app := application.New(application.Options{
		Name:        "progressor-todo-app",
		Description: "Progressor Todo App",
		Services: []application.Service{
			application.NewService(cardService),
			application.NewService(progressService),
			application.NewService(settingsService),
			application.NewService(projectService),
			application.NewService(skillService),
			application.NewService(notifications.New()),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	focusTimerService := service.NewFocusTimerService(cardService, settingsService, eventBus, app)
	focusTimerService.RegisterEventHandlers()

	app.OnShutdown(func() {
		focusTimerService.Shutdown()
		err := cardService.Cleanup()
		fmt.Println("Shutdown cleanup done")
		if err != nil {
			log.Printf("Error during card service cleanup on shutdown: %v", err)
		}
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "Progressor",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
		// Shortcuts:        shortcuts,
	})

	go func() {
		for {
			now := time.Now().Format(time.RFC1123)
			app.Event.Emit("time", now)
			time.Sleep(time.Second)
		}
	}()

	err = app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
