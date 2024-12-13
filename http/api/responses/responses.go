package responses

import (
	"time"
)

// The response content types
const (
	SignUpType         = "signup"
	FarmType           = "farm"
	BuyStocksType      = "buy-stocks"
	UpdateNameType     = "update-name"
	UpdatePasswordType = "update-password"
	BlockType          = "block"
	UnblockType        = "unblock"
	GetType            = "get"
	ErrorType          = "error"
)

type User struct {
	Id           uint64    `json:"id"`
	Name         string    `json:"name"`
	SolidBalance int64     `json:"solid_balance"`
	StockBalance int64     `json:"stock_balance"`
	IsBlocked    bool      `json:"is_blocked"`
	LastFarming  time.Time `json:"last_farming"`
	CreatedAt    time.Time `json:"created_at"`
}

type SignUp struct {
	User User `json:"user"`
}

type Farm struct {
	User   User  `json:"user"`
	Amount int64 `json:"amount"`
}

type BuyStocks struct {
	User User `json:"user"`
}

type UpdateName struct {
	User User `json:"user"`
}

type UpdatePassword struct {
	User User `json:"user"`
}

type Block struct {
	User User `json:"user"`
}

type Unblock struct {
	User User `json:"user"`
}

type Get struct {
	User User `json:"user"`
}
