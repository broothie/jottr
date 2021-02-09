package server

import (
	"net/http"
	"strings"
	"time"

	"github.com/broothie/jottr/logger"
)

const (
	jotIDsSep        = "|"
	jotIDsCookieName = "jot_ids"
)

func (s *Server) getJotIDs(r *http.Request) []string {
	cookie, _ := r.Cookie(jotIDsCookieName)
	if cookie == nil {
		s.log.Info("no cookie")
		return nil
	}

	s.log.Info("get cookie", logger.Field("value", cookie.Value))
	return strings.Split(cookie.Value, jotIDsSep)
}

func (s *Server) addJotID(w http.ResponseWriter, r *http.Request, id string) {
	ids := map[string]struct{}{id: {}}
	if cookie, _ := r.Cookie(jotIDsCookieName); cookie != nil {
		for _, id := range strings.Split(cookie.Value, jotIDsSep) {
			ids[id] = struct{}{}
		}
	}

	idSlice := make([]string, len(ids))
	i := 0
	for id := range ids {
		idSlice[i] = id
		i++
	}

	s.log.Info("set cookie", logger.Field("values", idSlice))
	http.SetCookie(w, &http.Cookie{
		Name:     jotIDsCookieName,
		Value:    strings.Join(idSlice, jotIDsSep),
		Path:     "/",
		Expires:  time.Now().Add(365 * 24 * time.Hour), // 1 year
		Secure:   s.config.Environment == "production",
		HttpOnly: false,
		SameSite: http.SameSiteStrictMode,
	})
}
