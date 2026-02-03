package consumer

import (
	"betapa-antik-service/configs"
	datasource "betapa-antik-service/internal/dataSource"
	rabbitmq "betapa-antik-service/pkg/constant/rabbitMQ"
	"betapa-antik-service/pkg/workers/payload"
	"context"
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

func DeleteImageConsumer(
	ctx context.Context,
	cloudinary datasource.CloudinaryService,
) {
	conn := configs.GetRabbitConn()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	// ✅ Batasi 1 message per worker
	_ = ch.Qos(1, 0, false)

	_, err = ch.QueueDeclare(
		rabbitmq.QueueDeleteImage,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(
		rabbitmq.QueueDeleteImage,
		"",
		false, // 🔥 MANUAL ACK
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	for msg := range msgs {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Println("[DeleteImageConsumer] panic recovered:", r)
					_ = msg.Nack(false, false)
				}
			}()

			var data payload.DeleteImagePayload
			if err := json.Unmarshal(msg.Body, &data); err != nil {
				log.Println("Invalid payload:", err)
				_ = msg.Nack(false, false)
				return
			}

			for _, url := range data.Path {
				if err := cloudinary.DeleteImageByURL(context.Background(), url); err != nil {
					log.Println("Failed delete image:", err)
					// optional: retry logic
				}
			}

			_ = msg.Ack(false)

			_ = configs.DeleteRedis(ctx, "materies:*")
			_ = configs.DeleteRedis(ctx, "materi:*")
		}()
	}
}

func deleteImage(
	msg amqp.Delivery,
	cloudinary datasource.CloudinaryService,
) {
	var data payload.DeleteImagePayload

	if err := json.Unmarshal(msg.Body, &data); err != nil {
		log.Println("Invalid payload:", err)
		_ = msg.Nack(false, false)
		return
	}

	for _, url := range data.Path {
		err := cloudinary.DeleteImageByURL(context.Background(), url)
		if err != nil {
			log.Println("Failed delete image:", err)
		}
	}

	_ = msg.Ack(false)
}
