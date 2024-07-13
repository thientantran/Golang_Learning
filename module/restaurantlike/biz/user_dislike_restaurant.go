package rstlikebiz

import (
	"Food-delivery/common"
	restaurantlikemodel "Food-delivery/module/restaurantlike/model"
	"context"
	"log"
	"time"
)

type UserDislikeRestaurantStore interface {
	Delete(ctx context.Context, userId, restaurantId int) error
}

type DecreaseLikeCountStore interface {
	DecreaseLikeCount(ctx context.Context, id int) error
}

type userDislikeRestaurantBiz struct {
	store    UserDislikeRestaurantStore
	decStore DecreaseLikeCountStore
}

func NewUserDisLikeRestaurantBiz(
	store UserDislikeRestaurantStore,
	decStore DecreaseLikeCountStore,
) *userDislikeRestaurantBiz {
	return &userDislikeRestaurantBiz{store: store, decStore: decStore}
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

	go func() {
		defer common.AppRecover()
		// ví dụ cái này tốn time để xử lý
		time.Sleep(3 * time.Second)
		if err := biz.decStore.DecreaseLikeCount(ctx, restaurantId); err != nil {
			log.Println(err)
			// ko panic gì hết,

			// muốn retry lại 3 lần mỗi khi gặp lỗi, nhưng lặp đi lặp lại rất mệt, và khoản cách mỗi lần retry muốn khác nhau
			for i := 0; i < 10; i++ {
				err := biz.decStore.DecreaseLikeCount(ctx, restaurantId)
				if err == nil {
					break
				}
				time.Sleep(2 * time.Second)
			}

		}
	}()
	return nil
}
