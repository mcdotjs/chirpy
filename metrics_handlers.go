package main

import (
	"fmt"
	"net/http"
)

func (c *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	count := c.fileserverHits.Load()
	fmt.Fprintf(w, "Hits: %d", count)
}

func (c *apiConfig) resetMetricsHandler(w http.ResponseWriter, r *http.Request) {
	c.db.DeleteAllUsers(r.Context())
	c.fileserverHits.Store(0)
	count := c.fileserverHits.Load()
	fmt.Fprintf(w, "Hits: %d", count)
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("incremeting")
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
