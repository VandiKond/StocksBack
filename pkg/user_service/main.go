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
)

var (
	FarmingLimit        = time.Hour
	StockCost    uint64 = 30
)

type SingUpUser struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type SingInUser struct {
	Id       uint64 `json:"id"`
	Password string `json:"password"`
}

func (u SingUpUser) SingUp(db db_cfg.DataBase) (*user_cfg.User, error) {
	id, err := db.GetLen()
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorGettingId, err, vanerrors.EmptyHandler)
	}
	usr, err := user_cfg.NewUser(u.Name, u.Password, id)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorCreatingUser, err, vanerrors.EmptyHandler)
	}
	err = db.NewUser(*usr)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorCreatingUser, err, vanerrors.EmptyHandler)
	}
	return usr, nil
}

func (u SingInUser) SingIn(db db_cfg.DataBase) (bool, *user_cfg.User, error) {
	usr, err := db.Select(u.Id)
	if err != nil {
		return false, nil, vanerrors.NewWrap(ErrorSelectingUser, err, vanerrors.EmptyHandler)
	}
	ok := usr.CheckPassword(u.Password)
	if ok {
		return ok, usr, nil
	}
	return ok, nil, nil
}

func Farm(id uint64, db db_cfg.DataBase) (uint64, *user_cfg.User, error) {
	usr, err := db.Select(id)
	if err != nil {
		return 0, nil, vanerrors.NewWrap(ErrorSelectingUser, err, vanerrors.EmptyHandler)
	}
	expected_time := time.Now().Add(-FarmingLimit)
	if usr.LastFarming.After(expected_time) {
		return 0, usr, vanerrors.NewSimple(ToEarlyFarming)
	}
	amount := rand.Uint64N(usr.SolidBalance / StockCost)

	usr.SolidBalance += amount
	usr.LastFarming = time.Now()
	if !usr.Valid() {
		return amount, usr, vanerrors.NewSimple(InvalidUser)
	}
	err = db.Update(*usr)
	if err != nil {
		return amount, usr, vanerrors.NewWrap(ErrorUpdatingUser, err, vanerrors.EmptyHandler)
	}
	return amount, usr, nil
}

func StockUpdate(db db_cfg.DataBase) ([]user_cfg.User, error) {
	users, err := db.SelectBy([]query.QuerySetting{{
		Separator: query.NOT_SEPARATOR,
		Type:      query.STOCK_BALANCE,
		Sing:      query.MORE,
		Y:         uint64(0),
	}})
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorSelectingUser, err, vanerrors.EmptyHandler)
	}
	for i := range users {
		usr := &users[i]
		usr.SolidBalance += usr.StockBalance
		if !usr.Valid() {
			return nil, vanerrors.NewSimple(InvalidUser)
		}
	}
	err = db.UpdateGroup(users)
	if err != nil {
		return users, vanerrors.NewWrap(ErrorUpdatingUser, err, vanerrors.EmptyHandler)
	}
	return users, nil
}

func BuyStocks(id uint64, num uint64, db db_cfg.DataBase) (*user_cfg.User, error) {
	usr, err := db.Select(id)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorSelectingUser, err, vanerrors.EmptyHandler)
	}
	cost := num * StockCost
	if usr.SolidBalance < cost {
		return usr, vanerrors.NewSimple(NotEnoughSolids, fmt.Sprintf("has %d, need %d", usr.SolidBalance, cost))
	}
	usr.SolidBalance -= cost
	usr.StockBalance += num
	if !usr.Valid() {
		return usr, vanerrors.NewSimple(InvalidUser)
	}
	err = db.Update(*usr)
	if err != nil {
		return usr, vanerrors.NewWrap(ErrorUpdatingUser, err, vanerrors.EmptyHandler)
	}
	return usr, nil
}

func UpdateName(id uint64, name string, db db_cfg.DataBase) (*user_cfg.User, error) {
	usr, err := db.Select(id)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorSelectingUser, err, vanerrors.EmptyHandler)
	}
	usr.Name = name
	if !usr.Valid() {
		return usr, vanerrors.NewSimple(InvalidUser)
	}
	err = db.Update(*usr)
	if err != nil {
		return usr, vanerrors.NewWrap(ErrorUpdatingUser, err, vanerrors.EmptyHandler)
	}
	return usr, nil
}

func UpdatePassword(id uint64, password string, db db_cfg.DataBase) (*user_cfg.User, error) {
	usr, err := db.Select(id)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorSelectingUser, err, vanerrors.EmptyHandler)
	}
	usr.NewPassword(password)
	if !usr.Valid() {
		return usr, vanerrors.NewSimple(InvalidUser)
	}
	err = db.Update(*usr)
	if err != nil {
		return usr, vanerrors.NewWrap(ErrorUpdatingUser, err, vanerrors.EmptyHandler)
	}
	return usr, nil
}

func Block(id uint64, db db_cfg.DataBase) (*user_cfg.User, error) {
	usr, err := db.Select(id)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorSelectingUser, err, vanerrors.EmptyHandler)
	}
	if usr.IsBlocked {
		return usr, vanerrors.NewSimple(UserIsBlocked)
	}
	usr.IsBlocked = true
	if !usr.Valid() {
		return usr, vanerrors.NewSimple(InvalidUser)
	}
	err = db.Update(*usr)
	if err != nil {
		return usr, vanerrors.NewWrap(ErrorUpdatingUser, err, vanerrors.EmptyHandler)
	}
	return usr, nil
}

func UnBlock(id uint64, db db_cfg.DataBase) (*user_cfg.User, error) {
	usr, err := db.Select(id)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorSelectingUser, err, vanerrors.EmptyHandler)
	}
	if !usr.IsBlocked {
		return usr, vanerrors.NewSimple(UserIsNotBlocked)
	}
	usr.IsBlocked = false
	if !usr.Valid() {
		return usr, vanerrors.NewSimple(InvalidUser)
	}
	err = db.Update(*usr)
	if err != nil {
		return usr, vanerrors.NewWrap(ErrorUpdatingUser, err, vanerrors.EmptyHandler)
	}
	return usr, nil
}
