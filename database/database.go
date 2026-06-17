package database

import (
	"database/sql"
	"embed"
	"fmt"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
	"zzz/internal/config"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func NewDB(cfg config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	log.Println("✅ Connected to PostgreSQL")

	if err := RunMigrations(sqlDB); err != nil {
		return nil, fmt.Errorf("migration failed: %v", err)
	}

	return db, nil
}

func RunMigrations(db *sql.DB) error {
	log.Println("🚀 Running migrations...")

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	log.Println("Running migrations...")
	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}
	log.Println("✅ Migrations applied successfully!")
	return nil
}
