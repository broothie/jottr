package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Jot struct {
	Body string `json:"body" firestore:"body"`
}

func (s *Server) newJot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := newID()
		if _, err := s.db.Collection("jots").Doc(id).Set(r.Context(), Jot{}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/jot/%s", id), http.StatusPermanentRedirect)
	}
}

func (s *Server) showJot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		s.render.HTML(w, http.StatusOK, "jots/show", map[string]interface{}{
			"id":            id,
			"body":          template.HTML(doc.Data()["body"].(string)),
			"save_delay_ms": s.config.SaveDelayMs,
		})
	}
}

func (s *Server) syncJot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := httprouter.ParamsFromContext(r.Context()).ByName("id")

		var jot Jot
		if err := json.NewDecoder(bufio.NewReader(r.Body)).Decode(&jot); err != nil {
			s.render.JSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
			return
		}

		if _, err := s.db.Collection("jots").Doc(id).Set(r.Context(), jot); err != nil {
			s.render.JSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
			return
		}

		s.render.JSON(w, http.StatusOK, map[string]string{"body": ""})
	}
}
