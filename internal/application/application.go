package application

import (
	"fmt"
	"time"

	"github.com/VandiKond/StocksBack/config/config"
	"github.com/VandiKond/StocksBack/config/db_cfg"
	"github.com/VandiKond/StocksBack/http/server"
	"github.com/VandiKond/StocksBack/pkg/cron"
	"github.com/VandiKond/StocksBack/pkg/db"
	"github.com/VandiKond/StocksBack/pkg/hash"
	"github.com/VandiKond/StocksBack/pkg/logger"
	"github.com/VandiKond/StocksBack/pkg/user_service"
	"github.com/VandiKond/vanerrors"
)

// The errors
const (
	ErrorUpdatingStocks = "error updating stocks"
)

// Thr application program
type Application struct {
	Duration  time.Duration
	IsService bool
	Logger    *logger.Logger
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

// Cron func for updating user
func CronFunc(db db_cfg.DataBase, logger *logger.Logger) func() error {
	return func() error {
		users, err := user_service.StockUpdate(db)
		for _, u := range users {
			logger.Println("%v got solids from stocks", u)
		}
		if err != nil {
			return vanerrors.NewWrap(ErrorUpdatingStocks, err, vanerrors.EmptyHandler)
		}
		return nil
	}
}

// Runs the application
func (a *Application) Run() {
	// Exiting in duration
	defer a.Logger.Fatalln("application stopped before timeout")
	go a.ExitTimeOut()

	// The program
	a.Logger.Println("Program started")
	// The unchangeable part with setting the program settings

	// Loading config

	cfg, err := config.LoadConfig("config/config.yml")
	if err != nil {
		a.Logger.Fatalln(err)
	}
	a.Logger.Println("got config")

	// Setting salt
	hash.SALT = cfg.Salt

	// Creating the data base
	db, err := db.New(cfg.Database, cfg.Key)
	if err != nil {
		a.Logger.Fatalln(err)
	}
	defer db.Close()
	// Creating the tables
	err = db.Init()
	if err != nil {
		a.Logger.Fatalln(err)
	}

	a.Logger.Println("database connected")

	// Running cron
	cr := cron.New(time.Hour*24, 21, CronFunc(db, a.Logger), a.Logger)
	cr.Run()

	handler := server.NewHandler(db, a.Logger)
	server := server.NewServer(handler, cfg.Port)
	server.Run()

	// The program end
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
