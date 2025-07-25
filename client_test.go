package sarvam

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	client := NewClient("test")
	assert.Equal(t, client.apiKey, "test")
}

func TestSetBaseURL(t *testing.T) {
	client := NewClient("test")
	client.SetBaseURL("https://api.sarvam.ai")
	assert.Equal(t, client.baseURL, "https://api.sarvam.ai")
}

func TestMakeHTTPRequest(t *testing.T) {
	testApiKey := "test" + time.Now().Format("20060102150405")
	httpTestServer := newMockServer(testApiKey)

	client := NewClient(testApiKey)
	client.SetBaseURL(httpTestServer.URL)

	response, err := client.makeHTTPRequest("GET", httpTestServer.URL+"/v1/test", nil)
	assert.NoError(t, err)
	assert.Equal(t, response.StatusCode, 200)
}

type handler struct {
	apiKey string
}

func newMockServer(apiKey string) *httptest.Server {
	h := &handler{apiKey: apiKey}
	s := httptest.NewServer(http.HandlerFunc(h.handleHTTPRequest))
	return s
}

func (h *handler) handleHTTPRequest(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("api-subscription-key")
	if authHeader != h.apiKey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func TestBool(t *testing.T) {
	b := Bool(true)
	assert.Equal(t, *b, true)
}

func TestFloat64(t *testing.T) {
	f := Float64(1.0)
	assert.Equal(t, *f, 1.0)
}

func TestInt(t *testing.T) {
	i := Int(1)
	assert.Equal(t, *i, 1)
}

func TestHTTPError(t *testing.T) {
	e := &HTTPError{
		StatusCode: http.StatusUnauthorized,
		Message:    "Unauthorized",
	}
	assert.Equal(t, e.Error(), "status code: 401, message: Unauthorized")
}
