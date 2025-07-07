package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌ Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName)

	DB, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("❌ Database is not reachable:", err)
	}

	// Run migrations automatically
	runMigrations()

	log.Println("✅ Connected to PostgreSQL!")
}

func runMigrations() {
	log.Println("🔄 Running database migrations...")

	migrationFiles := []string{
		"internal/db/migrations/001_create_logs_table.sql",
		"internal/db/migrations/002_create_metrics_table.sql",
		"internal/db/migrations/003_create_spans_and_traces_table.sql",
		"internal/db/migrations/004_add_performance_indexes.sql",
	}

	for _, file := range migrationFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			log.Printf("⚠️ Could not read migration %s: %v", file, err)
			continue
		}

		_, err = DB.Exec(string(content))
		if err != nil {
			log.Printf("⚠️ Migration %s failed: %v", file, err)
		} else {
			log.Printf("✅ Applied migration: %s", file)
		}
	}

	log.Println("✅ Database migrations complete!")
}
