package restaurantlikestorage

import (
	"Food-delivery/common"
	restaurantlikemodel "Food-delivery/module/restaurantlike/model"
	"context"
	"log"
)

func (s *sqlStore) Create(ctx context.Context, data *restaurantlikemodel.Like) error {
	db := s.db
	if err := db.Create(data).Error; err != nil {
		log.Println(err)
		return common.ErrDB(err)
	}
	return nil
}