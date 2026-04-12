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

		// Detect message type
		var raw map[string]any
		if err := json.Unmarshal(msgBytes, &raw); err != nil {
			c.hub.logger.Println("Invalid JSON")
			continue
		}
		//  Handled load_more
		if msgType, ok := raw["type"].(string); ok && msgType == "load_more" {
			offset := 0
			if val, ok := raw["offset"].(float64); ok {
				offset = int(val)
			}

			messages := c.hub.getRecentMessages(c.room, 20, offset)
			if len(messages) == 0 {
				response := map[string]string{
					"type": "no_more_messages",
				}

				respBytes, _ := json.Marshal(response)
				c.send <- respBytes
				continue
			}

			// send in correct order
			for i := len(messages) - 1; i >= 0; i-- {
				c.send <- messages[i]
			}
			continue
		}

		//  Normal chat message
		var msg models.Message
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			c.hub.logger.Println("Invalid message format")
			continue
		}

		// enrich from server
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
