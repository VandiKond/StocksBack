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
	InvalidUser        = "invalid user"
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
	FarmingLimit        = time.Hour // the farming limit
	StockCost    uint64 = 30        // the stock cost
)

// Sing up data
type SingUpUser struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// Sing in data
type SingInUser struct {
	Id       uint64 `json:"id"`
	Password string `json:"password"`
}

// Sing in data with key
type SingInKey struct {
	Key string `json:"key"`
	Id  uint64 `json:"id"`
}

// Creates a new user
func (u SingUpUser) SingUp(db db_cfg.DataBase) (*user_cfg.User, error) {
	// Gets the length of users
	id, err := db.GetLen()
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorGettingId, err, vanerrors.EmptyHandler)
	}

	// Created a new user
	usr, err := user_cfg.NewUser(u.Name, u.Password, id)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorCreatingUser, err, vanerrors.EmptyHandler)
	}

	// Creates the user in the data base
	err = db.NewUser(*usr)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorUpdatingUser, err, vanerrors.EmptyHandler)
	}

	return usr, nil
}

// Singing in with secret key
func (u SingInKey) SingInWithKey(db db_cfg.DataBase) (bool, *user_cfg.User, error) {
	// Checking key
	ok, err := db.CheckKey(u.Key)

	// Error
	if err != nil {
		return false, nil, vanerrors.NewWrap(ErrorCheckingKey, err, vanerrors.EmptyHandler)
	}

	// Wrong key
	if !ok {
		return false, nil, vanerrors.NewSimple(WrongKey)
	}

	// Getting user
	usr, err := db.Select(u.Id)
	if err != nil {
		return false, nil, vanerrors.NewWrap(ErrorSelectingUser, err, vanerrors.EmptyHandler)
	}

	return true, usr, nil
}

// Sings in
func (u SingInUser) SingIn(db db_cfg.DataBase) (bool, *user_cfg.User, error) {
	// Selects the user by id
	usr, err := db.Select(u.Id)
	if err != nil {
		return false, nil, vanerrors.NewWrap(ErrorSelectingUser, err, vanerrors.EmptyHandler)
	}

	// Checks the password
	ok := usr.CheckPassword(u.Password)
	if ok {
		return ok, usr, nil
	}

	return ok, nil, nil
}

// Farms
func Farm(id uint64, db db_cfg.DataBase) (uint64, *user_cfg.User, error) {
	// Selects the user by id
	usr, err := db.Select(id)
	if err != nil {
		return 0, nil, vanerrors.NewWrap(ErrorSelectingUser, err, vanerrors.EmptyHandler)
	}

	// Checks the limit
	expected_time := time.Now().Add(-FarmingLimit)
	if usr.LastFarming.After(expected_time) {
		return 0, usr, vanerrors.NewSimple(ToEarlyFarming)
	}

	// Gets the maximum value
	var max uint64 = usr.StockBalance
	if max <= StockCost {
		max = StockCost
	}

	// Gets the random value
	amount := rand.Uint64N(max)

	// Edits the user
	usr.SolidBalance += amount
	usr.LastFarming = time.Now()

	// Checks that user is valid
	if !usr.Valid() {
		return amount, usr, vanerrors.NewSimple(InvalidUser)
	}

	// Updates the user
	err = db.Update(*usr)
	if err != nil {
		return amount, usr, vanerrors.NewWrap(ErrorUpdatingUser, err, vanerrors.EmptyHandler)
	}

	return amount, usr, nil
}

// Updates all users with stacks
func StockUpdate(db db_cfg.DataBase) ([]user_cfg.User, error) {
	// Selects all users by query
	users, err := db.SelectBy(query.Query{
		{
			Separator: query.NOT_SEPARATOR,
			Type:      query.STOCK_BALANCE,
			Sing:      query.MORE,
			Y:         uint64(0),
		},
		{
			Separator: query.AND,
		}, {
			Separator: query.NOT_SEPARATOR,
			Type:      query.IS_BLOCKED,
			Sing:      query.EQUAL,
			Y:         false,
		},
	})
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorSelectingUser, err, vanerrors.EmptyHandler)
	}

	for i := range users {
		// Updates the user
		usr := &users[i]
		usr.SolidBalance += usr.StockBalance

		// Checks that user is valid
		if !usr.Valid() {
			return nil, vanerrors.NewSimple(InvalidUser)
		}
	}

	// Updates a group of users
	err = db.UpdateGroup(users)
	if err != nil {
		return users, vanerrors.NewWrap(ErrorUpdatingUser, err, vanerrors.EmptyHandler)
	}

	return users, nil
}

