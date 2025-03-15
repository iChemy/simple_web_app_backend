package infrastructure

import (
	"errors"
	"fmt"

	"github.com/iChemy/simple_web_app_backend/internal/domain/repository"
	"gorm.io/gorm"
)

func gormErrorHandling(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("%s:%w", repository.ErrRecordNotFound, err)
	} else if errors.Is(err, gorm.ErrDuplicatedKey) {
		return fmt.Errorf("%s:%w", repository.ErrDuplicatedKey, err)
	}

	return fmt.Errorf("%s:%w", repository.ErrUndefined, err)
}
