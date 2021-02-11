package server

import (
	"net/http"
)

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
