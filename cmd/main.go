package main

import (
	"time"

	"github.com/VandiKond/StocksBack/internal/application"
)

func main() {
	app := application.New(time.Hour)
	app.Run()
}
