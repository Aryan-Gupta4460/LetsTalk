package websocket

import (
	"encoding/json"

	"github.com/Aryan-Gupta4460/letstalk/internal/models"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan models.Message
	register   chan *Client
	unregister chan *Client
	history    [][]byte
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan models.Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {

		case client := <-h.register:
			h.clients[client] = true
			//  Send old messages
			for _, msg := range h.history {
				client.send <- msg
			}

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}

		case message := <-h.broadcast:
			//  Save to history
			msgBytes, _ := json.Marshal(message)
			h.history = append(h.history, msgBytes)

			// limit memory (last 100 messages)
			if len(h.history) > 100 {
				h.history = h.history[len(h.history)-100:]
			}

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
