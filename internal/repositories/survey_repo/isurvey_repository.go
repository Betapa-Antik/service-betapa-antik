package surveyrepo

import (
	"betapa-antik-service/internal/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ISurveyRepository interface {
	DB() *gorm.DB
	WithTx(tx *gorm.DB) ISurveyRepository
	GetSelectKeluarga(ctx context.Context, search string) ([]models.SelectKeluarga, error)
	CreateSurvey(ctx context.Context, data *models.Survey) error
	FindOrCreateSurveyLokasi(ctx context.Context, namaLokasi string, jenisSurvey string) (*models.SurveyLokasi, error)
	CreateSurveyItem(ctx context.Context, data *models.SurveyItem) error
	CreateSurveyFollowUpJentik(ctx context.Context, data *models.SurveyFollowUpJentik) error
	CreateSurveyFollowUpNyamuk(ctx context.Context, data *models.SurveyFollowUpNyamuk) error
	CreateGambar(ctx context.Context, data *models.Gambar) error
	CreatePivotSurveyGambar(ctx context.Context, data *models.SurveyGambar) error

	GetAllSurvey(ctx context.Context, limit int, offset int, search string, jenisSurvey string, startDate, endDate string, petugasId uuid.UUID) ([]models.Survey, int, error)
	GetSurveyByID(ctx context.Context, surveyId uuid.UUID) (*models.Survey, error)
	UpdateSurvey(ctx context.Context, surveyId uuid.UUID, updates map[string]interface{}) error
	DeleteSurvey(ctx context.Context, surveyId uuid.UUID) error
}
