package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/VandiKond/StocksBack/config/db_cfg"
	"github.com/VandiKond/StocksBack/config/headers"
	"github.com/VandiKond/StocksBack/config/responses"
	"github.com/VandiKond/StocksBack/config/user_cfg"
	"github.com/VandiKond/StocksBack/pkg/user_service"
	"github.com/VandiKond/vanerrors"
)

// The errors
const (
	UserNotFound          = "user not found"
	WrongPassword         = "wrong password"
	WrongKey              = "wrong key"
	NoAutorotationHeaders = "no authorization headers"
	InvalidHeader         = "invalid header"
)

// function with singing in
type SingInHandlerFunc func(w http.ResponseWriter, r *http.Request, u user_cfg.User, DB db_cfg.DataBase) error

// Sings in
func SingInMiddleware(next SingInHandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, DB db_cfg.DataBase) error {
		// Gets header
		key := r.Header.Get("Key")

		if key != "" {
			// Gets header data
			var keyData headers.KeyHeader
			err := json.Unmarshal([]byte(key), &keyData)

			if err != nil {
				// Set status
				w.WriteHeader(http.StatusBadRequest)

				// Creates an error
				resp := vanerrors.NewSimple(InvalidHeader)

				// Writes data
				err := json.NewEncoder(w).Encode(responses.SingInResponseError{
					ErrorResponse: ToErrorResponse(resp),
				})

				return err
			}
			ok, usr, err := keyData.SingInWithKey(DB)
			if err != nil {
				// Gets error
				resp := vanerrors.Get(err)
				if resp == nil {
					return err
				}

				// Checks error variants
				if resp.Name == user_service.ErrorSelectingUser || resp.Name == user_service.ErrorCheckingKey {
					w.WriteHeader(http.StatusInternalServerError)
				} else if resp.Name == user_service.WrongKey {
					w.WriteHeader(http.StatusBadRequest)
				}

				// Writes data
				err := json.NewEncoder(w).Encode(responses.SingInResponseError{
					ErrorResponse: ToErrorResponse(*resp),
				})

				return err
			}
			if !ok {
				// Set status
				w.WriteHeader(http.StatusUnauthorized)

				// Creates an error
				resp := vanerrors.NewSimple(WrongPassword)

				// Writes data
				err := json.NewEncoder(w).Encode(responses.SingInResponseError{
					ErrorResponse: ToErrorResponse(resp),
				})

				return err
			}
			return next(w, r, *usr, DB)
		}

		// Gets header
		key = r.Header.Get("Autorotation")

		if key == "" {
			// Set status
			w.WriteHeader(http.StatusUnauthorized)

			// Creates an error
			resp := vanerrors.NewSimple(NoAutorotationHeaders)

			// Writes data
			err := json.NewEncoder(w).Encode(responses.SingInResponseError{
				ErrorResponse: ToErrorResponse(resp),
			})

			return err
		}
		// Gets header data
		var authData headers.AuthorizationHeader
		err := json.Unmarshal([]byte(key), &authData)

		if err != nil {
			// Set status
			w.WriteHeader(http.StatusBadRequest)

			// Creates an error
			resp := vanerrors.NewSimple(InvalidHeader)

			// Writes data
			err := json.NewEncoder(w).Encode(responses.SingInResponseError{
				ErrorResponse: ToErrorResponse(resp),
			})

			return err
		}

		ok, usr, err := authData.SingIn(DB)
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

func CheckMethodMiddleware(method string, next HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, DB db_cfg.DataBase) error {
		// Checking method
		if r.Method != method {
			// Set status
			w.WriteHeader(http.StatusMethodNotAllowed)

			// Creates an error
			resp := vanerrors.NewSimple(WrongMethod, fmt.Sprintf("method %s is not allowed, allowed method: %s", r.Method, http.MethodPatch))

			// Writes data
			err := json.NewEncoder(w).Encode(responses.SingUpResponseError{
				ErrorResponse: ToErrorResponse(resp),
			})

			return err
		}
		return next(w, r, DB)
	}
}
