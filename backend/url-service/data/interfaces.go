package data

type URLInterface interface {
	Insert(url URL) (int64, error)
	Update(url URL) error
	GetByOriginalURL(originalURL string) (*URL, error)
	GetByShortCode(shortCode string) (*URL, error)
}
