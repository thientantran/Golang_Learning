package ginuser

import (
	"Food-delivery/common"
	"Food-delivery/component/appctx"
	"Food-delivery/component/hasher"
	"Food-delivery/component/tokenprovider/jwt"
	userbiz "Food-delivery/module/user/biz"
	usermodel "Food-delivery/module/user/model"
	userstorage "Food-delivery/module/user/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginUserData usermodel.UserLogin

		if err := c.ShouldBind(&loginUserData); err != nil {
			panic(common.ErrInValidRequest(err))
		}

		db := appCtx.GetMainDBConnection()

		tokenProvider := jwt.NewTokenJWTProvider(appCtx.SecretKey())

		store := userstorage.NewSQLStore(db)
		md5 := hasher.NewMd5Hash()

		business := userbiz.NewLoginBussiness(store, tokenProvider, md5, 60*60*24*30)
		account, err := business.Login(c.Request.Context(), &loginUserData)
		if err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(account))
	}
}
