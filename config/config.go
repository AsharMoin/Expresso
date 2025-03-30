package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AsharMoin/Expresso/sysinfo"
	"github.com/spf13/viper"
)

const APPLICATION_NAME = "Expresso"

type Config struct {
	openaikey string
	user      *sysinfo.User
}

// GetConfigDirectory returns the directory where the config file is stored
func GetConfigDirectory() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// fallback to current directory if home directory can't be found
		return "./.config"
	}
	return filepath.Join(homeDir, ".config", strings.ToLower(APPLICATION_NAME))
}

// GetConfigFilePath returns the full path to the config file
func GetConfigFilePath() string {
	return filepath.Join(GetConfigDirectory(), "config.yaml")
}

// InitConfig initializes and reads the configuration
func InitConfig() (*Config, error) {
	var config Config

	// always initialize the user regardless of config status
	config.user = sysinfo.NewUser()

	// set viper configuration
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(GetConfigDirectory())

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return &config, nil
		}
		fmt.Println("Error reading config:", err)
		return &config, err
	}

	config.openaikey = viper.GetString("openai_api_key")

	return &config, nil
}

func (c *Config) GetKey() string {
	return c.openaikey
}

func (c *Config) GetUser() *sysinfo.User {
	return c.user
}

// UpdateConfig updates the OpenAI API key in the configuration
func (c *Config) UpdateConfig(key string) error {
	c.openaikey = key
	viper.Set("openai_api_key", key)

	configDir := GetConfigDirectory()
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Println("Error creating config directory:", err)
		return err
	}

	configFile := GetConfigFilePath()

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if err := viper.SafeWriteConfigAs(configFile); err != nil {
			fmt.Println("Error writing new config file:", err)
			return err
		}
		return nil
	}

	fmt.Println("Updating existing config file at:", configFile)
	if err := viper.WriteConfig(); err != nil {
		fmt.Println("Error updating config file:", err)
		return err
	}

	return nil
}
