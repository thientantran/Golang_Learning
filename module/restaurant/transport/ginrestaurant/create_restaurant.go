package ginrestaurant

import (
	"Food-delivery/common"
	"Food-delivery/component/appctx"
	restaurantbiz "Food-delivery/module/restaurant/biz"
	restaurantmodel "Food-delivery/module/restaurant/model"
	restaurantstorage "Food-delivery/module/restaurant/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateRestaurant(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := appCtx.GetMainDBConnection()

		//go func() {
		//	defer common.AppRecover()
		//	// khi 1 routine bị lỗi thì chương trình đứng lun, phải đặt cái defer vào đây để bắt lỗi
		//	arr := []int{}
		//	log.Println(arr[0])
		//}()

		//arr := []int{}
		//log.Println(arr[0])

		var data restaurantmodel.RestaurantCreate

		//shouldbind là lấy data từ request và bind vào, dựa vào kiểu dữ liệu của data và struct đã khai báo ở trên
		if err := c.ShouldBind(&data); err != nil {
			//c.JSON(http.StatusBadRequest, common.ErrInValidRequest(err))
			//return
			panic(err) // chi duoc panic o tang transport, neu panic o tang biz thi se bi mat stack trace
		}
		store := restaurantstorage.NewSQLStore(db)
		biz := restaurantbiz.NewCreateRestaurantBiz(store)
		if err := biz.CreateRestaurant(c.Request.Context(), &data); err != nil {
			//c.JSON(http.StatusBadRequest, err)
			//return
			panic(err) // chi duoc panic o tang transport, neu panic o tang biz thi se bi mat stack trace
		}
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data.Id))
	}
}
