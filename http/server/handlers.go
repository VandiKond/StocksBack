package server

import (
	"encoding/json"
	"net/http"

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
func (h *Handler) SingUpHandler(w http.ResponseWriter, r *http.Request) {
	// Gets body
	req := requests.SingUpRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {

		// Creates an error
		resp := vanerrors.NewSimple(InvalidBody)

		// Writes data
		responses.SingUpResponseError{
			ErrorResponse: ToErrorResponse(resp),
		}.
			SendJson(w, http.StatusBadRequest)
		return
	}

	// Sings up
	usr, err := req.SingUp(h.db)

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

	// Converts user
	resp := ToResponseUser(*usr)

	// Sends data
	json.NewEncoder(w).Encode(responses.SingUpResponseOK{
		User: resp,
	})
}

// Farms
func (h *Handler) FarmHandler(w http.ResponseWriter, r *http.Request, u user_cfg.User) {
	// Farming
	amount, usr, err := user_service.Farm(u.Id, h.db)

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

	// Converts user
	resp := ToResponseUser(*usr)

	// Sends data
	json.NewEncoder(w).Encode(responses.FarmResponseOK{
		User:   resp,
		Amount: amount,
	})
}

// Buy stocks
func (h *Handler) BuyStocksHandler(w http.ResponseWriter, r *http.Request, u user_cfg.User) {
	// Gets body
	var req requests.BuyStocksRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {

		// Creates an error
		resp := vanerrors.NewSimple(InvalidBody)

		// Writes data
		responses.BuyStocksResponseError{
			ErrorResponse: ToErrorResponse(resp),
		}.
			SendJson(w, http.StatusBadRequest)
		return
	}

	// Buying stocks
	usr, err := user_service.BuyStocks(u.Id, req.Num, h.db)

	if err != nil {
		// Checks error variants
		var status = http.StatusBadRequest
		if user_service.IsServerError(err) {

			status = http.StatusInternalServerError
		}
		// Writes data
		responses.BuyStocksResponseError{
			ErrorResponse: ToErrorResponse(err),
		}.
			SendJson(w, status)
		return
	}

	// Converts user
	resp := ToResponseUser(*usr)

	// Sends data
	json.NewEncoder(w).Encode(responses.BuyStocksResponseOK{
		User: resp,
	})
}

// Update name
func (h *Handler) UpdateNameHandler(w http.ResponseWriter, r *http.Request, u user_cfg.User) {
	// Gets body
	var req requests.UpdateNameRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {

		// Creates an error
		resp := vanerrors.NewSimple(InvalidBody)

		// Writes data
		responses.UpdateNameResponseError{
			ErrorResponse: ToErrorResponse(resp),
		}.
			SendJson(w, http.StatusBadRequest)
		return
	}

	usr, err := user_service.UpdateName(u.Id, req.Name, h.db)

	if err != nil {
		// Checks error variants
		var status = http.StatusBadRequest
		if user_service.IsServerError(err) {

			status = http.StatusInternalServerError
		}
		// Writes data
		responses.UpdateNameResponseError{
			ErrorResponse: ToErrorResponse(err),
		}.
			SendJson(w, status)
		return
	}

	// Converts user
	resp := ToResponseUser(*usr)

	// Sends data
	json.NewEncoder(w).Encode(responses.UpdateNameResponseOK{
		User: resp,
	})
}

// Update password
func (h *Handler) UpdatePasswordHandler(w http.ResponseWriter, r *http.Request, u user_cfg.User) {
	// Gets body
	var req requests.UpdatePasswordRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {

		// Creates an error
		resp := vanerrors.NewSimple(InvalidBody)

		// Writes data
		responses.UpdatePasswordResponseError{
			ErrorResponse: ToErrorResponse(resp),
		}.
			SendJson(w, http.StatusBadRequest)
		return
	}

	usr, err := user_service.UpdatePassword(u.Id, req.Password, h.db)

	if err != nil {
		// Checks error variants
		var status = http.StatusBadRequest
		if user_service.IsServerError(err) {

			status = http.StatusInternalServerError
		}
		// Writes data
		responses.UpdatePasswordResponseError{
			ErrorResponse: ToErrorResponse(err),
		}.
			SendJson(w, status)
		return
	}

	// Converts user
	resp := ToResponseUser(*usr)

	// Sends data
	json.NewEncoder(w).Encode(responses.UpdatePasswordResponseOK{
		User: resp,
	})
}

// Block user
func (h *Handler) BlockHandler(w http.ResponseWriter, r *http.Request, u user_cfg.User) {
	// Gets body
	var req requests.BlockRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {

		// Creates an error
		resp := vanerrors.NewSimple(InvalidBody)

		// Writes data
		responses.BlockResponseError{
			ErrorResponse: ToErrorResponse(resp),
		}.
			SendJson(w, http.StatusBadRequest)
		return
	}

	usr, err := user_service.Block(u.Id, h.db)

	if err != nil {
		// Checks error variants
		var status = http.StatusBadRequest
		if user_service.IsServerError(err) {

			status = http.StatusInternalServerError
		}
		// Writes data
		responses.BlockResponseError{
			ErrorResponse: ToErrorResponse(err),
		}.
			SendJson(w, status)
		return
	}

	// Converts user
	resp := ToResponseUser(*usr)

	// Sends data
	json.NewEncoder(w).Encode(responses.BlockResponseOK{
		User: resp,
	})
}

// Unlock user
func (h *Handler) UnblockHandler(w http.ResponseWriter, r *http.Request, u user_cfg.User) {
	// Gets body
	var req requests.UnblockRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {

		// Creates an error
		resp := vanerrors.NewSimple(InvalidBody)

		// Writes data
		responses.UnblockResponseError{
			ErrorResponse: ToErrorResponse(resp),
		}.
			SendJson(w, http.StatusBadRequest)
		return
	}

	usr, err := user_service.Unblock(u.Id, h.db)

	if err != nil {
		// Checks error variants
		var status = http.StatusBadRequest
		if user_service.IsServerError(err) {

			status = http.StatusInternalServerError
		}
		// Writes data
		responses.UnblockResponseError{
			ErrorResponse: ToErrorResponse(err),
		}.
			SendJson(w, status)
		return
	}

	// Converts user
	resp := ToResponseUser(*usr)

	// Sends data
	json.NewEncoder(w).Encode(responses.UnblockResponseOK{
		User: resp,
	})
}
