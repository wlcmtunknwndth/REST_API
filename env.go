package REST_API

import (
	"github.com/joho/godotenv"
	"log"
)

func LoadEnv() {
	if err := godotenv.Load("local.env"); err != nil {
		log.Print("No .env file found")
	}
}
