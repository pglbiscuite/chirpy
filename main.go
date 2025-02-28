package main

import (
    "log"
    "net/http"
)

func main() {
    // Create a new ServeMux
    mux := http.NewServeMux()
    
    // Add a file server handler for the root path
    // This will serve files from the current directory
    fileServer := http.FileServer(http.Dir("."))
    mux.Handle("/app/", http.StripPrefix("/app", fileServer))
    mux.HandleFunc("/healthz", myHealthzHandler)
    
    
    
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


func myHealthzHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8") // 1. Sets the header
    w.WriteHeader(200)                                          // 2. Sets the status code
    w.Write([]byte("OK"))                                       // 3. Writes the body
}