package requests

import "github.com/VandiKond/StocksBack/pkg/user_service"

type SingUpRequest struct {
	User user_service.SingUpUser `json:"user"`
}

type SingInRequest struct {
	User user_service.SingInUser `json:"user"`
}

type FarmRequest struct {
	SingInRequest
}

type BuyStocksRequest struct {
	SingInRequest
	Num string `json:"num"`
}

type UpdateNameRequest struct {
	SingInRequest
	Name string `json:"name"`
}

type UpdatePasswordRequest struct {
	SingInRequest
	Password string `json:"password"`
}

type BlockRequest struct {
	Id  uint64 `json:"id"`
	Key string `json:"key"`
}

type UnlockRequest struct {
	Id  uint64 `json:"id"`
	Key string `json:"key"`
}
