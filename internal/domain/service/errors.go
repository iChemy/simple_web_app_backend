package service

import (
	"fmt"
	"net/http"
)

// アプリ内で統一的に扱うエラー型
type SrvError struct {
	Code    ErrorCode
	Message string
	Err     error
}

// エラーの種類を定義
type ErrorCode string

const (
	ErrNotFound     ErrorCode = "NOT_FOUND"
	ErrUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrForbidden    ErrorCode = "FORBIDDEN"
	ErrBadRequest   ErrorCode = "BAD_REQUEST"
	ErrConflict     ErrorCode = "CONFLICT"
	ErrInternal     ErrorCode = "INTERNAL"
)

// `error` インターフェースを満たす
func (e *SrvError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// エラーを作成するヘルパー関数
func customError(code ErrorCode, message string, err error) error {
	return &SrvError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// HTTP ステータスコードへの変換
func (e *SrvError) StatusCode() int {
	switch e.Code {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrForbidden:
		return http.StatusForbidden
	case ErrBadRequest:
		return http.StatusBadRequest
	case ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func (e *SrvError) Unwrap() error {
	return e.Err
}
