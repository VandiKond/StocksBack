package responses

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type MainResponse struct {
	Pages []string `json:"pages"`
}

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
	ErrorResponse
}

type BuyStocksResponseOK struct {
	User ResponseUser `json:"user"`
}

type BuyStocksResponseError struct {
	ErrorResponse
}

type UpdateNameResponseOK struct {
	User ResponseUser `json:"user"`
}

type UpdateNameResponseError struct {
	ErrorResponse
}

type UpdatePasswordResponseOK struct {
	User ResponseUser `json:"user"`
}

type UpdatePasswordResponseError struct {
	ErrorResponse
}

type BlockResponseOK struct {
	User ResponseUser `json:"user"`
}

type BlockResponseError struct {
	ErrorResponse
}

type UnblockResponseOK struct {
	User ResponseUser `json:"user"`
}

type UnblockResponseError struct {
	ErrorResponse
}

type GetResponseError struct {
	ErrorResponse
}

type GetResponseOK struct {
	User ResponseUser `json:"user"`
}
