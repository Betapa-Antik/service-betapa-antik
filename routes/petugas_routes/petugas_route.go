package petugasroutes

import (
	petugascontroller "betapa-antik-service/internal/controllers/petugas_controller"
	"betapa-antik-service/internal/middlewares"
	adminrepo "betapa-antik-service/internal/repositories/admin_repo"
	petugasrepo "betapa-antik-service/internal/repositories/petugas_repo"
	rolerepo "betapa-antik-service/internal/repositories/role_repo"
	petugasservice "betapa-antik-service/internal/services/petugas_service"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func PetugasRoutes(e *echo.Group, db *gorm.DB, rdb *redis.Client) {
	adminRepo := adminrepo.NewAdminRepositoryImpl(db)
	petugasRepo := petugasrepo.NewPetugasRepositoryImpl(db)
	roleRepo := rolerepo.NewRoleRepositoryImpl(db)
	petugasService := petugasservice.NewPetugasServiceImpl(petugasRepo, roleRepo, rdb)
	petugasCtrl := petugascontroller.NewPetugasController(petugasService)

	e.GET("/select-puskesmas", petugasCtrl.GetSelectPuskesmas)
	e.POST("/register", petugasCtrl.RegisterAkunPetugas)
	e.POST("/login", petugasCtrl.LoginPetugas)

	pm := e.Group("", middlewares.RequireRole("petugas puskesmas", adminRepo))
	pm.GET("/profile", petugasCtrl.GetProfilePetugas)
	pm.PUT("/profile/update", petugasCtrl.UpdateProfilePetugas)
	pm.POST("/logout", petugasCtrl.LogoutPetugas)
}
