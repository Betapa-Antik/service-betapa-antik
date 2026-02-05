package cache

import "betapa-antik-service/internal/models"

type CachePuskesmas struct {
	Puskesmas []models.PuskesmasWithTotal `json:"puskesmas"`
	Total     int                         `json:"total"`
}
