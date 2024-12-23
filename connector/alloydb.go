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
	// Version 当前版本号
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

// getEnvWithDefault 获取环境变量，如果未设置则返回默认值
func getEnvWithDefault(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// InitDB 初始化数据库连接
func InitDB() error {
	// 获取必需的环境变量
	instanceURI := mustGetenv("DB_HOST")   // AlloyDB实例URI
	username := mustGetenv("DB_USER")      // 数据库用户
	password := mustGetenv("DB_PASS")      // 数据库密码
	dbname := mustGetenv("DB_NAME")        // 数据库名称
	certpath := mustGetenv("DB_CERT_PATH") // 服务账号密钥文件路径

	// 获取可选的连接池配置
	maxOpenConns, _ := strconv.Atoi(getEnvWithDefault("DB_MAX_OPEN_CONNS", "0"))              // 默认0，无限制
	maxIdleConns, _ := strconv.Atoi(getEnvWithDefault("DB_MAX_IDLE_CONNS", "2"))              // 默认2
	connMaxLifetimeMinutes, _ := strconv.Atoi(getEnvWithDefault("DB_CONN_MAX_LIFETIME", "0")) // 默认0，无限制
	connMaxIdleMinutes, _ := strconv.Atoi(getEnvWithDefault("DB_CONN_MAX_IDLE_TIME", "0"))    // 默认0，无限制

	// 注册AlloyDB驱动
	_, err := pgxv5.RegisterDriver("alloydb", alloydbconn.WithCredentialsFile(certpath))
	if err != nil {
		return fmt.Errorf("failed to register alloydb driver: %w", err)
	}

	// 创建连接字符串
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		instanceURI, username, password, dbname)

	// 打开数据库连接
	sqlDB, err := sql.Open("alloydb", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// 仅当值不是默认值时才配置连接池
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

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// 使用现有连接创建GORM实例
	db, err = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to create gorm instance: %w", err)
	}

	return nil
}

// GetDB 获取GORM数据库实例
func GetDB() *gorm.DB {
	return db
}

// GetVersion 获取当前版本号
func GetVersion() string {
	return Version
}
