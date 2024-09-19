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
	ServerIp string
	Port     string
}

func LoadConfig() Config {
	u := flag.String("u", "anon", "Username")
	c := flag.String("c", "", "Color for your username, use -c=help for all colors")
	i := flag.String("ip", "", "IP address of server")
	p := flag.String("p", "", "Port")
	flag.Parse()

	if *c == "help" {
		util.Colors()
		return Config{}
	}

	if *c == "" {
		*c = styles.GenerateRandomANSIColor()
	}

	return Config{
		Username: *u,
		Color:    *c,
		ServerIp: *i,
		Port:     *p,
	}
}

func ValidateConfig(c Config) error {
	if c.ServerIp == "" || c.Port == "" {
		return fmt.Errorf("please provide IP and port of server and try again, use -h for help")
	}
	return nil
}

func GetUrl(c Config) string {
	scheme := "ws"
	host := c.ServerIp
	port := c.Port
	path := "/ws"

	// Construct URL
	u := url.URL{
		Scheme: scheme,
		Host:   fmt.Sprintf("%s:%s", host, port),
		Path:   path,
	}
	return u.String()
}
