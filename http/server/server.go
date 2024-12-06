// TODO Write comments
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
	logger logger.Logger
	db     db_cfg.DataBase
	funcs  map[string]http.HandlerFunc
}

func NewHandler(db db_cfg.DataBase, logger logger.Logger) *Handler {
	handler := Handler{
		logger: logger,
		db:     db,
	}
	handler.funcs = map[string]http.HandlerFunc{
		"/singup":     handler.CheckMethodMiddleware(http.MethodPost, handler.SingUpHandler),
		"/farm":       handler.CheckMethodMiddleware(http.MethodPatch, handler.SingInMiddleware(handler.FarmHandler)),
		"/buy_stocks": handler.CheckMethodMiddleware(http.MethodPatch, handler.SingInMiddleware(handler.BuyStocksHandler)),
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

// serve
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	fn, ok := h.funcs[r.URL.Path]
	if !ok {
		responses.ErrorResponse{
			ErrorName: "not found",
		}.
			SendJson(w, http.StatusNotFound)
		return
	}
	fn(w, r)
}

// Runs server
func (s *Server) Run() error {
	err := s.ListenAndServe()
	return err
}
