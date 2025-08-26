package models

import (
	"gorm.io/gorm"
)


type LogData struct {
    gorm.Model
    Email    string `gorm:"uniqueIndex"`
    Password string
}

type User struct{
	gorm.Model
	Email string `json:"email"`
	Password string `json:"password"`
}

type Entry struct {
	gorm.Model `json:"-"`
	ID         uint   `gorm:"primaryKey;autoIncrement"` // ID is auto-incremented
	UserID     uint   `gorm:"not null;index"` // Added index for better query performance
	Title      string `gorm:"size:100"`       // Added reasonable size limit
	Content    string `gorm:"type:text;not null"`
	IsPrivate  bool   `gorm:"default:true;not null"`
}
type EntryIn struct {          // ID is optional for creation, but useful for updates
	UserID     uint   `gorm:"not null;index"` // Added index for better query performance
	Title      string `gorm:"size:100"`       // Added reasonable size limit
	Content    string `gorm:"type:text;not null"`
	IsPrivate  bool   `gorm:"default:true;not null"`
}
