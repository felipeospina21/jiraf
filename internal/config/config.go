// Package config handles loading the mrjira TOML configuration file.
package config

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Board represents a Jira board entry in the config file.
type Board struct {
	Name string `mapstructure:"name"`
	ID   string `mapstructure:"id"`
	Key  string `mapstructure:"key"`
}

// Filter holds the board filter list from the config file.
type Filter struct {
	Boards []Board `mapstructure:"boards"`
}

// Config holds the application configuration.
type Config struct {
	BaseURL  string `mapstructure:"base_url"`
	APIToken string
	Filters  Filter `mapstructure:"filters"`
	DevMode  bool
}

var (
	GlobalConfig Config
	cmdName      = "mrjira"
)

// Load reads the config file and environment variables.
func Load(config *Config) error {
	config.DevMode = isDevMode()
	_ = godotenv.Load()                                          // cwd
	_ = godotenv.Load(os.ExpandEnv("$HOME/.config/mrjira/.env")) // config dir

	viper.SetConfigName(cmdName)
	viper.SetConfigType("toml")
	viper.AddConfigPath(fmt.Sprintf("$HOME/.config/%s/", cmdName))
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		if config.DevMode {
			config.BaseURL = "https://acme.atlassian.net"
			config.Filters.Boards = mockBoards
			return nil
		}
		return fmt.Errorf("config file: %w", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	if config.BaseURL == "" {
		return errors.New("base_url is required")
	}

	if config.DevMode {
		config.Filters.Boards = mockBoards
		return nil
	}

	return loadEnvVars(config)
}

func loadEnvVars(config *Config) error {
	viper.SetEnvPrefix(cmdName)

	_ = viper.BindEnv("token")

	token := viper.GetString("token")

	if !config.DevMode {
		if token == "" {
			return errors.New("MRJIRA_TOKEN not set")
		}
	}

	config.APIToken = token
	return nil
}

func isDevMode() bool {
	// Check env var first (for hub/embedded usage), then flag
	if os.Getenv("MRJIRA_DEV") == "1" || os.Getenv("MRJIRA_DEV") == "true" {
		return true
	}
	f := flag.Lookup("dev")
	if f != nil {
		return f.Value.String() == "true"
	}
	dev := flag.Bool("dev", false, "use mocked data instead of calling Jira API")
	if !flag.Parsed() {
		flag.Parse()
	}
	return *dev
}

var mockBoards = []Board{
	{Name: "UCP Board", ID: "9159", Key: "UCP"},
	{Name: "Platform Sprint", ID: "42", Key: "PLAT"},
	{Name: "Infrastructure", ID: "15", Key: "INFRA"},
}
