package consumer

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"betapa-antik-service/configs"
	datasource "betapa-antik-service/internal/dataSource"
	"betapa-antik-service/internal/models"
	adminrepo "betapa-antik-service/internal/repositories/admin_repo"
	queueConst "betapa-antik-service/pkg/constant/rabbitMQ"
	"betapa-antik-service/pkg/workers/payload"

	"gorm.io/gorm"
)

func StartPhotoConsumer(
	ctx context.Context,
	db *gorm.DB,
	cloud datasource.CloudinaryService,
) error {

	conn := configs.GetRabbitConn()
	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	_, err = ch.QueueDeclare(
		queueConst.AdminPhotoUploadQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		queueConst.AdminPhotoUploadQueue,
		"",
		false, // ✅ MANUAL ACK
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	adminRepo := adminrepo.NewAdminRepositoryImpl(db)

	go func() {
		for d := range msgs {
			var p payload.PhotoUploadPayload
			if err := json.Unmarshal(d.Body, &p); err != nil {
				log.Printf("invalid payload: %v", err)
				_ = d.Nack(false, false)
				continue
			}

			var lastPhotoURL string

			for _, file := range p.Files {
				f, err := os.Open(file.Path)
				if err != nil {
					log.Printf("file not found: %s", file.Path)
					continue
				}

				res, err := cloud.UploadFromReader(
					ctx,
					f,
					p.Folder,
					file.Filename,
				)
				f.Close()

				if err != nil {
					log.Printf("upload failed: %v", err)
					continue
				}

				_ = os.Remove(file.Path) // 🧹 cleanup tmp
				lastPhotoURL = res.URL

				log.Printf(
					"[PHOTO] uploaded file=%s url=%s",
					file.Filename,
					res.URL,
				)
			}

			if lastPhotoURL != "" {
				if err := adminRepo.Update(
					ctx,
					p.ID,
					&models.User{Foto: lastPhotoURL},
				); err != nil {
					log.Printf("db update failed: %v", err)
					_ = d.Nack(false, true)
					continue
				}

				_ = configs.DeleteRedis(ctx, "profile:"+p.ID.String())
			}

			d.Ack(false) // ✅ SUCCESS
		}
	}()

	return nil
}
