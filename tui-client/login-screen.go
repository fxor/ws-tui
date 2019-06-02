package tui

import (
	"strings"

	wssv "github.com/fxor/ws-tui/server"
	"github.com/jroimartin/gocui"
)

const (
	ViewLoginServer = "login-server"
	ViewLoginName   = "login-name"
	ViewLoginPW     = "login-pw"
)

func InitLoginScren(tui *TUI) {
	tui.viewData = append(tui.viewData, initLoginServerView(tui))
	tui.viewData = append(tui.viewData, initLoginNameView(tui))
	tui.viewData = append(tui.viewData, initLoginPWView(tui))
}

func initLoginServerView(tui *TUI) ViewData {
	return ViewData{
		Name: ViewLoginName,
		KeybindingsFunc: func(g *gocui.Gui) error {
			err := g.SetKeybinding(ViewLoginServer, gocui.KeyEnter, gocui.ModNone, nextView)
			return err
		},
		LayoutFunc: func(g *gocui.Gui) error {
			if user != "" {
				return nil
			}

			maxX, maxY := g.Size()
			v, err := g.SetView(ViewLoginServer, maxX/2-10, maxY/2-4, maxX/2+10, maxY/2+2-4)
			if err != nil {
				if err != gocui.ErrUnknownView {
					return err
				}
				v.Editable = true
				v.Title = "Server"
				if _, err := g.SetCurrentView(ViewLoginServer); err != nil {
					return err
				}
				defaultServer := "localhost:8080"
				v.Write([]byte(defaultServer))
				v.SetCursor(len(defaultServer), 0)
			}
			g.Cursor = true
			return nil
		},
	}
}

func initLoginNameView(tui *TUI) ViewData {
	return ViewData{
		Name: ViewLoginName,
		KeybindingsFunc: func(g *gocui.Gui) error {
			err := g.SetKeybinding(ViewLoginName, gocui.KeyEnter, gocui.ModNone, nextView)
			return err
		},
		LayoutFunc: func(g *gocui.Gui) error {
			if user != "" {
				return nil
			}

			maxX, maxY := g.Size()
			v, err := g.SetView(ViewLoginName, maxX/2-10, maxY/2-1, maxX/2+10, maxY/2+2-1)
			if err != nil {
				if err != gocui.ErrUnknownView {
					return err
				}
				v.Editable = true
				v.Title = "Username"
			}
			g.Cursor = true
			return nil
		},
	}
}

func initLoginPWView(tui *TUI) ViewData {
	return ViewData{
		Name: ViewLoginPW,
		KeybindingsFunc: func(g *gocui.Gui) error {
			err := g.SetKeybinding(ViewLoginPW, gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
				nameView, err := g.View(ViewLoginName)
				if err != nil {
					return err
				}
				serverView, err := g.View(ViewLoginServer)
				if err != nil {
					return err
				}
				server := strings.Trim(serverView.Buffer(), "\n")
				username := strings.Trim(nameView.Buffer(), "\n")
				pw := strings.Trim(v.Buffer(), "\n")
				wsc, err := wssv.NewClient(server, username, pw)
				if err != nil {
					return err
				}
				err = wsc.Connect()
				if err != nil {
					return err
				}
				v.Clear()
				v.SetCursor(0, 0)
				tui.wsc = wsc
				user = username
				return nil
			})
			return err
		},
		LayoutFunc: func(g *gocui.Gui) error {
			maxX, maxY := g.Size()
			if user != "" {
				return nil
			}
			v, err := g.SetView(ViewLoginPW, maxX/2-10, maxY/2+3-1, maxX/2+10, maxY/2+2+3-1)
			if err != nil {
				if err != gocui.ErrUnknownView {
					return err
				}
				v.Editable = true
				v.Title = "Password"
				v.Mask = '*'
			}
			g.Cursor = true
			return nil
		},
	}
}
