package server

import (
	"net/http"
	"strings"
	"time"

	"github.com/broothie/jottr/logger"
)

const (
	jotIDsCookieName = "jot_ids"
	jotIDsSeparator  = "|"
)

func (s *Server) getJotIDs(r *http.Request) []string {
	cookie, _ := r.Cookie(jotIDsCookieName)
	if cookie == nil {
		s.log.Info("no cookie")
		return nil
	}

	s.log.Info("get cookie", logger.Field("value", cookie.Value))
	return strings.Split(cookie.Value, jotIDsSeparator)
}

func (s *Server) addJotID(w http.ResponseWriter, r *http.Request, id string) {
	idSet := map[string]struct{}{id: {}}
	for _, id := range s.getJotIDs(r) {
		idSet[id] = struct{}{}
	}

	ids := make([]string, len(idSet))
	i := 0
	for id := range idSet {
		ids[i] = id
		i++
	}

	s.setJotIDsCookie(w, strings.Join(ids, jotIDsSeparator))
}

func (s *Server) removeJotID(w http.ResponseWriter, r *http.Request, id string) {
	var ids []string
	for _, cookieID := range s.getJotIDs(r) {
		if cookieID != id {
			ids = append(ids, cookieID)
		}
	}

	s.setJotIDsCookie(w, strings.Join(ids, jotIDsSeparator))
}

func (s *Server) setJotIDsCookie(w http.ResponseWriter, jotIDs string) {
	s.log.Info("set cookie", logger.Field("jot_ids", jotIDs))
	http.SetCookie(w, &http.Cookie{
		Name:     jotIDsCookieName,
		Value:    jotIDs,
		Path:     "/",
		Expires:  time.Now().Add(365 * 24 * time.Hour), // 1 year
		Secure:   s.config.Environment == "production",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}
