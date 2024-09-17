package main

// A simple program demonstrating the text input component from the Bubbles
// component library.

import (
	"flag"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/m1kkY8/gochat/src/connection"
	"github.com/m1kkY8/gochat/src/teamodel"
	"github.com/m1kkY8/gochat/src/util"
)

func main() {
	u := flag.String("u", "anon", "Username")
	c := flag.String("c", "202", "Color for your username, use -c=help for all colors")
	flag.Parse()

	if *c == "help" {
		util.Colors()
		return
	}

	url := "ws://139.162.132.8:1337/ws"

	conn, err := connection.ConnectToServer(url)
	if err != nil {
		log.Printf("Error creating websocket connection: %v", err)
		return
	}

	time.Sleep(time.Second * 1)

	teaModel := teamodel.New(*c, *u, conn)

	p := tea.NewProgram(teaModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Println(err)
	}
}
