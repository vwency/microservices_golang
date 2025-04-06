package repository

import (
	"fmt"

	"github.com/vwency/microservices_golang/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseInitRepository interface {
	RunMigrations() error
}

type databaseInitRepository struct {
	DB  *gorm.DB
	cfg config.ServiceConfig
}

func NewRepository(cfg config.ServiceConfig) (DatabaseInitRepository, error) {
	db, err := gorm.Open(postgres.Open(cfg.Database.URL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	return &databaseInitRepository{DB: db, cfg: cfg}, nil
}

func (r *databaseInitRepository) RunMigrations() error {
	if r.DB == nil {
		return fmt.Errorf("database connection is not established")
	}

	if err := r.DB.AutoMigrate(&User{}, &Column{}, &Card{}, &Comment{}).Error; err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
