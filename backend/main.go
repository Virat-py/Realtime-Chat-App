package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

var (
	clients   = make(map[*Client]bool)
	broadcast = make(chan []byte)
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func main() {
	curr_time:=time.Now().Format(time.DateTime)
	fmt.Println(curr_time)
	http.HandleFunc("/ws", handleConnections)
	go handleBroadcast()

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &Client{
		conn: conn,
		send: make(chan []byte),
	}

	clients[client] = true

	go client.write()
	go client.read()
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
