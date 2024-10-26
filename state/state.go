package state

import (
	"go_chi_pgx/repository"
	"sync"
)

type State struct {
	Config     *Config
	Repository repository.Repository
	Logger     *Logger
	Wg         sync.WaitGroup
}

func NewState(cfg *Config, db repository.Repository, logger *Logger) *State {
	return &State{
		Config:     cfg,
		Repository: db,
		Logger:     logger,
	}
}
