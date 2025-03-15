package repository

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

type SessionRepository interface {
	SaveSession(ctx context.Context, sessionID string, userID uuid.UUID, ttl time.Duration) error
	// セッション ID から userID を取得
	GetUserID(ctx context.Context, sessionID string) (uuid.UUID, error)
	DeleteSession(ctx context.Context, sessionID string) error
}
