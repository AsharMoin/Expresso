package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	openaikey string
	shell     string
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
		openaikey: viper.GetString("openai_api_key"),
		shell:     "powershell",
	}, nil
}

func (c *Config) GetKey() string {
	return c.openaikey
}

func (c *Config) GetShell() string {
	return c.shell
}
