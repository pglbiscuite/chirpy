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

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Increment the counter
        cfg.fileserverHits.Add(1)
        // Call the next handler
        next.ServeHTTP(w, r)
    })
}

// Add these methods outside of main
func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
    // Get the current count
    count := cfg.fileserverHits.Load()
    // Write the response
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(200)
    w.Write([]byte(fmt.Sprintf("Hits: %d", count)))
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
    // Reset the counter
    cfg.fileserverHits.Store(0)
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(200)
    w.Write([]byte("Counter reset"))
}

func myHealthzHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(200)
    w.Write([]byte("OK"))
}

func main() {
    // Create a new ServeMux
    mux := http.NewServeMux()

    // Create an instance of apiConfig
    apiCfg := apiConfig{}
    
    // Add a file server handler for the root path
    // This will serve files from the current directory
    fileServer := http.FileServer(http.Dir("."))
    mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))
    mux.HandleFunc("/healthz", myHealthzHandler)
    mux.HandleFunc("/metrics", apiCfg.metricsHandler)
    mux.HandleFunc("/reset", apiCfg.resetHandler)
    
    // Create an http.Server struct with the mux and port 8080
    server := http.Server{
        Addr:    ":8080",
        Handler: mux,
    }
    
    // Start the server with ListenAndServe
    log.Println("Server starting on port 8080...")
    err := server.ListenAndServe()
    if err != nil {
        log.Fatal("Server error:", err)
    }
}