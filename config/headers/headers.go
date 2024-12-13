package headers

import "github.com/vandi37/StocksBack/pkg/user_service"

type KeyHeader struct {
	user_service.SignInKey
}

type AuthorizationHeader struct {
	user_service.SignInUser
}
