package consumer

import (
	"betapa-antik-service/configs"
	datasource "betapa-antik-service/internal/dataSource"
	petugasrepo "betapa-antik-service/internal/repositories/petugas_repo"
	rabbitmq "betapa-antik-service/pkg/constant/rabbitMQ"
	"betapa-antik-service/pkg/workers/payload"
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

func PetugasPhotoConsumer(ctx context.Context, db *gorm.DB, cld datasource.CloudinaryService) error {
	petugasRepo := petugasrepo.NewPetugasRepositoryImpl(db)

	go func() {
		for {
			conn := configs.GetRabbitConn()
			if conn == nil || conn.IsClosed() {
				time.Sleep(5 * time.Second)
				continue
			}

			ch, err := conn.Channel()
			if err != nil {
				log.Printf("Failed to open channel: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			_ = ch.Qos(1, 0, false)

			msgs, err := ch.Consume(
				rabbitmq.PetugasUploadQueue,
				"", false, false, false, false, nil,
			)
			if err != nil {
				ch.Close()
				continue
			}

			log.Println("[RabbitMQ] petugas consumer started...")

			for d := range msgs {
				handlePetugasUpload(ctx, d, petugasRepo, cld)
			}

			log.Println("[RABBITMQ] channel closed, retrying in 5s...")
			ch.Close()
			time.Sleep(5 * time.Second)
		}
	}()
	return nil
}

func handlePetugasUpload(ctx context.Context, d amqp.Delivery, repo petugasrepo.IPetugasRepository, cld datasource.CloudinaryService) {
	var p payload.PhotoUploadPayload
	if err := json.Unmarshal(d.Body, &p); err != nil {
		log.Printf("Invalid payload: %v", err)
		_ = d.Nack(false, false)
		return
	}

	var lastPhotoURL string
	for _, file := range p.Files {
		f, err := os.Open(file.Path)
		if err != nil {
			log.Printf("File not found: %s", file.Path)
			continue
		}

		res, err := cld.UploadFromReader(ctx, f, p.Folder, file.Filename)
		f.Close()

		if err == nil {
			_ = os.Remove(file.Path)
			lastPhotoURL = res.URL
		}
	}

	updates := map[string]interface{}{}
	if lastPhotoURL != "" {
		updates["foto"] = lastPhotoURL
		err := repo.UpdateAkunPetugas(ctx, p.ID, updates)
		if err != nil {
			log.Printf("DB Update failed: %v", err)
			_ = d.Nack(false, true)
			return
		}

		_ = configs.DeleteRedis(ctx, "profile:"+p.ID.String())
	}

	_ = d.Ack(false)
}
