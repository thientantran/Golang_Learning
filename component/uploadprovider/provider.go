package uploadprovider

import (
	"Food-delivery/common"
	"context"
)

type UploadProvider interface {
	SaveFileUploaded(ctx context.Context, data []byte, path string) (*common.Image, error)
}
