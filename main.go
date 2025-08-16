package main

import (
	"embed"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/services/notifications"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env file (this may be normal if env vars are set externally):", err)
	}

	progressorApp := NewProgressorApp()

	log.Println("Starting Progressor Todo App...")

	wailsApp := application.New(application.Options{
		Name:        "progressor-todo-app",
		Description: "Progressor Todo App",
		Services: []application.Service{
			application.NewService(progressorApp),
			application.NewService(notifications.New()),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		OnShutdown: progressorApp.Shutdown,
	})

	progressorApp.Startup(wailsApp)

	wailsApp.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "Progressor",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
	})

	go func() {
		for {
			now := time.Now().Format(time.RFC1123)
			wailsApp.Event.Emit("time", now)
			time.Sleep(time.Second)
		}
	}()

	if err := wailsApp.Run(); err != nil {
		log.Fatal(err)
	}
}