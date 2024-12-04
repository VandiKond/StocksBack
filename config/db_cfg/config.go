package db_cfg

import (
	"io"

	"github.com/VandiKond/StocksBack/config/user_cfg"
	"github.com/VandiKond/StocksBack/pkg/query"
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
