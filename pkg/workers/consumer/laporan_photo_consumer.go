package consumer

import (
	"betapa-antik-service/configs"
	datasource "betapa-antik-service/internal/dataSource"
	"betapa-antik-service/internal/models"
	masyarakatrepo "betapa-antik-service/internal/repositories/masyarakat_repo"
	queueConst "betapa-antik-service/pkg/constant/rabbitMQ"
	"betapa-antik-service/pkg/utils"
	"betapa-antik-service/pkg/workers/payload"
	"bytes"
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

func StartLaporanPhotoConsumer(ctx context.Context, db *gorm.DB, cloud datasource.CloudinaryService) error {
	masyarakatRepo := masyarakatrepo.NewMasyarakatRepositoryImpl(db)

	go func() {
		for {
			conn := configs.GetRabbitConn()
			if conn == nil || conn.IsClosed() {
				log.Println("[RabbitMQ] Waiting for connection...")
				time.Sleep(5 * time.Second)
				continue
			}

			ch, err := conn.Channel()
			if err != nil {
				time.Sleep(5 * time.Second)
				continue
			}

			// PENTING: Batasi jumlah pesan yang diproses (1 per worker) agar tidak stuck
			_ = ch.Qos(2, 0, false)

			msgs, err := ch.Consume(
				queueConst.LaporanUploadQueue,
				"",    // consumer
				false, // auto-ack diganti ke FALSE agar aman
				false, // exclusive
				false, // no-local
				false, // no-wait
				nil,   // args
			)
			if err != nil {
				ch.Close()
				continue
			}

			log.Println("[RabbitMQ] Laporan Consumer is running...")

			for d := range msgs {
				// Proses satu per satu secara sekuensial atau gunakan worker pool terbatas
				// Menghapus 'go func()' di dalam loop file untuk menghindari race condition
				handleLaporanUpload(ctx, d, db, masyarakatRepo, cloud)
			}

			log.Println("[RabbitMQ] Channel closed, reconnecting...")
			ch.Close()
			time.Sleep(5 * time.Second)
		}
	}()

	return nil
}

func handleLaporanUpload(ctx context.Context, d amqp.Delivery, db *gorm.DB, laporanRepo masyarakatrepo.IMasyarakatRepository, cloud datasource.CloudinaryService) {
	var p payload.PhotoUploadPayload
	if err := json.Unmarshal(d.Body, &p); err != nil {
		log.Printf("Invalid payload: %v", err)
		_ = d.Ack(false)
		return
	}

	successCount := 0
	for _, file := range p.Files {
		// 1. Read File
		data, err := os.ReadFile(file.Path)
		if err != nil {
			log.Printf("Failed read file %s: %v", file.Path, err)
			continue
		}

		// 2. Upload to Cloudinary
		res, err := cloud.UploadFromReader(ctx, bytes.NewReader(data), p.Folder, file.Filename)
		if err != nil {
			log.Printf("Cloudinary error: %v", err)
			continue
		}

		// 3. DB Transaction
		err = utils.RunInTransaction(db, func(tx *gorm.DB) error {
			repo := laporanRepo.WithTx(tx)
			gambar := &models.Gambar{Path: res.URL}
			if err := repo.CreateGambar(ctx, gambar); err != nil {
				return err
			}
			return repo.CreateLaporanGambar(ctx, &models.LaporanGambar{
				LaporanID: p.ID,
				GambarID:  gambar.ID,
			})
		})

		if err == nil {
			_ = os.Remove(file.Path)
			successCount++
		}
	}

	// 5. Akhiri dengan Ack
	_ = d.Ack(false)
}
