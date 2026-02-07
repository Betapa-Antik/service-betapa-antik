package adminroutes

import (
	admincontroller "betapa-antik-service/internal/controllers/admin_controller"
	"betapa-antik-service/internal/middlewares"
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
	e.POST("/login", ctrl.Login)
	// protected profile routes
	// only users with role "admin" can access
	e.GET("/profile", ctrl.Profile, middlewares.RequireRole("admin", repo))
	e.PUT("/profile", ctrl.UpdateProfile, middlewares.RequireRole("admin", repo))
	e.PUT("/profile/photo", ctrl.UpdateProfilePhoto, middlewares.RequireRole("admin", repo))
	// logout
	e.POST("/logout", ctrl.Logout, middlewares.RequireRole("admin", repo))

	am := e.Group("", middlewares.RequireRole("admin", repo))
	//Manage Petugas
	am.GET("/petugas", ctrl.FindPetugas)
	am.PUT("/approve-or-reject-petugas/:petugasId", ctrl.ApproveOrRejecAkunPetugas)
	am.PUT("/active-or-nonactive-petugas/:petugasId", ctrl.ActiveOrNonActiveAkunPetugas)
}
