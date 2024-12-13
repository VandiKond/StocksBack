package requests

import "github.com/vandi37/StocksBack/pkg/user_service"

type SignUp struct {
	user_service.SignUpUser
}

type Farm struct{}

type BuyStocks struct {
	Num int64 `json:"num"`
}

type UpdateName struct {
	Name string `json:"name"`
}

type UpdatePassword struct {
	Password string `json:"password"`
}

type Block struct{}

type Unblock struct{}

type Get struct {
	Id uint64 `json:"id"`
}
