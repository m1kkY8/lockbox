package config

import (
	"errors"
	"net/url"
	"strings"

	"github.com/m1kkY8/lockbox/src/styles"
	"github.com/m1kkY8/lockbox/src/util"
)

// Predefined errors for better error handling
var (
	ErrConfigEmpty     = errors.New(ErrEmptyConfig)
	ErrHostEmpty       = errors.New(ErrEmptyHost)
	ErrColorListNeeded = errors.New(ErrColorList)
)

// Config holds the application configuration
type Config struct {
	Username string `yaml:"username" validate:"required,min=1,max=32"`
	Color    string `yaml:"color" validate:"color"`
	Host     string `yaml:"host" validate:"required"`
	Secure   string `yaml:"secure"`
}

// GetUrl constructs a WebSocket URL from the configuration
func GetUrl(c Config) url.URL {
	scheme := DefaultScheme
	if c.Secure != "" {
		scheme = SecureScheme
	}

	host := c.Host
	if !strings.Contains(host, ":") {
		host = host + ":" + DefaultPort
	}

	return url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   DefaultPath,
	}
}

// ValidateConfig validates the configuration and applies defaults
func ValidateConfig(c Config) error {
	if isEmpty(c) {
		return ErrConfigEmpty
	}

	if err := validateHost(c.Host); err != nil {
		return err
	}

	if err := validateColor(c.Color); err != nil {
		return err
	}

	return nil
}

// validateHost checks if the host configuration is valid
func validateHost(host string) error {
	if strings.TrimSpace(host) == "" {
		return ErrHostEmpty
	}
	return nil
}

// validateColor validates and handles color configuration
func validateColor(color string) error {
	if color == HelpColorCommand {
		util.Colors()
		return ErrColorListNeeded
	}
	return nil
}

// ApplyDefaults applies default values to empty configuration fields
func ApplyDefaults(c *Config) {
	if c.Color == "" {
		c.Color = styles.GenerateRandomANSIColor()
	}
}

// isEmpty checks if the configuration is completely empty
func isEmpty(c Config) bool {
	return strings.TrimSpace(c.Host) == "" &&
		strings.TrimSpace(c.Color) == "" &&
		strings.TrimSpace(c.Username) == "" &&
		strings.TrimSpace(c.Secure) == ""
}
