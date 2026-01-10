package handlers

import (
	"net/http"
	"log"
	"encoding/json"
	"backend/internal/db"
	"backend/internal/model"
	"backend/internal/auth"
	
	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	send chan []model.Message
}

var (
	clients   = make(map[*Client]bool)
	broadcast = make(chan []model.Message)
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func (h *Handler) handleWebsocker(w http.ResponseWriter, r *http.Request) {
	// auth user
	
	_,err:=auth.VerifyToken(token)
    // parse room_id
    // upgrade
    // create Client with DB
    
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		conn: conn,
		send: make(chan []model.Message),
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
		err=json.Unmarshal(msg,&currMessage)
		if err!=nil{
			log.Println(err)
			return
		}
		// db.AddMsg(h.)
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
