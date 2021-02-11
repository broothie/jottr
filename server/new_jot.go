package server

import (
	"fmt"
	"net/http"
)

func (s *Server) newJot(w http.ResponseWriter, r *http.Request) {
	id := newID()
	if _, err := s.db.Collection("jots").Doc(id).Set(r.Context(), Jot{ID: id}); err != nil {
		s.Error(w, err, "failed to create new jot", http.StatusInternalServerError)
		return
	}

	noCache(w)
	http.Redirect(w, r, fmt.Sprintf("/jot/%s", id), http.StatusPermanentRedirect)
}
