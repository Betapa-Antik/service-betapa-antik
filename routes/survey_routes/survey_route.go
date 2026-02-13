package surveyroutes

import (
	surveycontroller "betapa-antik-service/internal/controllers/survey_controller"
	"betapa-antik-service/internal/middlewares"
	adminrepo "betapa-antik-service/internal/repositories/admin_repo"
	surveyrepo "betapa-antik-service/internal/repositories/survey_repo"
	surveyservice "betapa-antik-service/internal/services/survey_service"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SurveyRoutes(e *echo.Group, db *gorm.DB, rdb *redis.Client) {
	adminRepo := adminrepo.NewAdminRepositoryImpl(db)
	surveyRepo := surveyrepo.NewSurveyRepositoryImpl(db)
	surveyService := surveyservice.NewSurveyServiceImpl(surveyRepo, rdb)
	surveyCtrl := surveycontroller.NewSurveyController(surveyService)

	sr := e.Group("", middlewares.RequireRole("petugas puskesmas", adminRepo))
	sr.GET("/select-keluarga", surveyCtrl.GetSelectKeluarga)
	sr.POST("/create", surveyCtrl.CreateSurvey)
	sr.GET("/all", surveyCtrl.GetAllSurvey)
	sr.GET("/:surveyId", surveyCtrl.GetSurveyByID)
	sr.PUT("/update/:surveyId", surveyCtrl.UpdateSurvey)
	sr.DELETE("/delete/:surveyId", surveyCtrl.DeleteSurvey)
}
