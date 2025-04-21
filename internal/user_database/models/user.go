package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID                 uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"user_id"`
	Username           string    `gorm:"type:varchar(100);unique" json:"username"`
	Email              *string   `gorm:"uniqueIndex;column:email" json:"email"`
	HashedRefreshToken string    `gorm:"type:text" json:"hashed_rt"`
	HashedPassword     string    `gorm:"type:text" json:"password"`
	HashedAccessToken  string    `gorm:"type:text" json:"hashed_at"`
}
