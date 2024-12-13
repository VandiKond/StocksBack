package db_cfg

import (
	"io"

	"github.com/vandi37/StocksBack/config/config"
	"github.com/vandi37/StocksBack/config/user_cfg"
	"github.com/vandi37/StocksBack/pkg/query"
)

// The data base interface should represent any one-table data base
// Ir should storage data based on the user_cfg.User signature
//
// - Create : Creates the table is it not exists
// - NewUser : Creates a new user
// - GetAll : Gets all users
// - Select : Selects a user by it's id
// - SelectBy : Selects all users by query
// - SelectNumBy : Selects users by query with num limit
// - SelectOneBy : Selects user by query
// - Update : Updates user data
// - UpdateGroup : Updates a group of users
// - GetLen : gets the total amount of users (it should get the last id of the user)
// - io.Closer : closes the data base
type DataBase interface {
	Init() error
	Create(user user_cfg.User) error
	GetAll() ([]user_cfg.User, error)
	GetAllBy(query query.Query) ([]user_cfg.User, error)
	GetNumBy(query query.Query, num int) ([]user_cfg.User, error)
	GetOneBy(query query.Query) (*user_cfg.User, error)
	GetOne(id uint64) (*user_cfg.User, error)
	UpdateSolids(id uint64, num int64) (*user_cfg.User, error)
	UpdateStocks(id uint64, num int64) (*user_cfg.User, error)
	UpdateName(id uint64, name string) (*user_cfg.User, error)
	UpdatePassword(id uint64, password string) (*user_cfg.User, error)
	UpdateBlock(id uint64, block bool) (*user_cfg.User, error)
	UpdateLastFarm(id uint64) (*user_cfg.User, error)
	Len() (uint64, error)
	CheckKey(key string) (bool, error)
	io.Closer
}

// It creates a new database
type Constructor interface {
	New(cfg config.DatabaseCfg, key string) (DataBase, error)
}
