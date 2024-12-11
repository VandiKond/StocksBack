package server

import (
	"fmt"
	"net/http"

	"github.com/VandiKond/StocksBack/config/db_cfg"
	"github.com/VandiKond/StocksBack/config/responses"
	"github.com/VandiKond/StocksBack/pkg/logger"
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
		// Main page
		"/": handler.MainHandler,
		// Sing uo
		"/singup": CheckMethodMiddleware(http.MethodPost, handler.SingUpHandler),

		// Stocks and solids
		"/buy_stocks": CheckMethodMiddleware(http.MethodPatch, handler.SingInMiddleware(handler.BuyStocksHandler)),
		"/farm":       CheckMethodMiddleware(http.MethodPatch, handler.SingInMiddleware(handler.FarmHandler)),

		// Name and password
		"/upd_name":     CheckMethodMiddleware(http.MethodPatch, handler.SingInMiddleware(handler.UpdateNameHandler)),
		"/upd_password": CheckMethodMiddleware(http.MethodPatch, handler.SingInMiddleware(handler.UpdatePasswordHandler)),

		// Block
		"/block":   CheckMethodMiddleware(http.MethodPatch, handler.SingInMiddleware(handler.BlockHandler)),
		"/unblock": CheckMethodMiddleware(http.MethodPatch, handler.SingInMiddleware(handler.UnblockHandler)),

		// Get
		"/get": CheckMethodMiddleware(http.MethodGet, handler.GetHandler),
	}

	return &handler
}

// The server
type Server struct {
	http.Server
}

// Creates a new server
func NewServer(handler http.Handler, port int) *Server {
	return &Server{http.Server{Addr: fmt.Sprint(":", port), Handler: handler}}
}

// Serve
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// The header
	w.Header().Add("Content-Type", "application/json")

	// Gets the handler func
	fn, ok := h.funcs[r.URL.Path]
	if !ok {

		// Not found
		responses.ErrorResponse{
			ErrorName: "not found",
		}.
			SendJson(w, http.StatusNotFound)

		return
	}

	// Runs the handler
	fn(w, r)
}

// Runs server
func (s *Server) Run() error {
	err := s.ListenAndServe()
	return err
}
