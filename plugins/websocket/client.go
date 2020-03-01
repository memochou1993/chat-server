package websocket

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/memochou1993/chat/helper"
)

// Client struct
type Client struct {
	ID   string
	Room *Room
	Conn *websocket.Conn
	Pool *Pool
}

// Message struct
type Message struct {
	RoomID   string `json:"roomId"`
	ClientID string `json:"clientId"`
	Type     int    `json:"type"`
	Body     string `json:"body"`
}

// NewClient func
func NewClient(pool *Pool, conn *websocket.Conn, room *Room, clientID string) *Client {
	return &Client{
		ID:   clientID,
		Room: room,
		Conn: conn,
		Pool: pool,
	}
}

// ReadMessage func
func (c *Client) ReadMessage() {
	defer func() {
		c.Pool.ClientUnregister <- c
		c.Conn.Close()
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()

		if err != nil {
			log.Println(err)
			return
		}

		message := Message{
			RoomID:   c.Room.ID,
			ClientID: c.ID,
			Type:     messageType,
			Body:     string(p),
		}

		c.Pool.Broadcast <- message
	}
}

// GetClientID func
func GetClientID(r *http.Request) string {
	id := fmt.Sprintf("%s: %s", helper.GetHost(r), helper.GetPlatform(r))

	return base64.StdEncoding.EncodeToString([]byte(id))
}
