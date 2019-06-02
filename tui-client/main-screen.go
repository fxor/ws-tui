package tui

import (
	"github.com/jroimartin/gocui"
)

const ViewInput = "input"
const ViewOutput = "output"
const ViewSide = "side"

func InitMainScren(tui *TUI) {
	tui.viewData = append(tui.viewData, initInputView(tui))
	tui.viewData = append(tui.viewData, initOutputView(tui))
	tui.viewData = append(tui.viewData, initSideView(tui))
}

func initInputView(tui *TUI) ViewData {
	return ViewData{
		Name: ViewInput,
		KeybindingsFunc: func(g *gocui.Gui) error {
			err := g.SetKeybinding(ViewInput, gocui.KeyEnter, gocui.ModNone, tui.sendMsg)
			if err != nil {
				return err
			}
			return nil
		},
		LayoutFunc: func(g *gocui.Gui) error {
			if user == "" {
				return nil
			}

			return nil
		},
	}
}

func initOutputView(tui *TUI) ViewData {
	return ViewData{
		Name: ViewOutput,
		KeybindingsFunc: func(g *gocui.Gui) error {
			err := g.SetKeybinding(ViewOutput, gocui.KeyArrowDown, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
				scrollView(v, 1)
				return nil
			})
			if err != nil {
				return err
			}
			err = g.SetKeybinding(ViewOutput, gocui.KeyArrowUp, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
				scrollView(v, -1)
				return nil
			})
			if err != nil {
				return err
			}
			return nil
		},
		LayoutFunc: func(g *gocui.Gui) error {
			if user == "" {
				return nil
			}
			return nil
		},
	}
}
func initSideView(tui *TUI) ViewData {
	return ViewData{
		Name: ViewSide,
		KeybindingsFunc: func(g *gocui.Gui) error {
			if err := g.SetKeybinding(ViewSide, gocui.KeyEnter, gocui.ModNone, getLine); err != nil {
				return err
			}
			if err := g.SetKeybinding(ViewSide, gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
				return err
			}
			if err := g.SetKeybinding(ViewSide, gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
				return err
			}
			return nil
		},
		LayoutFunc: func(g *gocui.Gui) error {
			if user == "" {
				return nil
			}
			return nil
		},
	}
}
