package app

import "context"

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Info(string)
	Error(string)
}

type Storage interface{}

func New(logger Logger, storage Storage) *App {
	return &App{logger: logger, storage: storage}
}

func (a *App) Ping(ctx context.Context) error {
	return nil
}
