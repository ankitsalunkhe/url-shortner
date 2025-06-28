//go:generate go tool oapi-codegen --package=api -o=api.gen.go ../spec/api.yaml
package api

import (
	"github.com/ankitsalunkhe/url-shortner/service"
	"github.com/labstack/echo/v4"
)

type API struct {
	echo    *echo.Echo
	service service.UrlShortnerService
}

func New(port int, basePath string, service service.UrlShortnerService) *echo.Echo {
	e := echo.New()

	a := &API{
		echo:    e,
		service: service,
	}

	RegisterHandlersWithBaseURL(e, a, basePath)

	return e
}

func (a *API) GetPing(c echo.Context) error {
	return c.String(200, "pong")
}
