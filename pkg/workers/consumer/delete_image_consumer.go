package consumer

import (
	"betapa-antik-service/configs"
	datasource "betapa-antik-service/internal/dataSource"
	rabbitmq "betapa-antik-service/pkg/constant/rabbitMQ"
	"betapa-antik-service/pkg/workers/payload"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/streadway/amqp"
)

func DeleteImageConsumer(
	ctx context.Context,
	cloudinary datasource.CloudinaryService,
) {
	// Kita jalankan dalam goroutine agar tidak mem-block startup aplikasi
	go func() {
		for {
			conn := configs.GetRabbitConn()
			if conn == nil || conn.IsClosed() {
				log.Println("[DeleteImageConsumer] Waiting for RabbitMQ connection...")
				time.Sleep(5 * time.Second)
				continue
			}

			ch, err := conn.Channel()
			if err != nil {
				log.Printf("[DeleteImageConsumer] Failed to open channel: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			// ✅ Batasi 1 message per worker agar tidak overload
			_ = ch.Qos(1, 0, false)

			_, err = ch.QueueDeclare(
				rabbitmq.QueueDeleteImage,
				true, false, false, false, nil,
			)
			if err != nil {
				ch.Close()
				continue
			}

			msgs, err := ch.Consume(
				rabbitmq.QueueDeleteImage,
				"",
				false, // ✅ MANUAL ACK
				false, false, false, nil,
			)
			if err != nil {
				ch.Close()
				continue
			}

			log.Println("[RABBITMQ] Delete Image Consumer is running...")

			// Loop untuk memproses pesan yang masuk
			for msg := range msgs {
				// Panggil fungsi handler terpisah
				handleDeleteProcess(ctx, msg, cloudinary)
			}

			// Jika loop msgs berhenti, berarti channel tutup
			log.Println("[DeleteImageConsumer] Channel closed, reconnecting...")
			ch.Close()
			time.Sleep(5 * time.Second)
		}
	}()
}

func handleDeleteProcess(_ context.Context, msg amqp.Delivery, cloudinary datasource.CloudinaryService) {
	// Gunakan defer untuk memastikan Nack dikirim jika terjadi panic
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[DeleteImageConsumer] Panic recovered: %v", r)
			_ = msg.Nack(false, true) // Requeue agar dicoba lagi
		}
	}()

	var data payload.DeleteImagePayload
	if err := json.Unmarshal(msg.Body, &data); err != nil {
		log.Println("[DeleteImageConsumer] Invalid payload:", err)
		_ = msg.Ack(false) // Ack saja karena data rusak tidak bisa diproses ulang
		return
	}

	// Eksekusi penghapusan di Cloudinary
	successDelete := false
	for _, url := range data.Path {
		if url == "" {
			continue
		}

		// Kita gunakan context background untuk penghapusan agar tidak terikat ctx parent yang mungkin dicancel
		if err := cloudinary.DeleteImageByURL(context.Background(), url); err != nil {
			log.Printf("[DeleteImageConsumer] Failed delete image %s: %v", url, err)
			// Lanjut ke URL berikutnya, jangan return dulu
		} else {
			successDelete = true
		}
	}

	// Invalidate Cache jika ada yang berhasil dihapus
	if successDelete {
		// _ = configs.DeleteRedis(ctx, "materies:*")
		// _ = configs.DeleteRedis(ctx, "materi:*")
	}

	// ✅ Akhiri dengan Ack
	_ = msg.Ack(false)
}
