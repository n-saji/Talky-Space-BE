package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"talky-space-be/config"
	"talky-space-be/daos"
	"talky-space-be/handlers"
	"talky-space-be/service"
	"talky-space-be/utils"
	"time"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {

	configuration, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading global configuration: %v\n", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	pool, err := config.ConnectToDb(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	err = config.PingDb(ctx, pool)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error pinging database: %v\n", err)
		os.Exit(1)
	}

	toRunGooseMigration(configuration.DatabaseURL)

	Dao := daos.NewPgxDao(pool)
	service := service.New(Dao)
	router := handlers.NewRouter(ctx, service, configuration.RequestTimeout)

	go utils.HubInstance.Run(ctx)

	server := http.Server{
		Addr:              ":" + configuration.ServerPort,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		Handler:           router,
	}

	go func() {
		fmt.Println("Starting server on :" + configuration.ServerPort + "...")
		if err := server.ListenAndServe(); err != nil {
			fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	fmt.Println("Shutting down server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		fmt.Fprintf(os.Stderr, "Server shutdown error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Server gracefully stopped")

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
	defer db.Close()
}
