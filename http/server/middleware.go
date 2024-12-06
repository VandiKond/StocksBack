package server

import (
	"encoding/json"
	"net/http"

	"github.com/VandiKond/StocksBack/config/db_cfg"
	"github.com/VandiKond/StocksBack/config/requests"
	"github.com/VandiKond/StocksBack/config/responses"
	"github.com/VandiKond/StocksBack/config/user_cfg"
	"github.com/VandiKond/vanerrors"
)

// The errors
const (
	UserNotFound  = "user not found"
	WrongPassword = "wrong password"
)

// function with singing in
type SingInHandlerFunc func(w http.ResponseWriter, r *http.Request, u user_cfg.User, DB db_cfg.DataBase) error

// Sings in
func SingInMiddleware(next SingInHandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, DB db_cfg.DataBase) error {
		req := requests.SingInRequest{}
		err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			// Set status
			w.WriteHeader(http.StatusBadRequest)

			// Creates an error
			resp := vanerrors.NewSimple(InvalidBody)

			// Writes data
			err := json.NewEncoder(w).Encode(responses.SingUpResponseError{
				ErrorResponse: ToErrorResponse(resp),
			})

			return err
		}

		ok, usr, err := req.User.SingIn(DB)
		if err != nil {
			// Set status
			w.WriteHeader(http.StatusInternalServerError)

			// Creates an error
			resp := vanerrors.NewSimple(UserNotFound)

			// Writes data
			err := json.NewEncoder(w).Encode(responses.SingUpResponseError{
				ErrorResponse: ToErrorResponse(resp),
			})

			return err
		}
		if !ok {
			// Set status
			w.WriteHeader(http.StatusUnauthorized)

			// Creates an error
			resp := vanerrors.NewSimple(WrongPassword)

			// Writes data
			err := json.NewEncoder(w).Encode(responses.SingUpResponseError{
				ErrorResponse: ToErrorResponse(resp),
			})

			return err
		}

		return next(w, r, *usr, DB)
	}
}
