package db

import (
	// "backend/internal/model"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectRoomDB() (*sql.DB,error) {
	db, err := sql.Open("sqlite3", "./app.db")
	if err != nil {
		return nil,err
	}

	if err := db.Ping(); err != nil {
		return nil,err
	}
	return db,nil
}

func CreateRoomTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS rooms (
			RoomID int PRIMARY KEY AUTOINCREMENT,
			RoomName TEXT NOT NULL
		);
		`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func AddRoom(db *sql.DB, RoomName string) (error){
	query:="INSERT INTO rooms VALUES RoomName=?"
	_,err:=db.Exec(query,RoomName)
	if err!=nil{
		return err
	}
	return nil
}

func GetRoomID(db *sql.DB, RoomName string) (int,error){
	query:="SELECT RoomID FROM rooms WHERE RoomName=?"
	var RoomID int
	err:=db.QueryRow(query,RoomName).Scan(&RoomID)
	if err!=nil{
		return -1,err
	}
	return RoomID,nil
}

func GetRoomName(db *sql.DB, RoomID int) (string,error){
	query:="SELECT RoomName FROM rooms WHERE RoomID=?"
	var RoomName string
	err:=db.QueryRow(query,RoomID).Scan(&RoomName)
	if err!=nil{
		return "",err
	}
	return RoomName,nil
}