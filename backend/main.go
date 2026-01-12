package main

import (
	"backend/internal/db"
	"backend/internal/handlers"
	"log"
	"net/http"
)

func main() {
	Database, err := db.ConnectDB()
	if err != nil {
		log.Println(err)
		return
	}

	if err := db.CreateUserTable(Database); err != nil {
		log.Fatal("failed to create users table:", err)
	}

	if err := db.CreateMsgTable(Database); err != nil {
		log.Fatal("failed to create messages table:", err)
	}
	
	if err := db.CreateRoomTable(Database); err != nil {
		log.Fatal("failed to create room table:", err)
	}

	defer Database.Close()

	h := &handlers.Handler{
		DB:Database,
	}
	
	go handlers.HandleBroadcast()

	http.HandleFunc("/register", h.RegisterUser)
	http.HandleFunc("/login", h.LoginUser)

	http.HandleFunc("/create_room",h.CreateRoom)
	http.HandleFunc("/get_rooms", h.GetAllRoomsHandler)
	http.HandleFunc("/room/", h.GetRoomData)

	http.HandleFunc("/ws", h.HandleWebSockets)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", enableCORS(http.DefaultServeMux)))
}

func enableCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
        w.Header().Set("Access-Control-Allow-Credentials", "true")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}

