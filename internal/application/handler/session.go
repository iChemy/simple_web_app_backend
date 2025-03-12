package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
)

const sessionCookieName = "session_id"

// コンテキストキー用の型を定義
type contextKey string

const userIDKey contextKey = "userID"

type SessionStore interface {
	SaveSession(ctx context.Context, sessionID string, userID uuid.UUID, ttl time.Duration) error
	GetUserID(ctx context.Context, sessionID string) (uuid.UUID, error)
	DeleteSession(ctx context.Context, sessionID string) error
}

type InMemorySessionStore struct {
	data map[string]uuid.UUID
	mu   sync.RWMutex
}

func NewInMemorySessionStore() *InMemorySessionStore {
	return &InMemorySessionStore{
		data: make(map[string]uuid.UUID),
	}
}

func (s *InMemorySessionStore) SaveSession(_ context.Context, sessionID string, userID uuid.UUID, _ time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[sessionID] = userID
	// TTL (有効期限) を考慮する場合は、時間経過後に削除する処理が必要
	return nil
}

func (s *InMemorySessionStore) GetUserID(_ context.Context, sessionID string) (uuid.UUID, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	userID, exists := s.data[sessionID]
	if !exists {
		return uuid.UUID{}, errors.New("session not found")
	}
	return userID, nil
}

func (s *InMemorySessionStore) DeleteSession(_ context.Context, sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, sessionID)
	return nil
}

func SessionMiddleware(sessionStore SessionStore) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Cookie からセッション ID を取得
			cookie, err := c.Cookie(sessionCookieName)
			if err != nil {
				// セッションがない場合はスルー
				return next(c)
			}

			// セッション ID から userID を取得
			userID, err := sessionStore.GetUserID(c.Request().Context(), cookie.Value)
			if err != nil {
				return next(c)
			}

			// `userID` を Context に格納
			ctx := context.WithValue(c.Request().Context(), userIDKey, userID)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

func generateSessionID() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal("failed to generate session ID")
	}
	return hex.EncodeToString(b)
}
