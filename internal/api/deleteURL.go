package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (a *API) DeleteURL(e echo.Context, shortUrl string) error {
	err := a.service.DeleteLongUrl(e.Request().Context(), shortUrl)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			fmt.Errorf("failed to delete long url: %w", err),
		)
	}

	return e.JSON(http.StatusOK, DeleteUrl{
		Message: "Ok",
	})
}
