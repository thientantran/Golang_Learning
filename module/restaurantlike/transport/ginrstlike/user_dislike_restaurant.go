package ginrstlike

import (
	"Food-delivery/common"
	"Food-delivery/component/appctx"
	restaurantstorage "Food-delivery/module/restaurant/storage"
	rstlikebiz "Food-delivery/module/restaurantlike/biz"
	restaurantlikestorage "Food-delivery/module/restaurantlike/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserDislikeRestaurant(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, err := common.FromBase58(c.Param("id"))

		if err != nil {
			panic(common.ErrInValidRequest(err))
		}

		requester := c.MustGet(common.CurrentUser).(common.Requester)

		store := restaurantlikestorage.NewSQLStore(appCtx.GetMainDBConnection())
		decStore := restaurantstorage.NewSQLStore(appCtx.GetMainDBConnection())
		biz := rstlikebiz.NewUserDisLikeRestaurantBiz(store, decStore)

		if err := biz.DislikeRestaurant(c.Request.Context(), requester.GetUserId(), int(uid.GetLocalID())); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
