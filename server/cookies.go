package server

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/broothie/jottr/set"
)

const (
	jotIDsSep        = ","
	jotIDsCookieName = "jot_ids"
)

func getJotIDs(r *http.Request) []string {
	cookie, _ := r.Cookie(jotIDsCookieName)
	if cookie == nil {
		return nil
	}

	log.Println("get cookie", cookie.Value)
	return strings.Split(cookie.Value, jotIDsSep)
}

func setJotID(w http.ResponseWriter, r *http.Request, id string) {
	ids := set.New(id)
	if cookie, _ := r.Cookie(jotIDsCookieName); cookie != nil {
		for _, id := range strings.Split(cookie.Value, jotIDsSep) {
			ids.Insert(id)
		}
	}

	idSlice := make([]string, len(ids))
	i := 0
	for id := range ids {
		idSlice[i] = id.(string)
		i++
	}

	log.Println("set cookie", jotIDsCookieName, idSlice)
	http.SetCookie(w, &http.Cookie{
		Name:     jotIDsCookieName,
		Value:    strings.Join(idSlice, jotIDsSep),
		Path:     "*",
		Expires:  time.Now().Add(30 * 24 * time.Hour), // 30 days
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}
