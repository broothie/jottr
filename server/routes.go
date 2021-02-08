package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s *Server) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(s.newJot)
	router.HandlerFunc(http.MethodGet, "/jot/:id", s.showJot)
	router.HandlerFunc(http.MethodPut, "/jot/:id/sync", s.syncJot)

	return router
}
