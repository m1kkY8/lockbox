package config

import (
	"fmt"
	"net/url"

	"github.com/m1kkY8/lockbox/src/styles"
	"github.com/m1kkY8/lockbox/src/util"
)

type Config struct {
	Username string
	Color    string
	Host     string
	Secure   string
}

func GetUrl(c Config) url.URL {
	var scheme string
	if c.Secure == "no" {
		scheme = "ws"
	} else {
		scheme = "wss"
	}

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

// OVO JE MEGA FUCKED OVO TREBA POPRAVITI
func LoadConfig(c Config) *Config {
	var config Config

	config.Host = c.Host
	config.Username = c.Username
	config.Color = c.Color
	config.Secure = c.Secure

	return &config
}

func ValidateConfig(c Config) error {
	if c.Color == "" {
		c.Color = styles.GenerateRandomANSIColor()
		return nil
	}
	if c.Color == "help" {
		util.Colors()
		return fmt.Errorf("Color list")
	}

	if c.Host == "" {
		return fmt.Errorf("please provide IP and port of server and try again, use -h for help")
	}
	return nil
}
