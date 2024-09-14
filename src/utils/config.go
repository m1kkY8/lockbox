package utils

import (
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Username string
	Server   []string
}

func CheckIfConfigExists() bool {
	usr, err := user.Current()
	if err != nil {
		log.Println(err)
	}

	if _, err := os.Stat(filepath.Join(usr.HomeDir, ".config", "gochat", "config.yaml")); os.IsNotExist(err) {
		return false
	}
	return true
}

func createConfigFolder() {
	usr, err := user.Current()
	if err != nil {
		log.Println(err)
	}

	configPath := filepath.Join(usr.HomeDir, ".config", "gochat")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		os.MkdirAll(configPath, 0700)
	}
}

func LoadConfig() Config {
	usr, err := user.Current()
	if err != nil {
		log.Println(err)
	}
	if _, err := os.Stat(filepath.Join(usr.HomeDir, ".config", "gochat", "config.yaml")); os.IsNotExist(err) {
		log.Println("Config file does not exist")
		createConfigFolder()
	}

	configPath := filepath.Join(usr.HomeDir, ".config", "gochat")

	viper.SetConfigName("config")
	viper.AddConfigPath(configPath)
	viper.SetConfigType("yaml")

	var config Config
	if err := viper.ReadInConfig(); err != nil {
		log.Println("Error reading config file")
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Println("Error unmarshalling config")
	}

	return Config{
		Username: config.Username,
		Server:   config.Server,
	}
}
