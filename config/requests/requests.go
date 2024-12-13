package requests

import "github.com/vandi37/StocksBack/pkg/user_service"

type SignUpRequest struct {
	user_service.SignUpUser
}

type FarmRequest struct{}

type BuyStocksRequest struct {
	Num int64 `json:"num"`
}

type UpdateNameRequest struct {
	Name string `json:"name"`
}

type UpdatePasswordRequest struct {
	Password string `json:"password"`
}

type BlockRequest struct{}

type UnblockRequest struct{}

type GetRequest struct {
	Id uint64 `json:"id"`
}
