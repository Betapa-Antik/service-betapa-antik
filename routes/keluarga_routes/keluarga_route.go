package keluargaroutes

import (
	keluargacontroller "betapa-antik-service/internal/controllers/keluarga_controller"
	"betapa-antik-service/internal/middlewares"
	adminrepo "betapa-antik-service/internal/repositories/admin_repo"
	keluargarepo "betapa-antik-service/internal/repositories/keluarga_repo"
	keluargaservice "betapa-antik-service/internal/services/keluarga_service"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func KeluargaRoutes(e *echo.Group, db *gorm.DB, rdb *redis.Client) {
	keluargaRepo := keluargarepo.NewKeluargaRepositoryImpl(db)
	adminRepo := adminrepo.NewAdminRepositoryImpl(db)
	keluargaSvc := keluargaservice.NewKeluargaServiceImpl(keluargaRepo, rdb)
	keluargaCtrl := keluargacontroller.NewKeluargaController(keluargaSvc)

	kr := e.Group("", middlewares.RequireRole("petugas puskesmas", adminRepo))
	kr.GET("/select-kecamatan", keluargaCtrl.GetSelectKecamatan)
	kr.GET("/select-kelurahan/:kecamatanId", keluargaCtrl.GetSelectKelurahan)
	kr.POST("/create", keluargaCtrl.CreateKeluarga)
	kr.GET("/all", keluargaCtrl.GetAllKeluarga)
	kr.GET("/:keluargaId", keluargaCtrl.GetKeluargaById)
	kr.PUT("/:keluargaId/edit", keluargaCtrl.UpdateKeluarga)
	kr.DELETE("/:keluargaId/delete", keluargaCtrl.DeleteKeluarga)
}
