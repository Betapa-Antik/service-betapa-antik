package cache

import "betapa-antik-service/internal/models"

type CacheKeluarga struct {
	Keluarga []models.Keluarga `json:"keluarga"`
	Total    int               `json:"total"`
}
