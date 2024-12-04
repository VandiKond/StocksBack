package db_cfg

import (
	"io"

	"github.com/VandiKond/StocksBack/config/user_cfg"
	"github.com/VandiKond/StocksBack/pkg/query"
)

type DataBase interface {
	Create() error
	NewUser(user_cfg.User) error
	GetAll() ([]user_cfg.User, error)
	Select(uint64) (*user_cfg.User, error)
	SelectBy(query.Query) ([]user_cfg.User, error)
	SelectNumBy(query.Query, int) ([]user_cfg.User, error)
	SelectOneBy(query.Query) (*user_cfg.User, error)
	Update(user_cfg.User) error
	UpdateGroup([]user_cfg.User) error
	GetLen() (uint64, error)
	io.Closer
}
