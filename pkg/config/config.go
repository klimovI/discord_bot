package config

import (
	"log"

	"github.com/spf13/viper"
)

type bot struct {
	Token string `mapstructure:"token"`
}

type config struct {
	Bot bot `mapstructure:"bot"`
}

var Config *config

func init() {
	viper.SetConfigName("local")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	if err := viper.Unmarshal(&Config); err != nil {
		log.Fatal(err)
	}
}
