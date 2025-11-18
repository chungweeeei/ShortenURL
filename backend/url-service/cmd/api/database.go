package main

import (
	"fmt"
	"log"
	"os"
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

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres sslmode=disable timezone=UTC connect_timeout=5",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to postgres database: %v", err)
		return
	}

	var exists bool
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	err = db.Raw(fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = '%s')", os.Getenv("DB_DATABASE"))).Scan(&exists).Error
	if err != nil {
		log.Printf("Failed to check database existence: %v", err)
		return
	}

	if !exists {
		err = db.Exec(fmt.Sprintf("CREATE DATABASE \"%s\"", os.Getenv("DB_DATABASE"))).Error
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

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DATABASE"),
	)
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
