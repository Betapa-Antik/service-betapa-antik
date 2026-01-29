package consumer

import (
	"context"
	"encoding/json"
	"io"
	"log"

	"betapa-antik-service/configs"
	datasource "betapa-antik-service/internal/dataSource"
	"betapa-antik-service/internal/models"
	adminrepo "betapa-antik-service/internal/repositories/admin_repo"
	queueConst "betapa-antik-service/pkg/constant/rabbitMQ"
	"betapa-antik-service/pkg/workers/payload"

	"gorm.io/gorm"
)

func StartPhotoConsumer(ctx context.Context, db *gorm.DB, cloud datasource.CloudinaryService) error {
	conn := configs.GetRabbitConn()
	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
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

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	adminRepo := adminrepo.NewAdminRepositoryImpl(db)

	go func() {
		for d := range msgs {
			var p payload.PhotoUploadPayload
			if err := json.Unmarshal(d.Body, &p); err != nil {
				log.Printf("failed to unmarshal photo payload: %v", err)
				continue
			}

			var photoURL string
			// upload all files in payload
			for _, file := range p.Files {
				res, err := cloud.UploadImageBytes(ctx, io.NopCloser(bytesReader(file.Data)), p.Folder, file.Filename)
				if err != nil {
					log.Printf("cloud upload failed for %s: %v", file.Filename, err)
					continue
				}
				photoURL = res.URL // keep last successful upload URL for user.Foto
				log.Printf("uploaded %s to %s, url=%s", file.Filename, p.Folder, res.URL)
			}

			// update user Foto field with last uploaded URL (for primary photo)
			if photoURL != "" {
				if err := adminRepo.Update(ctx, p.UserID, &models.User{Foto: photoURL}); err != nil {
					log.Printf("failed to update user foto: %v", err)
					continue
				}
				log.Printf("processed photo uploads for user %s, primary_url=%s", p.UserID, photoURL)
			}
		}
	}()

	return nil
}

// bytesReader provides an io.Reader from a byte slice
func bytesReader(b []byte) io.Reader {
	return &reader{b: b}
}

type reader struct{ b []byte }

func (r *reader) Read(p []byte) (int, error) {
	if len(r.b) == 0 {
		return 0, io.EOF
	}
	n := copy(p, r.b)
	r.b = r.b[n:]
	return n, nil
}
