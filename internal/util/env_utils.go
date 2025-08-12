package util

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// TODO: Setear de manera dinamica las variables de entorno
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
	}

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

func GetEnvTimeValues(envVar string, unit time.Duration) (time.Duration, error) {
	err := LoadEnv()

	if err != nil {
		log.Printf("Error loading .env file: %v", err)
		return 0, fmt.Errorf("error loading environment variables: %v", err)
	}

	value := os.Getenv(envVar)
	if value == "" {
		return 0, fmt.Errorf("environment variable %s not found or is empty", envVar)
	}

	converted, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("error converting environment variable %s to integer: %v", envVar, err)
	}

	return time.Duration(converted) * unit, nil
}
