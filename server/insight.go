package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/pamelag/pika/insight"
)

// Create an insight handler
// Create a router function
// Create the route handler function that decodes the request,
// calls the service and encodes the response or encodes error

type insightHandler struct {
	ins insight.Service
}

func (h *insightHandler) router() chi.Router {
	r := chi.NewRouter()

	r.Route("/trips", func(r chi.Router) {
		r.Get("/", h.getTrips)
		r.Post("/clearCache", h.clearCache)
	})

	return r
}

func (h *insightHandler) getTrips(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	medallionIDs := r.URL.Query().Get("medallionIDs")
	tripDate := r.URL.Query().Get("tripDate")
	cache := r.URL.Query().Get("ignoreCache")

	var ignoreCache bool

	medallions := strings.Split(medallionIDs, ",")
	if strings.ToLower(cache) == "true" {
		ignoreCache = true
	}

	ts, err := h.ins.GetTripCount(medallions, tripDate, ignoreCache)
	if err != nil {
		encodeError(ctx, err, w)
		return
	}

	var response = struct {
		Trips []insight.TripData `json:"trips"`
	}{
		Trips: ts,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		encodeError(ctx, err, w)
		return
	}
}

func (h *insightHandler) clearCache(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	h.ins.ClearCache()

	var response = struct {
		Cache string `json:"cache"`
	}{
		Cache: "cleared",
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		encodeError(ctx, err, w)
		return
	}
}
