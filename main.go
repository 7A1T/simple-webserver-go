package main

import (
    "fmt"
    "log"
    "net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, Welcome to my Go Web Server!")
}

func main() {
    // Define routes
    http.HandleFunc("/", helloHandler)

    // Specify the port to listen on
    port := ":8080"
    fmt.Printf("Server starting on port %s\n", port)
    
    // Start the server
    err := http.ListenAndServe(port, nil)
    if err != nil {
        log.Fatal("Error starting server: ", err)
    }
}
