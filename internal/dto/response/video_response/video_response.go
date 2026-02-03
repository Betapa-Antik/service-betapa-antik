package videoresponse

import (
	"betapa-antik-service/internal/models"
	"betapa-antik-service/pkg/utils"

	"github.com/google/uuid"
)

type VideoResponse struct {
	ID        uuid.UUID `json:"id"`
	Judul     string    `json:"judul"`
	Link      string    `json:"link"`
	Deskripsi string    `json:"deskripsi"`
	Status    string    `json:"status"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func ToVideoResponse(video models.Video) VideoResponse {
	return VideoResponse{
		ID:        video.ID,
		Judul:     video.Judul,
		Link:      video.Link,
		Deskripsi: video.Deskripsi,
		Status:    video.Status,
		CreatedAt: utils.FormatDate(video.CreatedAt),
		UpdatedAt: utils.FormatDate(video.UpdatedAt),
	}
}
