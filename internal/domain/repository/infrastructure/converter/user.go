package converter

import (
	"github.com/iChemy/simple_web_app_backend/internal/domain/entity"
	gormmodel "github.com/iChemy/simple_web_app_backend/internal/domain/repository/infrastructure/gorm_model"
)

/*
https://terasolunaorg.github.io/guideline/current/ja/Overview/ApplicationLayering.html#o-r-mapper
O/R Mapper
データベースと entity の相互マッピング
*/

func ConvertGormModelUserToEntityUser(gormModelUser gormmodel.User) entity.User {
	return entity.User{
		ID:           gormModelUser.ID,
		Name:         gormModelUser.Name,
		DisplayName:  gormModelUser.DisplayName,
		Bio:          gormModelUser.Bio,
		PasswordHash: gormModelUser.PasswordHash,
	}
}
