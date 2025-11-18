package main

import (
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

func initRedis() *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:6379", os.Getenv("REDIS_HOST")),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	if rdb == nil {
		log.Panic("Failed to connect to Redis.")
	}

	fmt.Println("Connected to Redis server successfully")

	return rdb
}
