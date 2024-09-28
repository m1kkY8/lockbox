package config

import (
	"flag"
	"fmt"
	"net/url"

	"github.com/m1kkY8/gochat/src/styles"
	"github.com/m1kkY8/gochat/src/util"
)

type Config struct {
	Username string
	Color    string
	Host     string
}

func LoadConfig() *Config {
	u := flag.String("u", "anon", "Username")
	c := flag.String("c", "", "Color for your username, use -c=help for all colors")
	h := flag.String("ip", "", "IP address of server")
	flag.Parse()

	var config Config

	config.Color = *c
	config.Username = *u
	config.Host = *h

	if config.Color == "" {
		*c = styles.GenerateRandomANSIColor()
	}

	return &config
}

func ValidateConfig(c Config) error {
	if c.Color == "help" {
		util.Colors()
		return fmt.Errorf("Color list")
	}

	if c.Host == "" {
		return fmt.Errorf("please provide IP and port of server and try again, use -h for help")
	}
	return nil
}

func GetUrl(c Config) url.URL {
	scheme := "wss"
	host := c.Host
	path := "/ws"

	// Construct URL
	u := url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   path,
	}
	return u
}