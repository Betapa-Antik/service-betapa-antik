package producers

import (
	"betapa-antik-service/configs"
	rabbitmq "betapa-antik-service/pkg/constant/rabbitMQ"
	"betapa-antik-service/pkg/workers/payload"
	"encoding/json"
	"log"
	"time"

	"github.com/streadway/amqp"
)

func PublishMateriPhotoUploadAsync(p payload.PhotoUploadPayload) {
	go func() {
		if err := publishMateriPhotoUpload(p); err != nil {
			log.Printf("[RABBIT] publish failed: %v", err)
		}
	}()
}

func publishMateriPhotoUpload(p payload.PhotoUploadPayload) error {
	ch, err := configs.GetRabbitChannel()
	if err != nil {
		return err
	}
	defer ch.Close() // ⬅️ PENTING

	q, err := ch.QueueDeclare(
		rabbitmq.MateriImageUploadQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	body, err := json.Marshal(p)
	if err != nil {
		return err
	}

	return ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			Body:         body,
		},
	)
}
