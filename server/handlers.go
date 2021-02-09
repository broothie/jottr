package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var whitespaceMatcher = regexp.MustCompile(`^\s+$`)

type Jot struct {
	ID       string      `json:"id" firestore:"id"`
	Body     string      `json:"body" firestore:"body"`
	Contents interface{} `json:"contents" firestore:"contents"`
}

func (j Jot) LinkText() string {
	if whitespaceMatcher.MatchString(j.Body) {
		return j.ID
	}

	return strings.Split(j.Body, "\n")[0]
}

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	cookieJotIDs := s.getJotIDs(r)
	docs, err := s.db.Collection("jots").
		Where("id", "in", cookieJotIDs).
		Documents(r.Context()).
		GetAll()
	if err != nil {
		s.log.Err(err, "failed to get recent jots")
	}

	jots := make([]Jot, 0, len(docs))
	foundJotIDs := make(map[string]struct{})
	for _, doc := range docs {
		var jot Jot
		if err := doc.DataTo(&jot); err != nil {
			s.log.Err(err, "failed to deserialize jot doc")
			continue
		}

		foundJotIDs[jot.ID] = struct{}{}
		jots = append(jots, jot)
	}

	for _, cookieJotID := range cookieJotIDs {
		if _, jotIDFound := foundJotIDs[cookieJotID]; !jotIDFound {
			s.removeJotID(w, r, cookieJotID)
		}
	}

	s.render.HTML(w, http.StatusOK, "home", jots)
}

func (s *Server) newJot(w http.ResponseWriter, r *http.Request) {
	id := newID()
	if _, err := s.db.Collection("jots").Doc(id).Set(r.Context(), Jot{ID: id}); err != nil {
		s.Error(w, err, "failed to create new jot", http.StatusInternalServerError)
		return
	}

	noCache(w)
	http.Redirect(w, r, fmt.Sprintf("/jot/%s", id), http.StatusPermanentRedirect)
}

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

func (s *Server) destroyJot(w http.ResponseWriter, r *http.Request) {
	id := httprouter.ParamsFromContext(r.Context()).ByName("id")

	if _, err := s.db.Collection("jots").Doc(id).Delete(r.Context()); err != nil {
		s.Error(w, err, "failed to delete jot", http.StatusInternalServerError)
		return
	}

	s.removeJotID(w, r, id)
}

func noCache(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
}
