package uploadbiz

import (
	"Food-delivery/common"
	"Food-delivery/component/uploadprovider"
	"bytes"
	"context"
	"fmt"

	"io"
	"log"
	"path/filepath"
	"strings"
	"time"
)
import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type CreateImageStorage interface {
	CreateImage(context context.Context, data *common.Image) error
}

type uploadBiz struct {
	provider uploadprovider.UploadProvider
	imgStore CreateImageStorage
}

func NewUploadBiz(provider uploadprovider.UploadProvider, imgStore CreateImageStorage) *uploadBiz {
	return &uploadBiz{
		provider: provider,
		imgStore: imgStore,
	}
}

func (biz *uploadBiz) Upload(ctx context.Context, data []byte, folder, fileName string) (*common.Image, error) {
	fileBytes := bytes.NewBuffer(data)
	w, h, err := getImageDimension(fileBytes)
	if err != nil {
		//return nil, uploadmodel.ErrFileIsNotImage(err)
		return nil, err
	}

	if strings.TrimSpace(folder) == "" {
		folder = "img"
	}

	fileExt := filepath.Ext(fileName) // lấy đuôi file có dấu chấm, ví dụ .png, .jpg
	fileName = fmt.Sprintf("%d%s", time.Now().Nanosecond(), fileExt)

	img, err := biz.provider.SaveFileUploaded(ctx, data, fmt.Sprintf("%s/%s", folder, fileName))
	if err != nil {
		//return nil, uploadmodel.ErrCannotSaveFile(err)
		return nil, err
	}

	img.Width = w
	img.Height = h
	img.Extension = fileExt

	//if err := biz.imgStore.CreateImage(ctx, img); err != nil {
	//	return nil, uploadmodel.ErrCannotSaveFile(err)
	//}
	return img, nil
}

func getImageDimension(reader io.Reader) (int, int, error) {
	img, format, err := image.DecodeConfig(reader)
	if err != nil {
		log.Printf("error decoding image: %v", err)
		return 0, 0, err
	}
	log.Printf("Image format: %s", format)
	return img.Width, img.Height, nil
}
