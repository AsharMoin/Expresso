package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	OpenAIKey string
	Shell     string
}

func InitConfig() (*Config, error) {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.ReadInConfig()

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file: %v", err)
	}

	return &Config{
		OpenAIKey: viper.GetString("openai_api_key"),
		Shell:     "powershell",
	}, nil
}

func (c *Config) GetKey() string {
	return c.OpenAIKey
}

func (c *Config) GetShell() string {
	return c.Shell
}
