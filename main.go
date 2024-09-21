package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/m1kkY8/gochat/src/config"
	"github.com/m1kkY8/gochat/src/connection"
	"github.com/m1kkY8/gochat/src/teamodel"
)

func main() {
	conf := config.LoadConfig()

	if err := config.ValidateConfig(conf); err != nil {
		log.Println(err)
		return
	}

	url := config.GetUrl(conf)

	conn, err := connection.ConnectToServer(url.String())
	if err != nil {
		log.Printf("Error creating websocket connection: %v", err)
		return
	}

	if err := connection.SendHandshake(conn, conf); err != nil {
		log.Println("error sending handshake")
	}

	teaModel := teamodel.New(conf.Color, conf.Username, conn)
	start(teaModel)
}

func start(teaModel *teamodel.Model) {
	p := tea.NewProgram(teaModel,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := p.Run(); err != nil {
		log.Println(err)
	}
}
