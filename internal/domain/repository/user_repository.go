package repository

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/iChemy/simple_web_app_backend/internal/domain/entity"
)

type UserRepository interface {
	GetUsers(ctx context.Context) ([]*entity.User, error)
	GetUser(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	CreateUser(ctx context.Context, args CreateUserArgs) (*entity.User, error)
	GetUserByName(ctx context.Context, userName string) (*entity.User, error)
	// `args` の全てのフィールドが更新される．
	// ゼロ値の場合も更新される．
	// 更新したくないフィールドがある場合，これを事前に取得する必要がある．
	UpdateUser(ctx context.Context, userID uuid.UUID, args UpdateUserArgs) (*entity.User, error)
}

type CreateUserArgs struct {
	ID           uuid.UUID
	Name         string
	DisplayName  string
	Bio          string
	PasswordHash string
}

type UpdateUserArgs struct {
	Name         string
	DisplayName  string
	Bio          string
	PasswordHash string
}
