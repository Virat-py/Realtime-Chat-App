package model

type User struct {
	UserID       string
	PasswordHash string
}

type Message struct {
	UserID   string `json:"user_id"`
	RoomID   int    `json:"room_id"`
	RoomName string `json:"room_name"`
	Message  string `json:"message"`
	Time     string `json:"time"`
}


type Room struct{
	RoomID int
	RoomName string
}
