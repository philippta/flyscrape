package api_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"flyscrape/api"

	"github.com/stretchr/testify/require"
)

func TestScrapeURL(t *testing.T) {
	svc := &api.ServiceMock{
		ScrapeURLFunc: func(url string, params map[string]any) (any, error) {
			return map[string]any{"foo": "bar"}, nil
		},
	}
	h := api.NewHandler(svc)

	r := httptest.NewRequest("POST", "/scrape", strings.NewReader(`{"url": "https://example.com", "data": {"foo":".foo"}}`))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	require.Equal(t, w.Result().StatusCode, http.StatusOK)
	require.Equal(t, w.Result().Header.Get("Content-Type"), "application/json")

	result := map[string]any{}
	require.NoError(t, json.NewDecoder(w.Result().Body).Decode(&result))
	require.Equal(t, result["url"].(string), "https://example.com")
	require.Equal(t, result["data"].(map[string]any)["foo"], "bar")
}

func TestScrapeURLInternalServerError(t *testing.T) {
	svc := &api.ServiceMock{
		ScrapeURLFunc: func(url string, params map[string]any) (any, error) {
			return nil, errors.New("whoops")
		},
	}
	h := api.NewHandler(svc)

	r := httptest.NewRequest("POST", "/scrape", strings.NewReader(`{"url": "https://example.com", "data": {"foo":".foo"}}`))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	require.Equal(t, w.Result().StatusCode, http.StatusInternalServerError)
	require.Equal(t, w.Result().Header.Get("Content-Type"), "application/json")
}

func TestScrapeURLBadRequest(t *testing.T) {
	svc := &api.ServiceMock{
		ScrapeURLFunc: func(url string, params map[string]any) (any, error) {
			return nil, errors.New("whoops")
		},
	}
	h := api.NewHandler(svc)

	r := httptest.NewRequest("POST", "/scrape", strings.NewReader(`{"}`))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	require.Equal(t, w.Result().StatusCode, http.StatusBadRequest)
	require.Equal(t, w.Result().Header.Get("Content-Type"), "application/json")
}
