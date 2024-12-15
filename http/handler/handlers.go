package handler

import (
	"encoding/json"
	"net/http"

	"github.com/vandi37/StocksBack/config/user_cfg"
	"github.com/vandi37/StocksBack/http/api"
	"github.com/vandi37/StocksBack/http/api/input/requests"
	"github.com/vandi37/StocksBack/http/api/responses"
	"github.com/vandi37/StocksBack/pkg/user_service"
	"github.com/vandi37/vanerrors"
)

// The errors
const (
	WrongMethod = "wrong method"
	InvalidBody = "invalid body"
	NotFound    = "not found"
)

// User to response user
func ToResponseUser(usr user_cfg.User) responses.User {
	return responses.User{
		Id:           usr.Id,
		Name:         usr.Name,
		SolidBalance: usr.SolidBalance,
		StockBalance: usr.StockBalance,
		IsBlocked:    usr.IsBlocked,
		LastFarming:  usr.LastFarming,
		CreatedAt:    usr.CreatedAt,
	}
}

// It creates a new user
func (h *Handler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	// Gets body
	req := requests.SignUp{}
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {

		// Creates an error
		resp := vanerrors.NewSimple(InvalidBody)

		// Writes data
		err = api.SendErrorResponse(w, http.StatusBadRequest, resp)
		if err != nil {
			h.logger.Errorln(err)
			return
		}

		return
	}

	// Signs up
	usr, err := req.SignUp(h.db)

	if err != nil {

		// Writes data
		err = api.SendErrorResponse(w, user_service.GetCode(err), err)
		if err != nil {
			h.logger.Errorln(err)
			return
		}

		h.logger.Warnf("unable to Sign up, reason: %v", err)
		return
	}

	// Converts user
	resp := ToResponseUser(*usr)

	// Sends data
	err = api.SendOkResponse(w, responses.SignUp{User: resp}, responses.SignUpType)
	if err != nil {
		h.logger.Errorln(err)
		return
	}

	h.logger.Printf("Signup: %v", *usr)
}

// Farms
func (h *Handler) FarmHandler(w http.ResponseWriter, r *http.Request, u user_cfg.User) {
	// Farming
	amount, usr, err := user_service.Farm(u.Id, h.db)

	if err != nil {
		// Writes data
		err = api.SendErrorResponse(w, user_service.GetCode(err), err)
		if err != nil {
			h.logger.Errorln(err)
			return
		}

		h.logger.Warnf("%v unable to farm, reason: %v", u, err)

		return
	}

	// Converts user
	resp := ToResponseUser(*usr)

	// Sends data
	err = api.SendOkResponse(w, responses.Farm{User: resp, Amount: amount}, responses.FarmType)
	if err != nil {
		h.logger.Errorln(err)
		return
	}

	h.logger.Printf("farm (%d) : %v", amount, *usr)
}

// Buy stocks
func (h *Handler) BuyStocksHandler(w http.ResponseWriter, r *http.Request, u user_cfg.User) {
	// Gets body
	var req requests.BuyStocks
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {

		// Creates an error
		resp := vanerrors.NewSimple(InvalidBody)

		// Writes data
		err = api.SendErrorResponse(w, http.StatusBadRequest, resp)
		if err != nil {
			h.logger.Errorln(err)
			return
		}

		return
	}

	// Buying stocks
	usr, err := user_service.BuyStocks(u.Id, req.Num, h.db)

	if err != nil {
		// Writes data
		err = api.SendErrorResponse(w, user_service.GetCode(err), err)
		if err != nil {
			h.logger.Errorln(err)
			return
		}

		h.logger.Warnf("%v unable to buy stocks, reason: %v", u, err)

		return
	}

	// Converts user
	resp := ToResponseUser(*usr)

	// Sends data
	err = api.SendOkResponse(w, responses.BuyStocks{
		User: resp,
	}, "buy-stocks")
	if err != nil {
		h.logger.Errorln(err)
		return
	}

	h.logger.Printf("buy stocks (%d) : %v", req.Num, *usr)
}

// Update name
func (h *Handler) UpdateNameHandler(w http.ResponseWriter, r *http.Request, u user_cfg.User) {
	// Gets body
	var req requests.UpdateName
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {

		// Creates an error
		resp := vanerrors.NewSimple(InvalidBody)

		// Writes data
		err = api.SendErrorResponse(w, http.StatusBadRequest, resp)
		if err != nil {
			h.logger.Errorln(err)
			return
		}

		return
	}

	usr, err := user_service.UpdateName(u.Id, req.Name, h.db)

	if err != nil {
		// Writes data
		err = api.SendErrorResponse(w, user_service.GetCode(err), err)
		if err != nil {
			h.logger.Errorln(err)
			return
		}

		h.logger.Warnf("%v unable to update name, reason: %v", u, err)

		return
	}

	// Converts user
	resp := ToResponseUser(*usr)

	// Sends data
	err = api.SendOkResponse(w, responses.UpdateName{User: resp}, responses.UpdateNameType)
	if err != nil {
		h.logger.Errorln(err)
		return
	}

	h.logger.Printf("update name (was %s) : %v", u.Name, *usr)

}

