package userstorage

import (
	"Food-delivery/common"
	usermodel "Food-delivery/module/user/model"
	"context"
	"go.opencensus.io/trace"
	"gorm.io/gorm"
)

func (s *sqlStore) FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*usermodel.User, error) {
	db := s.db.Table(usermodel.User{}.TableName())
	for i := range moreInfo {
		db = db.Preload(moreInfo[i])
	}

	var user usermodel.User

	_, span := trace.StartSpan(ctx, "store.user.find_user")
	// chú ý phải có defer để chắn chắn thoát
	defer span.End()

	if err := db.Where(conditions).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.RecordNotFound
		}
		return nil, common.ErrDB(err)
	}
	return &user, nil
}
