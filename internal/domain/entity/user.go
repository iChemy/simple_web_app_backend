package entity

import "github.com/gofrs/uuid"

type User struct {
	ID           uuid.UUID
	Name         string
	DisplayName  string
	Bio          string
	PasswordHash string
}
