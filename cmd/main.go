package main

import (
	"time"

	"github.com/VandiKond/StocksBack/internal/application"
)

func main() {
	// Creating a new application with a hour timeout
	app := application.New(time.Hour)

	// Running the app
	app.Run()
}
