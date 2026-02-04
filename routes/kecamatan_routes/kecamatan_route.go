package kecamatanroutes

import (
	kecamatancontroller "betapa-antik-service/internal/controllers/kecamatan_controller"
	"betapa-antik-service/internal/middlewares"
	adminrepo "betapa-antik-service/internal/repositories/admin_repo"
	kecamatanrepo "betapa-antik-service/internal/repositories/kecamatan_repo"
	kecamatanservice "betapa-antik-service/internal/services/kecamatan_service"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func KecamatanRoutes(e *echo.Group, db *gorm.DB, rdb *redis.Client) {
	adminRepo := adminrepo.NewAdminRepositoryImpl(db)
	kecamatanRepo := kecamatanrepo.NewKecamatanRepositoryImpl(db)
	kecamatanService := kecamatanservice.NewKecamatanServiceImpl(kecamatanRepo, rdb)
	kecamatanCtrl := kecamatancontroller.NewKecamatanController(kecamatanService)

	kr := e.Group("", middlewares.RequireRole("admin", adminRepo))
	kr.POST("/create", kecamatanCtrl.CreateKecamatan)
	kr.GET("/all", kecamatanCtrl.GetAllKecamatan)
	kr.GET("/:kecamatanId", kecamatanCtrl.GetKecamatanById)
	kr.PUT("/:kecamatanId/edit", kecamatanCtrl.UpdateKecamatan)
	kr.DELETE("/:kecamatanId/delete", kecamatanCtrl.DeleteKecamatan)
}
