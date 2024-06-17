package restaurantbiz

import (
	restaurantmodel "Food-delivery/module/restaurant/model"
	"context"
	"errors"
)

type CreateRestaurantStore interface {
	CreateRestaurant(context context.Context, data *restaurantmodel.RestaurantCreate) error
}

type createRestaurantBiz struct {
	store CreateRestaurantStore
}

func NewCreateRestaurantBiz(store CreateRestaurantStore) *createRestaurantBiz {
	return &createRestaurantBiz{store: store}
}

func (biz *createRestaurantBiz) CreateRestaurant(context context.Context, data *restaurantmodel.RestaurantCreate) error {
	if data.Name == "" {
		return errors.New("name of restaurant is required")

	}
	if err := biz.store.CreateRestaurant(context, data); err != nil {
		return err
	}

	return nil
}
