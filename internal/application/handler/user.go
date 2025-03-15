package handler

import (
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/iChemy/simple_web_app_backend/internal/application/presentation"
	"github.com/iChemy/simple_web_app_backend/internal/domain/service"
	"github.com/labstack/echo/v4"
)

/*
入力値バリデーションとサービスロジックの呼び出し
*/

type UserController struct {
	s  service.UserService
	v  *validator.Validate
	ss service.SessionService
}

type RegisterUserArgs struct {
	Name     string
	Password string
}

func (uc *UserController) Me(c echo.Context) error {
	ctx := c.Request().Context()
	// Context から `userID` を取得
	userID, err := uc.ss.GetUserID(ctx)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	u, err := uc.s.GetUser(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "user not found"})
	}

	res := presentation.UserRes{
		ID:          u.ID,
		Name:        u.Name,
		DisplayName: u.DisplayName,
		Bio:         u.Bio,
	}

	return c.JSON(http.StatusOK, res)
}

func (uc *UserController) RegisterUser(c echo.Context) error {
	req := presentation.RegisterUserReq{}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := uc.v.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	ctx := c.Request().Context()

	args := service.RegisterUserArgs{
		Name:     req.Name,
		Password: req.Password,
	}

	u, err := uc.s.RegisterUser(ctx, args)
	if err != nil {
		// err の内容を直接流して良いのかという問題あり
		return c.String(http.StatusInternalServerError, err.Error())
	}

	res := presentation.UserRes{
		ID:          u.ID,
		Name:        u.Name,
		DisplayName: u.DisplayName,
		Bio:         u.Bio,
	}

	return c.JSON(http.StatusOK, res)
}

// ログイン周りのエラーは全部 unauthorized が良い? セキュリティ的に
func (uc *UserController) LoginUser(c echo.Context) error {
	req := presentation.LoginUserReq{}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := uc.v.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	ctx := c.Request().Context()

	u, err := uc.s.LoginUser(ctx, req.Name, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
	}

	sessionID, err := uc.ss.SaveSession(ctx, u.ID, time.Hour)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create session"})
	}

	// Cookie をセット
	c.SetCookie(&http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionID,
		HttpOnly: true,
		Secure:   true, // HTTPS 必須
		Path:     "/",
		MaxAge:   3600,
	})

	return c.JSON(http.StatusOK, map[string]string{"message": "logged in"})
}
