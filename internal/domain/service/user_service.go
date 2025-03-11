package service

import (
	"context"
	"errors"

	"github.com/gofrs/uuid"
	"github.com/iChemy/simple_web_app_backend/internal/domain/entity"
	"github.com/iChemy/simple_web_app_backend/internal/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

/*
Service も UserService を interface として
userServiceImpl として実装と分離しても良いかもしれないが，
Repository におけるデータベースのように環境依存なものが Service には少ないように思え，
機能と実装が 1 対 1 で割り当てられるような気がしたので UserService を構造体として直接実装している．
*/

type RegisterUserArgs struct {
	Name     string
	Password string
}

type UserService struct {
	r repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return UserService{r: userRepository}
}

func (s UserService) RegisterUser(ctx context.Context, args RegisterUserArgs) (*entity.User, error) {
	userID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(args.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	createUserArgs := repository.CreateUserArgs{
		ID:           userID,
		Name:         args.Name,
		DisplayName:  args.Name,
		PasswordHash: string(hashedPass),
	}

	u, err := s.r.CreateUser(ctx, createUserArgs)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s UserService) LoginUser(ctx context.Context, userName string, password string) (*entity.User, error) {
	u, err := s.r.GetUserByName(ctx, userName)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return u, nil
}
