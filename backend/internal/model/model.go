package model

type User struct{
	UserID string
	PasswordHash string
}

type Message struct{
	UserID string
	RoomID int
	Message string
	Time string
}