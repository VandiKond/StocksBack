package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/VandiKond/StocksBack/config/headers"
	"github.com/VandiKond/StocksBack/config/responses"
	"github.com/VandiKond/StocksBack/config/user_cfg"
	"github.com/VandiKond/StocksBack/pkg/user_service"
	"github.com/VandiKond/vanerrors"
)

// The errors
const (
	UserNotFound           = "user not found"
	WrongPassword          = "wrong password"
	WrongKey               = "wrong key"
	NoAuthorizationHeaders = "no authorization headers"
	InvalidHeader          = "invalid header"
)

// function with Signing in
type HandlerFuncUser func(w http.ResponseWriter, r *http.Request, u user_cfg.User)

// Signs in
func (h *Handler) SignInMiddleware(next HandlerFuncUser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Gets header
		key := r.Header.Get("Key")

		if key != "" {
			// Gets header data
			var keyData headers.KeyHeader
			err := json.Unmarshal([]byte(key), &keyData)

			if err != nil {

				// Creates an error
				resp := vanerrors.NewSimple(InvalidHeader)

				// Writes data
				responses.SignInResponseError{
					ErrorResponse: ToErrorResponse(resp),
				}.
					SendJson(w, http.StatusBadRequest)

				return
			}
			usr, err := keyData.SignInWithKey(h.db)
			if err != nil {

				// Checks error variants
				var status = http.StatusBadRequest
				if user_service.IsServerError(err) {

					status = http.StatusInternalServerError
				}

				// Writes data
				responses.SignInResponseError{
					ErrorResponse: ToErrorResponse(err),
				}.
					SendJson(w, status)

				h.logger.Warnf("unable to login with key, reason: %v", err)

				return
			}
			next(w, r, *usr)
			return
		}

		// Gets header
		key = r.Header.Get("Authorization")

		if key == "" {

			// Creates an error
			resp := vanerrors.NewSimple(NoAuthorizationHeaders)

			// Writes data
			responses.SignInResponseError{
				ErrorResponse: ToErrorResponse(resp),
			}.
				SendJson(w, http.StatusUnauthorized)
			return
		}
		// Gets header data
		var authData headers.AuthorizationHeader
		err := json.Unmarshal([]byte(key), &authData)

		if err != nil {

			// Creates an error
			resp := vanerrors.NewSimple(InvalidHeader)

			// Writes data
			responses.SignInResponseError{
				ErrorResponse: ToErrorResponse(resp),
			}.
				SendJson(w, http.StatusBadRequest)
			return
		}

		ok, usr, err := authData.SignIn(h.db)
		if err != nil {
			// Checks error variants
			var status = http.StatusBadRequest
			if user_service.IsServerError(err) {

				status = http.StatusInternalServerError
			}

			// Writes data
			responses.SignUpResponseError{
				ErrorResponse: ToErrorResponse(err),
			}.
				SendJson(w, status)

			h.logger.Warnf("unable to login, reason: %v", err)

			return
		}
		if !ok {

			// Creates an error
			resp := vanerrors.NewSimple(WrongPassword)

			// Writes data
			responses.SignUpResponseError{
				ErrorResponse: ToErrorResponse(resp),
			}.
				SendJson(w, http.StatusUnauthorized)
			return
		}

		next(w, r, *usr)
	}
}

// Checks the method
func CheckMethodMiddleware(method string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Checking method
		if r.Method != method {

			// Creates an error
			resp := vanerrors.NewSimple(WrongMethod, fmt.Sprintf("method %s is not allowed, allowed method: %s", r.Method, method))

			// Writes data
			responses.SignUpResponseError{
				ErrorResponse: ToErrorResponse(resp),
			}.
				SendJson(w, http.StatusBadRequest)
			return
		}
		next(w, r)
	}
}
