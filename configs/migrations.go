package configs

import (
	"betapa-antik-service/internal/models"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {

	return db.AutoMigrate(
		&models.Role{},
		&models.Kecamatan{},
		&models.Puskesmas{},
		&models.Kelurahan{},
		&models.User{},
		&models.Gambar{},
		&models.Materi{},
		&models.MateriGambar{},
		&models.Video{},
		&models.LupaKataSandi{},
		&models.Keluarga{},
		&models.SurveyLokasi{},
		&models.Survey{},
		&models.SurveyItem{},
		&models.SurveyFollowUpNyamuk{},
		&models.SurveyFollowUpJentik{},
		&models.SurveyGambar{},
		&models.Laporan{},
		&models.LaporanGambar{},
	)
}
