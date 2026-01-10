package handlers

import (
	"backend/internal/auth"
	"backend/internal/model"
	"backend/internal/db"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn   *websocket.Conn
	send   chan []byte
	msgDB  *sql.DB
	userID string
	roomID int
}

var (
	clients   = make(map[*Client]bool)
	broadcast = make(chan []byte)
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func (h *Handler) HandleWebSockets(w http.ResponseWriter, r *http.Request) {
	// auth user
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "missing authorization header", http.StatusUnauthorized)
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, "invalid authorization format", http.StatusUnauthorized)
		return
	}

	tokenString := parts[1]
	claims, err := auth.VerifyToken(tokenString)
	if err != nil {
		log.Println(err)
		http.Error(w, "Token expired or invalid", http.StatusUnauthorized)
		return
	}
	// parse room_id
	roomIDStr := r.URL.Query().Get("room_id")
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		http.Error(w, "invalid room id", http.StatusBadRequest)
		return
	}
	// upgrade
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// create Client with DB

	client := &Client{
		conn:   conn,
		send:   make(chan []byte),
		msgDB:  h.MsgDB,
		userID: claims.UserID,
		roomID: roomID,
	}

	clients[client] = true

	go client.read()
	go client.write()

}

func (c *Client) read() {
	defer func() {
		delete(clients, c)
		c.conn.Close()
	}()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			return
		}
		var currMessage model.Message
		err = json.Unmarshal(msg, &currMessage)
		if err != nil {
			log.Println(err)
			return
		}
		currMessage.UserID=c.userID
		err = db.AddMsg(c.msgDB, currMessage)
		broadcast <- msg
	}
}

func (c *Client) write() {

	defer c.conn.Close()

	for msg := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}

func handleBroadcast() {
	for {
		msg := <-broadcast
		for client := range clients {
			select {
			case client.send <- msg:
			default:
				delete(clients, client)
				close(client.send)
			}
		}
	}
}
