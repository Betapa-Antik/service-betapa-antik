package routes

import (
	datasource "betapa-antik-service/internal/dataSource"
	adminroutes "betapa-antik-service/routes/admin_routes"
	kecamatanroutes "betapa-antik-service/routes/kecamatan_routes"
	kelurahanroute "betapa-antik-service/routes/kelurahan_route"
	materiroutes "betapa-antik-service/routes/materi_routes"
	petugasroutes "betapa-antik-service/routes/petugas_routes"
	puskesmasroutes "betapa-antik-service/routes/puskesmas_routes"
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
	//register kecamatan routes
	kecamatanroutes.KecamatanRoutes(v1.Group("/kecamatan"), db, rdb)
	//register kelurahan routes
	kelurahanroute.KelurahanRoutes(v1.Group("/kelurahan"), db, rdb)
	//register puskesmas routes
	puskesmasroutes.PuskesmasRoutes(v1.Group("/puskesmas"), db, rdb)
	//register petugas routes
	petugasroutes.PetugasRoutes(v1.Group("/petugas"), db, rdb)
}
