package main

import (
	"embed"
	"log"
	"runtime"
	"time"

	"github.com/joho/godotenv"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/icons"
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

	_ = wailsApp.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "Progressor",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		Windows: application.WindowsWindow{
			HiddenOnTaskbar: true,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
	})

	// System Tray Settings
	systemTray := wailsApp.SystemTray.New()

    // Support for template icons on macOS
    if runtime.GOOS == "darwin" {
        systemTray.SetTemplateIcon(icons.SystrayMacTemplate)
    } else {
        // Support for light/dark mode icons
        systemTray.SetDarkModeIcon(icons.SystrayDark)
        systemTray.SetIcon(icons.SystrayLight)
    }

    // Support for menu
    myMenu := wailsApp.NewMenu()
	myMenu.Add("Quit").OnClick(func(ctx *application.Context) {
		wailsApp.Quit()
	})
    systemTray.SetMenu(myMenu)

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