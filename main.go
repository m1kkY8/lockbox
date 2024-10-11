package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/m1kkY8/lockbox/src/config"
	"github.com/m1kkY8/lockbox/src/connection"
	"github.com/m1kkY8/lockbox/src/encryption"
	"github.com/m1kkY8/lockbox/src/login"
	"github.com/m1kkY8/lockbox/src/teamodel"
)

func main() {
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
	if err := config.ValidateConfig(conf); err != nil {
		return
	}
	// Novi rsa kljuc
	KeyPair, err := encryption.CreateRsaKey()
	if err != nil {
		log.Println(err)
		return
	}

	url := config.GetUrl(conf)

	conn, err := connection.ConnectToServer(url.String())
	if err != nil {
		fmt.Println("Failed establishing connection")
		return
	}

	if err := connection.SendHandshake(conn, conf, KeyPair.PublicKey); err != nil {
		fmt.Println("Failed sending initial handshake")
		return
	}

	// pokreni glavnu formu
	teaModel := teamodel.New(conf, conn, KeyPair)
	p = tea.NewProgram(teaModel,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := p.Run(); err != nil {
		log.Println(err)
	}
}
