package application

import (
	"fmt"
	"time"

	"github.com/VandiKond/StocksBack/config/config"
	"github.com/VandiKond/StocksBack/http/server"
	"github.com/VandiKond/StocksBack/pkg/file_db"
	"github.com/VandiKond/StocksBack/pkg/hash"
	"github.com/VandiKond/StocksBack/pkg/logger"
)

// Thr application program
type Application struct {
	Duration  time.Duration
	IsService bool
	Logger    logger.Logger
}

// Creates a new service application
func NewService() *Application {
	return &Application{
		IsService: true,
		Logger:    logger.New(),
	}
}

// Creates a new application
func New(d time.Duration) *Application {
	return &Application{
		Duration: d,
		Logger:   logger.New(),
	}
}

// Runs the application
func (a *Application) Run() error {
	// Exiting in duration
	defer a.Logger.Fatalln("application stopped before timeout")
	go a.ExitTimeOut()

	// The program

	// The unchangeable part with setting the program settings

	// Loading config
	cfg, err := config.LoadConfig("config/config.yml")
	if err != nil {
		a.Logger.Fatalln(err)
	}

	// Setting salt
	hash.SALT = cfg.Salt

	// Creating the data base
	db, err := file_db.NewFileDB(cfg.Database.Name, cfg.Key)
	if err != nil {
		a.Logger.Fatalln(err)
	}
	defer db.Close()
	// Creating the tables
	err = db.Create()
	if err != nil {
		a.Logger.Fatalln(err)
	}

	handler := server.NewHandler(db, a.Logger)
	server := server.NewServer(handler, cfg.Port)
	server.Run()

	// The program end

	// Returning without error
	return nil
}

// Exit in duration, if the application isn't in service mode
func (a *Application) ExitTimeOut() {
	// Checking service mod
	if a.IsService {
		return
	}

	// Waiting duration seconds
	time.Sleep(a.Duration)

	// Exiting after timeout
	fmt.Println("")
	a.Logger.Fatalf("timeout %s has passed. Ending the program", a.Duration)
}
