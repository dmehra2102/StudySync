package database

import (
	"fmt"

	"github.com/dmehra2102/StudySync/internal/config"
	"github.com/dmehra2102/StudySync/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgres(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&domain.User{},
		&domain.Notification{},
		&domain.Task{},
		&domain.StudySession{},
	); err != nil {
		return nil, err
	}

	return db, nil
}
