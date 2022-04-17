package internals

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func NewServer() *echo.Echo {
	e := echo.New()
	HideBanner(e)
	e.Use(NoCache())
	return e
}

func HideBanner(e *echo.Echo) {
	e.HideBanner = true
	e.HidePort = true
}

func NoCache() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			res := c.Response()
			res.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			res.Header().Set("Pragma", "no-cache")
			res.Header().Set("Expires", "0")
			return next(c)
		}
	}
}

func DefaultJsonContentTypeMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			switch c.Request().Method {
			case http.MethodPost:
				fallthrough
			case http.MethodPut:
				fallthrough
			case http.MethodPatch:
				contentType := c.Request().Header.Get("Content-Type")
				if len(contentType) == 0 || contentType == "application/x-www-form-urlencoded" {
					c.Request().Header.Set("Content-Type", "application/json")
				}
			}

			if next == nil {
				return nil
			}
			return next(c)
		}
	}
}
