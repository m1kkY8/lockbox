package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/m1kkY8/lockbox/src/config"
	"github.com/m1kkY8/lockbox/src/connection"
	"github.com/m1kkY8/lockbox/src/encryption"
	"github.com/m1kkY8/lockbox/src/login"
	"github.com/m1kkY8/lockbox/src/teamodel"
)

func main() {
	// Set up context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle OS signals for graceful shutdown
	go handleSignals(cancel)

	if err := run(ctx); err != nil {
		log.Printf("Application error: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	// Create login configuration
	var conf config.Config
	loginPrompt := login.New(&conf)

	// Run login form
	p := tea.NewProgram(loginPrompt,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("login form error: %w", err)
	}

	// Apply defaults and validate configuration
	config.ApplyDefaults(&conf)
	if err := config.ValidateConfig(conf); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	// Generate RSA key pair
	keyPair, err := encryption.CreateRsaKey()
	if err != nil {
		return fmt.Errorf("failed to create RSA keys: %w", err)
	}

	// Get server URL
	url := config.GetUrl(conf)

	// Connect to server
	conn, err := connection.ConnectToServer(url.String())
	if err != nil {
		return fmt.Errorf("%s: %w", config.ErrConnectionFailed, err)
	}
	defer conn.Close()

	// Send initial handshake
	if err := connection.SendHandshake(conn, conf, keyPair.PublicKey); err != nil {
		return fmt.Errorf("%s: %w", config.ErrHandshakeFailed, err)
	}

	// Create and run main TUI
	teaModel := teamodel.New(conf, conn, keyPair)
	mainProgram := tea.NewProgram(teaModel,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := mainProgram.Run(); err != nil {
		return fmt.Errorf("main application error: %w", err)
	}

	return nil
}

func handleSignals(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	cancel()
}
