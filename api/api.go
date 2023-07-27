package api

import (
	"encoding/json"
	"net/http"

	"github.com/alexedwards/flow"
)

type ScrapeRequest struct {
	URL  string         `json:"url"`
	Data map[string]any `json:"data"`
}

type ScrapeResponse struct {
	URL  string `json:"url"`
	Data any    `json:"data"`
}

//go:generate moq -out api_service_mock_test.go . Service
type Service interface {
	ScrapeURL(url string, params map[string]any) (any, error)
}

func NewHandler(svc Service) http.Handler {
	h := &Handler{
		router: flow.New(),
		svc:    svc,
	}
	h.routes()
	return h
}

type Handler struct {
	router *flow.Mux
	svc    Service
}

func (h *Handler) routes() {
	h.router.HandleFunc("/scrape", h.handleScrape, "POST")
}

func (h *Handler) handleScrape(w http.ResponseWriter, r *http.Request) {
	var req ScrapeRequest
	if err := decodeRequest(r, &req); err != nil {
		respondErr(w, http.StatusBadRequest, err)
		return
	}

	result, err := h.svc.ScrapeURL(req.URL, req.Data)
	if err != nil {
		respondErr(w, http.StatusInternalServerError, err)
		return
	}

	respond(w, ScrapeResponse{
		URL:  req.URL,
		Data: result,
	})
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func decodeRequest(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func respond(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(v)
}

func respondErr(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	})
}
