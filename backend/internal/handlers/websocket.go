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
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn     *websocket.Conn
	send     chan model.Message
	DB       *sql.DB
	userID   string
	roomID   int
	roomName string
}

type incomingMsg struct {
	Message string `json:"message"`
}

var (
	clients      = make(map[*Client]bool)
	broadcast    = make(chan model.Message)
	roomToClient = make(map[int][]*Client)

	clientsMu sync.Mutex
	roomsMu   sync.Mutex

	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func RemoveClientFromRoom(roomID int, clientToRemove *Client) {

	roomsMu.Lock()
	defer roomsMu.Unlock()

	clients_in_room, exists := roomToClient[roomID]
	if !exists {
		return // room doesn't exist
	}

	// Find and remove the client
	for i, c := range clients_in_room {
		if c == clientToRemove { // compare pointers
			// Remove by slicing
			roomToClient[roomID] = append(clients_in_room[:i], clients_in_room[i+1:]...)
			return
		}
	}
	// client not found in this room → nothing to do
}

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
	//
	roomName, err := db.GetRoomName(h.DB, roomID)

	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		conn:     conn,
		send:     make(chan model.Message),
		DB:       h.DB,
		userID:   claims.UserID,
		roomID:   roomID,
		roomName: roomName,
	}

	clientsMu.Lock()
	clients[client] = true
	clientsMu.Unlock()

	roomsMu.Lock()
	roomToClient[roomID] = append(roomToClient[roomID], client)
	roomsMu.Unlock()

	go client.read()
	go client.write()

}

func (c *Client) read() {
	defer func() {

		clientsMu.Lock()
		delete(clients, c)
		clientsMu.Unlock()

		RemoveClientFromRoom(c.roomID, c)

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
		currMsgText := currMsgTextJSON.Message
		// don't trust client send fields

		var currMessage model.Message
		currMessage.Message = currMsgText
		currMessage.UserID = c.userID
		currMessage.RoomID = c.roomID
		ist := time.FixedZone("IST", 5*60*60+30*60) // +5 hours 30 minutes
		nowIST := time.Now().In(ist)
		currMessage.Time = nowIST.Format(time.DateTime)
		currMessage.RoomName = c.roomName

		err = db.AddMsg(c.DB, currMessage)
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

func HandleBroadcast() {
	for {
		msg := <-broadcast

		for _, client := range roomToClient[msg.RoomID] {
			select {
			case client.send <- msg:
			default:
				delete(clients, client)
				close(client.send)
			}
		}
	}
}
