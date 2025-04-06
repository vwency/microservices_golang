package repository

type User struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"type:varchar(100)"`
	Email string `gorm:"type:varchar(100);unique"`
}

type Column struct {
	ID          string  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"column_id"`
	UserID      string  `gorm:"type:uuid" json:"user_id"`
	ColumnName  string  `gorm:"unique;type:varchar(255)" json:"column_name"`
	Description *string `gorm:"type:varchar(255)" json:"description"`

	Cards    []Card    `gorm:"foreignKey:ColumnID;references:ID" json:"cards"`
	Comments []Comment `gorm:"foreignKey:ColumnID;references:ID" json:"comments"`
	User     User      `gorm:"foreignKey:UserID;references:ID" json:"user"`
}

type Card struct {
	ID          string  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"card_id"`
	UserID      string  `gorm:"type:uuid" json:"user_id"`
	ColumnID    string  `gorm:"type:uuid" json:"column_id"`
	CardName    string  `gorm:"unique;type:varchar(255)" json:"card_name"`
	Description *string `gorm:"type:varchar(255)" json:"description"`

	User     User      `gorm:"foreignKey:UserID;references:ID" json:"user"`
	Column   Column    `gorm:"foreignKey:ColumnID;references:ID" json:"column"`
	Comments []Comment `gorm:"foreignKey:CardID;references:ID" json:"comments"`
}

type Comment struct {
	ID          string  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"comment_id"`
	UserID      string  `gorm:"type:uuid" json:"user_id"`
	ColumnID    string  `gorm:"type:uuid" json:"column_id"`
	CardID      string  `gorm:"type:uuid" json:"card_id"`
	CommentName string  `gorm:"unique;type:varchar(255)" json:"comment_name"`
	Description *string `gorm:"type:varchar(255)" json:"description"`

	User   User   `gorm:"foreignKey:UserID;references:ID" json:"user"`
	Column Column `gorm:"foreignKey:ColumnID;references:ID" json:"column"`
	Card   Card   `gorm:"foreignKey:CardID;references:ID" json:"card"`
}
