package main

import (
	"flag"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/m1kkY8/gochat/src/connection"
	"github.com/m1kkY8/gochat/src/styles"
	"github.com/m1kkY8/gochat/src/teamodel"
	"github.com/m1kkY8/gochat/src/util"
)

func main() {
	u := flag.String("u", "anon", "Username")
	c := flag.String("c", "", "Color for your username, use -c=help for all colors")
	flag.Parse()

	if *c == "help" {
		util.Colors()
		return
	}

	if *c == "" {
		*c = styles.GenerateRandomANSIColor()
	}

	url := "ws://139.162.132.8:1337/ws"

	conn, err := connection.ConnectToServer(url)
	if err != nil {
		log.Printf("Error creating websocket connection: %v", err)
		return
	}

	var handshake connection.Handshake

	handshake.Username = *u
	handshake.ClientId = connection.GenerateUUID()
	handshake.PublicKey = "kljuc"

	err = connection.SendHandshake(conn, handshake)
	if err != nil {
		log.Println("error sending handshake")
	}

	teaModel := teamodel.New(*c, *u, conn)
	start(teaModel)
}

func start(teaModel *teamodel.Model) {
	p := tea.NewProgram(teaModel,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		log.Println(err)
	}
}
