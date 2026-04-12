package main

import (
	"log"
	"net/http"

	"github.com/Aryan-Gupta4460/letstalk/internal/db"
	migration "github.com/Aryan-Gupta4460/letstalk/internal/db"
	"github.com/Aryan-Gupta4460/letstalk/internal/websocket"
	"github.com/Aryan-Gupta4460/letstalk/pkg/config"
	"github.com/Aryan-Gupta4460/letstalk/pkg/logger"
)

func main() {
	cfg := config.LoadConfig("C:\\Users\\user\\Desktop\\LetsTalk\\configs\\config.yaml")
	db := db.InitDB(cfg)
	migration.RunMigrations(db)
	logger.Init()
	logger.Logger.Println("Database initialized:")
	hub := websocket.NewHub(db, logger.Logger)
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWS(hub, w, r)
	})

	logger.Logger.Println("Server started  on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
