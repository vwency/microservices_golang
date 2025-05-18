package main

import (
	"github.com/vwency/microservices_golang/internal/user_database/repository"
	"gorm.io/gorm"
)

func newRepository(db *gorm.DB) repository.Repository {
	repo := repository.NewRepository(db)
	return *repo
}
