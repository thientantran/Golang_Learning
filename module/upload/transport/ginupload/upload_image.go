package ginupload

import (
	"Food-delivery/common"
	"Food-delivery/component/appctx"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UploadImage(appCtx appctx.AppContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		fileHeader, err := c.FormFile("file")
		if err != nil {
			panic(err)
		}

		if err := c.SaveUploadedFile(fileHeader, fmt.Sprintf("static/%s", fileHeader.Filename)); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(common.Image{
			Id:        1,
			Url:       "http://localhost:8080/static/" + fileHeader.Filename,
			Width:     100,
			Height:    100,
			CloudName: "local",
			Extension: "png",
		}))
	}
}
