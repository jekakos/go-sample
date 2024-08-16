package infra

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	if _, err := os.Stat("../.env"); err == nil {
		err := godotenv.Load("../.env")
		if err != nil {
			panic("Error loading .env file")
		}
	}
	return nil
}
