package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func InitDB() {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Printf("⚠️  Warning loading .env file: %v", err)
		} else {
			log.Println("📄 Loaded configuration from .env file")
		}
	} else {
		log.Println("🐳 Running in container mode - using environment variables")
	}

	config := getDatabaseConfig()

	if err := validateConfig(config); err != nil {
		log.Fatal("❌ Invalid database configuration:", err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.User, config.Password, config.Host, config.Port, config.Name)

	log.Printf("🔗 Connecting to database at %s:%s/%s as %s",
		config.Host, config.Port, config.Name, config.User)

	if err := connectWithRetry(dsn, 10, 2*time.Second); err != nil {
		log.Fatal("❌ Failed to connect to database after retries:", err)
	}

	configureConnectionPool()
	runMigrations()
	log.Println("✅ Database initialized successfully!")
}

type DatabaseConfig struct {
	User     string
	Password string
	Name     string
	Host     string
	Port     string
}

func getDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		User:     getEnvWithDefault("DB_USER", "postgres"),
		Password: getEnvWithDefault("DB_PASS", ""),
		Name:     getEnvWithDefault("DB_NAME", "go_insight"),
		Host:     getEnvWithDefault("DB_HOST", "localhost"),
		Port:     getEnvWithDefault("DB_PORT", "5432"),
	}
}

func validateConfig(config DatabaseConfig) error {
	if config.Password == "" {
		return fmt.Errorf("DB_PASS environment variable is required")
	}
	if config.Host == "" {
		return fmt.Errorf("DB_HOST environment variable is required")
	}
	if config.Name == "" {
		return fmt.Errorf("DB_NAME environment variable is required")
	}
	return nil
}

func connectWithRetry(dsn string, maxRetries int, delay time.Duration) error {
	var err error

	for i := range maxRetries {
		DB, err = sql.Open("pgx", dsn)
		if err != nil {
			log.Printf("⏳ Database connection attempt %d/%d failed: %v", i+1, maxRetries, err)
			time.Sleep(delay)
			continue
		}

		if err = DB.Ping(); err != nil {
			log.Printf("⏳ Database ping attempt %d/%d failed: %v", i+1, maxRetries, err)
			DB.Close()
			time.Sleep(delay)
			continue
		}

		return nil
	}

	return fmt.Errorf("failed to connect after %d attempts: %v", maxRetries, err)
}

func configureConnectionPool() {
	maxOpenConns := getEnvIntWithDefault("DB_MAX_OPEN_CONNS", 25)
	maxIdleConns := getEnvIntWithDefault("DB_MAX_IDLE_CONNS", 10)
	connMaxLifetime := getEnvIntWithDefault("DB_CONN_MAX_LIFETIME", 300) // 5 minutes

	DB.SetMaxOpenConns(maxOpenConns)
	DB.SetMaxIdleConns(maxIdleConns)
	DB.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)

	log.Printf("🔧 Connection pool configured: max_open=%d, max_idle=%d, max_lifetime=%ds",
		maxOpenConns, maxIdleConns, connMaxLifetime)
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntWithDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		log.Printf("⚠️  Invalid integer value for %s: %s, using default: %d", key, value, defaultValue)
	}
	return defaultValue
}

func runMigrations() {
	log.Println("🔄 Running database migrations...")

	migrationFiles := []string{
		"internal/db/migrations/001_create_logs_table.sql",
		"internal/db/migrations/002_create_metrics_table.sql",
		"internal/db/migrations/003_create_spans_and_traces_table.sql",
		"internal/db/migrations/004_add_performance_indexes.sql",
	}

	successCount := 0
	for _, file := range migrationFiles {
		if err := runSingleMigration(file); err != nil {
			log.Printf("⚠️ Migration %s failed: %v", file, err)
		} else {
			log.Printf("✅ Applied migration: %s", file)
			successCount++
		}
	}

	log.Printf("✅ Database migrations complete! (%d/%d successful)", successCount, len(migrationFiles))
}

func runSingleMigration(file string) error {
	content, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("could not read migration file: %w", err)
	}

	if _, err := DB.Exec(string(content)); err != nil {
		return fmt.Errorf("migration execution failed: %w", err)
	}

	return nil
}
