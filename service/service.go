package service

import (
	"context"
	"fmt"

	"github.com/ankitsalunkhe/url-shortner/db"
	base62shortner "github.com/ankitsalunkhe/url-shortner/shortner"
)

type UrlShortnerService interface {
	UpsertShortUrl(context.Context, string) (string, error)
	GetLongUrl(context.Context, string) (string, error)
}

type urlShortnerService struct {
	db db.Database
}

var _ UrlShortnerService = (*urlShortnerService)(nil)

func New(db db.Database) urlShortnerService {
	return urlShortnerService{
		db: db,
	}
}

func (s *urlShortnerService) UpsertShortUrl(ctx context.Context, longUrl string) (string, error) {
	shortUrl := base62shortner.New().Generate(100000000000)

	err := s.db.UpsertUrl(ctx, db.Url{
		LongUrl:  longUrl,
		ShortUrl: shortUrl,
	})
	if err != nil {
		return "", fmt.Errorf("save item in db: %w", err)
	}

	return shortUrl, nil
}

func (s *urlShortnerService) GetLongUrl(ctx context.Context, shortUrl string) (string, error) {
	longUrl, err := s.db.GetUrl(ctx, db.Url{ShortUrl: shortUrl})
	if err != nil {
		return "", fmt.Errorf("get item from db: %w", err)
	}

	return longUrl, nil
}
