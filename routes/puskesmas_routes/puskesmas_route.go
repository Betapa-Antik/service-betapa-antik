package puskesmasroutes

import (
	puskesmascontroller "betapa-antik-service/internal/controllers/puskesmas_controller"
	"betapa-antik-service/internal/middlewares"
	adminrepo "betapa-antik-service/internal/repositories/admin_repo"
	puskesmasrepo "betapa-antik-service/internal/repositories/puskesmas_repo"
	puskesmasservice "betapa-antik-service/internal/services/puskesmas_service"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func PuskesmasRoutes(e *echo.Group, db *gorm.DB, rdb *redis.Client) {
	adminRepo := adminrepo.NewAdminRepositoryImpl(db)
	puskesmasRepo := puskesmasrepo.NewPuskesmasRepositoryImpl(db)
	puskesmasService := puskesmasservice.NewPuskesmasServiceImpl(puskesmasRepo, rdb)
	puskesmasCtrl := puskesmascontroller.NewPuskesmasController(puskesmasService)

	pr := e.Group("", middlewares.RequireRole("admin", adminRepo))
	pr.POST("/create", puskesmasCtrl.CreatePuskesmas)
	pr.GET("/all", puskesmasCtrl.GetAllPuskesmas)
	pr.GET("/:puskesmasId", puskesmasCtrl.GetPuskesmasById)
	pr.PUT("/:puskesmasId/edit", puskesmasCtrl.UpdatePuskesmas)
	pr.DELETE("/:puskesmasId/delete", puskesmasCtrl.DeletePuskesmas)

	pr.GET("/select-kecamatan", puskesmasCtrl.GetSelectKecamatan)
	pr.GET("/select-kelurahan/:kecamatanId", puskesmasCtrl.GetSelectKelurahan)
	pr.GET("/petugas/:puskesmasId", puskesmasCtrl.GetPetugasByPuskesmasId)
}
