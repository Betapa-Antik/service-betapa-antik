package cache

import "betapa-antik-service/internal/models"

type SurveyCache struct {
	Survey []models.Survey `json:"survey"`
	Total  int             `json:"total"`
}
