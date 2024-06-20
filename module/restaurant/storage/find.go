package restaurantstorage

import (
	"Food-delivery/common"
	restaurantmodel "Food-delivery/module/restaurant/model"
	"context"
	"gorm.io/gorm"
)

func (s *sqlStore) FindDataWithCondition(
	context context.Context,
	condition map[string]interface{}, //map where the keys are strings and the values are of type interface{}, which means they can be any type
	moreKeys ...string,
) (*restaurantmodel.Restaurant, error) {
	var data restaurantmodel.Restaurant

	if err := s.db.Where(condition).First(&data).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.RecordNotFound
			// gorm cũng cung cấp một số lỗi cơ bản, trong trường hợp không tìm thấy record thì gorm sẽ trả về lỗi ErrRecordNotFound, nhưng mình muốn tự đồng bộ với bắt lỗi của mình lun
		}
		return nil, common.ErrDB(err)
	}
	return &data, nil
}
