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
	UserNotFound          = "user not found"
	WrongPassword         = "wrong password"
	WrongKey              = "wrong key"
	NoAutorotationHeaders = "no authorization headers"
	InvalidHeader         = "invalid header"
)

// function with singing in
type HandlerFuncUser func(w http.ResponseWriter, r *http.Request, u user_cfg.User)

// Sings in
func (h *Handler) SingInMiddleware(next HandlerFuncUser) http.HandlerFunc {
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
				responses.SingInResponseError{
					ErrorResponse: ToErrorResponse(resp),
				}.
					SendJson(w, http.StatusBadRequest)

				return
			}
			usr, err := keyData.SingInWithKey(h.db)
			if err != nil {

				// Checks error variants
				var status = http.StatusBadRequest
				if user_service.IsServerError(err) {

					status = http.StatusInternalServerError
				}

				// Writes data
				responses.SingInResponseError{
					ErrorResponse: ToErrorResponse(err),
				}.
					SendJson(w, status)
				return
			}
			next(w, r, *usr)
			return
		}

		// Gets header
		key = r.Header.Get("Autorotation")

		if key == "" {

			// Creates an error
			resp := vanerrors.NewSimple(NoAutorotationHeaders)

			// Writes data
			responses.SingInResponseError{
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
			responses.SingInResponseError{
				ErrorResponse: ToErrorResponse(resp),
			}.
				SendJson(w, http.StatusBadRequest)
			return
		}

		ok, usr, err := authData.SingIn(h.db)
		if err != nil {
			// Checks error variants
			var status = http.StatusBadRequest
			if user_service.IsServerError(err) {

				status = http.StatusInternalServerError
			}

			// Writes data
			responses.SingUpResponseError{
				ErrorResponse: ToErrorResponse(err),
			}.
				SendJson(w, status)
			return
		}
		if !ok {

			// Creates an error
			resp := vanerrors.NewSimple(WrongPassword)

			// Writes data
			responses.SingUpResponseError{
				ErrorResponse: ToErrorResponse(resp),
			}.
				SendJson(w, http.StatusUnauthorized)
			return
		}

		next(w, r, *usr)
	}
}

// Checks the method
func (h *Handler) CheckMethodMiddleware(method string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Checking method
		if r.Method != method {

			// Creates an error
			resp := vanerrors.NewSimple(WrongMethod, fmt.Sprintf("method %s is not allowed, allowed method: %s", r.Method, http.MethodPatch))

			// Writes data
			responses.SingUpResponseError{
				ErrorResponse: ToErrorResponse(resp),
			}.
				SendJson(w, http.StatusBadRequest)
			return
		}
		next(w, r)
	}
}
