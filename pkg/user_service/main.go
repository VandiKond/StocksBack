package user_service

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/VandiKond/StocksBack/config/db_cfg"
	"github.com/VandiKond/StocksBack/config/user_cfg"
	"github.com/VandiKond/StocksBack/pkg/query"
	"github.com/VandiKond/vanerrors"
)

// The errors
const (
	ErrorGettingId     = "error getting id"
	ErrorCreatingUser  = "error creating user"
	ErrorSelectingUser = "error selecting user"
	ToEarlyFarming     = "to early farming"
	ErrorUpdatingUser  = "error updating user"
	NotEnoughSolids    = "not enough solids"
	UserIsBlocked      = "user is blocked"
	UserIsNotBlocked   = "user isn't blocked"
	ErrorCheckingKey   = "error getting key"
	WrongKey           = "wrong key"
)

// Checks the error
func IsServerError(err error) bool {
	s := vanerrors.GetName(err)
	return s == ErrorGettingId || s == ErrorSelectingUser || s == ErrorUpdatingUser || s == ErrorCheckingKey
}

// Global variables
var (
	FarmingLimit       = time.Hour // the farming limit
	StockCost    int64 = 30        // the stock cost
)

// Sign up data
type SignUpUser struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// Sign in data
type SignInUser struct {
	Id       uint64 `json:"id"`
	Password string `json:"password"`
}

// Sign in data with key
type SignInKey struct {
	Key string `json:"key"`
	Id  uint64 `json:"id"`
}

// Creates a new user
func (u SignUpUser) SignUp(db db_cfg.DataBase) (*user_cfg.User, error) {
	// Gets the length of users
	id, err := db.Len()
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorGettingId, err, vanerrors.EmptyHandler)
	}

	// Created a new user
	usr, err := user_cfg.NewUser(u.Name, u.Password, id)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorCreatingUser, err, vanerrors.EmptyHandler)
	}

	// Creates the user in the data base
	err = db.Create(*usr)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorUpdatingUser, err, vanerrors.EmptyHandler)
	}

	return usr, nil
}

// Signing in with secret key
func (u SignInKey) SignInWithKey(db db_cfg.DataBase) (*user_cfg.User, error) {
	// Checking key
	ok, err := db.CheckKey(u.Key)

	// Error
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorCheckingKey, err, vanerrors.EmptyHandler)
	}

	// Wrong key
	if !ok {
		return nil, vanerrors.NewSimple(WrongKey)
	}

	// Getting user
	usr, err := db.GetOne(u.Id)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorSelectingUser, err, vanerrors.EmptyHandler)
	}

	return usr, nil
}

// Signs in
func (u SignInUser) SignIn(db db_cfg.DataBase) (bool, *user_cfg.User, error) {
	// Selects the user by id
	usr, err := Get(u.Id, db)
	if err != nil {
		return false, nil, err
	}

	// Checks the password
	ok := usr.CheckPassword(u.Password)
	if ok {
		return true, usr, nil
	}

	return false, nil, nil
}

// Gets user
func Get(id uint64, db db_cfg.DataBase) (*user_cfg.User, error) {
	// Selects the user by id
	usr, err := db.GetOne(id)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorSelectingUser, err, vanerrors.EmptyHandler)
	}
	return usr, err
}

// Farms
func Farm(id uint64, db db_cfg.DataBase) (int64, *user_cfg.User, error) {
	// Selects the user by id
	usr, err := Get(id, db)
	if err != nil {
		return 0, nil, err
	}

	// Checks the limit
	expected_time := time.Now().Add(-FarmingLimit)
	if usr.LastFarming.After(expected_time) {
		return 0, usr, vanerrors.NewSimple(ToEarlyFarming)
	}

	// Gets the maximum value
	var max int64 = usr.StockBalance
	if max <= StockCost {
		max = StockCost
	}

	// Gets the random value
	amount := rand.Int64N(max)

	// Edits the user

	usr, err = db.UpdateLastFarm(id)
	if err != nil {
		return amount, usr, vanerrors.NewWrap(ErrorUpdatingUser, err, vanerrors.EmptyHandler)
	}

	usr, err = db.UpdateSolids(id, amount)
	if err != nil {
		return amount, usr, vanerrors.NewWrap(ErrorUpdatingUser, err, vanerrors.EmptyHandler)
	}

	return amount, usr, nil
}

