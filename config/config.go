package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	Postgres_User         string
	Postgres_Db_Name      string
	Postgres_Password     string
	Postgres_Port         int64
	Postgres_Host         string
	Port                  string
	DB_URL                string
	FRONTEND_URL          string
	ACCESS_SECRET_KEY  string
	REFRESH_SECRET_KEY string
)

func Init() {

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file - ", err)
	}

	Postgres_User = os.Getenv("POSTGRES_USER")
	Postgres_Db_Name = os.Getenv("POSTGRES_DB_NAME")
	Postgres_Password = os.Getenv("POSTGRES_PASSWORD")
	Postgres_Port, _ = strconv.ParseInt(os.Getenv("POSTGRES_PORT"), 10, 64)
	Postgres_Host = os.Getenv("POSTGRES_HOST")
	Port = ":" + os.Getenv("PORT")
	SSL_MODE := os.Getenv("SSL_MODE")
	FRONTEND_URL = os.Getenv("FRONTEND_URL")
	ACCESS_SECRET_KEY = os.Getenv("ACCESS_TOKEN_SECRET")
	REFRESH_SECRET_KEY = os.Getenv("REFRESH_TOKEN_SECRET")
	if ACCESS_SECRET_KEY == "" || REFRESH_SECRET_KEY == "" {
		log.Fatal("ACCESS_SECRET_KEY and REFRESH_SECRET_KEY must be set in environment variables")
	}


	os.Setenv("db_url", fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s", Postgres_Host, Postgres_User, Postgres_Password, Postgres_Db_Name, Postgres_Port, SSL_MODE))
	DB_URL = os.Getenv("db_url")
}
