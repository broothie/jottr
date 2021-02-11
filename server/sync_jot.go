package server

import (
	"bufio"
	"encoding/json"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/julienschmidt/httprouter"
)

func (s *Server) syncJot(w http.ResponseWriter, r *http.Request) {
	id := httprouter.ParamsFromContext(r.Context()).ByName("id")

	var jot Jot
	if err := json.NewDecoder(bufio.NewReader(r.Body)).Decode(&jot); err != nil {
		s.render.JSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}

	updates := []firestore.Update{{Path: "body", Value: jot.Body}, {Path: "contents", Value: jot.Contents}}
	if _, err := s.db.Collection("jots").Doc(id).Update(r.Context(), updates); err != nil {
		s.render.JSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	s.render.JSON(w, http.StatusOK, jot)
}
