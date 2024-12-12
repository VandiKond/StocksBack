package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/VandiKond/StocksBack/internal/application"
	"github.com/VandiKond/StocksBack/pkg/db"
)

func main() {
	// Creating a new application with a hour timeout
	app := application.New(time.Second, db.Constructor{})

	// Adding graceful  shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Running the app
	app.Run(ctx)
}
