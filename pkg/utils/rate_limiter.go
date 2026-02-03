package utils

import (
	"betapa-antik-service/configs"
	"context"
	"time"
)

const (
	LoginAttemptKeyPrefix = "login:attempt:"
	LoginAttemptLimit     = 3
	LoginAttemptTTL       = time.Minute
)

// IncrementLoginAttempt increments failed login attempts and returns current attempts and TTL remaining
func IncrementLoginAttempt(ctx context.Context, email string) (int, time.Duration, error) {
	key := LoginAttemptKeyPrefix + email
	val, err := configs.RDB.Incr(ctx, key).Result()
	if err != nil {
		return 0, 0, err
	}
	// Start rate-limit TTL only when attempts reach the configured limit
	if val >= LoginAttemptLimit {
		if err := configs.RDB.Expire(ctx, key, LoginAttemptTTL).Err(); err != nil {
			return int(val), 0, err
		}
	}
	ttl, err := configs.RDB.TTL(ctx, key).Result()
	if err != nil {
		return int(val), 0, err
	}
	return int(val), ttl, nil
}

// ResetLoginAttempt resets attempts for email
func ResetLoginAttempt(ctx context.Context, email string) error {
	key := LoginAttemptKeyPrefix + email
	return configs.RDB.Del(ctx, key).Err()
}

// IsLoginAllowed checks whether login allowed; returns allowed, attempts, ttlLeft
func IsLoginAllowed(ctx context.Context, email string) (bool, int, time.Duration, error) {
	key := LoginAttemptKeyPrefix + email
	val, err := configs.RDB.Get(ctx, key).Int()
	if err != nil {
		// if key not found, allowed
		return true, 0, 0, nil
	}
	if val >= LoginAttemptLimit {
		ttl, err := configs.RDB.TTL(ctx, key).Result()
		if err != nil {
			return false, val, 0, err
		}
		return false, val, ttl, nil
	}
	return true, val, 0, nil
}
