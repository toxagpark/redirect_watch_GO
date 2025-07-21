package redis

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	RDB *redis.Client
)

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Ошибка подключения к Redis: %v", err)
	}

	fmt.Println("Успешное подключение к Redis")
}

func AddURLRedis(original_url, short_url string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := RDB.Set(ctx, short_url, original_url, time.Hour).Err()
	if err != nil {
		return fmt.Errorf("не удалось сохранить URL: %w", err)
	}

	return nil
}

func GetURLRedis(short_url string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	val, err := RDB.Get(ctx, short_url).Result()
	if errors.Is(err, redis.Nil) {
		return ""
	} else if err != nil {
		log.Fatalf("Ошибка Редис: %v", err)
	}
	return val
}
