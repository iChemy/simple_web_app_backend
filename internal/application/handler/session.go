package handler

import (
	"github.com/iChemy/simple_web_app_backend/internal/domain/service"
	"github.com/labstack/echo/v4"
)

const sessionCookieName = "session_id"

type SessionController struct {
	ss service.SessionService
}

func (sc *SessionController) SessionMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Cookie からセッション ID を取得
			cookie, err := c.Cookie(sessionCookieName)
			if err != nil {
				// セッションがない場合はスルー
				return next(c)
			}

			ctx, err := sc.ss.GetAndPreserveUserID(c.Request().Context(), cookie.Value)
			if err != nil {
				return next(c)
			}

			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}