// Byes stocks
func BuyStocks(id uint64, num uint64, db db_cfg.DataBase) (*user_cfg.User, error) {
	// Selects the user by id
	usr, err := db.Select(id)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorSelectingUser, err, vanerrors.EmptyHandler)
	}

	// Gets the cost
	cost := num * StockCost

	// Checks user balance
	if usr.SolidBalance < cost {
		return usr, vanerrors.NewSimple(NotEnoughSolids, fmt.Sprintf("has %d, need %d", usr.SolidBalance, cost))
	}

	// Updates the user
	usr.SolidBalance -= cost
	usr.StockBalance += num
	// Checks that user is valid
	if !usr.Valid() {
		return usr, vanerrors.NewSimple(InvalidUser)
	}

	// Updates user in data base
	err = db.Update(*usr)
	if err != nil {
		return usr, vanerrors.NewWrap(ErrorUpdatingUser, err, vanerrors.EmptyHandler)
	}

	return usr, nil
}

// Updates the user name
func UpdateName(id uint64, name string, db db_cfg.DataBase) (*user_cfg.User, error) {
	// Selects the user by id
	usr, err := db.Select(id)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorSelectingUser, err, vanerrors.EmptyHandler)
	}

	// Updates the user
	usr.Name = name

	// Checks that user is valid
	if !usr.Valid() {
		return usr, vanerrors.NewSimple(InvalidUser)
	}

	// Updates user in data base
	err = db.Update(*usr)
	if err != nil {
		return usr, vanerrors.NewWrap(ErrorUpdatingUser, err, vanerrors.EmptyHandler)
	}

	return usr, nil
}

// Updates the user password
func UpdatePassword(id uint64, password string, db db_cfg.DataBase) (*user_cfg.User, error) {
	// Selects the user by id
	usr, err := db.Select(id)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorSelectingUser, err, vanerrors.EmptyHandler)
	}

	// Updates the password
	usr.NewPassword(password)

	// Checks that user is valid
	if !usr.Valid() {
		return usr, vanerrors.NewSimple(InvalidUser)
	}

	// Updates user in data base
	err = db.Update(*usr)
	if err != nil {
		return usr, vanerrors.NewWrap(ErrorUpdatingUser, err, vanerrors.EmptyHandler)
	}

	return usr, nil
}

// Blocks the user
func Block(id uint64, db db_cfg.DataBase) (*user_cfg.User, error) {
	// Selects the user by id
	usr, err := db.Select(id)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorSelectingUser, err, vanerrors.EmptyHandler)
	}

	// Checks that the user is not blocked
	if usr.IsBlocked {
		return usr, vanerrors.NewSimple(UserIsBlocked)
	}

	// Blocks user
	usr.IsBlocked = true

	// Checks that user is valid
	if !usr.Valid() {
		return usr, vanerrors.NewSimple(InvalidUser)
	}

	// Updates user in data base
	err = db.Update(*usr)
	if err != nil {
		return usr, vanerrors.NewWrap(ErrorUpdatingUser, err, vanerrors.EmptyHandler)
	}

	return usr, nil
}

// Unblocks the user
func Unblock(id uint64, db db_cfg.DataBase) (*user_cfg.User, error) {
	// Selects the user by id
	usr, err := db.Select(id)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorSelectingUser, err, vanerrors.EmptyHandler)
	}

	// Checks that the user is blocked
	if !usr.IsBlocked {
		return usr, vanerrors.NewSimple(UserIsNotBlocked)
	}

	// Unblocks user
	usr.IsBlocked = false

	// Checks that user is valid\
	if !usr.Valid() {
		return usr, vanerrors.NewSimple(InvalidUser)
	}

	// Updates user in data base

	err = db.Update(*usr)
	if err != nil {
		return usr, vanerrors.NewWrap(ErrorUpdatingUser, err, vanerrors.EmptyHandler)
	}

	return usr, nil
}
