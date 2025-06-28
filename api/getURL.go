package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (a *API) GetURL(e echo.Context, shortUrl string) error {
	longUrl, err := a.service.GetLongUrl(e.Request().Context(), shortUrl)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			fmt.Errorf("failed to get long url: %w", err),
		)
	}

	return e.Redirect(http.StatusFound, longUrl)
}
