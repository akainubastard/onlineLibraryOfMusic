package details

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type config struct{
	LOCAL_PORT string
	POSTGRES_HOST string
	POSTGRES_PORT string
	POSTGRES_USER string
	POSTGRES_PASSWORD string
	POSTGRES_DBNAME string
}
var Config *config

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	Config = &config{
        LOCAL_PORT: os.Getenv("LOCAL_PORT"),
        POSTGRES_HOST: os.Getenv("POSTGRES_HOST"),
		POSTGRES_PORT: os.Getenv("POSTGRES_PORT"),
		POSTGRES_USER: os.Getenv("POSTGRES_USER"),
		POSTGRES_PASSWORD: os.Getenv("POSTGRES_PASSWORD"),
		POSTGRES_DBNAME: os.Getenv("POSTGRES_DBNAME"),
    }
}


