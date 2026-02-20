package configs

import (
	"context"
	"errors"
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
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,

		DialTimeout:  10 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,

		PoolSize:     20,
		MinIdleConns: 5,

		MaxRetries:      5,
		MinRetryBackoff: 500 * time.Millisecond,
		MaxRetryBackoff: 3 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RDB.Ping(ctx).Err(); err != nil {
		log.Fatalf("❌ Redis ping failed: %v", err)
	}

	log.Println("✅ Redis connected")
	return RDB
}

func SetRedis(ctx context.Context, key string, value interface{}, duration time.Duration) error {
	return RDB.Set(ctx, key, value, duration).Err()
}

// Get value berdasarkan key
func GetRedis(ctx context.Context, key string) (string, error) {
	if RDB == nil {
		return "", errors.New("redis not initialized")
	}
	return RDB.Get(ctx, key).Result()
}

// Delete key dari Redis
func DeleteRedis(ctx context.Context, key string) error {
	return RDB.Del(ctx, key).Err()
}
func DeleteByPattern(ctx context.Context, pattern string) error {
	iter := RDB.Scan(ctx, 0, pattern, 0).Iterator()

	for iter.Next(ctx) {
		if err := RDB.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}

	return iter.Err()
}
