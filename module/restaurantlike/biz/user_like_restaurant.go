package rstlikebiz

import (
	"Food-delivery/common"
	restaurantlikemodel "Food-delivery/module/restaurantlike/model"
	"Food-delivery/pubsub"
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
	//incStore IncreaseLikeCountStore
	ps pubsub.Pubsub
}

func NewUserLikeRestaurantBiz(
	store UserLikeRestaurantStore,
	//incStore IncreaseLikeCountStore,
	ps pubsub.Pubsub,
) *userLikeRestaurantBiz {
	return &userLikeRestaurantBiz{store: store,
		//incStore: incStore,
		ps: ps}
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

	// send message
	if err := biz.ps.Publish(ctx, common.TopicUserLikeRestaurant, pubsub.NewMessage(data)); err != nil {
		log.Println(err)
	}

	////side effect
	//j := asyncjob.NewJob(func(ctx2 context.Context) error {
	//	return biz.incStore.IncreaseLikeCount(ctx2, data.RestaurantId)
	//})
	//if err := asyncjob.NewGroup(true, j).Run(ctx); err != nil {
	//	log.Println(err)
	//}

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
