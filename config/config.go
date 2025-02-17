package config

import (
	"fmt"

	"github.com/AsharMoin/Expresso/sysinfo"
	"github.com/spf13/viper"
)

type Config struct {
	openaikey string
	user      *sysinfo.User
}

func InitConfig() (*Config, error) {
	var config Config

	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.ReadInConfig()

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file: %v", err)
	}

	config.openaikey = viper.GetString("openai_api_key")
	config.user = sysinfo.NewUser()

	return &config, nil
}

func (c *Config) GetKey() string {
	return c.openaikey
}

func (c *Config) GetUser() *sysinfo.User {
	return c.user
}
