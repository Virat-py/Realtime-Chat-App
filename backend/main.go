package main

import (
	"backend/internal/db"
	"backend/internal/handlers"
	"log"
	"net/http"
)

func main() {
	userDatabase, err := db.ConnectUserDB()
	if err != nil {
		log.Println(err)
		return
	}

	if err := db.CreateUserTable(userDatabase); err != nil {
		log.Fatal("failed to create users table:", err)
	}

	msgDatabase, err := db.ConnectMsgDB()
	if err != nil {
		log.Println(err)
		return
	}

	if err := db.CreateMsgTable(msgDatabase); err != nil {
		log.Fatal("failed to create messages table:", err)
	}
	
	roomDatabase, err := db.ConnectRoomDB()
	if err != nil {
		log.Println(err)
		return
	}

	if err := db.CreateRoomTable(msgDatabase); err != nil {
		log.Fatal("failed to create messages table:", err)
	}

	defer userDatabase.Close()
	defer msgDatabase.Close()
	defer roomDatabase.Close()

	h := &handlers.Handler{
		UserDB: userDatabase,
		MsgDB:  msgDatabase,
		RoomDB: roomDatabase,
	}

	http.HandleFunc("/register", h.RegisterUser)
	http.HandleFunc("/login", h.LoginUser)

	http.HandleFunc("/create_room",h.CreateRoom)
	http.HandleFunc("/get_rooms", h.GetAllRoomsHandler)
	http.HandleFunc("/room/", h.GetRoomData)

	http.HandleFunc("/ws", h.HandleWebSockets)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
