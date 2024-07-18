package main

import (
	"Food-delivery/component/appctx"
	"Food-delivery/memcache"
	"Food-delivery/middleware"
	"Food-delivery/module/restaurant/transport/ginrestaurant"
	"Food-delivery/module/restaurantlike/transport/ginrstlike"
	"Food-delivery/module/upload/transport/ginupload"
	userstorage "Food-delivery/module/user/storage"
	"Food-delivery/module/user/transport/ginuser"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func setupRoute(appContext appctx.AppContext, v1 *gin.RouterGroup) {
	userStore := userstorage.NewSQLStore(appContext.GetMainDBConnection())
	userCachingStore := memcache.NewUserCaching(memcache.NewCaching(), userStore)

	v1.POST("/upload", ginupload.UploadImage(appContext))

	v1.POST("/register", ginuser.Register(appContext))
	v1.POST("/authenticate", ginuser.Login(appContext))
	v1.GET("/profile", middleware.RequireAuth(appContext, userCachingStore), ginuser.Profile(appContext))
	restaurants := v1.Group("/restaurants", middleware.RequireAuth(appContext, userCachingStore))
	restaurants.POST("", ginrestaurant.CreateRestaurant(appContext))

	// GET a restaurant
	restaurants.GET("/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		var data Restaurant
		appContext.GetMainDBConnection().Where("id = ?", id).First(&data)
		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})

	})

	// GET all restaurants
	restaurants.GET("", ginrestaurant.ListRestaurant(appContext))

	// UPDATE a restaurant
	restaurants.PATCH("/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
		var data RestaurantUpdate

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		appContext.GetMainDBConnection().Where("id = ?", id).Updates(&data)

		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	})

	restaurants.DELETE("/:id", ginrestaurant.DeleteRestaurant(appContext))
	restaurants.POST("/:id/liked-user", ginrstlike.UserLikeRestaurant(appContext))
	restaurants.DELETE("/:id/liked-user", ginrstlike.UserDislikeRestaurant(appContext))
	restaurants.GET("/:id/liked-user", ginrstlike.ListUserLikedRestaurant(appContext))
}
