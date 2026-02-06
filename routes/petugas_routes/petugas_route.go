package petugasroutes

import (
	petugascontroller "betapa-antik-service/internal/controllers/petugas_controller"
	petugasrepo "betapa-antik-service/internal/repositories/petugas_repo"
	rolerepo "betapa-antik-service/internal/repositories/role_repo"
	petugasservice "betapa-antik-service/internal/services/petugas_service"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func PetugasRoutes(e *echo.Group, db *gorm.DB, rdb *redis.Client) {
	petugasRepo := petugasrepo.NewPetugasRepositoryImpl(db)
	roleRepo := rolerepo.NewRoleRepositoryImpl(db)
	petugasService := petugasservice.NewPetugasServiceImpl(petugasRepo, roleRepo, rdb)
	petugasCtrl := petugascontroller.NewPetugasController(petugasService)

	e.GET("/select-puskesmas", petugasCtrl.GetSelectPuskesmas)
	e.POST("/register", petugasCtrl.RegisterAkunPetugas)
}
