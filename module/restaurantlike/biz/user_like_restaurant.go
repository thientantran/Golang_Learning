package rstlikebiz

import (
	"Food-delivery/component/asyncjob"
	restaurantlikemodel "Food-delivery/module/restaurantlike/model"
	"context"
	"log"
)

type UserLikeRestaurantStore interface {
	Create(ctx context.Context, data *restaurantlikemodel.Like) error
}

type IncreaseLikeCountStore interface {
	IncreaseLikeCount(ctx context.Context, id int) error
}

type userLikeRestaurantBiz struct {
	store    UserLikeRestaurantStore
	incStore IncreaseLikeCountStore
}

func NewUserLikeRestaurantBiz(
	store UserLikeRestaurantStore,
	incStore IncreaseLikeCountStore,
) *userLikeRestaurantBiz {
	return &userLikeRestaurantBiz{store: store, incStore: incStore}
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
	j := asyncjob.NewJob(func(ctx2 context.Context) error {
		return biz.incStore.IncreaseLikeCount(ctx2, data.RestaurantId)
	})
	if err := asyncjob.NewGroup(true, j).Run(ctx); err != nil {
		log.Println(err)
	}
	//// cái này ko quan trọng, nên có thể để vào đây cho nó chạy ngầm
	//go func() {
	//	defer common.AppRecover()
	//	// ví dụ cái này tốn time để xử lý
	//	time.Sleep(3 * time.Second)
	//	if err := biz.incStore.IncreaseLikeCount(ctx, data.RestaurantId); err != nil {
	//		log.Println(err)
	//		// ko panic gì hết,
	//	}
	//}()

	return nil
}
