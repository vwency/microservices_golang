package user

import (
	"github.com/vwency/microservices_golang/internal/database/models"
	"gorm.io/gorm"
)

// RunMigrations выполняет миграции для пользователей, создавая таблицы в базе данных
func RunMigrations(db *gorm.DB) error {
	// Создаем таблицу пользователей, если она еще не существует
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return err
	}
	return nil
}
