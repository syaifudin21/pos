package redis

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

var Rdb *redis.Client

func InitRedis() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := os.Getenv("REDIS_DB")

	var db int
	if redisDB != "" {
		var err error
		db, err = strconv.Atoi(redisDB)
		if err != nil {
			log.Printf("Warning: Invalid REDIS_DB value '%s', using default DB 0: %v", redisDB, err)
			db = 0 // Default to DB 0 on error
		}
	}

	Rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := Rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	log.Println("Connected to Redis!")
}

func CloseRedis() {
	if Rdb != nil {
		if err := Rdb.Close(); err != nil {
			log.Printf("Error closing Redis client: %v", err)
		}
	}
}
