package routes

import (
	datasource "betapa-antik-service/internal/dataSource"
	adminroutes "betapa-antik-service/routes/admin_routes"
	materiroutes "betapa-antik-service/routes/materi_routes"
	roleroutes "betapa-antik-service/routes/role_routes"
	videoroutes "betapa-antik-service/routes/video_routes"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Routes(e *echo.Echo, db *gorm.DB, rdb *redis.Client, cldSvc *datasource.CloudinaryService) {
	v1 := e.Group("/betapa-antik/api/v1")
	// register role routes
	roleroutes.RoleRoutes(v1.Group("/role"), db)
	// register admin routes
	adminroutes.AdminRoutes(v1.Group("/admin"), db)
	// register materi routes
	materiroutes.MateriRoutes(v1.Group("/materi"), db, rdb)
	//register video routes
	videoroutes.VideoRoutes(v1.Group("/video"), db, rdb)
}
