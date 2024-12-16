package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "sync"
    "time"

    "github.com/gorilla/mux" // more sophisticated routing
)

// User struct to demonstrate JSON serialization
type User struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

// In-memory user storage with thread-safe operations
type UserStore struct {
    mu    sync.RWMutex
    users map[int]User
    nextID int
}

// Methods for thread-safe user operations
func (us *UserStore) AddUser(user User) int {
    us.mu.Lock()
    defer us.mu.Unlock()
    
    user.ID = us.nextID
    user.CreatedAt = time.Now()
    us.users[us.nextID] = user
    us.nextID++
    
    return user.ID
}

func (us *UserStore) GetUser(id int) (User, bool) {
    us.mu.RLock()
    defer us.mu.RUnlock()
    
    user, exists := us.users[id]
    return user, exists
}

func main() {
    // Initialize user store
    userStore := &UserStore{
        users:  make(map[int]User),
        nextID: 1,
    }

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
