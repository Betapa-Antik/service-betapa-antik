package consumer

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"betapa-antik-service/configs"
	datasource "betapa-antik-service/internal/dataSource"
	"betapa-antik-service/internal/models"
	adminrepo "betapa-antik-service/internal/repositories/admin_repo"
	queueConst "betapa-antik-service/pkg/constant/rabbitMQ"
	"betapa-antik-service/pkg/workers/payload"

	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

func StartPhotoConsumer(
	ctx context.Context,
	db *gorm.DB,
	cloud datasource.CloudinaryService,
) error {
	adminRepo := adminrepo.NewAdminRepositoryImpl(db)

	go func() {
		for {
			conn := configs.GetRabbitConn()
			if conn == nil || conn.IsClosed() {
				log.Println("[RABBITMQ] Admin photo consumer waiting for connection...")
				time.Sleep(5 * time.Second)
				continue
			}

			ch, err := conn.Channel()
			if err != nil {
				log.Printf("[RABBITMQ] Failed to open channel: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			// Batasi jumlah pesan yang diproses sekaligus agar tidak memakan RAM/CPU
			_ = ch.Qos(1, 0, false)

			_, err = ch.QueueDeclare(
				queueConst.AdminPhotoUploadQueue,
				true, false, false, false, nil,
			)
			if err != nil {
				ch.Close()
				continue
			}

			msgs, err := ch.Consume(
				queueConst.AdminPhotoUploadQueue,
				"",    // consumer
				false, // ✅ Manual Ack
				false, false, false, nil,
			)
			if err != nil {
				ch.Close()
				continue
			}

			log.Println("[RABBITMQ] Admin Photo Consumer started...")

			for d := range msgs {
				// Proses pesan
				handleAdminPhotoProcess(ctx, d, adminRepo, cloud)
			}

			log.Println("[RABBITMQ] Admin channel closed, retrying in 5s...")
			ch.Close()
			time.Sleep(5 * time.Second)
		}
	}()

	return nil
}

func handleAdminPhotoProcess(ctx context.Context, d amqp.Delivery, repo adminrepo.IAdminRepository, cloud datasource.CloudinaryService) {
	var p payload.PhotoUploadPayload
	if err := json.Unmarshal(d.Body, &p); err != nil {
		log.Printf("Invalid payload: %v", err)
		_ = d.Ack(false) // Buang pesan yang formatnya salah
		return
	}

	var lastPhotoURL string
	for _, file := range p.Files {
		f, err := os.Open(file.Path)
		if err != nil {
			log.Printf("File not found: %s", file.Path)
			continue
		}

		res, err := cloud.UploadFromReader(ctx, f, p.Folder, file.Filename)
		f.Close()

		if err != nil {
			log.Printf("Upload failed for %s: %v", file.Filename, err)
			continue
		}

		_ = os.Remove(file.Path) // 🧹 Cleanup
		lastPhotoURL = res.URL
	}

	if lastPhotoURL != "" {
		// Gunakan Updates() atau Field spesifik di Repo untuk menghindari unique constraint error
		if err := repo.Update(ctx, p.ID, &models.User{Foto: lastPhotoURL}); err != nil {
			log.Printf("DB update failed: %v", err)
			_ = d.Nack(false, true) // Requeue agar dicoba lagi nanti
			return
		}
		_ = configs.DeleteRedis(ctx, "profile:"+p.ID.String())
	}

	_ = d.Ack(false) // ✅ Selesai dengan sukses
}
