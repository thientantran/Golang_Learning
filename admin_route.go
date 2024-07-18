package main

import (
	"Food-delivery/component/appctx"
	"Food-delivery/memcache"
	"Food-delivery/middleware"
	userstorage "Food-delivery/module/user/storage"
	"Food-delivery/module/user/transport/ginuser"
	"github.com/gin-gonic/gin"
)

func setupAdminRoute(appContext appctx.AppContext, v1 *gin.RouterGroup) {
	userStore := userstorage.NewSQLStore(appContext.GetMainDBConnection())
	userCachingStore := memcache.NewUserCaching(memcache.NewCaching(), userStore)
	admin := v1.Group("/admin", middleware.RequireAuth(appContext, userCachingStore), middleware.RoleRequired(appContext, "admin", "mod"))

	{
		admin.GET("profile", ginuser.Profile(appContext))
	}
}
