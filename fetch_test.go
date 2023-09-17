package flyscrape_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/philippta/flyscrape"
	"github.com/stretchr/testify/require"
)

func TestFetchFetch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("foobar"))
	}))

	fetch := flyscrape.Fetch()

	req, _ := http.NewRequest("GET", srv.URL, nil)
	html, err := fetch(req)
	require.NoError(t, err)
	require.Equal(t, html, "foobar")
}

func TestFetchCachedFetch(t *testing.T) {
	numcalled := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		numcalled++
		w.Write([]byte("foobar"))
	}))

	fetch := flyscrape.CachedFetch(flyscrape.Fetch())

	req, _ := http.NewRequest("GET", srv.URL, nil)
	html, err := fetch(req)
	require.NoError(t, err)
	require.Equal(t, html, "foobar")

	html, err = fetch(req)
	require.NoError(t, err)
	require.Equal(t, html, "foobar")

	require.Equal(t, 1, numcalled)
}

func TestFetchProxiedFetch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, r.URL.String(), "http://example.com/foo")
		w.Write([]byte("foobar"))
	}))

	fetch := flyscrape.ProxiedFetch(srv.URL)

	req, _ := http.NewRequest("GET", "http://example.com/foo", nil)
	html, err := fetch(req)
	require.NoError(t, err)
	require.Equal(t, html, "foobar")
}
