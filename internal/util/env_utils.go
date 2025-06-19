package util

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func LoadEnv() error {
	envLocations := []string{
		".env",
		filepath.Join("..", ".env"),
		filepath.Join("..", "..", ".env"),
		filepath.Join("..", "..", "..", ".env"),
		filepath.Join(os.Getenv("HOME"), ".env"),
	}

	for _, location := range envLocations {
		err := godotenv.Load(location)
		if err == nil {
			return nil
		}
		// Ignoramos el error y continuamos intentando con la siguiente ubicación
	}

	// No encontramos ningún archivo .env, pero no es un error
	// Las variables podrían estar configuradas directamente en el entorno
	return nil
}

func GetTimeoutFromEnv(envVar string, unit time.Duration) time.Duration {
	err := LoadEnv()

	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	value := os.Getenv(envVar)
	if value == "" {
		log.Printf("Searching in %s/.env the variable", envVar)
	}

	converted, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}

	return time.Duration(converted) * unit
}
