package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct{
	DbUrl string
}

func Read()Config{
	if err := godotenv.Load(); err!=nil{
		log.Fatal("Error loading .env file")
	}
	
	dbUrl := os.Getenv("DB_URL")

	return Config{
		DbUrl: dbUrl,
	}
}