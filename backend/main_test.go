package main

import (
	"backend/internal/auth"
	"backend/internal/db"
	"backend/internal/handlers"
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestDB(t *testing.T) (*sql.DB, *sql.DB, func()) {
	t.Helper()

	// Use an in-memory SQLite database for testing
	userDatabase, err := db.ConnectUserDB()
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.CreateUserTable(userDatabase); err != nil {
		t.Fatalf("failed to create users table: %v", err)
	}

	msgDatabase, err := db.ConnectUserDB()
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.CreateMsgTable(msgDatabase); err != nil {
		t.Fatalf("failed to create messages table: %v", err)
	}

	// Teardown function to close the database connection and clean up tables
	teardown := func() {
		_, err := userDatabase.Exec("DELETE FROM users")
		if err != nil {
			t.Fatalf("failed to clean users table: %v", err)
		}
		_, err = msgDatabase.Exec("DELETE FROM messages")
		if err != nil {
			t.Fatalf("failed to clean messages table: %v", err)
		}
		userDatabase.Close()
		msgDatabase.Close()
	}

	return userDatabase, msgDatabase, teardown
}

func TestLoginUser(t *testing.T) {
	userDB, _, teardown := setupTestDB(t)
	defer teardown()

	h := &handlers.Handler{
		UserDB: userDB,
	}

	// First, register a user to login with
	password := "password123"
	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	err = db.AddUser(userDB, "testuser", hash)
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	loginData := map[string]string{
		"UserID":   "testuser",
		"Password": password,
	}
	body, _ := json.Marshal(loginData)

	req, err := http.NewRequest("POST", "/login", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.LoginUser)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body for a token
	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("could not parse response body: %v", err)
	}

	if _, ok := response["token"]; !ok {
		t.Errorf("handler did not return a token")
	}
}
