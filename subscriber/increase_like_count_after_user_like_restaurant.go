package subscriber

import (
	"Food-delivery/common"
	"Food-delivery/component/appctx"
	restaurantstorage "Food-delivery/module/restaurant/storage"
	"context"
	"log"
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

// thên vào, ví dụ thông báo cho người dùng hay restaurant có người like, blabla
func PushNotificationWhenUserLikeRestaurant(appCtx appctx.AppContext, ctx context.Context) {
	c, _ := appCtx.GetPubSub().Subscribe(ctx, common.TopicUserLikeRestaurant)

	go func() {
		defer common.AppRecover()
		for {
			msg := <-c
			likeData := msg.Data().(HasRestaurantId)
			log.Println("Push notification when user like restaurant: ", likeData.GetRestaurantId())
		}
	}()
}
