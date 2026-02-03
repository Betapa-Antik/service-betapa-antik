package materiroutes

import (
	matericontroller "betapa-antik-service/internal/controllers/materi_controller"
	"betapa-antik-service/internal/middlewares"
	adminrepo "betapa-antik-service/internal/repositories/admin_repo"
	gambarrepo "betapa-antik-service/internal/repositories/gambar_repo"
	materigambarrepo "betapa-antik-service/internal/repositories/materI_gambar_repo"
	materirepo "betapa-antik-service/internal/repositories/materi_repo"
	materiservice "betapa-antik-service/internal/services/materi_service"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func MateriRoutes(e *echo.Group, db *gorm.DB, rdb *redis.Client) {
	repo := adminrepo.NewAdminRepositoryImpl(db)
	materiRepo := materirepo.NewMateriRepositoryImpl(db)
	gambarRepo := gambarrepo.NewGambarRepositoryImpl(db)
	materiGambarRepo := materigambarrepo.NewMateriGambarRepositoryImpl(db)
	svc := materiservice.NewMateriServiceImpl(materiRepo, gambarRepo, materiGambarRepo, rdb)
	ctrl := matericontroller.NewMateriController(svc)

	rm := e.Group("", middlewares.RequireRole("admin", repo))
	rm.POST("/create", ctrl.CreateMateri)
	rm.GET("/all", ctrl.GetAllMateri)
	rm.GET("/:materiId", ctrl.GetByID)
	rm.PUT("/:materiId/edit", ctrl.UpdateMateri)
	rm.PATCH("/:materiId/update-status", ctrl.UpdateStatusMateri)
	rm.DELETE("/:materiId/delete", ctrl.DeleteMateri)
}
