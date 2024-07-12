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

	// không được làm vậy vì cái này ở nghiệp vụ restaurant like, ko tác động đến restaurant, -> single resposibility
	//db.Exec("UPDATE restaurants SET liked_count = liked_count + 1 WHERE id = ?", data.RestaurantId)
	return nil
}
