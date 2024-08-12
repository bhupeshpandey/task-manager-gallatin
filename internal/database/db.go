package database

import (
	"database/sql"
	"fmt"
	"github.com/bhupeshpandey/task-manager-gallatin/internal/models"
	_ "github.com/lib/pq"
	"log"
)

func NewDatabase(cfg models.DatabaseConfig) models.TaskRepository {
	var db *sql.DB
	var err error

	var dbRepo models.TaskRepository
	switch cfg.Type {
	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DBName)
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Fatalf("failed to connect to PostgreSQL: %v", err)
		}
		dbRepo = NewPostgresRepository(db)

	case "sqlite":
		db, err = sql.Open("sqlite3", cfg.SQLite.Filepath)
		if err != nil {
			log.Fatalf("failed to connect to SQLite: %v", err)
		}

	default:
		log.Fatalf("unsupported database type: %s", cfg.Type)
	}

	// Verify the connection
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping the database: %v", err)
	}

	log.Println("Database connection successfully established")
	return dbRepo
}
