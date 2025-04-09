package repository

import (
	"fmt"
	"log"

	"github.com/vwency/microservices_golang/internal/database/models"
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
	if cfg.Database.URL == "" {
		return nil, fmt.Errorf("database URL is empty")
	}

	log.Printf("Attempting to connect to database with URL: %s", cfg.Database.URL)

	db, err := gorm.Open(postgres.Open(cfg.Database.URL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	// Verify the connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Try pinging the database to verify connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established successfully.")
	repo := &databaseInitRepository{DB: db, cfg: cfg}
	log.Printf("Repository initialized: %+v", repo)
	return repo, nil
}

// RunMigrations сбрасывает таблицы, а затем создает их заново
func (r *databaseInitRepository) RunMigrations() error {
	if r == nil || r.DB == nil {
		return fmt.Errorf("repository or DB connection is nil")
	}

	migrator := r.DB.Migrator()
	if migrator == nil {
		return fmt.Errorf("failed to get database migrator")
	}

	log.Println("Starting migration process.")

	// Проверка моделей перед миграцией
	models := []interface{}{
		&models.User{},
	}

	for _, model := range models {
		if err := migrator.DropTable(model); err != nil {
			return fmt.Errorf("failed to drop table %T: %w", model, err)
		}
	}

	// Создание таблиц
	if err := r.DB.AutoMigrate(models...); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Migrations completed successfully.")
	return nil
}
