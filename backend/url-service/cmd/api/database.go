package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {

	ensureDatabaseExists()

	conn := connectToDB()
	if conn == nil {
		log.Panic("Can not connect to database.")
	}

	return conn
}

func ensureDatabaseExists() {

	dsn := "host=localhost user=root password=root dbname=postgres sslmode=disable timezone=UTC connect_timeout=5"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to postgres database: %v", err)
		return
	}

	var exists bool
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	err = db.Raw("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = 'URLs')").Scan(&exists).Error
	if err != nil {
		log.Printf("Failed to check database existence: %v", err)
		return
	}

	if !exists {
		err = db.Exec("CREATE DATABASE \"URLs\"").Error
		if err != nil {
			log.Printf("Failed to create database: %v", err)
		} else {
			fmt.Println("Database 'URLs' created successfully")
		}
	} else {
		fmt.Println("Database 'URLs' already exists")
	}
}

func connectToDB() *gorm.DB {

	count := 0

	dsn := "host=localhost user=root password=root dbname=URLs sslmode=disable timezone=UTC connect_timeout=5"

	for {
		connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			fmt.Println("Postgres not yet ready, retrying...")
		} else {
			DB, err := connection.DB()
			if err != nil {
				fmt.Println("Failed connect to database")
				return nil
			}

			DB.SetMaxIdleConns(5)
			DB.SetConnMaxLifetime(30 * time.Minute)

			fmt.Println("Connected to Postgres database successfully")
			return connection
		}

		if count > 10 {
			return nil
		}

		fmt.Println("Backing off for 1 second")
		time.Sleep(1 * time.Second)
		count++
	}
}
