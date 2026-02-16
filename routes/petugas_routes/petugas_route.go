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
	pm.PUT("/ubah-kata-sandi", petugasCtrl.UbahKataSandiPetugas)

	e.POST("/lupa-kata-sandi-request", petugasCtrl.LupaKataSandi)
	e.GET("/status-lupa-kata-sandi/:logId", petugasCtrl.StatusVerifikasiLupaKataSandi)
	e.PUT("/atur-ulang-kata-sandi/:petugasId", petugasCtrl.AturUlangKataSandi)

	pm.GET("/laporan/all", petugasCtrl.GetAllLaporan)
	pm.GET("/laporan/:laporanId", petugasCtrl.GetLaporanByID)
	pm.PUT("/laporan/:laporanId/update-status", petugasCtrl.UpdateStatusLaporan)

	pm.GET("/dashboard", petugasCtrl.GetDashboard)
	pm.GET("/laporan-latest", petugasCtrl.GetLatestLaporanByPuskesmasID)
	pm.GET("/survey-latest", petugasCtrl.GetLatestSurveyByPetugasID)
}
