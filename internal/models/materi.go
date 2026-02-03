package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	MateriStatusPublished = "published"
	MateriStatusDraft     = "draft"
)

// Materi represents educational content which can have multiple images via a pivot table.
type Materi struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Judul           string         `gorm:"type:varchar(255);not null" json:"judul"`
	Deskripsi       string         `gorm:"type:text;not null" json:"deskripsi"`
	Status          string         `gorm:"type:varchar(50);not null" json:"status"`
	CatatanTambahan *string        `gorm:"type:text;default:null" json:"catatan_tambahan,omitempty"`
	Gambar          []Gambar       `gorm:"many2many:materi_gambar;constraint:OnDelete:CASCADE;" json:"gambar,omitempty"`
	MateriGambars   []MateriGambar `gorm:"foreignKey:MateriID;constraint:OnDelete:CASCADE;" json:"-"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Materi) TableName() string {
	return "materi"
}

// MateriGambar is the pivot model between Materi and Gambar.
type MateriGambar struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	MateriID  uuid.UUID `gorm:"type:uuid;not null;index" json:"materi_id"`
	GambarID  uuid.UUID `gorm:"type:uuid;not null;index" json:"gambar_id"`
	Materi    Materi    `gorm:"foreignKey:MateriID;references:ID" json:"materi,omitempty"`
	Gambar    Gambar    `gorm:"foreignKey:GambarID;references:ID" json:"gambar,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (MateriGambar) TableName() string {
	return "materi_gambar"
}
