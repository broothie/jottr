package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_setJotID(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()
	setJotID(recorder, request, "new-dddid")

	var jotCookie *http.Cookie
	for _, cookie := range recorder.Result().Cookies() {
		if cookie.Name == jotIDsCookieName {
			jotCookie = cookie
			return
		}
	}

	assert.NotNil(t, jotCookie)
	assert.Equal(t, "new-id", jotCookie.Value)
}
