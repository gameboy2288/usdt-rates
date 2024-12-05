package config

import (
	"flag"
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Port       string
	DbUser     string
	DbPassword string
	DbName     string
	DbHost     string
	DbPort     string
}

// LoadConfig загружает конфигурацию из флагов или переменных окружения
func LoadConfig() (*Config, error) {
	// Загружаем переменные из .env файла
	err := godotenv.Load()
	if err != nil {
		log.Printf("No .env file found or unable to load it: %v", err)
	}

	// Чтение переменных окружения с помощью viper
	viper.AutomaticEnv()

	// Определяем флаги
	port := flag.String("port", viper.GetString("PORT"), "Port for the application")
	dbUser := flag.String("db-user", viper.GetString("DB_USER"), "Database user")
	dbPassword := flag.String("db-password", viper.GetString("DB_PASSWORD"), "Database password")
	dbName := flag.String("db-name", viper.GetString("DB_NAME"), "Database name")
	dbHost := flag.String("db-host", viper.GetString("DB_HOST"), "Database host")
	dbPort := flag.String("db-port", viper.GetString("DB_PORT"), "Database port")

	// Парсим флаги
	flag.Parse()

	// Формируем и возвращаем конфигурацию
	config := &Config{
		Port:       fallback(*port, viper.GetString("PORT"), "50051"),
		DbUser:     fallback(*dbUser, viper.GetString("DB_USER"), "user"),
		DbPassword: fallback(*dbPassword, viper.GetString("DB_PASSWORD"), "password"),
		DbName:     fallback(*dbName, viper.GetString("DB_NAME"), "usdt_rates"),
		DbHost:     fallback(*dbHost, viper.GetString("DB_HOST"), "localhost"),
		DbPort:     fallback(*dbPort, viper.GetString("DB_PORT"), "5432"),
	}

	return config, nil
}

// fallback возвращает значение флага, если оно есть, иначе — переменную окружения или значение по умолчанию
func fallback(flagValue, envValue, defaultValue string) string {
	if flagValue != "" {
		return flagValue
	}
	if envValue != "" {
		return envValue
	}
	return defaultValue
}
