package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s *Server) destroyJot(w http.ResponseWriter, r *http.Request) {
	id := httprouter.ParamsFromContext(r.Context()).ByName("id")

	if _, err := s.db.Collection("jots").Doc(id).Delete(r.Context()); err != nil {
		s.Error(w, err, "failed to delete jot", http.StatusInternalServerError)
		return
	}

	s.removeJotID(w, r, id)
}
