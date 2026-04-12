package websocket

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/Aryan-Gupta4460/letstalk/internal/models"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan models.Message
	register   chan *Client
	unregister chan *Client
	db         *sql.DB
	logger     *log.Logger
}

func NewHub(db *sql.DB, logger *log.Logger) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan models.Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		db:         db,
		logger:     logger,
	}
}

func (h *Hub) Run() {
	for {
		select {

		case client := <-h.register:
			h.clients[client] = true

			// Load history from DB
			messages := h.getRecentMessages(client.room, 20, 0)
			if len(messages) == 0 {
				h.logger.Println("No more messages available")
			}
			for i := len(messages) - 1; i >= 0; i-- {
				client.send <- messages[i]
			}

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}

		case message := <-h.broadcast:
			//  Save to DB
			_, err := h.db.Exec(
				"INSERT INTO messages (username, content, room, timestamp) VALUES ($1, $2, $3, $4)",
				message.Username, message.Content, message.Room, message.Timestamp,
			)
			if err != nil {
				h.logger.Println("DB insert error:", err)
			}
			//  Save to history
			msgBytes, _ := json.Marshal(message)

			for client := range h.clients {
				// FILTER BY ROOM
				if client.room != message.Room {
					continue
				}
				select {
				case client.send <- msgBytes:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
func (h *Hub) getRecentMessages(room string, limit int, offset int) [][]byte {
	rows, err := h.db.Query(
		`SELECT username, content, room, timestamp 
		 FROM messages 
		 WHERE room=$1 
		 ORDER BY timestamp DESC 
		 LIMIT $2 OFFSET $3`,
		room, limit, offset,
	)
	if err != nil {
		h.logger.Println("DB fetch error:", err)
		return nil
	}
	defer rows.Close()

	var messages [][]byte

	for rows.Next() {
		var msg models.Message

		err := rows.Scan(&msg.Username, &msg.Content, &msg.Room, &msg.Timestamp)
		if err != nil {
			h.logger.Println("Row scan error:", err)
			continue
		}

		msgBytes, _ := json.Marshal(msg)
		messages = append(messages, msgBytes)
	}

	return messages
}
