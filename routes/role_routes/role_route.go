package roleroutes

import (
	rolecontroller "betapa-antik-service/internal/controllers/role_controller"
	rolerepo "betapa-antik-service/internal/repositories/role_repo"
	roleservice "betapa-antik-service/internal/services/role_service"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RoleRoutes(e *echo.Group, db *gorm.DB) {
	// initialize repository, service, controller
	repo := rolerepo.NewRoleRepositoryImpl(db)
	svc := roleservice.NewRoleServiceImpl(repo)
	ctrl := rolecontroller.NewRoleController(svc)

	e.POST("/create", ctrl.Create)
	e.GET("/all", ctrl.FindAll)
	e.GET("/:roleId", ctrl.FindByID)
	e.PUT("/:roleId/edit", ctrl.Update)
	e.DELETE("/:roleId/delete", ctrl.Delete)
}
