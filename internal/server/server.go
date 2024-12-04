// TODO Write comments
package server

import (
	"net/http"

	"github.com/VandiKond/StocksBack/config/db_cfg"
)

type Server struct {
	DB    db_cfg.DataBase
	Funcs map[string]http.HandlerFunc
}

func NewServer(DB db_cfg.DataBase) *Server {
	return &Server{
		DB: DB,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fn, ok := s.Funcs[r.URL.Path]
	if !ok {
		http.NotFound(w, r)
		return
	}
	fn(w, r)
}

func RunServer(server *Server, port string) {
	http.ListenAndServe(":"+port, server)
}
