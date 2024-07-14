package subscriber

import (
	"Food-delivery/common"
	"Food-delivery/component/appctx"
	restaurantstorage "Food-delivery/module/restaurant/storage"
	"context"
)

type HasRestaurantId interface {
	GetRestaurantId() int
	//GetUserId() int
}

func IncreaseLikeCountAfterUserLikeRestaurant(appCtx appctx.AppContext, ctx context.Context) {
	c, _ := appCtx.GetPubSub().Subscribe(ctx, common.TopicUserLikeRestaurant)
	store := restaurantstorage.NewSQLStore(appCtx.GetMainDBConnection())
	go func() {
		defer common.AppRecover()
		for {
			msg := <-c
			//neu ko co hasRestaurantId thi ép kiểu
			//likeData := msg.Data().(*restaurantlikemodel.Like)
			//_ = store.IncreaseLikeCount(ctx, likeData.RestaurantId)
			likeData := msg.Data().(HasRestaurantId) // ko cần con trỏ do interface đã là con trỏ rồi
			_ = store.IncreaseLikeCount(ctx, likeData.GetRestaurantId())
		}
	}()
}
