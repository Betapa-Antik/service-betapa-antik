package models

import (
	"time"

	"github.com/google/uuid"
)

// Gambar represents an uploaded image record.
type Gambar struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Path      string    `gorm:"type:varchar(255);not null" json:"path"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Gambar) TableName() string {
	return "gambar"
}
