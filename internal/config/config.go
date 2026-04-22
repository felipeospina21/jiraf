// Package config handles loading the jiraf TOML configuration file.
package config

import (
	"errors"
	"flag"
	"fmt"

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

// ThemeOverrides holds optional color overrides for the theme.
// Each field is a hex color string (e.g. "#FF5733"). Nil means use default.
type ThemeOverrides struct {
	Primary         *string `mapstructure:"primary"`
	PrimaryBright   *string `mapstructure:"primary_bright"`
	PrimaryFg       *string `mapstructure:"primary_fg"`
	PrimaryDim      *string `mapstructure:"primary_dim"`
	Info            *string `mapstructure:"info"`
	InfoBright      *string `mapstructure:"info_bright"`
	Success         *string `mapstructure:"success"`
	SuccessBright   *string `mapstructure:"success_bright"`
	Danger          *string `mapstructure:"danger"`
	DangerBright    *string `mapstructure:"danger_bright"`
	Warning         *string `mapstructure:"warning"`
	WarningBright   *string `mapstructure:"warning_bright"`
	Caution         *string `mapstructure:"caution"`
	Text            *string `mapstructure:"text"`
	TextInverse     *string `mapstructure:"text_inverse"`
	TextDimmed      *string `mapstructure:"text_dimmed"`
	Muted           *string `mapstructure:"muted"`
	Dim             *string `mapstructure:"dim"`
	Border          *string `mapstructure:"border"`
	ModalBorder     *string `mapstructure:"modal_border"`
	SurfaceDim      *string `mapstructure:"surface_dim"`
	SelectionBorder *string `mapstructure:"selection_border"`
	StatusText      *string `mapstructure:"status_text"`
	StatusNormal    *string `mapstructure:"status_normal"`
	StatusLoading   *string `mapstructure:"status_loading"`
	StatusError     *string `mapstructure:"status_error"`
	StatusDev       *string `mapstructure:"status_dev"`
	StatusAccent1   *string `mapstructure:"status_accent1"`
	StatusAccent2   *string `mapstructure:"status_accent2"`
}

// Config holds the application configuration.
type Config struct {
	BaseURL  string         `mapstructure:"base_url"`
	APIToken string
	Filters  Filter         `mapstructure:"filters"`
	Theme    ThemeOverrides `mapstructure:"theme"`
	DevMode  bool
}

var (
	GlobalConfig Config
	cmdName      = "jiraf"
)

// Load reads the config file and environment variables.
func Load(config *Config) error {
	config.DevMode = isDevMode()

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
			return errors.New("JIRAF_TOKEN not set")
		}
	}

	config.APIToken = token
	return nil
}

var devFlag = flag.Bool("dev", false, "use mocked data instead of calling Jira API")

func isDevMode() bool {
	if !flag.Parsed() {
		flag.Parse()
	}
	return *devFlag
}

var mockBoards = []Board{
	{Name: "Project Board", ID: "1", Key: "PROJ"},
	{Name: "Platform Sprint", ID: "42", Key: "PLAT"},
	{Name: "Infrastructure", ID: "15", Key: "INFRA"},
}
