package responses

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ResponseUser struct {
	Id           uint64    `json:"id"`
	Name         string    `json:"name"`
	SolidBalance int64     `json:"solid_balance"`
	StockBalance int64     `json:"stock_balance"`
	IsBlocked    bool      `json:"is_blocked"`
	LastFarming  time.Time `json:"last_farming"`
	CreatedAt    time.Time `json:"created_at"`
}

type ErrorResponse struct {
	ErrorName string `json:"error_name"`
	Error     string `json:"error"`
}

func (e ErrorResponse) SendJson(w http.ResponseWriter, code int) error {
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}
	w.WriteHeader(code)
	fmt.Fprint(w, string(b))
	return err
}

type SingUpResponseOK struct {
	User ResponseUser `json:"user"`
}

type SingUpResponseError struct {
	ErrorResponse
}

type SingInResponseError struct {
	ErrorResponse
}

type FarmResponseOK struct {
	User   ResponseUser `json:"user"`
	Amount int64        `json:"amount"`
}

type FarmResponseError struct {
	User ResponseUser `json:"user"`
	ErrorResponse
}

type BuyStocksResponseOK struct {
	User ResponseUser `json:"user"`
}

type BuyStocksResponseError struct {
	User ResponseUser `json:"user"`
	ErrorResponse
}

type UpdateNameResponseOK struct {
	User ResponseUser `json:"user"`
}

type UpdateNameResponseError struct {
	User ResponseUser `json:"user"`
	ErrorResponse
}

type UpdatePasswordResponseOK struct {
	User ResponseUser `json:"user"`
}

type UpdatePasswordResponseError struct {
	User ResponseUser `json:"user"`
	ErrorResponse
}

type BlockResponseOK struct {
	User ResponseUser `json:"user"`
}

type BlockResponseError struct {
	User ResponseUser `json:"user"`
	ErrorResponse
}

type UnblockResponseOK struct {
	User ResponseUser `json:"user"`
}

type UnblockResponseError struct {
	User ResponseUser `json:"user"`
	ErrorResponse
}
