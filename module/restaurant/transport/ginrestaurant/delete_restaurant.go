package ginrestaurant

import (
	"Food-delivery/common"
	"Food-delivery/component/appctx"
	restaurantbiz "Food-delivery/module/restaurant/biz"
	restaurantstorage "Food-delivery/module/restaurant/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DeleteRestaurant(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := appCtx.GetMainDBConnection()
		//id, err := strconv.Atoi(c.Param("id"))
		uid, err := common.FromBase58(c.Param("id"))

		if err != nil {
			//c.JSON(http.StatusBadRequest, gin.H{
			//	"error": err.Error(),
			//})
			//return
			panic(err)
		}
		store := restaurantstorage.NewSQLStore(db)
		biz := restaurantbiz.NewDeleteRestaurantBiz(store)
		if err := biz.DeleteRestaurant(c.Request.Context(), int(uid.GetLocalID())); err != nil {
			//c.JSON(http.StatusBadRequest, gin.H{
			//	"error": err.Error(),
			//})
			//return
			panic(err)
		}
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
