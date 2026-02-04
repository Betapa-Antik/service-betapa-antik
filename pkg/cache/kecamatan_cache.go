package cache

import "betapa-antik-service/internal/models"

type CacheKecamatan struct {
	Kecamatans []models.KecamatanWithTotal `json:"kecamatans"`
	Total      int                         `json:"total"`
}
