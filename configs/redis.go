package configs

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()
var RDB *redis.Client

func InitRedis() *redis.Client {
	RDB = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Username: "default",
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	_, err := RDB.Ping(context.Background()).Result()
	if err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}
	log.Println("Redis Connected!")
	return RDB
}

func SetRedis(ctx context.Context, key string, value interface{}, duration time.Duration) error {
	return RDB.Set(ctx, key, value, duration).Err()
}

// Get value berdasarkan key
func GetRedis(ctx context.Context, key string) (string, error) {
	return RDB.Get(ctx, key).Result()
}

// Delete key dari Redis
func DeleteRedis(ctx context.Context, key string) error {
	return RDB.Del(ctx, key).Err()
}
