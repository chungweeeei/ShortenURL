package data

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type URL struct {
	ID           int64     `json:"id" gorm:"primaryKey;BIGSERIAL"`
	ShortCode    string    `json:"short_code" gorm:"unique;not null;VARCHAR(10)"`
	OriginalURL  string    `json:"original_url" gorm:"unique;not null;TEXT"`
	CreatedAt    time.Time `json:"created_at" gorm:"type:timestamp without time zone"`
	LastAccessAt time.Time `json:"last_access_at" gorm:"type:timestamp without time zone"`
}

func (u *URL) Insert(url URL) (int64, error) {

	result := db.Create(&url)

	if result.Error != nil {
		return 0, errors.New("failed to insert url")
	}

	return url.ID, nil
}

func (u *URL) Update(url URL) error {

	result := db.Save(&url)

	if result.Error != nil {
		return errors.New("failed to update url")
	}

	return nil
}

func (u *URL) GetByOriginalURL(originalURL string) (*URL, error) {

	url := URL{}
	result := db.Where("original_url = ?", originalURL).First(&url)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("url with original url %s not found", originalURL)
		}
		return nil, errors.New("failed to retrieve url by original url")
	}

	return &url, nil
}

func (u *URL) GetByShortCode(shortCode string) (*URL, error) {

	url := URL{}
	result := db.Where("short_code = ?", shortCode).First(&url)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("url with short code %s not found", shortCode)
		}
		return nil, errors.New("failed to retrieve url by short code")
	}

	return &url, nil
}
