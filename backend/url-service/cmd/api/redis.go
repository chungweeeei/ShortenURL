package main

import (
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func initRedis() *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	if rdb == nil {
		log.Panic("Failed to connect to Redis.")
	}

	fmt.Println("Connected to Redis server successfully")

	return rdb
}
