package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vandi37/StocksBack/config/user_cfg"
	"github.com/vandi37/StocksBack/http/api"
	"github.com/vandi37/StocksBack/http/api/input/headers"
	"github.com/vandi37/StocksBack/pkg/user_service"
	"github.com/vandi37/vanerrors"
)

// The errors
const (
	UserNotFound           = "user not found"
	WrongPassword          = "wrong password"
	WrongKey               = "wrong key"
	NoAuthorizationHeaders = "no authorization headers"
	InvalidHeader          = "invalid header"
	NotAllowed             = "not allowed"
)

// function with Signing in
type HandlerFuncUser func(w http.ResponseWriter, r *http.Request, u user_cfg.User)

// Checks the method
func (h *Handler) CheckMethodMiddleware(method string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Checking method
		if r.Method != method {

			// Creates an error
			resp := vanerrors.NewSimple(WrongMethod, fmt.Sprintf("method %s is not allowed, allowed method: %s", r.Method, method))

			// Writes data
			err := api.SendErrorResponse(w, http.StatusMethodNotAllowed, resp)
			if err != nil {
				h.logger.Errorln(err)
				return
			}
			return
		}
		next(w, r)
	}
}

// Signs in
func (h *Handler) AuthorizationMiddleware(checkBlock bool, next HandlerFuncUser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Gets header
		key := r.Header.Get("Key")

		if key != "" {
			// Gets header data
			var keyData headers.Key
			err := json.Unmarshal([]byte(key), &keyData)

			if err != nil {

				// Creates an error
				resp := vanerrors.NewSimple(InvalidHeader)

				// Writes data
				err = api.SendErrorResponse(w, http.StatusBadRequest, resp)
				if err != nil {
					h.logger.Errorln(err)
					return
				}

				return
			}
			usr, err := keyData.SignInWithKey(h.db)
			if err != nil {

				// Writes data
				err = api.SendErrorResponse(w, user_service.GetCode(err), err)
				if err != nil {
					h.logger.Errorln(err)
					return
				}

				h.logger.Warnf("unable to login with key, reason: %v", err)

				return
			}
			if checkBlock && usr.IsBlocked {
				err = api.SendErrorResponse(w, http.StatusForbidden, vanerrors.NewSimple(NotAllowed, "user is blocked"))
				if err != nil {
					h.logger.Errorln(err)
					return
				}
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
			err := api.SendErrorResponse(w, http.StatusUnauthorized, resp)
			if err != nil {
				h.logger.Errorln(err)
				return
			}
			return
		}
		// Gets header data
		var authData headers.Authorization
		err := json.Unmarshal([]byte(key), &authData)

		if err != nil {

			// Creates an error
			resp := vanerrors.NewSimple(InvalidHeader)

			// Writes data
			err = api.SendErrorResponse(w, http.StatusBadRequest, resp)
			if err != nil {
				h.logger.Errorln(err)
				return
			}
			return
		}

		ok, usr, err := authData.SignIn(h.db)
		if err != nil {
			// Writes data
			err = api.SendErrorResponse(w, user_service.GetCode(err), err)
			if err != nil {
				h.logger.Errorln(err)
				return
			}

			h.logger.Warnf("unable to login, reason: %v", err)

			return
		}
		if !ok {

			// Creates an error
			resp := vanerrors.NewSimple(WrongPassword)

			// Writes data
			err := api.SendErrorResponse(w, http.StatusUnauthorized, resp)
			if err != nil {
				h.logger.Errorln(err)
				return
			}
			return
		}

		if checkBlock && usr.IsBlocked {
			err = api.SendErrorResponse(w, http.StatusForbidden, vanerrors.NewSimple(NotAllowed, "user is blocked"))
			if err != nil {
				h.logger.Errorln(err)
				return
			}
			return
		}
		next(w, r, *usr)
	}
}

// Checks admin
func (h *Handler) KeyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Checking the key header (without it not allowed)
		key := r.Header.Get("Key")

		if key == "" {
			// Creates an error
			resp := vanerrors.NewSimple(InvalidHeader)

			// Writes data
			err := api.SendErrorResponse(w, http.StatusForbidden, resp)
			if err != nil {
				h.logger.Errorln(err)
				return
			}

			return
		}
	}
}
