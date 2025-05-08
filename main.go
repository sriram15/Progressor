package main

import (
	"context"
	"embed"
	"log"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/sriram15/progressor-todo-app/internal/connection"
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

func main() {
	// Load .env file at the very beginning
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env file (this may be normal if env vars are set externally):", err)
	}

	var startupCtx context.Context

	db, err := connection.OpenDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	// defer db.Close()

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
				log.Println(err.Error())
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
