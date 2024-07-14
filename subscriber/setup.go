package subscriber

import (
	"Food-delivery/component/appctx"
	"context"
)

func Setup(appCtx appctx.AppContext, ctx context.Context) {
	IncreaseLikeCountAfterUserLikeRestaurant(appCtx, ctx)
	DecreaseLikeCountAfterUserDisLikeRestaurant(appCtx, ctx)
	PushNotificationWhenUserLikeRestaurant(appCtx, ctx)
}
