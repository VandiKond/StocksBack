package db_cfg

import (
	"io"

	"github.com/VandiKond/StocksBack/config/user_cfg"
)

type DataBase interface {
	Create() error
	NewUser(user_cfg.User) error
	GetAll() ([]user_cfg.User, error)
	Select(uint64) (*user_cfg.User, error)
	SelectBy(string) ([]user_cfg.User, error)
	SelectOneBy(string) (*user_cfg.User, error)
	io.Closer
}
