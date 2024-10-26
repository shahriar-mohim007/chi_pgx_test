package main

import (
	"github.com/joho/godotenv"
	"go_chi_pgx/cmd/httpserver"
	"go_chi_pgx/repository"
	"go_chi_pgx/state"
	"os"
)

func main() {

	logger := state.New(os.Stdout, state.LevelInfo)
	_ = godotenv.Load()
	cfg, err := state.NewConfig()

	if err != nil {
		logger.PrintError(err, map[string]string{
			"context": "Error loading env value",
		})
	}

	db, err := repository.NewPgRepository(cfg.DatabaseUrl)

	if err != nil {
		logger.PrintError(err, map[string]string{
			"context": "Error initializing the database",
		})
		os.Exit(1)
	}

	appState := state.NewState(cfg, db, logger)

	err = httpserver.Serve(appState)
	if err != nil {
		logger.PrintError(err, map[string]string{
			"context": "serving the application",
		})
	}

}
