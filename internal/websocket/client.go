package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Aryan-Gupta4460/letstalk/internal/models"
	"github.com/gorilla/websocket"
)

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	username string
	room     string
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, msgBytes, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		var msg models.Message
		err = json.Unmarshal(msgBytes, &msg)
		if err != nil {
			log.Println("Invalid message format")
			continue
		}

		// Add server timestamp
		msg.Username = c.username
		msg.Room = c.room
		msg.Timestamp = time.Now()
		c.hub.broadcast <- msg
	}
}

func (c *Client) writePump() {
	for message := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
