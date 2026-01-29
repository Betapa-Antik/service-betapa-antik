package configs

import (
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
)

var RabbitConn *amqp.Connection

func InitRabbitMQ() {
	amqpURL := os.Getenv("AMQP_URL")
	if amqpURL == "" {
		amqpURL = "amqp://guest:guest@localhost:5672/"
	}

	for i := 0; i < 5; i++ {
		conn, err := amqp.Dial(amqpURL)
		if err == nil {
			RabbitConn = conn
			log.Println("RabbitMQ connected successfully")
			return
		}

		log.Printf("Retrying RabbitMQ connection (%d): %v", i+1, err)
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	log.Fatal("Failed to connect RabbitMQ")
}

func GetRabbitConn() *amqp.Connection {
	if RabbitConn == nil || RabbitConn.IsClosed() {
		log.Println("RabbitMQ reconnecting...")
		InitRabbitMQ()
	}
	return RabbitConn
}

func CloseConnections() {
	if RabbitConn != nil {
		_ = RabbitConn.Close()
	}
}
