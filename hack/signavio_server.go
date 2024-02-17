// Small hack that implements the Signavio API
package main

import (
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func main() {
	logger.Info("Starting signavio-mock-server")

	mux := http.NewServeMux()

	mux.HandleFunc("/p/login", handleLogin)
	mux.HandleFunc("/p/user", handleGetUsers)
	mux.HandleFunc("/api/v2/user-jobs/delete", handleDeleteUser)
	mux.HandleFunc("/api/v2/user-jobs/", handleJobStatus)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	logger.Info("Login request", "method", r.Method, "url", r.URL.String(), "user", r.FormValue("name"), "tenant", r.FormValue("tenant"))
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Check the credentials from www-form-urlencoded body
	// If they are correct, set the cookie
	if r.FormValue("name") != "test@test.com" || r.FormValue("password") != "test" || r.FormValue("tenant") != "test" || r.FormValue("tokenonly") != "true" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// Set the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "1234",
		Path:     "/",
		HttpOnly: false,
	})
	w.WriteHeader(http.StatusOK)
}

func handleGetUsers(w http.ResponseWriter, r *http.Request) {
	logger.Info("GetUsers request", "method", r.Method, "url", r.URL.String())
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Check the cookie
	cookie, err := r.Cookie("token")
	if err != nil || cookie.Value != "1234" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// Return the users, no GZIP for simplicity
	w.Write([]byte(`[
		{
		  "rel": "user",
		  "href": "/user/0123456789",
		  "rep": {
			"lastName": "Doe",
			"mail": "john.doe@test.com",
			"language": "de",
			"isAlwaysGrantedAccess": false,
			"title": "",
			"principal": "john.doe@test.com",
			"isFirstStart": true,
			"firstName": "John",
			"deleted": false,
			"isGuestUser": false,
			"phone": "",
			"name": "John Doe",
			"company": "",
			"state": "active"
		  }
		},  { "rel": "info", "href": "/user", "rep": { "size": 1 } }
		]`))
}

func handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Check the cookie
	cookie, err := r.Cookie("token")
	if err != nil || cookie.Value != "1234" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// Return the job ID
	w.Write([]byte(`"0123456789"`))
}

func handleJobStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Check the cookie
	cookie, err := r.Cookie("token")
	if err != nil || cookie.Value != "1234" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// Return the job status in a randomized way
	rand.New(rand.NewSource(time.Now().UnixNano()))
	status := []string{"RUNNING", "SCHEDULED", "COMPLETED", "ERROR"}
	w.Write([]byte(`"` + status[rand.Intn(4)] + `"`))
}
