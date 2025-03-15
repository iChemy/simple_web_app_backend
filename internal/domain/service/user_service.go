package service

import (
	"context"
	"errors"
	"unicode"

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

type UserService interface {
	GetUser(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	RegisterUser(ctx context.Context, args RegisterUserArgs) (*entity.User, error)
	LoginUser(ctx context.Context, userName string, password string) (*entity.User, error)
}

type userService struct {
	r repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{r: userRepository}
}

func (s *userService) GetUser(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	return s.r.GetUser(ctx, userID)
}

func (s *userService) RegisterUser(ctx context.Context, args RegisterUserArgs) (*entity.User, error) {
	if !isComplexPassword(args.Password) {
		// これは BadRequest として処理されるべき
		return nil, errors.New("password must contain at least one number and one letter")
	}

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
		// ユーザーネームの重複は BadRequest として処理されるべき
		return nil, err
	}

	return u, nil
}

func (s *userService) LoginUser(ctx context.Context, userName string, password string) (*entity.User, error) {
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

func isComplexPassword(s string) bool {
	return hasNumber(s) && hasLetter(s)
}

func hasNumber(s string) bool {
	for _, r := range s {
		if unicode.IsNumber(r) {
			return true
		}
	}
	return false
}

func hasLetter(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) {
			return true
		}
	}
	return false
}
