package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "sync"
    "time"

    "github.com/gorilla/mux"
    httpSwagger "github.com/swaggo/http-swagger"
    _ "github.com/7A1T/simple-webserver-go/docs"
)

// User represents a user in the system
// @Description User model with basic information
type User struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

// UserStore manages user data with thread-safe operations
type UserStore struct {
    mu      sync.RWMutex
    users   map[int]User
    nextID  int
}

// AddUser creates a new user and returns its ID
// @Summary Create a new user
// @Description Adds a new user to the system and returns the assigned user ID
// @Tags users
// @Accept json
// @Produce json
// @Param user body User true "User to create"
// @Success 201 {object} map[string]int "User created successfully"
// @Failure 400 {string} string "Invalid user data"
// @Router /users [post]
func (us *UserStore) AddUser(user User) int {
    us.mu.Lock()
    defer us.mu.Unlock()
    
    user.ID = us.nextID
    user.CreatedAt = time.Now()
    us.users[us.nextID] = user
    us.nextID++
    
    return user.ID
}

// GetUser retrieves a user by ID
// @Summary Get a user by ID
// @Description Retrieves a specific user from the system
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} User "User found"
// @Failure 404 {string} string "User not found"
// @Router /users/{id} [get]
func (us *UserStore) GetUser(id int) (User, bool) {
    us.mu.RLock()
    defer us.mu.RUnlock()
    
    user, exists := us.users[id]
    return user, exists
}

// @title User Management API
// @version 1.0
// @description A simple user management API with Swagger documentation
// @host localhost:8080
// @BasePath /
func main() {
    // Initialize user store
    userStore := &UserStore{
        users:  make(map[int]User),
        nextID: 1,
    }
    
    // Use gorilla mux for routing
    r := mux.NewRouter()
    
    // Middleware
    r.Use(loggingMiddleware)
    
    // Health check endpoint
    // @Summary Health Check
    // @Description Checks if the server is running
    // @Tags system
    // @Success 200 {string} string "Server is operational"
    // @Router /health [get]
    r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("server is okely dokely!"))
    }).Methods("GET")
    
    // Swagger UI endpoint
    r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
    
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
    
    // Start the server
    port := ":8080"
    fmt.Printf("Server starting on port %s\n", port)
    fmt.Printf("Swagger UI available at http://localhost%s/swagger/index.html\n", port)
    log.Fatal(http.ListenAndServe(port, r))
}

// Logging middleware to log each request
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s %s", r.Method, r.URL.Path, time.Now().Format(time.RFC3339))
        next.ServeHTTP(w, r)
    })
}
