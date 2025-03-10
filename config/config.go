package config

import (
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
		return nil, err
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

func (c *Config) UpdateConfig(key string) error {
	viper.Set("openai_api_key", key)

	// Try to write to existing config
	err := viper.WriteConfig()
	if err != nil {
		// If config file doesn't exist, create it
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return viper.SafeWriteConfig()
		}
		return err
	}

	return nil
}
