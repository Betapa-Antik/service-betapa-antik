package videoroutes

import (
	videocontroller "betapa-antik-service/internal/controllers/video_controller"
	"betapa-antik-service/internal/middlewares"
	adminrepo "betapa-antik-service/internal/repositories/admin_repo"
	videorepo "betapa-antik-service/internal/repositories/video_repo"
	videoservice "betapa-antik-service/internal/services/video_service"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func VideoRoutes(e *echo.Group, db *gorm.DB, rdb *redis.Client) {
	videoRepo := videorepo.NewVideoRepositoryImpl(db)
	repo := adminrepo.NewAdminRepositoryImpl(db)
	videoService := videoservice.NewVideoServiceImpl(videoRepo, rdb)
	videoCtrl := videocontroller.NewVideoController(videoService)

	vr := e.Group("", middlewares.RequireRole("admin", repo))
	vr.POST("/create", videoCtrl.CreateVideo)
	vr.GET("/all", videoCtrl.GetAllVideo)
	vr.GET("/:videoId", videoCtrl.GetVideoById)
	vr.PUT("/:videoId/edit", videoCtrl.UpdateVideo)
	vr.PATCH("/:videoId/edit-status", videoCtrl.UpdateStatusVideo)
	vr.DELETE("/:videoId/delete", videoCtrl.DeleteVideo)
}
