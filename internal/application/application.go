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

type Application struct {
	Duration  time.Duration
	IsService bool
}

func NewService() *Application {
	return &Application{
		IsService: true,
	}
}

func New(d time.Duration) *Application {
	return &Application{
		Duration: d,
	}
}

func (a *Application) Run() error {
	// Exiting in duration
	defer log.Printf("application stopped before timeout")
	go a.ExitTimeOut()

	// The program
	log.Println("the program is working")
	cfg, err := config.LoadConfig("config/config.yml")
	if err != nil {
		panic(err)
	}
	hash.SALT = cfg.Salt
	db, err := file_db.NewFileDB("users.json")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Create()
	if err != nil {
		panic(err)
	}
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

	fmt.Println(user_service.StockUpdate(db))

	// The program end

	// Returning without error
	return nil
}

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
