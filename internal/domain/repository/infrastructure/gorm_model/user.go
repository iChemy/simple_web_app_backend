package gormmodel

import (
	"github.com/gofrs/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"primaryKey;type:char(36)"`
	Name         string    `gorm:"uniqueIndex;not null;type:varchar(128)"`
	DisplayName  string    `gorm:"not null;type:varchar(256)"`
	Bio          string    `gorm:"not null;type:text"`
	PasswordHash string    `gorm:"not null;type:text"`
}
