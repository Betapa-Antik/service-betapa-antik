package videorepo

import (
	"betapa-antik-service/internal/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IVideoRepository interface {
	DB() *gorm.DB
	WithTx(tx *gorm.DB) IVideoRepository
	CreateVideo(ctx context.Context, data *models.Video) error
	GetAllVideo(ctx context.Context, limit, offset int, search string) ([]*models.Video, int, error)
	GetVideoById(ctx context.Context, videoId uuid.UUID) (*models.Video, error)
	UpdateVideo(ctx context.Context, videoId uuid.UUID, updates map[string]interface{}) error
	UpdateStatusVideo(ctx context.Context, videoId uuid.UUID, status string) error
	DeleteVideoById(ctx context.Context, videoId uuid.UUID) error
}
