package main

import (
    "log"
    "net/http"
)

func main() {
    // Create a new ServeMux
    mux := http.NewServeMux()
    
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