package api

import (
	"context"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// RegisterHealth registers the health check route.
// If client is non-nil, GET /health pings MongoDB and returns 503 when unhealthy.
func RegisterHealth(mux *http.ServeMux, client *mongo.Client) {
	if client != nil {
		mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
			defer cancel()
			if err := client.Ping(ctx, readpref.Primary()); err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte("unhealthy"))
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
		return
	}
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
}
