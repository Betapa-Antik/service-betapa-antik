package models

import (
	"time"

	"github.com/google/uuid"
)

type SurveyFollowUpNyamuk struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	SurveyID         uuid.UUID `gorm:"type:uuid;not null" json:"survey_id"`
	Survey           *Survey   `gorm:"foreignKey:SurveyID;references:ID;constraint:OnDelete:CASCADE;" json:"-"`
	DitemukanAedes   bool      `gorm:"type:bool" json:"ditemukan_aedes"`
	TingkatInfestasi string    `gorm:"type:varchar(100)" json:"tingkat_infestasi"`
	FoogingStatus    string    `gorm:"type:varchar(100)" json:"fooging_status"`
	EdukasiOrAbate   bool      `gorm:"type:bool" json:"edukasi_or_abate"`
	Catatan          string    `gorm:"type:text" json:"catatan"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type SurveyFollowUpJentik struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	SurveyID     uuid.UUID `gorm:"type:uuid;not null" json:"survey_id"`
	Survey       *Survey   `gorm:"foreignKey:SurveyID;references:ID;constraint:OnDelete:CASCADE;" json:"-"`
	EdukasiPSN   bool      `gorm:"type:bool" json:"edukasi_psn"`
	TindakLanjut string    `gorm:"type:varchar(255)" json:"tindak_lanjut"`
	Catatan      string    `gorm:"type:text" json:"catatan"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (SurveyFollowUpNyamuk) TableName() string {
	return "survey_followup_nyamuk"
}

func (SurveyFollowUpJentik) TableName() string {
	return "survey_followup_jentik"
}
