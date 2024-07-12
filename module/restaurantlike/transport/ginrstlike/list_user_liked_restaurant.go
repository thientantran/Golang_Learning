package ginrstlike

import (
	"Food-delivery/common"
	"Food-delivery/component/appctx"
	rstlikebiz "Food-delivery/module/restaurantlike/biz"
	restaurantlikemodel "Food-delivery/module/restaurantlike/model"
	restaurantlikestorage "Food-delivery/module/restaurantlike/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ListUserLikedRestaurant(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, err := common.FromBase58(c.Param("id"))
		if err != nil {
			panic(common.ErrInValidRequest(err))
		}

		filter := restaurantlikemodel.Filter{
			RestaurantId: int(uid.GetLocalID()),
		}

		var paging common.Paging

		if err := c.ShouldBind(&paging); err != nil {
			panic(common.ErrInValidRequest(err))
		}

		paging.Fulfill()

		store := restaurantlikestorage.NewSQLStore(appCtx.GetMainDBConnection())
		biz := rstlikebiz.NewListUserLikeRestaurantBiz(store)

		result, err := biz.ListUsers(c.Request.Context(), &filter, &paging)

		if err != nil {
			panic(err)
		}

		for i := range result {
			result[i].Mask(false)
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(result, paging, filter))
	}
}
