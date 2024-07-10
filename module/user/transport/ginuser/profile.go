package ginuser

import (
	"Food-delivery/common"
	"Food-delivery/component/appctx"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Profile(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		data := c.MustGet(common.CurrentUser).(common.Requester)
		//data.GetUserId()
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data))
	}
}
