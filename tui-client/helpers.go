package tui

import (
	"fmt"
	"log"

	wssv "github.com/fxor/ws-tui/server"
	"github.com/jroimartin/gocui"
)

var user string

var mainViews = []string{
	"input",
	"output",
	"side",
}
var mainViewsIndex int

var loginViews = []string{
	"login-server",
	"login-name",
	"login-pw",
}
var loginViewsIndex int

func nextView(g *gocui.Gui, v *gocui.View) error {
	var views []string
	var viewIndex *int
	if user == "" {
		views = loginViews
		viewIndex = &loginViewsIndex
	} else {
		views = mainViews
		viewIndex = &mainViewsIndex
	}
	*viewIndex = (*viewIndex + 1) % len(views)
	nextView := views[*viewIndex]
	_, err := g.SetCurrentView(nextView)
	return err
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func getLine(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	maxX, maxY := g.Size()
	if v, err := g.SetView("msg", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, l)
		if _, err := g.SetCurrentView("msg"); err != nil {
			return err
		}

	}
	return nil
}

func delMsg(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("msg"); err != nil {
		return err
	}
	if _, err := g.SetCurrentView("side"); err != nil {
		return err
	}
	return nil
}

func (tui *TUI) sendMsg(g *gocui.Gui, v *gocui.View) error {
	msg := v.Buffer()
	if msg != "" {
		err := tui.wsc.Send(msg)
		if err != nil {
			fmt.Println("Coduln't send msg")
		}
		v.Clear()
		v.SetCursor(0, 0)
		outputView, err := g.View("output")
		if err != nil {
			return err
		}
		outputView.Autoscroll = true
	}
	return nil
}

func pushToOutput(g *gocui.Gui, msg string) error {
	outputView, err := g.View("output")
	if err != nil {
		return err
	}
	_, err = outputView.Write([]byte(msg))
	if err != nil {
		return err
	}
	g.Update(func(*gocui.Gui) error { return nil })
	return nil
}

func autoscroll(g *gocui.Gui, v *gocui.View) error {
	v.Autoscroll = true
	return nil
}

func scrollView(v *gocui.View, dy int) error {
	if v != nil {
		v.Autoscroll = false
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy+dy); err != nil {
			return err
		}
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func listener(conn *wssv.Client, gui *gocui.Gui) {
	for {
		if conn.IsClosed() {
			return
		}
		msg, err := conn.Receive()
		if err != nil {
			log.Fatal("Listener:", err)
		}
		err = pushToOutput(gui, msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}
