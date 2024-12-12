package application

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/VandiKond/StocksBack/config/config"
	"github.com/VandiKond/StocksBack/config/db_cfg"
	"github.com/VandiKond/StocksBack/http/server"
	"github.com/VandiKond/StocksBack/pkg/closer"
	"github.com/VandiKond/StocksBack/pkg/cron"
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
	duration    time.Duration
	isService   bool
	logger      *logger.Logger
	constructor db_cfg.Constructor
}

// Creates a new service application
func NewService(constr db_cfg.Constructor) *Application {
	return &Application{
		isService:   true,
		logger:      logger.New(),
		constructor: constr,
	}
}

// Creates a new application
func New(d time.Duration, constr db_cfg.Constructor) *Application {
	return &Application{
		duration:    d,
		logger:      logger.New(),
		constructor: constr,
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
func (a *Application) Run(ctx context.Context) {
	// Exiting in duration
	if !a.isService {
		var stop context.CancelFunc
		ctx, stop = context.WithTimeout(ctx, a.duration)
		defer stop()
	}

	// The program
	a.logger.Println("Program started")
	// The unchangeable part with setting the program settings

	// Loading config

	cfg, err := config.LoadConfig("config/config.yml")
	if err != nil {
		a.logger.Fatalln(err)
	}
	a.logger.Println("got config")

	// Creating closer
	closer := closer.New(a.logger)

	// Setting salt
	hash.SALT = cfg.Salt

	// Creating the data base
	db, err := a.constructor.New(cfg.Database, cfg.Key)
	if err != nil {
		a.logger.Fatalln(err)
	}
	closer.Add(db.Close)

	// Creating the tables
	err = db.Init()
	if err != nil {
		a.logger.Fatalln(err)
	}

	a.logger.Println("database connected")

	// Running cron
	cr := cron.New(time.Hour*24, 21, CronFunc(db, a.logger), a.logger)
	cr.Run()

	handler := server.NewHandler(db, a.logger)
	server := server.NewServer(handler, cfg.Port)
	closer.Add(server.Close)

	go server.Run()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	closer.Close(ctx)

	// The program end
	a.logger.Println("application stopped")

	os.Exit(http.StatusTeapot)
}
