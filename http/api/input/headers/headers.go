package headers

import "github.com/vandi37/StocksBack/pkg/user_service"

type Key struct {
	user_service.SignInKey
}

type Authorization struct {
	user_service.SignInUser
}
