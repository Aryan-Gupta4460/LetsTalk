package models

import "time"

type Message struct {
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Room      string    `json:"room"`
	Timestamp time.Time `json:"timestamp"`
}
