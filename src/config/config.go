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
	scheme := "wss"
	host := c.Host
	path := "/chat"

	// Construct URL
	u := url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   path,
	}
	return u
}

func ValidateConfig(c Config) error {
	if isEmpty(c) {
		return fmt.Errorf("Empty config")
	}
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

func isEmpty(c Config) bool {
	if c.Host == "" && c.Color == "" && c.Username == "" && c.Secure == "" {
		return true
	}
	return false
}
