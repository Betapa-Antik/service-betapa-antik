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

func PublishPuskesmasPhotoUpload(p payload.PhotoUploadPayload) error {
	conn := configs.GetRabbitConn()
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		rabbitmq.PuskesmasUploadQueue,
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

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
		},
	)
	if err != nil {
		return err
	}
	log.Printf("Published photo upload message for puskesmas %s with %d file(s) to folder %s", p.ID, len(p.Files), p.Folder)
	return nil
}
