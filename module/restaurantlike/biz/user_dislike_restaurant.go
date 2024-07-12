package rstlikebiz

import (
	restaurantlikemodel "Food-delivery/module/restaurantlike/model"
	"context"
)

type UserDislikeRestaurantStore interface {
	Delete(ctx context.Context, userId, restaurantId int) error
}

type userDislikeRestaurantBiz struct {
	store UserDislikeRestaurantStore
}

func NewUserDisLikeRestaurantBiz(
	store UserDislikeRestaurantStore,
) *userDislikeRestaurantBiz {
	return &userDislikeRestaurantBiz{store: store}
}

func (biz *userDislikeRestaurantBiz) DislikeRestaurant(
	ctx context.Context,
	userId,
	restaurantId int,
) error {
	err := biz.store.Delete(ctx, userId, restaurantId)
	if err != nil {
		return restaurantlikemodel.ErrCannotDisLikeRestaurant(err)
	}

	return nil
}
