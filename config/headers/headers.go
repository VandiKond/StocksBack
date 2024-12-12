package headers

import "github.com/VandiKond/StocksBack/pkg/user_service"

type KeyHeader struct {
	user_service.SignInKey
}

type AuthorizationHeader struct {
	user_service.SignInUser
}
