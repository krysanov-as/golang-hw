package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/krysanov-as/hw12_13_14_15_calendar/internal/app" //nolint:depguard
	"github.com/krysanov-as/hw12_13_14_15_calendar/internal/logger" //nolint:depguard
	internalhttp "github.com/krysanov-as/hw12_13_14_15_calendar/internal/server/http" //nolint:depguard
	memorystorage "github.com/krysanov-as/hw12_13_14_15_calendar/internal/storage/memory" //nolint:depguard
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	cfg, err := NewConfig(configFile)
	if err != nil {
		panic(err)
	}

	logg := logger.New(cfg.Logger.Level)
	storage := memorystorage.New()
	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(
		logg,
		calendar,
		cfg.HTTP.Host,
		cfg.HTTP.Port,
	)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		_ = server.Stop(shutdownCtx)
	}()

	logg.Info("calendar is running")

	if err := server.Start(ctx); err != nil {
		logg.Error(err.Error())
		os.Exit(1)
	}
}
