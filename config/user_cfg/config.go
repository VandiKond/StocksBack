package user_cfg

import (
	"fmt"
	"regexp"
	"time"

	"github.com/VandiKond/StocksBack/pkg/hash"
	"github.com/VandiKond/vanerrors"
)

// The errors
const (
	InvalidName     = "invalid name"     // invalid name
	InvalidPassword = "invalid password" // invalid password
)

// The user structure
type User struct {
	Id           uint64    `json:"id"`
	Name         string    `json:"name"`
	Password     string    `json:"password"`
	SolidBalance uint64    `json:"solid_balance"`
	StockBalance uint64    `json:"stock_balance"`
	IsBlocked    bool      `json:"is_blocked"`
	LastFarming  time.Time `json:"last_farming"`
	CreatedAt    time.Time `json:"created_at"`
}

// Sets the user to string
func (u User) String() string {
	var blocked string
	if u.IsBlocked {
		blocked = "[BLOCKED] "
	}
	return fmt.Sprintf("%suser %d (%s). balance: solids - %d, stocks - %d. Created at: %s", blocked, u.Id, u.Name, u.SolidBalance, u.StockBalance, u.CreatedAt.Format("01.02.2006 15:04:05"))
}

// Creates a new user
func NewUser(name string, password string, id uint64) (*User, error) {
	if ok := validName(name); !ok {
		return nil, vanerrors.NewSimple(InvalidName, fmt.Sprintf("name %s has not allowed symbols", name))
	}
	if ok := validPassword(password); !ok {
		return nil, vanerrors.NewSimple(InvalidPassword, fmt.Sprintf("password %s has not allowed symbols", password))
	}
	hashed_password, err := hash.HashPassword(password)
	if err != nil {
		return nil, err
	}
	return &User{
		Id:           id,
		Name:         name,
		Password:     hashed_password,
		SolidBalance: 0,
		StockBalance: 0,
		IsBlocked:    false,
		CreatedAt:    time.Now(),
	}, nil
}

func validName(name string) bool {
	matched, err := regexp.MatchString(`^[a-zA-Z0-9 _-]*$`, name)
	if err != nil {
		return false
	}
	return matched
}

func validPassword(password string) bool {
	matched, err := regexp.MatchString(`^[a-zA-Z0-9!#$*_&-]*$`, password)
	if err != nil {
		return false
	}
	return matched
}

func (u User) Valid() bool {
	return validPassword(u.Password) && validName(u.Name)
}

func (u User) CheckPassword(password string) bool {
	ok, err := hash.CompareHash(password, u.Password)
	return ok && err == nil
}
