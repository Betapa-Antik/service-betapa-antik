package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	VideoStatusPublished = "published"
	VideoStatusDraft     = "draft"
)

type Video struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Judul     string    `gorm:"type:varchar(255);not null" json:"judul"`
	Link      string    `gorm:"type:text;not null" json:"link"`
	Deskripsi string    `gorm:"type:text;not null" json:"deskripsi"`
	Status    string    `gorm:"type:varchar(50);not null" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Video) TableName() string {
	return "video"
}
