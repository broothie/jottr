package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) showJot(w http.ResponseWriter, r *http.Request) {
	id := httprouter.ParamsFromContext(r.Context()).ByName("id")

	doc, err := s.db.Collection("jots").Doc(id).Get(r.Context())
	if err != nil {
		if status.Code(err) == codes.NotFound {
			s.render.HTML(w, http.StatusNotFound, "jots/404", nil)
			return
		}

		s.render.HTML(w, http.StatusInternalServerError, "error", nil)
		return
	}

	noCache(w)
	s.addJotID(w, r, id)
	s.render.HTML(w, http.StatusOK, "jots/show", doc.Data())
}
