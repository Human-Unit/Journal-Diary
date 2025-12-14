package models

import (
	"gorm.io/gorm"
)

type Entry struct {
	gorm.Model `json:"-"`
	UserID     uint64   `gorm:"not null;index"` // Added index for better query performance
	Title      string `gorm:"size:100"`       // Added reasonable size limit
	Content    string `gorm:"type:text;not null"`
	IsPrivate  bool   `gorm:"default:true;not null"`
}
type EntryIn struct {          // ID is optional for creation, but useful for updates
	UserID     uint64   `gorm:"not null;index"` // Added index for better query performance
	Title      string `gorm:"size:100"`       // Added reasonable size limit
	Content    string `gorm:"type:text;not null"`
	IsPrivate  bool   `gorm:"default:true;not null"`
}
type EntryResponse struct {
    ID      uint64 `json:"id"`
    UserID  uint64 `json:"user_id"`
    Title   string `json:"title"`
    Content string `json:"content"`
}
