package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var users = map[string]string{
	"shay1": "password",
	"shay2": "password",
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	expectedPassword, userExists := users[creds.Username]
	if !userExists || expectedPassword != creds.Password {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	log.Printf("User %s signed in", creds.Username)

	w.WriteHeader(http.StatusOK)
}

func main() {
	err := InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v\n", err)
	}
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/signin", signinHandler)
	r.HandleFunc("/ws", handleConnections)
	r.HandleFunc("/save_area", saveAreaHandler).Methods("POST")
	r.HandleFunc("/update_area", updateAreaHandler).Methods("POST")
	r.HandleFunc("/areas_in_bounds", getAreasInBoundsHandler).Methods("GET")
	// r.Handle("/", enableCORS(r))

	go handleMessages()

	log.Println("Server started on :8000")
	err = http.ListenAndServe(":8000", r)
	if err != nil {
		log.Fatalf("Server error: %v\n", err)
	}
}

// func enableCORS(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Access-Control-Allow-Origin", "*")
// 		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
// 		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

// 		if r.Method == "OPTIONS" {
// 			w.WriteHeader(http.StatusOK)
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }
