package ginupload

import (
	"Food-delivery/common"
	"Food-delivery/component/appctx"
	uploadbiz "Food-delivery/module/upload/biz"
	"github.com/gin-gonic/gin"
	"net/http"
)

// viêt upload chỉ 1 file, trong trường hợp upload nhiều files thì sẽ chia tải ra trên các node server khác nhau
func UploadImage(appCtx appctx.AppContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		fileHeader, err := c.FormFile("file")
		if err != nil {
			panic(common.ErrInValidRequest(err))
		}

		folder := c.DefaultPostForm("folder", "static")
		file, err := fileHeader.Open()

		if err != nil {
			panic(common.ErrInValidRequest(err))
		}

		defer file.Close()

		dataBytes := make([]byte, fileHeader.Size)
		if _, err := file.Read(dataBytes); err != nil {
			panic(common.ErrInValidRequest(err))
		}
		//imgStore:= uploadstorage.NewSQLStore(db)
		biz := uploadbiz.NewUploadBiz(appCtx.UploadProvider(), nil)
		img, err := biz.Upload(c.Request.Context(), dataBytes, folder, fileHeader.Filename)

		if err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(img))
	}
}
