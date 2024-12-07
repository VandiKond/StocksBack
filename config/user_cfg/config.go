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
	// Adding block prefix
	var blocked string
	if u.IsBlocked {
		blocked = "[BLOCKED] "
	}

	// Returning the user data
	return fmt.Sprintf("%suser %d (%s). balance: solids - %d, stocks - %d", blocked, u.Id, u.Name, u.SolidBalance, u.StockBalance)
}

// Creates a new user
func NewUser(name string, password string, id uint64) (*User, error) {
	// Checking the name
	if ok := validName(name); !ok {
		return nil, vanerrors.NewSimple(InvalidName, fmt.Sprintf("name %s has not allowed symbols", name))
	}

	// Checks the password
	if ok := validPassword(password); !ok {
		return nil, vanerrors.NewSimple(InvalidPassword, fmt.Sprintf("password %s has not allowed symbols", password))
	}

	// Hash password
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

// Updates the password
func (u *User) NewPassword(password string) error {
	// Checks the password
	if ok := validPassword(password); !ok {
		return vanerrors.NewSimple(InvalidPassword, fmt.Sprintf("password %s has not allowed symbols", password))
	}

	// Hash password
	hashed_password, err := hash.HashPassword(password)
	if err != nil {
		return err
	}

	// Updating the password
	u.Password = hashed_password

	return nil
}

// valid name
func validName(name string) bool {
	matched, err := regexp.MatchString(`^[a-zA-Z0-9 _-]*$`, name)
	if err != nil {
		return false
	}
	return matched
}

// valid password
func validPassword(password string) bool {
	matched, err := regexp.MatchString(`^[a-zA-Z0-9!#$*_&-]*$`, password)
	if err != nil {
		return false
	}
	return matched
}

// Checks if the user valid
func (u User) Valid() bool {
	return validName(u.Name)
}

// comperes password
func (u User) CheckPassword(password string) bool {
	ok, err := hash.CompareHash(password, u.Password)
	return ok && err == nil
}
