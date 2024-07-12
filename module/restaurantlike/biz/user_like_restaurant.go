package rstlikebiz

import (
	restaurantlikemodel "Food-delivery/module/restaurantlike/model"
	"context"
	"log"
)

type UserLikeRestaurantStore interface {
	Create(ctx context.Context, data *restaurantlikemodel.Like) error
}

//type IncreaseLikeCountStore interface {
//	IncreaseLikeCount(ctx context.Context, id int) error
//}

type userLikeRestaurantBiz struct {
	store UserLikeRestaurantStore
}

func NewUserLikeRestaurantBiz(
	store UserLikeRestaurantStore,
) *userLikeRestaurantBiz {
	return &userLikeRestaurantBiz{store: store}
}

func (biz *userLikeRestaurantBiz) LikeRestaurant(
	ctx context.Context,
	data *restaurantlikemodel.Like,
) error {
	err := biz.store.Create(ctx, data)
	log.Println("loi:", err)
	if err != nil {
		return restaurantlikemodel.ErrCannotLikeRestaurant(err)
	}
	return nil
}
