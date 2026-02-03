package producers

import (
	"betapa-antik-service/configs"
	rabbitmq "betapa-antik-service/pkg/constant/rabbitMQ"
	"betapa-antik-service/pkg/workers/payload"
	"encoding/json"

	"github.com/streadway/amqp"
)

func PublishDeleteImageAsync(imageURL []string) {
	ch, err := configs.GetRabbitChannel()
	if err != nil {
		return
	}

	// ✅ DECLARE QUEUE
	_, err = ch.QueueDeclare(
		rabbitmq.QueueDeleteImage,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return
	}

	payload := payload.DeleteImagePayload{
		Path: imageURL,
	}

	body, _ := json.Marshal(payload)

	_ = ch.Publish(
		"",
		rabbitmq.QueueDeleteImage,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
