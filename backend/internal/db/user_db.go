package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectDB() (*sql.DB,error) {
	db, err := sql.Open("sqlite3", "./app.db")
	if err != nil {
		return nil,err
	}

	if err := db.Ping(); err != nil {
		return nil,err
	}
	return db,nil
}

func CreateUserTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			UserID TEXT PRIMARY KEY,
			PasswordHash TEXT NOT NULL
		);
		`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func AddUser(db *sql.DB, UserID string, PasswordHash string) error {
	query := "INSERT INTO users (UserID,PasswordHash) VALUES (?,?)"

	_, err := db.Exec(query, UserID, PasswordHash)
	if err != nil {
		return err
	}
	return nil
}

func CheckUserExists(db *sql.DB, UserID string) (bool,error) {
	query := `SELECT 1 FROM users WHERE UserID = ? LIMIT 1`

	var exists int
	err := db.QueryRow(query, UserID).Scan(&exists)

	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetUserHash(db *sql.DB, UserID string) (string, error) {
	query := "SELECT PasswordHash FROM users WHERE UserID=?"
	var hash string

	err := db.QueryRow(query, UserID).Scan(&hash)
	if err != nil {
		return "", err
	}
	return hash, nil
}
