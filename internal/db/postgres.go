package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Aryan-Gupta4460/letstalk/pkg/config"
	_ "github.com/lib/pq"
)

func InitDB(cfg *config.Config) *sql.DB {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("DB connection error:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("DB ping failed:", err)
	}

	log.Println("Connected to PostgreSQL")

	return db
}

func RunMigrations(db *sql.DB) {
	query, err := os.ReadFile("C:\\Users\\user\\Desktop\\LetsTalk\\internal\\db\\migrations.sql")
	if err != nil {
		log.Fatal("Failed to read migration file:", err)
	}

	_, err = db.Exec(string(query))
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Migration executed successfully")
}
