package models

import (
	"time"

	"github.com/google/uuid"
)

type SurveyLokasi struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	JenisSurvey string    `gorm:"type:varchar(50);not null" json:"jenis_survey"`
	NamaLokasi  string    `gorm:"type:varchar(100);not null" json:"nama_lokasi"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (SurveyLokasi) TableName() string {
	return "survey_lokasi"
}
