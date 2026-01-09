package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"backend/internal/auth"
	"backend/internal/db"
)

type Handler struct {
	DB *sql.DB
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// check if request is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var credentials struct {
		UserID   string `json:"UserID"`
		Password string `json:"Password"`
	}
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}
	// check if user already exists
	userAlreadyExists, err := db.CheckUserExists(h.DB, credentials.UserID)
	if err != nil {
		http.Error(w, "Server Error with DB", http.StatusInternalServerError)
		return
	}
	if userAlreadyExists == true {
		http.Error(w, "User Already Exists", http.StatusBadRequest)
		return
	}
	// create hash
	passwordHash, err := auth.HashPassword(credentials.Password)
	if err != nil {
		http.Error(w, "Server Error with Password", http.StatusInternalServerError)
		return
	}
	// store creds in DB
	err = db.AddUser(h.DB, credentials.UserID, passwordHash)
	if err != nil {
		http.Error(w, "Server Error with DB storage", http.StatusInternalServerError)
		return
	}
	token, err := auth.GenerateToken(passwordHash)
	if err != nil {
		http.Error(w, "Server Error with generating token", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	json.NewEncoder(w).Encode(map[string]string{
			"token": token,
	})
	
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	// check if request is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var credentials struct {
		UserID   string `json:"UserID"`
		Password string `json:"Password"`
	}
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}
	// check if user already exists
	userAlreadyExists, err := db.CheckUserExists(h.DB, credentials.UserID)
	if err != nil {
		http.Error(w, "Server Error with DB", http.StatusInternalServerError)
		return
	}
	if userAlreadyExists == false {
		http.Error(w, "User Doesn't Exists", http.StatusBadRequest)
		return
	}
	// create hash
	passwordHash, err := auth.HashPassword(credentials.Password)
	if err != nil {
		http.Error(w, "Server Error with Password", http.StatusInternalServerError)
		return
	}
	// check creds with DB
	correct_hash,err:=db.GetUserHash(h.DB,credentials.UserID)
	if err != nil {
		http.Error(w, "Server Error checking DB creds", http.StatusInternalServerError)
		return
	}
	if passwordHash!=correct_hash{
		http.Error(w, "Wrong Password for given UserID", http.StatusBadRequest)
		return
	}
	
	//Generate token given password is correct
	token, err := auth.GenerateToken(passwordHash)
	if err != nil {
		http.Error(w, "Server Error with generating token", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	json.NewEncoder(w).Encode(map[string]string{
			"token": token,
	})
	
}
