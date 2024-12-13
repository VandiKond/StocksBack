package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/vandi37/StocksBack/internal/application"
	"github.com/vandi37/StocksBack/pkg/file_db"
)

func main() {
	// Creating a new application with a hour timeout
	app := application.New(time.Hour, file_db.Constructor{})

	// Adding graceful  shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Running the app
	app.Run(ctx)
}
