package consumer

import (
	"betapa-antik-service/configs"
	datasource "betapa-antik-service/internal/dataSource"
	kecamatanrepo "betapa-antik-service/internal/repositories/kecamatan_repo"
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

func KecamatanPhotoConsumer(ctx context.Context, db *gorm.DB, cld datasource.CloudinaryService) error {
	kecamatanRepo := kecamatanrepo.NewKecamatanRepositoryImpl(db)

	// Bungkus dalam fungsi untuk memudahkan retry saat koneksi hilang
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

			// Set QoS agar consumer tidak kewalahan (stuck)
			_ = ch.Qos(1, 0, false)

			msgs, err := ch.Consume(
				rabbitmq.KecamatanUploadQueue,
				"", false, false, false, false, nil,
			)
			if err != nil {
				ch.Close()
				continue
			}

			log.Println("[RABBITMQ] Kecamatan Consumer started...")

			for d := range msgs {
				// Proses payload...
				handleUpload(ctx, d, kecamatanRepo, cld)
			}

			// Jika loop msgs berhenti, artinya channel/conn bermasalah
			log.Println("[RABBITMQ] Channel closed, retrying in 5s...")
			ch.Close()
			time.Sleep(5 * time.Second)
		}
	}()

	return nil
}

// Pisahkan logic agar bersih
func handleUpload(ctx context.Context, d amqp.Delivery, repo kecamatanrepo.IKecamatanRepository, cld datasource.CloudinaryService) {
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
		// PERBAIKAN: Pastikan di Repository menggunakan:
		// db.Model(&Kecamatan{}).Where("id = ?", id).Update("foto", lastPhotoURL)
		updates["foto"] = lastPhotoURL
		err := repo.Update(ctx, p.ID, updates)
		if err != nil {
			log.Printf("DB Update failed: %v", err)
			_ = d.Nack(false, true) // Requeue agar dicoba lagi
			return
		}
		_ = configs.DeleteRedis(ctx, "profile:"+p.ID.String())
	}

	_ = d.Ack(false)
}
