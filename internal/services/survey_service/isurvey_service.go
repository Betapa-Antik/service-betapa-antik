package surveyservice

import (
	surveyrequest "betapa-antik-service/internal/dto/request/survey_request"
	"betapa-antik-service/internal/models"
	"context"

	"github.com/google/uuid"
)

type ISurveyService interface {
	GetSelectKeluarga(ctx context.Context, search string) ([]models.SelectKeluarga, error)
	CreateSurvey(ctx context.Context, petugasId uuid.UUID, req surveyrequest.CreateSurveyRequest) error
	GetAllSurvey(ctx context.Context, req surveyrequest.GetAllSurveyRequest, petugasId uuid.UUID) ([]models.Survey, int, error)
	GetSurveyByID(ctx context.Context, surveyId uuid.UUID) (*models.Survey, error)
	UpdateSurvey(ctx context.Context, surveyId uuid.UUID, req surveyrequest.UpdateSurveyRequest) error
	DeleteSurvey(ctx context.Context, surveyId uuid.UUID) error
}
