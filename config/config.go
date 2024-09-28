package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	DB_HOST        string `mapstructure:"DB_HOST"`
	DB_PORT        string `mapstructure:"DB_PORT"`
	DB_USER        string `mapstructure:"DB_USER"`
	DB_PASSWORD    string `mapstructure:"DB_PASSWORD"`
	DB_NAME        string `mapstructure:"DB_NAME"`
	SERVER_PORT    string `mapstructure:"SERVER_PORT"`
	ACCESS_SECRET  string `mapstructure:"ACCESS_SECRET"`
	REFRESH_SECRET string `mapstructure:"REFRESH_SECRET"`
}

func LoadConfig() (config Config) {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}
	return
}
