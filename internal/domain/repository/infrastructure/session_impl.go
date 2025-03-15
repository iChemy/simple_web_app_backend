package infrastructure

import (
	"context"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"github.com/iChemy/simple_web_app_backend/internal/domain/repository"
)

// インメモリデータベース
type sessionRepository struct {
	data map[string]uuid.UUID
	mu   sync.RWMutex
}

func NewSessionRepository() repository.SessionRepository {
	return &sessionRepository{
		data: make(map[string]uuid.UUID),
	}
}

func (s *sessionRepository) SaveSession(_ context.Context, sessionID string, userID uuid.UUID, _ time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[sessionID] = userID
	// TTL (有効期限) を考慮する場合は、時間経過後に削除する処理が必要
	return nil
}

func (s *sessionRepository) GetUserID(_ context.Context, sessionID string) (uuid.UUID, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	userID, exists := s.data[sessionID]
	if !exists {
		return uuid.UUID{}, repository.ErrSessionNotFound
	}
	return userID, nil
}

func (s *sessionRepository) DeleteSession(_ context.Context, sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, sessionID)
	return nil
}
