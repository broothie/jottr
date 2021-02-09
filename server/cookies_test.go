package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_setJotID(t *testing.T) {
	t.Run("single cookie", func(t *testing.T) {
		// Set up
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()

		// Exercise
		setJotID(recorder, request, "new-id")

		// Verify
		var jotCookie *http.Cookie
		for _, cookie := range recorder.Result().Cookies() {
			if cookie.Name == jotIDsCookieName {
				jotCookie = cookie
				break
			}
		}

		require.NotNil(t, jotCookie)
		assert.Equal(t, "new-id", jotCookie.Value)
	})

	t.Run("multiple cookies", func(t *testing.T) {
		// Set up
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		request.AddCookie(&http.Cookie{Name: jotIDsCookieName, Value: "first-id"})
		recorder := httptest.NewRecorder()

		// Exercise
		setJotID(recorder, request, "second-id")

		// Verify
		var jotCookie *http.Cookie
		for _, cookie := range recorder.Result().Cookies() {
			if cookie.Name == jotIDsCookieName {
				jotCookie = cookie
				break
			}
		}

		require.NotNil(t, jotCookie)
		assert.Equal(t, "second-id|first-id", jotCookie.Value)
	})
}

func Test_getJotIDs(t *testing.T) {

}
