package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectUserDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./app.db")
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}

func CreateUserTable(db *sql.DB) (error) {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			user_id INTEGER PRIMARY KEY,
			password_hash TEXT NOT NULL
		);
		`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func AddUser(db *sql.DB, user_id string, password_hash string) (error) {
	query := "INSERT INTO users (user_id,password_hash) VALUES (?,?)"
	
	_, err := db.Exec(query, user_id, password_hash)
	if err != nil {
		return err
	}
	return nil
}

func GetUserHash(db *sql.DB,user_id string) (string,error) {
	query:="SELECT password_hash FROM users WHERE user_id=?"
	var hash string
	
	err:=db.QueryRow(query,user_id).Scan(&hash)
	if err!=nil{
		return "",err
	}
	return hash,nil
}
