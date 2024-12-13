package handler

import (
	"net/http"

	"github.com/vandi37/StocksBack/config/db_cfg"
	"github.com/vandi37/StocksBack/http/api"
	"github.com/vandi37/StocksBack/pkg/logger"
	"github.com/vandi37/vanerrors"
)

// The handler func
type HandlerFunc func(w http.ResponseWriter, r *http.Request, DB db_cfg.DataBase) int

// The handler
type Handler struct {
	logger *logger.Logger
	db     db_cfg.DataBase
	funcs  map[string]http.HandlerFunc
}

// Created a new handler
func NewHandler(db db_cfg.DataBase, logger *logger.Logger) *Handler {
	// Creating handler
	handler := Handler{
		logger: logger,
		db:     db,
	}

	// Adding functions
	handler.funcs = map[string]http.HandlerFunc{
		// Sign uo
		"/signup": handler.CheckMethodMiddleware(http.MethodPost, handler.SignUpHandler),

		// Stocks and solids
		"/buy":  handler.CheckMethodMiddleware(http.MethodPatch, handler.AuthorizationMiddleware(handler.BuyStocksHandler)),
		"/farm": handler.CheckMethodMiddleware(http.MethodPatch, handler.AuthorizationMiddleware(handler.FarmHandler)),

		// Name and password
		"/change/name":     handler.CheckMethodMiddleware(http.MethodPatch, handler.AuthorizationMiddleware(handler.UpdateNameHandler)),
		"/change/password": handler.CheckMethodMiddleware(http.MethodPatch, handler.AuthorizationMiddleware(handler.UpdatePasswordHandler)),

		// Block
		"/block":   handler.KeyMiddleware(handler.CheckMethodMiddleware(http.MethodPatch, handler.AuthorizationMiddleware(handler.BlockHandler))),
		"/unblock": handler.KeyMiddleware(handler.CheckMethodMiddleware(http.MethodPatch, handler.AuthorizationMiddleware(handler.UnblockHandler))),

		// Get
		"/get": handler.CheckMethodMiddleware(http.MethodGet, handler.GetHandler),
	}

	return &handler
}

// Serve
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// The header
	w.Header().Add("Content-Type", "application/json")

	// Gets the handler func
	fn, ok := h.funcs[r.URL.Path]
	if !ok {

		// Not found
		err := api.SendErrorResponse(w, http.StatusNotFound, vanerrors.NewSimple(NotFound))
		if err != nil {
			h.logger.Errorln(err)
			return
		}

		return
	}

	// Runs the handler
	fn(w, r)
}
