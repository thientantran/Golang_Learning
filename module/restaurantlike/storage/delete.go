package restaurantlikestorage

import (
	"Food-delivery/common"
	restaurantlikemodel "Food-delivery/module/restaurantlike/model"
	"context"
)

func (s *sqlStore) Delete(ctx context.Context, userId, restaurantId int) error {
	db := s.db
	if err := db.Table(restaurantlikemodel.Like{}.TableName()).
		Where("user_id = ? and restaurant_id = ?", userId, restaurantId).
		Delete(nil).
		Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}
