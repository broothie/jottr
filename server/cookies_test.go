package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/broothie/jottr/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer_addJotID(t *testing.T) {
	t.Run("single cookie", func(t *testing.T) {
		// Set up
		server := &Server{log: logger.New()}
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()

		// Exercise
		server.addJotID(recorder, request, "new-id")

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
		server := &Server{log: logger.New()}
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		request.AddCookie(&http.Cookie{Name: jotIDsCookieName, Value: "first-id"})
		recorder := httptest.NewRecorder()

		// Exercise
		server.addJotID(recorder, request, "second-id")

		// Verify
		var jotCookie *http.Cookie
		for _, cookie := range recorder.Result().Cookies() {
			if cookie.Name == jotIDsCookieName {
				jotCookie = cookie
				break
			}
		}

		require.NotNil(t, jotCookie)
		assert.Contains(t, jotCookie.Value, "first-id")
		assert.Contains(t, jotCookie.Value, "second-id")
	})
}

func TestServer_removeJotID(t *testing.T) {
	// Set up
	server := &Server{log: logger.New()}
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.AddCookie(&http.Cookie{Name: jotIDsCookieName, Value: "first-id|second-id|third-id"})
	recorder := httptest.NewRecorder()

	// Exercise
	server.removeJotID(recorder, request, "second-id")

	// Verify
	var jotCookie *http.Cookie
	for _, cookie := range recorder.Result().Cookies() {
		if cookie.Name == jotIDsCookieName {
			jotCookie = cookie
			break
		}
	}

	require.NotNil(t, jotCookie)
	assert.NotContains(t, jotCookie.Value, "second-id")
}

func TestServer_getJotIDs(t *testing.T) {
	// Set up
	server := &Server{log: logger.New()}
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.AddCookie(&http.Cookie{Name: jotIDsCookieName, Value: "first-id|second-id|third-id"})

	// Exercise
	jotIDs := server.getJotIDs(request)

	// Verify
	assert.Len(t, jotIDs, 3)
	assert.Contains(t, jotIDs, "first-id")
	assert.Contains(t, jotIDs, "second-id")
	assert.Contains(t, jotIDs, "third-id")
}
