package presentation

import "github.com/gofrs/uuid"

type RegisterUserReq struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginUserReq struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password"`
}

type UserRes struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"displayName"`
	Bio         string    `json:"bio"`
}
