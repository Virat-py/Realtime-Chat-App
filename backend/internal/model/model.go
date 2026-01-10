package model

type User struct {
	UserID       string
	PasswordHash string
}

type Message struct {
	UserID  string
	RoomID  int
	RoomName string
	Message string
	Time    string
}

type Room struct{
	RoomID int
	RoomName string
}
