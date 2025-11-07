package models

import "time"

type Link struct {
	ID        uint      `gorm:"primaryKey"`
	Shortcode string    `gorm:"size:10;uniqueIndex;not null"`
	LongURL   string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// Link représente un lien raccourci dans la base de données.
// Les tags `gorm:"..."` définissent comment GORM doit mapper cette structure à une table SQL.
// ID qui est une primaryKey
// Shortcode : doit être unique, indexé pour des recherches rapide (voir doc), taille max 10 caractères
// LongURL : doit pas être null
// CreateAt : Horodatage de la créatino du lien
