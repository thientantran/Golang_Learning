package restaurantstorage

import (
	restaurantmodel "Food-delivery/module/restaurant/model"
	"context"
)

func (s *sqlStore) FindDataWithCondition(
	context context.Context,
	condition map[string]interface{}, //map where the keys are strings and the values are of type interface{}, which means they can be any type
	moreKeys ...string,
) (*restaurantmodel.Restaurant, error) {
	var data restaurantmodel.Restaurant

	if err := s.db.Where(condition).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}
