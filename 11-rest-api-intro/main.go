package main

import (
	"encoding/json"
	"net/http"
	"sync"
)

// User struct represents a user in the system
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var (
	users    = make(map[int]User)
	nextID   = 1
	usersMtx sync.RWMutex
)

// Seed data
func init() {
	users[1] = User{ID: 1, Name: "John Doe", Email: "john@example.com"}
	users[2] = User{ID: 2, Name: "Jane Doe", Email: "jane@example.com"}
	users[3] = User{ID: 3, Name: "Alice Smith", Email: "alice@example.com"}
}

// Get all users
func getUsers(w http.ResponseWriter, r *http.Request) {
	usersMtx.RLock()
	defer usersMtx.RUnlock()

	userList := make([]User, 0, len(users))
	for _, user := range users {
		userList = append(userList, user)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userList)
}

// Get user by ID
func getUser(w http.ResponseWriter, r *http.Request) {
	usersMtx.RLock()
	defer usersMtx.RUnlock()

	id := r.URL.Query().Get("id")
	for _, user := range users {
		if strconv.Itoa(user.ID) == id {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(user)
			return
		}
	}
	http.Error(w, "User not found", http.StatusNotFound)
}

// Create user
func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	usersMtx.Lock()
	defer usersMtx.Unlock()

	user.ID = nextID
	nextID++
	users[user.ID] = user

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Update user
func updateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	usersMtx.Lock()
	defer usersMtx.Unlock()

	if _, exists := users[user.ID]; !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	users[user.ID] = user
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// Delete user
func deleteUser(w http.ResponseWriter, r *http.Request) {
	usersMtx.Lock()
	defer usersMtx.Unlock()

	id := r.URL.Query().Get("id")
	for userID := range users {
		if strconv.Itoa(userID) == id {
			delete(users, userID)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "User not found", http.StatusNotFound)
}

// Main function
func main() {
	http.HandleFunc("/users", getUsers)
	http.HandleFunc("/user", getUser)
	http.HandleFunc("/user/create", createUser)
	http.HandleFunc("/user/update", updateUser)
	http.HandleFunc("/user/delete", deleteUser)
	
	// Start the server
	http.ListenAndServe(":8080", nil)
}