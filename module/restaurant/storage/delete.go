package restaurantstorage

import (
	"Food-delivery/common"
	restaurantmodel "Food-delivery/module/restaurant/model"
	"context"
)

func (s *sqlStore) Delete(
	context context.Context,
	id int,
) error {
	// hàm delete chỉ làm đúng bản chất nhiệm vụ của nó là delete, không cần quan tâm đến việc xử lý lỗi
	if err := s.db.Table(restaurantmodel.Restaurant{}.TableName()).
		Where("id = ?", id).
		Updates(map[string]interface{}{"status": 0}).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}
