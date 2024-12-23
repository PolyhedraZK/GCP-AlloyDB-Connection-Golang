package connector

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/alloydbconn"
	"cloud.google.com/go/alloydbconn/driver/pgxv5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	// Version represents the current version number
	Version = "v1.0.0"
)

var db *gorm.DB

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("Fatal Error: %s environment variable not set.", k)
	}
	return v
}

// getEnvWithDefault retrieves environment variable with a default value
func getEnvWithDefault(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// InitDB initializes the database connection
func InitDB() error {
	// Get required environment variables
	instanceURI := mustGetenv("DB_HOST")   // AlloyDB instance URI
	username := mustGetenv("DB_USER")      // Database user
	password := mustGetenv("DB_PASS")      // Database password
	dbname := mustGetenv("DB_NAME")        // Database name
	certpath := mustGetenv("DB_CERT_PATH") // Service account key file path

	// Get optional connection pool configuration
	maxOpenConns, _ := strconv.Atoi(getEnvWithDefault("DB_MAX_OPEN_CONNS", "0"))              // Default 0: unlimited
	maxIdleConns, _ := strconv.Atoi(getEnvWithDefault("DB_MAX_IDLE_CONNS", "2"))              // Default 2
	connMaxLifetimeMinutes, _ := strconv.Atoi(getEnvWithDefault("DB_CONN_MAX_LIFETIME", "0")) // Default 0: unlimited
	connMaxIdleMinutes, _ := strconv.Atoi(getEnvWithDefault("DB_CONN_MAX_IDLE_TIME", "0"))    // Default 0: unlimited

	// Register AlloyDB driver
	_, err := pgxv5.RegisterDriver("alloydb", alloydbconn.WithCredentialsFile(certpath))
	if err != nil {
		return fmt.Errorf("failed to register alloydb driver: %w", err)
	}

	// Create connection string
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		instanceURI, username, password, dbname)

	// Open database connection
	sqlDB, err := sql.Open("alloydb", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool only when values are different from defaults
	if maxOpenConns != 0 {
		sqlDB.SetMaxOpenConns(maxOpenConns)
	}
	if maxIdleConns != 2 {
		sqlDB.SetMaxIdleConns(maxIdleConns)
	}
	if connMaxLifetimeMinutes != 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(connMaxLifetimeMinutes) * time.Minute)
	}
	if connMaxIdleMinutes != 0 {
		sqlDB.SetConnMaxIdleTime(time.Duration(connMaxIdleMinutes) * time.Minute)
	}

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Create GORM instance with existing connection
	db, err = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to create gorm instance: %w", err)
	}

	return nil
}

// GetDB returns the GORM database instance
func GetDB() *gorm.DB {
	return db
}

// GetVersion returns the current version number
func GetVersion() string {
	return Version
}
