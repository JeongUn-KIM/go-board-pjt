package domain

import "gorm.io/gorm"

type Board struct {
	gorm.Model
	Title   string `gorm:"not null" json:"title"`
	Author  string `gorm:"not null" json:"author"`
	Content string `gorm:"not null" json:"content"`
}

type WritePost struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Content string `json:"content"`
}
