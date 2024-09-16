package main

// A simple program demonstrating the text input component from the Bubbles
// component library.

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gorilla/websocket"
	"github.com/m1kkY8/gochat/src/teamodel"
)

type message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
	To       string `json:"to"`
}

var (
	color    string
	username string
)

func CreateWebSocketConnection(url string) (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func main() {
	u := flag.String("u", "anon", "Username")
	c := flag.String("c", "202", "Color for your username")
	flag.Parse()

	url := "ws://139.162.132.8:1337/ws"
	conn, err := CreateWebSocketConnection(url)
	if err != nil {
		log.Printf("Error creating websocket connection: %v", err)
		return
	}

	defer conn.Close()

	color = *c
	username = *u

	time.Sleep(time.Second * 1)

	teaModel := teamodel.New(color)
	teaModel.Conn = conn
	teaModel.Username = username
	teaModel.UserColor = color

	// go teaModel.HandleIncomingMessage()

	fmt.Println("routine1\n", runtime.NumGoroutine())

	p := tea.NewProgram(teaModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Println(err)
	}

	fmt.Println("routine2\n", runtime.NumGoroutine())
	teaModel.CloseConnection()
}
