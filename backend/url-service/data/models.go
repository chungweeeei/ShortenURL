package data

import (
	"log"

	"gorm.io/gorm"
)

var db *gorm.DB

func New(dbPool *gorm.DB) Models {
	db = dbPool

	// Do auto migration
	err := db.AutoMigrate(&URL{})
	if err != nil {
		log.Println("Failed to auto migrate urls table")
	}

	return Models{
		URL: &URL{},
	}
}

type Models struct {
	URL URLInterface
}
