package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL string
	Port        string
}

func LoadConfig() *Config {
	viper.AutomaticEnv() // Читаем переменные окружения

	return &Config{
		DatabaseURL: viper.GetString("DATABASE_URL"),
		Port:        viper.GetString("PORT"),
	}
}
