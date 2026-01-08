package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectMsgDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./app.db")
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}

func CreateMsgTable(db *sql.DB) (error) {
	query := `
		CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTO INCREMENT,
			user_id TEXT NOT NULL,
			room_id INT NOT NULL,
			message TEXT NOT NULL,
			time TEXT NOT NULL
		);
		`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func AddMsg(db *sql.DB, user_id string,room_id int, message string) (error) {
	
	query := "INSERT INTO messages (user_id,room_id,message,time) VALUES (?,?,?,?)"
	curr_time:=time.Now().Format(time.DateTime)
	
	_, err := db.Exec(query, user_id,room_id, message,curr_time)
	if err != nil {
		return err
	}
	return nil
}


