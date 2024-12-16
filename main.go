package main

import (
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/gorilla/mux" // more sophisticated routing
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, Welcome to my Go Web Server!")
}

func main() {

    // Use gorilla mux for more advanced routing
    r := mux.NewRouter()
    
    // middleware
    r.Use(loggingMiddleware)
    
    // Health check endpoint
    r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
      w.WriteHeader(http.StatusOK)
      w.Write([]byte("server is okely dokely!"))
    }).Methods("GET") // restricted to get requests only

    // start the server
    port := ":8080"
    fmt.Printf("Server starting on port %s\n", port)
    log.Fatal(http.ListenAndServe(port, r))
}

// Logging middleware to log each request
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s %s", r.Method, r.URL.Path, time.Now().Format(time.RFC3339))
        next.ServeHTTP(w, r)
    })
}
