package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/m1kkY8/lockbox/src/config"
	"github.com/m1kkY8/lockbox/src/connection"
	"github.com/m1kkY8/lockbox/src/encryption"
	"github.com/m1kkY8/lockbox/src/login"
	"github.com/m1kkY8/lockbox/src/teamodel"
)

func main() {
	start()
}

// Start the program
func start() {
	// Pravi novi login
	var conf config.Config
	loginPrompt := login.New(&conf)

	p := tea.NewProgram(loginPrompt,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		log.Println(err)
	}
	// Ovde se zavrsava forma al ne znam dal ce da ucita config
	// Novi rsa kljuc
	KeyPair, err := encryption.CreateRsaKey()
	if err != nil {
		log.Println(err)
		return
	}

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

	if err := connection.SendHandshake(conn, conf, KeyPair.PublicKey); err != nil {
		log.Println("error sending handshake")
		return
	}

	teaModel := teamodel.New(conf, conn, KeyPair)

	p = tea.NewProgram(teaModel,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := p.Run(); err != nil {
		log.Println(err)
	}
}
