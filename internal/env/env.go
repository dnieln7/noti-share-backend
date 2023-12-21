package env

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func GetEnvProperties() *EvnProperties {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Could not get .env")
	}

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("PORT not found")
	}

	return &EvnProperties{
		Port: port,
	}
}
