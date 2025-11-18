package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/chungweeeei/ShortenURL/data"
	"github.com/gin-gonic/gin"
)

type ShortenRequest struct {
	URL string `json:"url" binding:"required,url" example:"https://example.com"`
}

func (app *Config) ShortenURL(c *gin.Context) {

	var request ShortenRequest
	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"short_url": "",
			"message":   "Failed to parse request body.",
		})
		return
	}

	// before generate shorten URL, check if the original URL already exists
	existingURL, _ := app.Model.URL.GetByOriginalURL(request.URL)
	if existingURL != nil {
		c.JSON(http.StatusCreated, gin.H{
			"short_url": fmt.Sprintf("http://localhost:%s/%s", forwardingPort, existingURL.ShortCode),
			"message":   "URL shortened successfully.",
		})
		return
	}

	if app.Generator == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"short_url": "",
			"message":   "ID generator has not initialized.",
		})
		return
	}

	IdResponse, err := app.Generator.GenerateShortURL()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"short_url": "",
			"message":   "Failed to generate short URL.",
		})
		return
	}

	// Insert data into database
	currentTime := time.Now()
	_, err = app.Model.URL.Insert(data.URL{
		ID:           IdResponse.ID,
		ShortCode:    IdResponse.ShortURL,
		OriginalURL:  request.URL,
		CreatedAt:    currentTime,
		LastAccessAt: currentTime,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"short_url": "",
			"message":   "Failed to store URL mapping.",
		})
		return
	}

	// Insert data to Redis cache
	expiredTime := 5 * time.Minute
	err = app.RDB.Set(c, IdResponse.ShortURL, request.URL, expiredTime).Err()
	if err != nil {
		app.Errorlog.Printf("Failed to cache URL in Redis: %v", err)
	}

	c.JSON(http.StatusCreated, gin.H{
		"short_url": fmt.Sprintf("http://localhost:%s/%s", forwardingPort, IdResponse.ShortURL),
		"message":   "URL shortened successfully.",
	})
}

func (app *Config) RedirectLongURL(c *gin.Context) {

	shortCode := c.Param("shortCode")

	// check Redis cache first
	cachedURL, err := app.RDB.Get(c, shortCode).Result()
	if err == nil {
		// cache hit
		app.InfoLog.Printf("Cache hit for short code: %s", shortCode)

		// update latest access time in database asynchronously
		go func() {
			url, err := app.Model.URL.GetByShortCode(shortCode)
			if err != nil {
				app.Errorlog.Printf("Short URL not found during async update: %s, error: %v", shortCode, err)
				return
			}
			url.LastAccessAt = time.Now()
			err = app.Model.URL.Update(*url)
			if err != nil {
				app.Errorlog.Printf("Failed to update last access time for short code %s: %v", shortCode, err)
			}
		}()

		c.Redirect(http.StatusFound, cachedURL)
		return
	}

	// cache miss, query from database
	url, err := app.Model.URL.GetByShortCode(shortCode)
	if err != nil {
		app.Errorlog.Printf("Short URL not found: %s, error: %v", shortCode, err)
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Short URL not found.",
		})
		return
	}

	// update latest access time asynchronously
	go func() {

		// update Redis cache
		expiredTime := 5 * time.Minute
		err = app.RDB.Set(c, shortCode, url.OriginalURL, expiredTime).Err()
		if err != nil {
			app.Errorlog.Printf("Failed to cache URL in Redis: %v", err)
		}

		// update latest access time in database
		url.LastAccessAt = time.Now()
		err = app.Model.URL.Update(*url)
		if err != nil {
			app.Errorlog.Printf("Failed to update last access time for short code %s: %v", shortCode, err)
		}
	}()

	// 301 永久重定向 或 302 臨時重定向
	c.Redirect(http.StatusFound, url.OriginalURL)
}

func (app *Config) GetLongURL(c *gin.Context) {

	shortCode := c.Param("shortCode")

	url, err := app.Model.URL.GetByShortCode(shortCode)
	if err != nil {
		app.Errorlog.Printf("Short URL not found: %s, error: %v", shortCode, err)
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Short URL not found.",
		})
		return
	}

	// update latest access time
	url.LastAccessAt = time.Now()
	app.Model.URL.Update(*url)

	app.InfoLog.Printf("Redirecting %s to %s", shortCode, url.OriginalURL)

	c.JSON(http.StatusOK, gin.H{
		"original_url": url.OriginalURL,
		"short_code":   url.ShortCode,
		"message":      "Original URL retrieved successfully.",
	})
}
