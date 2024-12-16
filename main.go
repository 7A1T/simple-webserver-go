package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
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

    // User creation endpoint
    r.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
        var user User
        err := json.NewDecoder(r.Body).Decode(&user)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        id := userStore.AddUser(user)
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(map[string]int{"id": id})
    }).Methods("POST")

    // User retrieval endpoint
    r.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id, _ := strconv.Atoi(vars["id"])
        
        user, exists := userStore.GetUser(id)
        if !exists {
            http.Error(w, "User not found", http.StatusNotFound)
            return
        }

        json.NewEncoder(w).Encode(user)
    }).Methods("GET")
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
