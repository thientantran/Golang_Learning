package ginrestaurant

import (
	"Food-delivery/common"
	"Food-delivery/component/appctx"
	restaurantbiz "Food-delivery/module/restaurant/biz"
	restaurantmodel "Food-delivery/module/restaurant/model"
	restaurantstorage "Food-delivery/module/restaurant/storage"
	restaurantlikestorage "Food-delivery/module/restaurantlike/storage"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func ListRestaurant(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := appCtx.GetMainDBConnection()
		var pagingData common.Paging

		if err := c.ShouldBind(&pagingData); err != nil {
			//c.JSON(http.StatusBadRequest, gin.H{
			//	"error": err.Error(),
			//})
			//return
			panic(common.ErrInValidRequest(err))
		}

		pagingData.Fulfill()
		log.Println(pagingData)
		var filter restaurantmodel.Filter

		if err := c.ShouldBind(&filter); err != nil {
			//c.JSON(http.StatusBadRequest, gin.H{
			//	"error": err.Error(),
			//})
			//return
			panic(common.ErrInValidRequest(err))
		}

		filter.Status = []int{1}

		var result []restaurantmodel.Restaurant
		store := restaurantstorage.NewSQLStore(db)
		likeStore := restaurantlikestorage.NewSQLStore(db)
		biz := restaurantbiz.NewListRestaurantBiz(store, likeStore)
		result, err := biz.ListRestaurant(c.Request.Context(), &filter, &pagingData)
		if err != nil {
			//c.JSON(http.StatusBadRequest, gin.H{
			//	"error": err.Error(),
			//})
			//return
			panic(err)
		}

		for i := range result {
			result[i].Mask(false)
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(result, pagingData, filter))
	}
}
