package main

import (
	"flag"
	"fmt"
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
	i := flag.String("ip", "", "IP address of server")
	p := flag.String("p", "", "Port")
	flag.Parse()

	if *c == "help" {
		util.Colors()
		return
	}

	if *c == "" {
		*c = styles.GenerateRandomANSIColor()
	}

	if *i == "" || *p == "" {
		fmt.Println("Please provide IP and port of server and try again, use -h for help")
		return
	}

	url := fmt.Sprintf("ws://%s:%s/ws", *i, *p)

	conn, err := connection.ConnectToServer(url)
	if err != nil {
		log.Printf("Error creating websocket connection: %v", err)
		return
	}

	var handshake connection.Handshake

	handshake.Username = *u
	handshake.ClientId = connection.GenerateUUID()
	handshake.PublicKey = "kljuc"
	handshake.Color = *c

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
