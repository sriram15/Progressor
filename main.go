package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"runtime"

	"github.com/sriram15/progressor-todo-app/internal"
	"github.com/sriram15/progressor-todo-app/internal/database"
	"github.com/sriram15/progressor-todo-app/internal/service"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed all:internal/database/migrations
var migrations embed.FS

func main() {

	var startupCtx context.Context

	entries, err := migrations.ReadDir("internal/database/migrations")
	if err != nil {
		log.Fatalf("Failed to read migrations directory: %v", err)
	}

	for _, entry := range entries {
		fmt.Println(entry.Name())
	}

	db, err := internal.OpenDB()
	if err != nil {
		log.Fatalf("Failed to auto-migrate: %v", err)
		panic("There is an issue with auto-migrate")
	}
	// TODO: Find a place to close the db

	queries := database.New(db)
	projectService := service.NewProjectService()
	taskCompletionService := service.NewTaskCompletionService(db, queries)
	cardService := service.NewCardService(db, queries, projectService, taskCompletionService)
	progressService := service.NewProgressService(queries, taskCompletionService)
	settingsService := service.NewSettingService()
	// shortcuts := internal.NewShortcut()

	// Create menu
	appMenu := menu.NewMenu()
	fileMenu := appMenu.AddSubmenu("File")

	fileMenu.AddText("Quit", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
		wailsRuntime.Quit(startupCtx)
	})
	fileMenu.AddText("Command Prompt", keys.CmdOrCtrl("k"), func(_ *menu.CallbackData) {
		wailsRuntime.EventsEmit(startupCtx, "globalMenu:CommandPrompt")
	})

	if runtime.GOOS == "darwin" {
		appMenu.Append(menu.EditMenu()) // on macos platform, we should append EditMenu to enable Cmd+C,Cmd+V,Cmd+Z... shortcut
	}

	// Create application with options
	wailsApp := &options.App{
		Title:      "Progressor",
		Fullscreen: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		// StartHidden:      true,
		OnStartup: func(ctx context.Context) {
			startupCtx = ctx
			// shortcuts.Startup(ctx)
		},
		OnBeforeClose: func(ctx context.Context) (prevent bool) {
			err := cardService.Cleanup()
			if err != nil {
				log.Printf(err.Error())
			}
			dialog, err := wailsRuntime.MessageDialog(ctx, wailsRuntime.MessageDialogOptions{
				Type:    wailsRuntime.QuestionDialog,
				Title:   "Quit?",
				Message: "Are you sure you want to quit?",
			})
			if err != nil {
				return false
			}
			return dialog != "Yes"
		},
		Menu:             appMenu,
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		Bind: []interface{}{
			cardService,
			progressService,
			settingsService,
		},
	}

	err = wails.Run(wailsApp)

	if err != nil {
		println("Error:", err.Error())
	}
}
