package user

import (
	"github.com/vwency/microservices_golang/internal/database/models"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	if db.Migrator().HasTable(&models.User{}) {
		// Удаляем таблицу, если она существует
		if err := db.Migrator().DropTable(&models.User{}); err != nil {
			return err
		}
	}
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return err
	}
	return nil
}
