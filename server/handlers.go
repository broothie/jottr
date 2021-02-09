package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"regexp"

	"cloud.google.com/go/firestore"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var breakRegexp = regexp.MustCompile(`<br/?>`)

type Jot struct {
	ID   string `json:"id" firestore:"id"`
	Body string `json:"body" firestore:"body"`
}

func (j Jot) LinkText() string {
	if j.Body == "" {
		return j.ID
	}

	lines := breakRegexp.Split(j.Body, -1)
	firstLine := lines[0]
	firstLineStripped := strip.StripTags(firstLine)
	return firstLineStripped
}

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	docs, err := s.db.Collection("jots").
		Where("id", "in", s.getJotIDs(r)).
		Documents(r.Context()).
		GetAll()
	if err != nil {
		s.log.Err(err, "failed to get recent jots")
	}

	jots := make([]Jot, 0, len(docs))
	for _, doc := range docs {
		var jot Jot
		if err := doc.DataTo(&jot); err != nil {
			s.log.Err(err, "failed to deserialize jot doc")
			continue
		}

		jots = append(jots, jot)
	}

	s.render.HTML(w, http.StatusOK, "home", jots)
}

func (s *Server) newJot(w http.ResponseWriter, r *http.Request) {
	id := newID()
	if _, err := s.db.Collection("jots").Doc(id).Set(r.Context(), Jot{ID: id}); err != nil {
		s.Error(w, err, "failed to create new jot", "failed to create new jot", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
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

	s.addJotID(w, r, id)
	s.render.HTML(w, http.StatusOK, "jots/show", map[string]interface{}{
		"id":            id,
		"body":          template.HTML(doc.Data()["body"].(string)),
		"save_delay_ms": s.config.SaveDelayMs,
	})
}

func (s *Server) syncJot(w http.ResponseWriter, r *http.Request) {
	id := httprouter.ParamsFromContext(r.Context()).ByName("id")

	var requestPayload map[string]string
	if err := json.NewDecoder(bufio.NewReader(r.Body)).Decode(&requestPayload); err != nil {
		s.render.JSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}

	updates := []firestore.Update{{Path: "body", Value: requestPayload["body"]}}
	if _, err := s.db.Collection("jots").Doc(id).Update(r.Context(), updates); err != nil {
		s.render.JSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	s.render.JSON(w, http.StatusOK, requestPayload)
}
