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
		// Sign uo
		"/signup": CheckMethodMiddleware(http.MethodPost, handler.SignUpHandler),

		// Stocks and solids
		"/buy": CheckMethodMiddleware(http.MethodPatch, handler.SignInMiddleware(handler.BuyStocksHandler)),
		"/farm":       CheckMethodMiddleware(http.MethodPatch, handler.SignInMiddleware(handler.FarmHandler)),

		// Name and password
		"/change/name":     CheckMethodMiddleware(http.MethodPatch, handler.SignInMiddleware(handler.UpdateNameHandler)),
		"/change/password": CheckMethodMiddleware(http.MethodPatch, handler.SignInMiddleware(handler.UpdatePasswordHandler)),

		// Block
		"/block":   CheckMethodMiddleware(http.MethodPatch, handler.SignInMiddleware(handler.BlockHandler)),
		"/unblock": CheckMethodMiddleware(http.MethodPatch, handler.SignInMiddleware(handler.UnblockHandler)),

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
