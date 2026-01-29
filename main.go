package main

import (
	"betapa-antik-service/configs"
	datasource "betapa-antik-service/internal/dataSource"
	"betapa-antik-service/pkg/utils"
	"betapa-antik-service/pkg/workers/consumer"
	"betapa-antik-service/routes"
	"context"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	configs.LoadEnv()

	db := configs.InitDB()
	rdb := configs.InitRedis()

	configs.RunMigrations(db)

	e := echo.New()
	e.Validator = utils.NewCustomValidator()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH},
	}))

	for _, r := range e.Routes() {
		log.Printf("ROUTE %s %s", r.Method, r.Path)
	}

	cloudinarySvc, err := datasource.NewCloudinaryService()
	if err != nil {
		log.Fatalf("Failed to initialize Cloudinary service: %v", err)
	}

	configs.InitRabbitMQ()
	defer configs.CloseConnections()

	// start photo upload consumer worker
	go func() {
		if err := consumer.StartPhotoConsumer(context.Background(), db, cloudinarySvc); err != nil {
			log.Printf("Failed to start photo consumer: %v", err)
		}
	}()

	routes.Routes(e, db, rdb, &cloudinarySvc)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(e.Start(":" + port))
}
