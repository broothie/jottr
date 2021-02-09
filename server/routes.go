package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s *Server) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.RedirectHandler("/home", http.StatusPermanentRedirect)
	router.Handler(http.MethodGet, "/", http.RedirectHandler("/new", http.StatusPermanentRedirect))
	router.HandlerFunc(http.MethodGet, "/home", s.home)
	router.HandlerFunc(http.MethodGet, "/new", s.newJot)
	router.HandlerFunc(http.MethodGet, "/jot/:id", s.showJot)
	router.HandlerFunc(http.MethodPut, "/jot/:id/sync", s.syncJot)

	return router
}
