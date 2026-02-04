package kelurahanroute

import (
	kelurahancontroller "betapa-antik-service/internal/controllers/kelurahan_controller"
	"betapa-antik-service/internal/middlewares"
	adminrepo "betapa-antik-service/internal/repositories/admin_repo"
	kelurahanrepo "betapa-antik-service/internal/repositories/kelurahan_repo"
	kelurahanservice "betapa-antik-service/internal/services/kelurahan_service"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func KelurahanRoutes(e *echo.Group, db *gorm.DB, rdb *redis.Client) {
	adminRepo := adminrepo.NewAdminRepositoryImpl(db)
	kelRepo := kelurahanrepo.NewKelurahanRepositoryImpl(db)
	kelSvc := kelurahanservice.NewKelurahanServiceImpl(kelRepo, rdb)
	kelCtrl := kelurahancontroller.NewKelurahanController(kelSvc)

	kr := e.Group("", middlewares.RequireRole("admin", adminRepo))
	kr.POST("/create", kelCtrl.CreateKelurahan)
	kr.GET("/all/:kecamatanId", kelCtrl.GetAllKelurahan)
	kr.GET("/:kelurahanId", kelCtrl.GetKelurahanById)
	kr.PUT("/:kelurahanId/edit", kelCtrl.UpdateKelurahan)
	kr.DELETE("/:kelurahanId/delete", kelCtrl.DeleteKelurahan)
}
