package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/VandiKond/StocksBack/config/db_cfg"
	"github.com/VandiKond/StocksBack/config/requests"
	"github.com/VandiKond/StocksBack/config/responses"
	"github.com/VandiKond/StocksBack/config/user_cfg"
	"github.com/VandiKond/StocksBack/pkg/user_service"
	"github.com/VandiKond/vanerrors"
)

// The errors
const (
	WrongMethod = "wrong method"
	InvalidBody = "invalid body"
)

// User to response user
func ToResponseUser(usr user_cfg.User) responses.ResponseUser {
	return responses.ResponseUser{
		Id:           usr.Id,
		Name:         usr.Name,
		SolidBalance: usr.SolidBalance,
		StockBalance: usr.StockBalance,
		IsBlocked:    usr.IsBlocked,
		LastFarming:  usr.LastFarming,
		CreatedAt:    usr.CreatedAt,
	}
}

// Vanerror to response error
func ToErrorResponse(err vanerrors.VanError) responses.ErrorResponse {
	return responses.ErrorResponse{
		ErrorName: err.Name,
		Error:     err.Message,
	}
}

// It creates a new user
func SingUpHandler(w http.ResponseWriter, r *http.Request, DB db_cfg.DataBase) error {

	// Checking method
	if r.Method != http.MethodPost {
		// Set status
		w.WriteHeader(http.StatusBadRequest)

		// Creates an error
		resp := vanerrors.NewSimple(WrongMethod, fmt.Sprintf("method %s, allowed method %s", r.Method, http.MethodPost))

		// Writes data
		err := json.NewEncoder(w).Encode(responses.SingUpResponseError{
			ErrorResponse: ToErrorResponse(resp),
		})

		return err
	}

	// Gets body
	req := requests.SingUpRequest{}
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

	// Sings up
	usr, err := req.User.SingUp(DB)

	if err != nil {
		// Gets error
		resp := vanerrors.Get(err)
		if resp == nil {
			return err
		}

		// Checks error variants
		if resp.Name == user_service.ErrorGettingId || resp.Name == user_service.ErrorUpdatingUser {

			w.WriteHeader(http.StatusInternalServerError)
		} else if resp.Name == user_service.ErrorCreatingUser {

			w.WriteHeader(http.StatusBadRequest)
		}

		// Writes data
		err := json.NewEncoder(w).Encode(responses.SingUpResponseError{
			ErrorResponse: ToErrorResponse(*resp),
		})

		return err
	}

	// Converts user
	resp := ToResponseUser(*usr)

	// Sends data
	err = json.NewEncoder(w).Encode(responses.SingUpResponseOK{
		User: resp,
	})

	return err
}

func FarmHandler(w http.ResponseWriter, r *http.Request, u user_cfg.User, DB db_cfg.DataBase) error {
	amount, usr, err := user_service.Farm(u.Id, DB)

	if err != nil {
		// Gets error
		resp := vanerrors.Get(err)
		if resp == nil {
			return err
		}

		// Checks error variants
		if resp.Name == user_service.ErrorSelectingUser || resp.Name == user_service.ErrorUpdatingUser {
			w.WriteHeader(http.StatusInternalServerError)
		} else if resp.Name == user_service.ToEarlyFarming {

			w.WriteHeader(http.StatusTooManyRequests)
		}

		// Writes data
		err := json.NewEncoder(w).Encode(responses.FarmResponseError{
			ErrorResponse: ToErrorResponse(*resp),
		})

		return err
	}

	// Converts user
	resp := ToResponseUser(*usr)

	// Sends data
	err = json.NewEncoder(w).Encode(responses.FarmResponseOK{
		User:   resp,
		Amount: amount,
	})

	return err
}
