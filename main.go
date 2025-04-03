package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (c *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	count := c.fileserverHits.Load()
	fmt.Fprintf(w, "Hits: %d", count)
}

func (c *apiConfig) resetMetricsHandler(w http.ResponseWriter, r *http.Request) {
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

func main() {

	apiCfg := &apiConfig{}
	mux := http.NewServeMux()
	port := "8080"
	filepathRoot := "."

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))

	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(200)

		html := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, apiCfg.fileserverHits.Load())
		fmt.Fprintf(w, html)
	})
	mux.HandleFunc("POST /admin/reset", apiCfg.resetMetricsHandler)
	mux.HandleFunc("GET /api/metrics", apiCfg.metricsHandler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	//log.Fatal(srv.ListenAndServe())
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("listen", err)
	}
}
