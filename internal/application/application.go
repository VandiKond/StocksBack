package application

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/vandi37/StocksBack/config/config"
	"github.com/vandi37/StocksBack/config/db_cfg"
	"github.com/vandi37/StocksBack/config/db_cfg/constructors"
	"github.com/vandi37/StocksBack/http/server"
	"github.com/vandi37/StocksBack/pkg/closer"
	"github.com/vandi37/StocksBack/pkg/cron"
	"github.com/vandi37/StocksBack/pkg/hash"
	"github.com/vandi37/StocksBack/pkg/logger"
	"github.com/vandi37/StocksBack/pkg/user_service"
	"github.com/vandi37/vanerrors"
)

// The errors
const (
	ErrorUpdatingStocks  = "error updating stocks"
	ErrorParsingDuration = "error parsing duration"
)

// Thr application program
type Application struct {
}

// Creates a new application
func New() *Application {
	return &Application{}
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
	// Creates logger
	logger := logger.New()

	// Loading config
	cfg, err := config.LoadConfig("config/config.yml")
	if err != nil {
		logger.Fatalln(err)
	}
	logger.Println("got config")

	// Getting database constructor
	constructor, err := constructors.Get(cfg.Database.Type)
	if err != nil {
		logger.Fatalln(err)
	}

	// Getting duration
	d, err := time.ParseDuration(cfg.App.Duration)
	if err != nil {
		logger.Fatalln(ErrorParsingDuration)
	}

	// Setting context
	if !cfg.App.IsService {
		var stop context.CancelFunc
		ctx, stop = context.WithTimeout(ctx, d)
		defer stop()
	}

	// The program

	// Creating closer
	closer := closer.New(logger)

	// Setting salt
	hash.SALT = cfg.Salt

	// Creating the data base
	db, err := constructor.New(cfg.Database, cfg.Key)
	if err != nil {
		logger.Fatalln(err)
	}
	closer.Add(db.Close)

	// Creating the tables
	err = db.Init()
	if err != nil {
		logger.Fatalln(err)
	}

	logger.Println("database connected")

	// Running cron
	cr := cron.New(time.Hour*24, 21, CronFunc(db, logger), logger)
	cr.Run()

	handler := server.NewHandler(db, logger)
	server := server.NewServer(handler, cfg.Port)
	closer.Add(server.Close)

	go server.Run()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	closer.Close(ctx)

	// The program end
	logger.Println("application stopped")

	os.Exit(http.StatusTeapot)
}
