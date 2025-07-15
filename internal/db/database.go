package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/NathanSanchezDev/go-insight/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func InitDB(cfg *config.Config) {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Printf("⚠️  Warning loading .env file: %v", err)
		} else {
			log.Println("📄 Loaded configuration from .env file")
		}
	} else {
		log.Println("🐳 Running in container mode - using environment variables")
	}

	dbConfig := getDatabaseConfig(cfg)

	if err := validateConfig(dbConfig); err != nil {
		log.Fatal("❌ Invalid database configuration:", err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name)

	log.Printf("🔗 Connecting to database at %s:%d/%s as %s",
		dbConfig.Host, dbConfig.Port, dbConfig.Name, dbConfig.User)

	if err := connectWithRetry(dsn, 10, 2*time.Second); err != nil {
		log.Fatal("❌ Failed to connect to database after retries:", err)
	}

	configureConnectionPool(cfg)
	runMigrations()
	log.Println("✅ Database initialized successfully!")
}

type DatabaseConfig struct {
	User     string
	Password string
	Name     string
	Host     string
	Port     int
}

func getDatabaseConfig(cfg *config.Config) DatabaseConfig {
	return DatabaseConfig{
		User:     cfg.Database.User,
		Password: getEnvWithDefault("DB_PASS", ""), // Keep password in ENV
		Name:     cfg.Database.Name,
		Host:     getEnvWithDefault("DB_HOST", "localhost"), // Keep host in ENV
		Port:     cfg.Database.Port,
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
		return fmt.Errorf("DB_NAME cannot be empty")
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

func configureConnectionPool(cfg *config.Config) {
	maxIdleConns := cfg.Database.MaxConnections / 2 // Half of max as idle
	connMaxLifetime := 5 * time.Minute

	DB.SetMaxOpenConns(cfg.Database.MaxConnections)
	DB.SetMaxIdleConns(maxIdleConns)
	DB.SetConnMaxLifetime(connMaxLifetime)

	log.Printf("🔧 Connection pool configured: max_open=%d, max_idle=%d, max_lifetime=%v",
		cfg.Database.MaxConnections, maxIdleConns, connMaxLifetime)
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
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
