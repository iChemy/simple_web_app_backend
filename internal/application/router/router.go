package router

import (
	"fmt"

	"github.com/iChemy/simple_web_app_backend/internal/application/handler"
	"github.com/labstack/echo/v4"
)

type Router struct {
	e *echo.Echo
	h *handler.Handlers
}

func (r *Router) Setup() {
	r.e = echo.New()
	{
		pingGroup := r.e.Group("/ping")
		pingGroup.GET("/", r.h.Ping.Ping)
	}
}

func (r *Router) Start(port int) {
	r.e.Logger.Fatal(r.e.Start(fmt.Sprintf(":%d", port)))
}
