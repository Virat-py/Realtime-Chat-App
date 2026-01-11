package handlers

import (
	"backend/internal/auth"
	"backend/internal/db"
	"backend/internal/model"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn   *websocket.Conn
	send   chan model.Message
	msgDB  *sql.DB
	roomDB *sql.DB
	userID string
	roomID int
}

type incomingMsg struct {
	Message string `json:"message"`
}

var (
	clients   = make(map[*Client]bool)
	broadcast = make(chan model.Message)
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
		send:   make(chan model.Message),
		msgDB:  h.MsgDB,
		roomDB: h.RoomDB,
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

		var currMsgTextJSON incomingMsg
		err = json.Unmarshal(msg, &currMsgTextJSON)
		if err != nil {
			log.Println(err)
			return
		}
		currMsgText:=currMsgTextJSON.Message
		// don't trust client send fields
		
		var currMessage model.Message
		currMessage.Message=currMsgText
		currMessage.UserID = c.userID
		currMessage.RoomID = c.roomID
		ist := time.FixedZone("IST", 5*60*60+30*60) // +5 hours 30 minutes
		nowIST := time.Now().In(ist)
		currMessage.Time = nowIST.Format(time.DateTime)
		roomName,err:=db.GetRoomName(c.roomDB,c.roomID)
		currMessage.RoomName=roomName

		err = db.AddMsg(c.msgDB, currMessage)
		if err != nil {
			log.Println(err)
			return
		}
		broadcast <- currMessage
	}
}

func (c *Client) write() {

	defer c.conn.Close()

	for currMessage := range c.send {
		data, err := json.Marshal(currMessage)
		err = c.conn.WriteMessage(websocket.TextMessage, data)
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
