package masyarakatroutes

import (
	masyarakatcontroller "betapa-antik-service/internal/controllers/masyarakat_controller"
	masyarakatrepo "betapa-antik-service/internal/repositories/masyarakat_repo"
	masyarakatservice "betapa-antik-service/internal/services/masyarakat_service"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func MasyarakatRoutes(e *echo.Group, db *gorm.DB) {
	masyarakatRepo := masyarakatrepo.NewMasyarakatRepositoryImpl(db)
	masyarakatSvc := masyarakatservice.NewMasyarakatServiceImpl(masyarakatRepo)
	masyarakatCtrl := masyarakatcontroller.NewMasyarakatController(masyarakatSvc)
	e.GET("/materi/latest", masyarakatCtrl.GetLatestMaterial)
	e.GET("/video/latest", masyarakatCtrl.GetLatestVideo)
	e.GET("/materi/all", masyarakatCtrl.GetAllPublicMateri)
	e.GET("/video/all", masyarakatCtrl.GetAllPublicVideo)
	e.GET("/materi/:materiId", masyarakatCtrl.GetPublicMateriByID)
	e.GET("/video/:videoId", masyarakatCtrl.GetPublicVideoByID)
	e.POST("/laporan/create", masyarakatCtrl.CreateLaporan)
	e.GET("/location-resolve", masyarakatCtrl.GetLocationResolve)
}
