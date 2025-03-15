package router

import (
	"fmt"
	"net/http"

	"github.com/iChemy/simple_web_app_backend/internal/application/handler"
	"github.com/iChemy/simple_web_app_backend/internal/domain/service"
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

	r.e.HTTPErrorHandler = newHTTPErrorHandler(r.e)
}

func (r *Router) Start(port int) {
	r.e.Logger.Fatal(r.e.Start(fmt.Sprintf(":%d", port)))
}

func newHTTPErrorHandler(e *echo.Echo) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		var (
			code int
			herr *echo.HTTPError
		)

		srvErr, ok := err.(*service.SrvError)
		if !ok {
			code = http.StatusInternalServerError
			herr = echo.NewHTTPError(code, http.StatusText(code)).SetInternal(err)
		} else {
			code = srvErr.StatusCode()
			herr = echo.NewHTTPError(code, http.StatusText(code)).SetInternal(srvErr)
		}
		e.DefaultHTTPErrorHandler(herr, c)
	}
}