// Update password
func (h *Handler) UpdatePasswordHandler(w http.ResponseWriter, r *http.Request, u user_cfg.User) {
	// Gets body
	var req requests.UpdatePassword
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {

		// Creates an error
		resp := vanerrors.NewSimple(InvalidBody)

		// Writes data
		err = api.SendErrorResponse(w, http.StatusBadRequest, resp)
		if err != nil {
			h.logger.Errorln(err)
			return
		}

		return
	}

	usr, err := user_service.UpdatePassword(u.Id, req.Password, h.db)

	if err != nil {
		// Writes data
		err = api.SendErrorResponse(w, user_service.GetCode(err), err)
		if err != nil {
			h.logger.Errorln(err)
			return
		}

		h.logger.Warnf("%v unable to update password, reason: %v", u, err)

		return
	}

	// Converts user
	resp := ToResponseUser(*usr)

	// Sends data
	err = api.SendOkResponse(w, responses.UpdatePassword{User: resp}, responses.UpdatePasswordType)
	if err != nil {
		h.logger.Errorln(err)
		return
	}

	h.logger.Printf("update password: %v", *usr)
}

// Block user
func (h *Handler) BlockHandler(w http.ResponseWriter, r *http.Request, u user_cfg.User) {
	// Gets body
	var req requests.Block
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		// Creates an error
		resp := vanerrors.NewSimple(InvalidBody)

		// Writes data
		err = api.SendErrorResponse(w, http.StatusBadRequest, resp)
		if err != nil {
			h.logger.Errorln(err)
			return
		}

		return
	}

	usr, err := user_service.Block(u.Id, h.db)

	if err != nil {
		// Writes data
		err = api.SendErrorResponse(w, user_service.GetCode(err), err)
		if err != nil {
			h.logger.Errorln(err)
			return
		}

		h.logger.Warnf("%v unable to block, reason: %v", u, err)

		return
	}

	// Converts user
	resp := ToResponseUser(*usr)

	// Sends data
	err = api.SendOkResponse(w, responses.Block{User: resp}, responses.BlockType)
	if err != nil {
		h.logger.Errorln(err)
		return
	}

	h.logger.Printf("block: %v", *usr)
}

// Unlock user
func (h *Handler) UnblockHandler(w http.ResponseWriter, r *http.Request, u user_cfg.User) {
	// Gets body
	var req requests.Unblock
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {

		// Creates an error
		resp := vanerrors.NewSimple(InvalidBody)

		// Writes data
		err = api.SendErrorResponse(w, http.StatusBadRequest, resp)
		if err != nil {
			h.logger.Errorln(err)
			return
		}
		return
	}

	usr, err := user_service.Unblock(u.Id, h.db)

	if err != nil {
		// Writes data
		err = api.SendErrorResponse(w, user_service.GetCode(err), err)
		if err != nil {
			h.logger.Errorln(err)
			return
		}

		h.logger.Warnf("%v unable to unblock, reason: %v", u, err)

		return
	}

	// Converts user
	resp := ToResponseUser(*usr)

	// Sends data
	err = api.SendOkResponse(w, responses.Unblock{User: resp}, responses.UnblockType)
	if err != nil {
		h.logger.Errorln(err)
		return
	}

	h.logger.Printf("unblock: %v", *usr)
}

// Get's user
func (h *Handler) GetHandler(w http.ResponseWriter, r *http.Request) {
	// Gets body
	var req requests.Get
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {

		// Creates an error
		resp := vanerrors.NewSimple(InvalidBody)

		// Writes data
		err = api.SendErrorResponse(w, http.StatusBadRequest, resp)
		if err != nil {
			h.logger.Errorln(err)
			return
		}

		return
	}

	usr, err := user_service.Get(req.Id, h.db)

	if err != nil {
		// Writes data
		err = api.SendErrorResponse(w, user_service.GetCode(err), err)
		if err != nil {
			h.logger.Errorln(err)
			return
		}

		h.logger.Warnf("user %d not got, reason:", req.Id, err)

		return
	}

	// Converts user
	resp := ToResponseUser(*usr)

	// Sends data
	err = api.SendOkResponse(w, responses.Get{User: resp}, responses.GetType)
	if err != nil {
		h.logger.Errorln(err)
		return
	}

	h.logger.Printf("sended user: %v", *usr)
}
