package sarvam

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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

	response, err := client.makeJsonHTTPRequest("GET", httpTestServer.URL+"/v1/test", nil)
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

func TestHTTPErrorWithAPIError(t *testing.T) {
	e := &HTTPError{
		StatusCode: http.StatusBadRequest,
		Message:    "body.model : Input should be 'saarika:v1', 'saarika:v2', 'saarika:v2.5' or 'saarika:flash'",
		Code:       "invalid_request_error",
		RequestID:  "20250801_1f62ae94-d102-4513-a31d-bbb0718052dc",
	}
	expected := "status code: 400, code: invalid_request_error, message: body.model : Input should be 'saarika:v1', 'saarika:v2', 'saarika:v2.5' or 'saarika:flash', request_id: 20250801_1f62ae94-d102-4513-a31d-bbb0718052dc"
	assert.Equal(t, e.Error(), expected)
}

func TestParseAPIError(t *testing.T) {
	// Test with valid API error response
	apiErrorJSON := `{
		"error": {
			"message": "body.model : Input should be 'saarika:v1', 'saarika:v2', 'saarika:v2.5' or 'saarika:flash'",
			"code": "invalid_request_error",
			"request_id": "20250801_1f62ae94-d102-4513-a31d-bbb0718052dc"
		}
	}`

	resp := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       io.NopCloser(strings.NewReader(apiErrorJSON)),
	}

	err := parseAPIError(resp)
	httpErr, ok := err.(*HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpErr.StatusCode)
	assert.Equal(t, "body.model : Input should be 'saarika:v1', 'saarika:v2', 'saarika:v2.5' or 'saarika:flash'", httpErr.Message)
	assert.Equal(t, "invalid_request_error", httpErr.Code)
	assert.Equal(t, "20250801_1f62ae94-d102-4513-a31d-bbb0718052dc", httpErr.RequestID)
}
