package main

import (
	"database/sql"
	"embed"
	"talky-space-be/config"
	"talky-space-be/handlers"
	"talky-space-be/utils"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/robfig/cron"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {

	config.Init()
	db := config.DBInit()
	toRunGooseMigration(config.DB_URL)
	defer config.CloseDB(db)

	handlerConnection := handlers.New(db)

	defer config.CloseDB(db)
	s := cron.New()

	go utils.HubInstance.Run()

	// _, err := s.AddFunc("@every 10m", jobs.RunDailyMigrations)
	// if err != nil {
	// 	log.Println("Error scheduling RunDailyMigrations:", err)
	// }

	// _, err = s.AddFunc("@every 10s", jobs.SendMessages)
	// if err != nil {
	// 	log.Println("Error scheduling SendMessages:", err)
	// }

	s.Start()

	r := handlerConnection.GetRouter()
	main_err := r.Run(config.Port)
	if main_err != nil {
		return
	}
	s.Stop()

}

func toRunGooseMigration(url string) {

	db, err := sql.Open("postgres", url)
	if err != nil {
		panic(err)
	}
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}
	if err := goose.Up(db, "migrations", goose.WithAllowMissing()); err != nil {
		panic(err)
	}
	db.Close()
}
