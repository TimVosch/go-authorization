package main

import (
	"context"
	"flag"
	"fmt"
	"goauthz/app"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var (
	// Migrate flag
	migrate = flag.Bool("migrate", false, "Migrate database")
)

func init() {
	flag.Parse()
}

func main() {
	if err := Run(); err != nil {
		fmt.Printf("Error occured: %v", err)
	}
}

func connectDB() (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", "data/data.db")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func migrateDB(db *sqlx.DB) error {
	// Load migrations
	migration, err := os.ReadFile("./migration.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(migration))
	if err != nil {
		return err
	}

	fmt.Println("Applied migrations")
	return nil
}

func Run() error {
	fmt.Println("Starting server...")
	db, err := connectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// Migrate database
	if *migrate {
		if err := migrateDB(db); err != nil {
			return err
		}
	}

	svc := app.New(db)

	router := svc.Routes()
	srv := &http.Server{
		Addr:    "0.0.0.0:3000",
		Handler: router,
	}

	errChan := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()
	fmt.Println("Server listening...")

	// catch interrupts
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Wait for error or exit
	err = nil
	select {
	case err = <-errChan:
	case <-sigChan:
	}
	srv.Shutdown(context.Background())

	return err
}
