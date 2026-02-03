package configs

import (
	"log"
	"os"
)

type CloudinaryConfig struct {
	CloudName string
	APIKey    string
	APISecret string
}

func LoadCloudinaryConfig() *CloudinaryConfig {
	cfg := &CloudinaryConfig{
		CloudName: os.Getenv("CLOUDINARY_CLOUD_NAME"),
		APIKey:    os.Getenv("CLOUDINARY_API_KEY"),
		APISecret: os.Getenv("CLOUDINARY_API_SECRET"),
	}

	if cfg.CloudName == "" || cfg.APIKey == "" || cfg.APISecret == "" {
		log.Fatal("Cloudinary environment variables are not set")
	}

	return cfg
}
