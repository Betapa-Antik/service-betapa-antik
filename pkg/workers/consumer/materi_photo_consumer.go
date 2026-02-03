package consumer

import (
	"betapa-antik-service/configs"
	datasource "betapa-antik-service/internal/dataSource"
	"betapa-antik-service/internal/models"
	materirepo "betapa-antik-service/internal/repositories/materi_repo"
	queueConst "betapa-antik-service/pkg/constant/rabbitMQ"
	"betapa-antik-service/pkg/utils"
	"betapa-antik-service/pkg/workers/payload"
	"bytes"
	"context"
	"encoding/json"
	"os"

	"gorm.io/gorm"
)

func StartMateriPhotoConsumer(ctx context.Context, db *gorm.DB, cloud datasource.CloudinaryService) error {
	conn := configs.GetRabbitConn()
	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		queueConst.MateriImageUploadQueue,
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

	materiRepo := materirepo.NewMateriRepositoryImpl(db)

	go func() {
		for d := range msgs {
			var p payload.PhotoUploadPayload
			if err := json.Unmarshal(d.Body, &p); err != nil {
				continue
			}

			for _, file := range p.Files {
				go func(f payload.PhotoFile) {
					data, err := os.ReadFile(f.Path)
					if err != nil {
						return
					}

					res, err := cloud.UploadFromReader(
						ctx,
						bytes.NewReader(data),
						p.Folder,
						f.Filename,
					)
					if err != nil {
						return
					}

					_ = os.Remove(f.Path) // 🔥 bersihin tmp

					// DB TRANSACTION
					_ = utils.RunInTransaction(db, func(tx *gorm.DB) error {
						repo := materiRepo.WithTx(tx)

						gambar := &models.Gambar{Path: res.URL}
						if err := repo.CreateGambar(ctx, gambar); err != nil {
							return err
						}

						return repo.CreatePivotMateriGambar(ctx, &models.MateriGambar{
							MateriID: p.ID,
							GambarID: gambar.ID,
						})
					})
				}(file)
				_ = configs.DeleteRedis(ctx, "materies:*")
				_ = configs.DeleteRedis(ctx, "materi:"+p.ID.String())
			}
		}

	}()

	return nil
}
