package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Port        string
	Db_user     string
	Db_password string
	Db_name     string
	Db_host     string
	Db_port     string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	viper.AutomaticEnv() // Читаем переменные окружения

	return &Config{
		Port:        viper.GetString("PORT"),
		Db_user:     viper.GetString("DB_USER"),
		Db_password: viper.GetString("DB_PASSWORD"),
		Db_name:     viper.GetString("DB_NAME"),
		Db_host:     viper.GetString("DB_HOST"),
		Db_port:     viper.GetString("DB_PORT"),
	}, nil
}
