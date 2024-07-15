package subscriber

import (
	"Food-delivery/component/appctx"
	restaurantstorage "Food-delivery/module/restaurant/storage"
	"Food-delivery/pubsub"
	"context"
	"log"
)

type HasRestaurantId interface {
	GetRestaurantId() int
	GetUserId() int
}

//func IncreaseLikeCountAfterUserLikeRestaurant(appCtx appctx.AppContext, ctx context.Context) {
//	c, _ := appCtx.GetPubSub().Subscribe(ctx, common.TopicUserLikeRestaurant)
//	store := restaurantstorage.NewSQLStore(appCtx.GetMainDBConnection())
//	go func() {
//		defer common.AppRecover()
//		for {
//			msg := <-c
//			//neu ko co hasRestaurantId thi ép kiểu
//			//likeData := msg.Data().(*restaurantlikemodel.Like)
//			//_ = store.IncreaseLikeCount(ctx, likeData.RestaurantId)
//			likeData := msg.Data().(HasRestaurantId) // ko cần con trỏ do interface đã là con trỏ rồi
//			_ = store.IncreaseLikeCount(ctx, likeData.GetRestaurantId())
//		}
//	}()
//}

func IncreaseLikeCountAfterUserLikeRestaurant(appCtx appctx.AppContext) consumerJob {
	return consumerJob{
		Title: "Increase Like Count After User Like Restaurant",
		Hld: func(ctx context.Context, message *pubsub.Message) error {
			store := restaurantstorage.NewSQLStore(appCtx.GetMainDBConnection())
			likeData := message.Data().(HasRestaurantId)

			return store.IncreaseLikeCount(ctx, likeData.GetRestaurantId())
		},
	}
}

//// thên vào, ví dụ thông báo cho người dùng hay restaurant có người like, blabla
//func PushNotificationWhenUserLikeRestaurant(appCtx appctx.AppContext, ctx context.Context) {
//	c, _ := appCtx.GetPubSub().Subscribe(ctx, common.TopicUserLikeRestaurant)
//
//	go func() {
//		defer common.AppRecover()
//		for {
//			msg := <-c
//			likeData := msg.Data().(HasRestaurantId)
//			log.Println("Push notification when user like restaurant: ", likeData.GetRestaurantId())
//		}
//	}()
//}

func PushNotificationWhenUserLikeRestaurant(appCtx appctx.AppContext) consumerJob {
	return consumerJob{
		Title: "Push Notification when User Like Restaurant",
		Hld: func(ctx context.Context, message *pubsub.Message) error {
			likeData := message.Data().(HasRestaurantId)
			log.Println("Push notification when user like restaurant: ", likeData.GetRestaurantId())

			return nil
		},
	}
}

func EmitRealtimeAfterUserLikeRestaurant(appCtx appctx.AppContext) consumerJob {
	return consumerJob{
		Title: "Realtime emit After User Like Restaurant",
		Hld: func(ctx context.Context, message *pubsub.Message) error {
			likeData := message.Data().(HasRestaurantId)
			//vì realtime engine bỏ vào AppContext, nên lấy ra rồi emit được
			appCtx.GetRealTimeEngine().EmitToUser(likeData.GetUserId(), string(message.Channel()), likeData)
			return nil
		},
	}
}
