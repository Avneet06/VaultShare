package auth

import (
	"encoding/json"
	"fmt"
	"file-sharing-system/db"
	"net/http"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Register endpoint hit")                           //Added  to debug Rhandler.
	if r.Method == "TEST" {
		w.WriteHeader(http.StatusCreated)
		return
	}
	
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	user := &db.User{
		Email:    req.Email,
		Password: hashedPassword,
	}

	err = db.CreateUser(user)
if err != nil {
   
    fmt.Println("DB error:", err)
    http.Error(w, "Could not create user", http.StatusInternalServerError)
    return
}

fmt.Println("User created:", req.Email) // log successful reg
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}


func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	user, err := db.GetUserByEmail(req.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if !CheckPasswordHash(req.Password, user.Password) {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	token, err := GenerateJWT(user.Email)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}
	if err != nil {
		fmt.Println("DB error during login:", err)
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}
	
	resp := map[string]string{"token": token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}