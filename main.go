package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

type MessageRequest struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

var validToken = "secret-token"

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/log", logHandler)
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		log.Println("Failed to decode request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Println("Received login request:", loginReq)

	ok, data, err := AuthUsingLDAP(loginReq.Username, loginReq.Password)

	if !ok {
		log.Println("LDAP authentication failed for user:", loginReq.Username)
		http.Error(w, "Invalid username/password", http.StatusUnauthorized)
		return
	}
	if err != nil {
		log.Println("LDAP authentication error:", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	tokenString, err := CreateToken(data.FullName)
	if err != nil {
		log.Println("Failed to create token:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	resp := LoginResponse{
		Token:   tokenString,
		Message: fmt.Sprintf("Welcome %s", data.FullName),
	}

	log.Println("Sending response:", resp)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("Failed to encode response:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var msgReq MessageRequest
	if err := json.NewDecoder(r.Body).Decode(&msgReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err := verifyToken(msgReq.Token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Invalid token")
		return
	}

	//if msgReq.Token != tokens {
	//	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	//	return
	//}

	log.Println("Received message from client:", msgReq.Message)
	fmt.Fprintf(w, "Logged message: %s", msgReq.Message)
}
