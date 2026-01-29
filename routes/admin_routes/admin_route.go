package adminroutes

import (
	admincontroller "betapa-antik-service/internal/controllers/admin_controller"
	adminrepo "betapa-antik-service/internal/repositories/admin_repo"
	rolerepo "betapa-antik-service/internal/repositories/role_repo"
	adminservice "betapa-antik-service/internal/services/admin_service"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func AdminRoutes(e *echo.Group, db *gorm.DB) {
	repo := adminrepo.NewAdminRepositoryImpl(db)
	rrepo := rolerepo.NewRoleRepositoryImpl(db)
	svc := adminservice.NewAdminServiceImpl(repo, rrepo, db)
	ctrl := admincontroller.NewAdminController(svc)

	e.POST("/register", ctrl.Register)
}
