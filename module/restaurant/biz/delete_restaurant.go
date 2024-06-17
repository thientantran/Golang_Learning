package restaurantbiz

import (
	restaurantmodel "Food-delivery/module/restaurant/model"
	"context"
	"errors"
)

type DeleteRestaurantStore interface {
	Delete(context context.Context, id int) error
	FindDataWithCondition(
		context context.Context,
		condition map[string]interface{}, //map where the keys are strings and the values are of type interface{}, which means they can be any type
		moreKeys ...string,
	) (*restaurantmodel.Restaurant, error)
}

type deleteRestaurantBiz struct {
	store DeleteRestaurantStore
}

func NewDeleteRestaurantBiz(store DeleteRestaurantStore) *deleteRestaurantBiz {
	return &deleteRestaurantBiz{store: store}
}

func (biz *deleteRestaurantBiz) DeleteRestaurant(context context.Context, id int) error {
	oldData, err := biz.store.FindDataWithCondition(context, map[string]interface{}{"id": id})
	if err != nil {
		return err
	}
	if oldData.Status == 0 {
		return errors.New("restaurant has been deleted")
	}

	if err := biz.store.Delete(context, id); err != nil {
		return err
	}
	return nil
}
