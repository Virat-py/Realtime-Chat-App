package db

import (
	"backend/internal/model"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectMsgDB() (*sql.DB,error) {
	db, err := sql.Open("sqlite3", "./app.db")
	if err != nil {
		return nil,err
	}

	if err := db.Ping(); err != nil {
		return nil,err
	}
	return db,nil
}

func CreateMsgTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS messages (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			UserID TEXT NOT NULL,
			RoomID INT NOT NULL,
			RoomName TEXT NOT NULL,
			Message TEXT NOT NULL,
			Time TEXT NOT NULL
		);
		`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func AddMsg(db *sql.DB, msg model.Message) error {

	query := "INSERT INTO messages (UserID,RoomID,RoomName,Message,Time) VALUES (?,?,?,?,?)"

	_, err := db.Exec(query, msg.UserID, msg.RoomID, msg.RoomName, msg.Message, msg.Time)
	if err != nil {
		return err
	}
	return nil
}

func GetMsgOfRoom(db *sql.DB, RoomID int) ([]model.Message, error) {

	query := "SELECT UserID,RoomID,Message,Time FROM messages WHERE RoomID=? ORDER BY ID ASC"

	rows, err := db.Query(query, RoomID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var messages []model.Message

	for rows.Next() {
		var msg model.Message
		err := rows.Scan(&msg.UserID, &msg.RoomID, &msg.Message, &msg.Time)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}


func GetRooms(db *sql.DB) ([]model.Room, error) {

	query := `
	SELECT RoomID, RoomName
	FROM messages
	GROUP BY RoomID, RoomName
	ORDER BY MIN(ID) ASC;
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var rooms []model.Room

	for rows.Next() {
		var room model.Room
		err := rows.Scan(&room.RoomID,&room.RoomName)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return rooms, nil
}
