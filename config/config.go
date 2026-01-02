package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	Port               string
	DB_URL             string
	FRONTEND_URL       string
	ACCESS_SECRET_KEY  string
	REFRESH_SECRET_KEY string
)

func Init() {

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file - ", err)
	}

	Port = ":" + os.Getenv("PORT")
	FRONTEND_URL = os.Getenv("FRONTEND_URL")
	ACCESS_SECRET_KEY = os.Getenv("ACCESS_TOKEN_SECRET")
	REFRESH_SECRET_KEY = os.Getenv("REFRESH_TOKEN_SECRET")
	if ACCESS_SECRET_KEY == "" || REFRESH_SECRET_KEY == "" {
		log.Fatal("ACCESS_SECRET_KEY and REFRESH_SECRET_KEY must be set in environment variables")
	}

	DB_URL = os.Getenv("DB_URL")
}
