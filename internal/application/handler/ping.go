package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type PingController struct{}

func (pc *PingController) Ping(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
