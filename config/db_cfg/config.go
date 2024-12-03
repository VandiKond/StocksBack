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
	SelectBy(any) ([]user_cfg.User, error)
	SelectNumBy(any, int) ([]user_cfg.User, error)
	SelectOneBy(any) (*user_cfg.User, error)
	Update(user_cfg.User) error
	UpdateGroup([]user_cfg.User) error
	GetLast() (uint64, error)
	io.Closer
}
