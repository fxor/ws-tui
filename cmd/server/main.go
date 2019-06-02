package main

import (
	"log"

	wssv "github.com/fxor/ws-tui/server"
)

func main() {
	sv := wssv.NewServer(":8080")
	log.Fatal(sv.Start())
}
