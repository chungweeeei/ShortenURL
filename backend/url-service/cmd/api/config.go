package main

import (
	"log"

	"github.com/chungweeeei/ShortenURL/data"
	"github.com/chungweeeei/ShortenURL/helpers"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Config struct {
	DB            *gorm.DB
	RDB           *redis.Client
	Generator     *helpers.Generator
	Model         data.Models
	InfoLog       *log.Logger
	Errorlog      *log.Logger
	ErrorChan     chan error
	ErrorDoneChan chan bool
}
