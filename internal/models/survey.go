package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	JenisSurveyNyamuk = "Survey Nyamuk"
	JenisSurveyJentik = "Survey Jentik"
)

type Survey struct {
	ID               uuid.UUID             `gorm:"type:uuid;primaryKey" json:"id"`
	KeluargaID       uuid.UUID             `gorm:"type:uuid;not null" json:"keluarga_id"`
	Keluarga         Keluarga              `gorm:"foreignKey:KeluargaID;references:ID;constraint:OnDelete:CASCADE;" json:"keluarga,omitempty"`
	PetugasID        uuid.UUID             `gorm:"type:uuid;not null" json:"petugas_id"`
	Petugas          User                  `gorm:"foreignKey:PetugasID;references:ID;constraint:OnDelete:CASCADE;" json:"user,omitempty"`
	Tanggal          time.Time             `gorm:"type:date;not null" json:"tanggal"`
	JenisSurvey      string                `gorm:"type:varchar(50);not null" json:"jenis_survey"`
	CreatedAt        time.Time             `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time             `gorm:"autoUpdateTime" json:"updated_at"`
	Items            []SurveyItem          `gorm:"foreignKey:SurveyID;constraint:OnDelete:CASCADE;" json:"items,omitempty"`
	SurveyGambar     []SurveyGambar        `gorm:"foreignKey:SurveyID;constraint:OnDelete:CASCADE;" json:"survey_gambar,omitempty"`
	FollowUpNyamuk   *SurveyFollowUpNyamuk `gorm:"foreignKey:SurveyID;constraint:OnDelete:CASCADE;" json:"survey_followup_nyamuk,omitempty"`
	FollowUpJentik   *SurveyFollowUpJentik `gorm:"foreignKey:SurveyID;constraint:OnDelete:CASCADE;" json:"survey_followup_jentik,omitempty"`
	SurveyPSN        *SurveyPSN            `gorm:"foreignKey:SurveyID;constraint:OnDelete:CASCADE;" json:"psn,omitempty"`
	SurveyNyamukInfo *SurveyNyamukInfo     `gorm:"foreignKey:SurveyID;constraint:OnDelete:CASCADE;" json:"nyamuk_info,omitempty"`

	Status string `gorm:"-" json:"status"`
}

func (Survey) TableName() string {
	return "survey"
}
