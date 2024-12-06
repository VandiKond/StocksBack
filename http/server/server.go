// TODO Write comments
package server

import (
	"net/http"

	"github.com/VandiKond/StocksBack/config/db_cfg"
	"github.com/VandiKond/StocksBack/pkg/logger"
)

// The handler func
type HandlerFunc func(w http.ResponseWriter, r *http.Request, DB db_cfg.DataBase) error

// The server
type Server struct {
	DB     db_cfg.DataBase
	Logger logger.Logger
	Funcs  map[string]HandlerFunc
}

// Creates a new server
func NewServer(Logger logger.Logger, DB db_cfg.DataBase) *Server {
	return &Server{
		DB:     DB,
		Logger: Logger,
		Funcs: map[string]HandlerFunc{
			"/singup": SingUpHandler,
			"/farm":   SingInMiddleware(FarmHandler),
		},
	}
}

// serve
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fn, ok := s.Funcs[r.URL.Path]
	if !ok {
		http.NotFound(w, r)
		return
	}
	fn(w, r, s.DB)
}

// Runs server
func (s *Server) Run(port string) {
	http.ListenAndServe(":"+port, s)
}
