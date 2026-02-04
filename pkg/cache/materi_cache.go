package cache

import "betapa-antik-service/internal/models"

type CacheMateri struct {
	Materies []*models.Materi `json:"materies"`
	Total    int              `json:"total"`
}
