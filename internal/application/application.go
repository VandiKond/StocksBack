package application

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/VandiKond/StocksBack/config/config"
	"github.com/VandiKond/StocksBack/pkg/file_db"
	"github.com/VandiKond/StocksBack/pkg/hash"
	"github.com/VandiKond/StocksBack/pkg/user_service"
)

// Thr application program
type Application struct {
	Duration  time.Duration
	IsService bool
}

// Creates a new service application
func NewService() *Application {
	return &Application{
		IsService: true,
	}
}

// Creates a new application
func New(d time.Duration) *Application {
	return &Application{
		Duration: d,
	}
}

// Runs the application
func (a *Application) Run() error {
	// Exiting in duration
	defer log.Printf("application stopped before timeout")
	go a.ExitTimeOut()

	// The program

	// The unchangeable part with setting the program settings

	// Loading config
	cfg, err := config.LoadConfig("config/config.yml")
	if err != nil {
		panic(err)
	}

	// Setting salt
	hash.SALT = cfg.Salt

	// Creating the data base
	db, err := file_db.NewFileDB(cfg.Database.Name)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// Creating the tables
	err = db.Create()
	if err != nil {
		panic(err)
	}

	// The loading part exit

	// Creating test users

	// length, err := db.GetLen()
	// if err != nil {
	// 	panic(err)
	// }
	// usr, err := user_cfg.NewUser("usr1", "pass", length)
	// if err != nil {
	// 	panic(err)
	// }
	// usr.StockBalance = 15
	// db.NewUser(*usr)
	// length, err = db.GetLen()
	// if err != nil {
	// 	panic(err)
	// }
	// usr2, err := user_cfg.NewUser("usr2", "pass", length)
	// if err != nil {
	// 	panic(err)
	// }
	// usr2.StockBalance = 30
	// db.NewUser(*usr2)

	// Updating stocks
	fmt.Println(user_service.StockUpdate(db))

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
	log.Printf("timeout %s has passed. Ending the program", a.Duration)
	os.Exit(http.StatusTeapot)
}
