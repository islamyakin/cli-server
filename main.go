package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var validUsername = "admin"
var validPassword = "password"

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
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if loginReq.Username == validUsername && loginReq.Password == validPassword {
		sapaLogin := "Login berhasil, Halo" + " " + loginReq.Username
		resp := LoginResponse{Token: validToken, Message: sapaLogin}
		json.NewEncoder(w).Encode(resp)
		log.Println("Successful login fon user:", loginReq.Username)
	} else {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
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

	if msgReq.Token != validToken {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	log.Println("Received message from client:", msgReq.Message)
	fmt.Fprintf(w, "Logged message: %s", msgReq.Message)
}
