package start

import (
	datasource "betapa-antik-service/internal/dataSource"
	"betapa-antik-service/pkg/workers/consumer"
	"context"

	"gorm.io/gorm"
)

func StartWorker(db *gorm.DB, cldSvc datasource.CloudinaryService) {
	go consumer.StartPhotoConsumer(context.Background(), db, cldSvc)
	go consumer.StartMateriPhotoConsumer(context.Background(), db, cldSvc)
	go consumer.DeleteImageConsumer(context.Background(), cldSvc)

	select {}
}
