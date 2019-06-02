package tui

/*
	See examples in: ~/git/src/github.com/jroimartin/gocui/_examples
*/

import (
	"fmt"
	"log"

	wssv "github.com/fxor/ws-tui/server"
	"github.com/jroimartin/gocui"
)

type TUI struct {
	gui      *gocui.Gui
	wsc      *wssv.Client
	viewData []ViewData
}

type ViewData struct {
	Name            string
	LayoutFunc      func(*gocui.Gui) error
	KeybindingsFunc func(*gocui.Gui) error
}

func (tui *TUI) initScreens() {
	InitLoginScren(tui)
	InitMainScren(tui)
}

func New() (*TUI, error) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return nil, err
	}
	g.Cursor = true
	g.Highlight = true
	g.SelFgColor = gocui.ColorBlue

	ctui := &TUI{
		gui:      g,
		viewData: []ViewData{},
	}

	ctui.initScreens()
	ctui.setLayouts()
	err = ctui.setKeybindings()
	if err != nil {
		log.Fatal(err)
	}

	return ctui, nil
}

func (tui *TUI) MainLoop() error {
	err := tui.gui.MainLoop()
	if err != nil && err != gocui.ErrQuit {
		return err
	}
	return nil
}

func (tui *TUI) Close() {
	if tui.gui != nil {
		tui.gui.Close()
	}
	if tui.wsc != nil {
		tui.wsc.Close()
	}
}

func (tui *TUI) setKeybindings() error {
	g := tui.gui
	// general keybindings
	err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit)
	if err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		return err
	}

	// msg view kb
	if err := g.SetKeybinding("msg", gocui.KeyEnter, gocui.ModNone, delMsg); err != nil {
		return err
	}

	// per view keybindings
	for _, vd := range tui.viewData {
		vd.KeybindingsFunc(g)
	}

	return nil
}

func (tui *TUI) setLayouts() {
	layoutMain := func(g *gocui.Gui) error {
		if user == "" {
			return nil
		}
		maxX, maxY := g.Size()
		currentView := g.CurrentView()
		if currentView != nil {
			if currentView.Name() == "input" {
				g.Cursor = true
			} else {
				g.Cursor = false
			}
		}
		v, err := g.SetView("side", 0, 0, 29, maxY-1)
		if err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Highlight = true
			v.SelBgColor = gocui.ColorGreen
			v.SelFgColor = gocui.ColorBlack
			fmt.Fprintln(v, "Item 1")
			fmt.Fprintln(v, "Item 2")
			fmt.Fprintln(v, "Item 3")
			fmt.Fprint(v, "\rWill be")
			fmt.Fprint(v, "deleted\rItem 4\nItem 5")
			v.Title = "Side"
		}
		v, err = g.SetView("output", 30, 0, maxX-1, maxY-4)
		if err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Wrap = true
			v.Autoscroll = true
			v.Title = "Item 1"
			go func() {
				listener(tui.wsc, g)
			}()
		}
		v, err = g.SetView("input", 30, maxY-3, maxX-1, maxY-1)
		if err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Editable = true
			v.Wrap = true
			if _, err := g.SetCurrentView("input"); err != nil {
				return err
			}
			v.Autoscroll = true
		}
		return nil
	}

	layout := func(*gocui.Gui) error { return nil }
	for _, vd := range tui.viewData {
		layout = appendFuncs(layout, vd.LayoutFunc)
	}
	layout = appendFuncs(layout, layoutMain)
	tui.gui.SetManagerFunc(layout)
}

func appendFuncs(f1, f2 func(g *gocui.Gui) error) func(g *gocui.Gui) error {
	return func(g *gocui.Gui) error {
		err := f1(g)
		if err != nil {
			return err
		}
		err = f2(g)
		if err != nil {
			return err
		}
		return nil
	}
}
