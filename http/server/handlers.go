package server

import (
	"encoding/json"
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
func ToErrorResponse(err error) responses.ErrorResponse {
	return responses.ErrorResponse{
		ErrorName: vanerrors.GetName(err),
		Error:     vanerrors.GetMessage(err),
	}
}

// It creates a new user
func SingUpHandler(w http.ResponseWriter, r *http.Request, DB db_cfg.DataBase) error {
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
	usr, err := req.SingUp(DB)

	if err != nil {
		// Checks error variants
		if user_service.IsServerError(err) {

			w.WriteHeader(http.StatusInternalServerError)
		} else {

			w.WriteHeader(http.StatusBadRequest)
		}

		// Writes data
		err := json.NewEncoder(w).Encode(responses.SingUpResponseError{
			ErrorResponse: ToErrorResponse(err),
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
	// Farming
	amount, usr, err := user_service.Farm(u.Id, DB)

	if err != nil {
		// Checks error variants
		if user_service.IsServerError(err) {

			w.WriteHeader(http.StatusInternalServerError)
		} else {

			w.WriteHeader(http.StatusBadRequest)
		}

		// Writes data
		err := json.NewEncoder(w).Encode(responses.SingUpResponseError{
			ErrorResponse: ToErrorResponse(err),
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

func BuyStocksHandler(w http.ResponseWriter, r *http.Request, u user_cfg.User, DB db_cfg.DataBase) error {
	// Gets body
	var req requests.BuyStocksRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		// Set status
		w.WriteHeader(http.StatusBadRequest)

		// Creates an error
		resp := vanerrors.NewSimple(InvalidBody)

		// Writes data
		err := json.NewEncoder(w).Encode(responses.BlockResponseError{
			ErrorResponse: ToErrorResponse(resp),
		})

		return err
	}

	// Buying stocks
	usr, err := user_service.BuyStocks(u.Id, req.Num, DB)

	if err != nil {
		// Checks error variants
		if user_service.IsServerError(err) {

			w.WriteHeader(http.StatusInternalServerError)
		} else {

			w.WriteHeader(http.StatusBadRequest)
		}

		// Writes data
		err := json.NewEncoder(w).Encode(responses.SingUpResponseError{
			ErrorResponse: ToErrorResponse(err),
		})

		return err
	}

	// Converts user
	resp := ToResponseUser(*usr)

	// Sends data
	err = json.NewEncoder(w).Encode(responses.BuyStocksResponseOK{
		User: resp,
	})

	return err
}
