package db

import (
	"backend/internal/model"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func CreateRoomTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS rooms (
			RoomID INTEGER PRIMARY KEY AUTOINCREMENT,
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
	query:="INSERT INTO rooms (RoomName) VALUES (?)"
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

func GetRooms(db *sql.DB) ([]model.Room, error) {

	query := `
	SELECT RoomID, RoomName
	FROM rooms
	ORDER BY RoomID ASC;
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