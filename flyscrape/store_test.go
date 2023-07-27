package flyscrape_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"flyscrape/flyscrape"

	"github.com/stretchr/testify/require"
)

func TestFetch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(404)
		w.Write([]byte(`{"error": "not found"}`))
	}))
	defer srv.Close()

	result, err := flyscrape.Fetch(http.DefaultClient, srv.URL)
	require.NoError(t, err)
	require.Equal(t, result.URL, srv.URL)
	require.Equal(t, result.StatusCode, 404)
	require.Equal(t, result.Header.Get("content-type"), "application/json")
	require.Equal(t, result.Body, []byte(`{"error": "not found"}`))
}

func TestScrapeStore(t *testing.T) {
	store, err := flyscrape.NewScrapeStore("test.db")
	require.NoError(t, err)
	defer store.Close()

	err = store.Migrate()
	require.NoError(t, err)

	result := &flyscrape.ScrapeResult{
		URL:        "http://example.com/page",
		Body:       []byte(`<html><body>Example</body></html>`),
		StatusCode: 200,
		Header:     http.Header{"Content-Type": []string{"text/html"}},
		Timestamp:  time.Now().UTC(),
	}

	err = store.InsertScrapeResult(result)
	require.NoError(t, err)
}
