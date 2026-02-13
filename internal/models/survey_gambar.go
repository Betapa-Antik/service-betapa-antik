package models

import (
	"time"

	"github.com/google/uuid"
)

type SurveyGambar struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	SurveyID  uuid.UUID `gorm:"type:uuid;not null;index" json:"survey_id"`
	GambarID  uuid.UUID `gorm:"type:uuid;not null;index" json:"gambar_id"`
	Survey    *Survey   `gorm:"foreignKey:SurveyID;references:ID;constraint:OnDelete:CASCADE;" json:"-"`
	Gambar    Gambar    `gorm:"foreignKey:GambarID;references:ID;constraint:OnDelete:CASCADE;" json:"gambar"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (SurveyGambar) TableName() string {
	return "survey_gambar"
}
