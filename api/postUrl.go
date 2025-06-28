package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Url struct {
	ShortUrl string `dynamodbav:"ShortUrl"`
	LongUrl  string `dynamodbav:"LongUrl"`
}

func (a *API) PostURL(e echo.Context) error {
	var request Request
	if err := e.Bind(&request); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			fmt.Errorf("invalid body: %w", err),
		)
	}

	shortUrl, err := a.service.UpsertShortUrl(e.Request().Context(), request.Url)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			fmt.Errorf("failed to generate short url: %w", err),
		)
	}

	return e.JSON(200, CreateUrl{
		Url: shortUrl,
	})
}
