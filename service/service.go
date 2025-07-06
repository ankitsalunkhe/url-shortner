package service

import (
	"context"
	"fmt"

	"github.com/ankitsalunkhe/url-shortner/db"
	"github.com/ankitsalunkhe/url-shortner/retriever"
	base62shortner "github.com/ankitsalunkhe/url-shortner/shortner"
)

type UrlShortnerService interface {
	UpsertShortUrl(context.Context, string) (string, error)
	GetLongUrl(context.Context, string) (string, error)
	DeleteLongUrl(context.Context, string) error
}

type urlShortnerService struct {
	db  db.Database
	ret retriever.Retriever
}

var _ UrlShortnerService = (*urlShortnerService)(nil)

func New(db db.Database, ret retriever.Retriever) urlShortnerService {
	return urlShortnerService{
		db:  db,
		ret: ret,
	}
}

func (s *urlShortnerService) UpsertShortUrl(ctx context.Context, longUrl string) (string, error) {
	shortUrl, err := s.db.GetShortUrl(ctx, longUrl)
	if err != nil {
		return "", fmt.Errorf("checking if long url already exists: %w", err)
	}

	if shortUrl != "" {
		return shortUrl, nil
	}

	base, err := s.ret.GetBase()
	if err != nil {
		return "", fmt.Errorf("unbale to get new range: %w", err)
	}

	shortUrl = base62shortner.New().Generate(base)

	err = s.db.UpsertUrl(ctx, db.Url{
		LongUrl:  longUrl,
		ShortUrl: shortUrl,
	})
	if err != nil {
		return "", fmt.Errorf("save item in db: %w", err)
	}

	return shortUrl, nil
}

func (s *urlShortnerService) GetLongUrl(ctx context.Context, shortUrl string) (string, error) {
	longUrl, err := s.db.GetLongUrl(ctx, db.Url{ShortUrl: shortUrl})
	if err != nil {
		return "", fmt.Errorf("get item from db: %w", err)
	}

	if longUrl == "" {
		return "", fmt.Errorf("no shortUrl found: %s", shortUrl)
	}

	return longUrl, nil
}

func (s *urlShortnerService) DeleteLongUrl(ctx context.Context, shortUrl string) error {
	err := s.db.DeletUrl(ctx, db.Url{ShortUrl: shortUrl})
	if err != nil {
		return fmt.Errorf("get item from db: %w", err)
	}

	return nil
}
