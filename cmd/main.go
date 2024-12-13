package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/vandi37/StocksBack/internal/application"
)

func main() {
	// Creating a new application with a hour timeout
	app := application.New()

	// Adding graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Running the app
	app.Run(ctx)
}
