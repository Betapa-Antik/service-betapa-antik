package models

type CurrenLocation struct {
	Wilayah string  `json:"wilayah"`
	Scope   string  `json:"scope"`
	Df      float64 `json:"df"`
	Status  string  `json:"status"`
}