// Updates all users with stacks
func StockUpdate(db db_cfg.DataBase) ([]user_cfg.User, error) {
	// Selects all users by query
	users, err := db.GetAllBy(query.Query{
		{
			Separator: query.NOT_SEPARATOR,
			Type:      query.STOCK_BALANCE,
			Sign:      query.MORE,
			Y:         uint64(0),
		},
		{
			Separator: query.AND,
		}, {
			Separator: query.NOT_SEPARATOR,
			Type:      query.IS_BLOCKED,
			Sign:      query.EQUAL,
			Y:         false,
		},
	})
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorSelectingUser, err, vanerrors.EmptyHandler)
	}

	for i := range users {

		usr := &users[i]

		// Updates the user
		usr, err = db.UpdateSolids(usr.Id, usr.StockBalance)
		if err != nil {
			return users, vanerrors.NewWrap(ErrorUpdatingUser, err, vanerrors.EmptyHandler)
		}
		users[i] = *usr
	}

	return users, nil
}

// Byes stocks
func BuyStocks(id uint64, num int64, db db_cfg.DataBase) (*user_cfg.User, error) {
	// Selects the user by id
	usr, err := Get(id, db)
	if err != nil {
		return nil, err
	}

	// Gets the cost
	cost := num * StockCost

	// Checks user balance
	if usr.SolidBalance < cost {
		return usr, vanerrors.NewSimple(NotEnoughSolids, fmt.Sprintf("has %d, need %d", usr.SolidBalance, cost))
	}

	// Updates the user

	usr, err = db.UpdateSolids(usr.Id, -cost)
	if err != nil {
		return usr, vanerrors.NewWrap(ErrorUpdatingUser, err, vanerrors.EmptyHandler)
	}

	usr, err = db.UpdateStocks(usr.Id, num)
	if err != nil {
		return usr, vanerrors.NewWrap(ErrorUpdatingUser, err, vanerrors.EmptyHandler)
	}

	return usr, nil
}

// Updates the user name
func UpdateName(id uint64, name string, db db_cfg.DataBase) (*user_cfg.User, error) {
	// Selects the user by id
	usr, err := Get(id, db)
	if err != nil {
		return nil, err
	}

	// Updates the user
	usr, err = db.UpdateName(usr.Id, name)
	if err != nil {
		return usr, vanerrors.NewWrap(ErrorUpdatingUser, err, vanerrors.EmptyHandler)
	}

	return usr, nil
}

// Updates the user password
func UpdatePassword(id uint64, password string, db db_cfg.DataBase) (*user_cfg.User, error) {
	// Selects the user by id
	usr, err := Get(id, db)
	if err != nil {
		return nil, err
	}

	// Updates the password
	usr.NewPassword(password)
	usr, err = db.UpdatePassword(usr.Id, usr.Password)
	if err != nil {
		return usr, vanerrors.NewWrap(ErrorUpdatingUser, err, vanerrors.EmptyHandler)
	}

	return usr, nil
}

// Blocks the user
func Block(id uint64, db db_cfg.DataBase) (*user_cfg.User, error) {
	// Selects the user by id
	usr, err := Get(id, db)
	if err != nil {
		return nil, err
	}

	// Checks that the user is not blocked
	if usr.IsBlocked {
		return usr, vanerrors.NewSimple(UserIsBlocked)
	}

	// Blocks user
	usr, err = db.UpdateBlock(usr.Id, true)
	if err != nil {
		return usr, vanerrors.NewWrap(ErrorUpdatingUser, err, vanerrors.EmptyHandler)
	}

	return usr, nil
}

// Unblocks the user
func Unblock(id uint64, db db_cfg.DataBase) (*user_cfg.User, error) {
	// Selects the user by id
	usr, err := Get(id, db)
	if err != nil {
		return nil, err
	}

	// Checks that the user is blocked
	if !usr.IsBlocked {
		return usr, vanerrors.NewSimple(UserIsNotBlocked)
	}

	// Blocks user
	usr, err = db.UpdateBlock(usr.Id, false)
	if err != nil {
		return usr, vanerrors.NewWrap(ErrorUpdatingUser, err, vanerrors.EmptyHandler)
	}

	return usr, nil
}
