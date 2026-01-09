package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"log"

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
		log.Println(err)
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}
	// check if user already exists
	userAlreadyExists, err := db.CheckUserExists(h.DB, credentials.UserID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if userAlreadyExists == true {
		log.Printf("UserID:%v already exists",credentials.UserID)
		http.Error(w, "UserID Already Exists", http.StatusBadRequest)
		return
	}
	// create hash
	passwordHash, err := auth.HashPassword(credentials.Password)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// store creds in DB
	err = db.AddUser(h.DB, credentials.UserID, passwordHash)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	token, err := auth.GenerateToken(credentials.UserID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if userAlreadyExists == false {
		http.Error(w, "Invalid Credentials", http.StatusBadRequest)
		return
	}
	
	// find hash of userID in DB and check if its correct
	userHash,err:=db.GetUserHash(h.DB,credentials.UserID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	
	passwordCorrect:=auth.CheckPasswordHash(credentials.Password,userHash)
	if passwordCorrect!=true{
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}
	//Generate token given password is correct
	token, err := auth.GenerateToken(credentials.UserID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	json.NewEncoder(w).Encode(map[string]string{
			"token": token,
	})
	
}
