package middleware

import (
	"Food-delivery/common"
	"Food-delivery/component/appctx"
	"github.com/gin-gonic/gin"
)

func Recover(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.Header("Content-Type", "application/json")
				// neu co dang appErr thi tra ve appErr, neu khong (loi cua golong) co thi tra ve ErrInternal
				if appErr, ok := err.(*common.AppError); ok {
					c.JSON(appErr.StatusCode, appErr)
					panic(err)
					return
				}

				appErr := common.ErrInternal(err.(error))
				c.AbortWithStatusJSON(appErr.StatusCode, appErr)
				panic(err) // panic to log stack trace vi trong gin cung co middleware recover trong gin defaul . engine,do đó sau khi creash thì sẽ up stack, đến recover của mình, sau đó đến recover của gin
				return
			}
		}()
		c.Next()
	}
}
