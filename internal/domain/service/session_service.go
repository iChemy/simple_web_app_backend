package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"time"

	"github.com/gofrs/uuid"
	"github.com/iChemy/simple_web_app_backend/internal/domain/repository"
)

type SessionService interface {
	// Cookie の値から session 情報を取得し
	// どのユーザーのセッションかを特定して， ctx に保持する．
	// session 情報が取得できなかった場合は error を返す
	GetAndPreserveUserID(ctx context.Context, cookieValue string) (context.Context, error)
	GetUserID(ctx context.Context) (uuid.UUID, error)
	// セッションID を生成し返す
	SaveSession(ctx context.Context, userID uuid.UUID, ttl time.Duration) (string, error)
	generateSessionID() string
}

type sessionServiceImpl struct {
	r repository.SessionRepository
}

// コンテキストキー用の型を定義
type contextKey string

const userIDKey contextKey = "userID"

func NewSessionService(sessionRepository repository.SessionRepository) SessionService {
	return &sessionServiceImpl{
		r: sessionRepository,
	}
}

func (s *sessionServiceImpl) GetAndPreserveUserID(ctx context.Context, cookieValue string) (context.Context, error) {
	userID, err := s.r.GetUserID(ctx, cookieValue)
	if err != nil {
		return nil, err
	}

	// `userID` を Context に格納
	return context.WithValue(ctx, userIDKey, userID), nil
}

func (s *sessionServiceImpl) generateSessionID() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal("failed to generate session ID")
	}
	return hex.EncodeToString(b)
}

func (s *sessionServiceImpl) SaveSession(ctx context.Context, userID uuid.UUID, ttl time.Duration) (string, error) {
	sessionID := s.generateSessionID()
	return sessionID, s.r.SaveSession(ctx, sessionID, userID, ttl)
}

func (s *sessionServiceImpl) GetUserID(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, errors.New("ctx has not been preserved userID")
	}

	return userID, nil
}
