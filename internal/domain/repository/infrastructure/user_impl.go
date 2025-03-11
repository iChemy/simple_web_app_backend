package infrastructure

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/iChemy/simple_web_app_backend/internal/domain/entity"
	"github.com/iChemy/simple_web_app_backend/internal/domain/repository"
	gormmodel "github.com/iChemy/simple_web_app_backend/internal/domain/repository/infrastructure/gorm_model"
	"gorm.io/gorm"
)

// https://terasolunaorg.github.io/guideline/current/ja/ImplementationAtEachLayer/DomainLayer.html#id11
// によると Repository と RepositoryImpl は同じパッケージに置くとのこと (なので repository/infrastructure においている)
type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return userRepositoryImpl{db}
}

func (r userRepositoryImpl) GetUsers(ctx context.Context) ([]*entity.User, error) {
	us := make([]*gormmodel.User, 0)

	if err := r.db.WithContext(ctx).Find(&us).Error; err != nil {
		return nil, err
	}

	ret := make([]*entity.User, 0, len(us))

	for _, u := range us {
		ret = append(ret, &entity.User{
			ID:           u.ID,
			Name:         u.Name,
			DisplayName:  u.DisplayName,
			Bio:          u.Bio,
			PasswordHash: u.PasswordHash,
		})
	}

	return ret, nil
}

func (r userRepositoryImpl) GetUser(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	u := new(gormmodel.User)

	if err := r.db.WithContext(ctx).Where(&gormmodel.User{ID: userID}).First(u).Error; err != nil {
		return nil, err
	}

	return &entity.User{
		ID:           u.ID,
		Name:         u.Name,
		DisplayName:  u.DisplayName,
		Bio:          u.Bio,
		PasswordHash: u.PasswordHash,
	}, nil
}

func (r userRepositoryImpl) CreateUser(ctx context.Context, args repository.CreateUserArgs) (*entity.User, error) {
	ctxDB := r.db.WithContext(ctx)

	user := gormmodel.User{
		ID:           args.ID,
		Name:         args.Name,
		DisplayName:  args.DisplayName,
		Bio:          args.Bio,
		PasswordHash: args.PasswordHash,
	}

	err := ctxDB.Create(&user).Error
	if err != nil {
		return nil, err
	}

	u := new(gormmodel.User)

	err = ctxDB.Where(&gormmodel.User{ID: args.ID}).First(u).Error

	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:           u.ID,
		Name:         u.Name,
		DisplayName:  u.DisplayName,
		Bio:          u.Bio,
		PasswordHash: u.PasswordHash,
	}, nil
}

func (r userRepositoryImpl) GetUserByName(ctx context.Context, userName string) (*entity.User, error) {
	u := new(gormmodel.User)

	if err := r.db.WithContext(ctx).Where(&gormmodel.User{Name: userName}).First(u).Error; err != nil {
		return nil, err
	}

	return &entity.User{
		ID:           u.ID,
		Name:         u.Name,
		DisplayName:  u.DisplayName,
		Bio:          u.Bio,
		PasswordHash: u.PasswordHash,
	}, nil
}

func (r userRepositoryImpl) UpdateUser(ctx context.Context, userID uuid.UUID, args repository.UpdateUserArgs) (*entity.User, error) {
	ctxDB := r.db.WithContext(ctx)

	user := gormmodel.User{
		ID:           userID,
		Name:         args.Name,
		DisplayName:  args.DisplayName,
		Bio:          args.Bio,
		PasswordHash: args.PasswordHash,
	}

	err := ctxDB.Save(&user).Error
	if err != nil {
		return nil, err
	}

	u := new(gormmodel.User)

	err = ctxDB.Where(&gormmodel.User{ID: userID}).First(u).Error
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:           u.ID,
		Name:         u.Name,
		DisplayName:  u.DisplayName,
		Bio:          u.Bio,
		PasswordHash: u.PasswordHash,
	}, nil
}
