package internal

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.design/x/hotkey"
	"golang.design/x/hotkey/mainthread"
)

type Shortcut struct {
	ctx context.Context
}

func NewShortcut() *Shortcut {
	return &Shortcut{}
}

func (s *Shortcut) Startup(ctx context.Context) {
	s.ctx = ctx
	mainthread.Init(s.BindShortcuts)
}

func (s *Shortcut) BindShortcuts() {
	registerHoykeys(s)
}

func registerHoykeys(s *Shortcut) {

	// the actual shortcut keybind - Ctrl + Shift + P
	// for more info - refer to the golang.design/x/hotkey documentation
	hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyP)
	err := hk.Register()
	if err != nil {
		return
	}

	<-hk.Keyup()
	// do anything you want on Key up event
	fmt.Printf("hotkey: %v is up\n", hk)

	runtime.Hide(s.ctx)
	runtime.Show(s.ctx)

	hk.Unregister()

	// reattach listener
	registerHoykeys(s)
}
