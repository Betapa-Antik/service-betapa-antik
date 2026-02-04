package cache

import "betapa-antik-service/internal/models"

type CacheKelurahan struct {
	Kelurahans []models.Kelurahan `json:"kelurahans"`
	Total      int                `json:"total"`
}
