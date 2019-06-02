package main

import (
	"log"

	ctui "github.com/fxor/ws-tui/tui-client"
)

func main() {
	tui, err := ctui.New()
	if err != nil {
		log.Fatal(err)
	}
	defer tui.Close()

	err = tui.MainLoop()
	if err != nil {
		log.Fatal(err)
	}
}
