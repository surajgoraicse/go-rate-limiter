package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v5"
	"github.com/surajgoraicse/go-rate-limiter/rate"
)

func getPort() string {
	if len(os.Args) > 1 {
		return ":" + os.Args[1]
	}
	return ":8000"
}

func main() {

	e := echo.New()
	rateLimiter := rate.NewRateLimiter(rate.RateConfig{
		Cap:  10,
		Rate: 1,
	})
	e.Use(rateLimiter)

	e.GET("/", func(c *echo.Context) error {
		
		return c.String(http.StatusOK, "Hello, World!")
	})

	if err := e.Start(getPort()); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	} else {
		e.Logger.Info("server is runnign at http://localhost:" + getPort())
	}
}
