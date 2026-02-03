package videoservice

import (
	videoreqeuest "betapa-antik-service/internal/dto/request/video_reqeuest"
	"betapa-antik-service/internal/models"
	"context"

	"github.com/google/uuid"
)

type IVideoService interface {
	CreateVideo(ctx context.Context, req videoreqeuest.CreateVideoRequest) error
	GetAllVideo(ctx context.Context, req videoreqeuest.GetAllVideoRequest) ([]*models.Video, int, error)
	GetVideoById(ctx context.Context, videoId uuid.UUID) (*models.Video, error)
	UpdateVideo(ctx context.Context, videoId uuid.UUID, req videoreqeuest.UpdateVideoRequest) error
	UpdateStatusVideo(ctx context.Context, videoId uuid.UUID, req videoreqeuest.UpdateStatusVideoRequest) error
	DeleteVideo(ctx context.Context, videoId uuid.UUID) error
}
