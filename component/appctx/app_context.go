package appctx

import (
	"Food-delivery/component/uploadprovider"
	"Food-delivery/pubsub"
	"Food-delivery/skio"
	"gorm.io/gorm"
)

type AppContext interface {
	GetMainDBConnection() *gorm.DB
	UploadProvider() uploadprovider.UploadProvider
	SecretKey() string
	GetPubSub() pubsub.Pubsub
	GetRealTimeEngine() skio.RealTimeEngine
}

type appCtx struct {
	db             *gorm.DB
	uploadProvider uploadprovider.UploadProvider
	secretKey      string
	PS             pubsub.Pubsub
	rtEngine       skio.RealTimeEngine
}

func NewAppContext(db *gorm.DB, uploadProvider uploadprovider.UploadProvider, secretKey string, PS pubsub.Pubsub) *appCtx {
	return &appCtx{db: db, uploadProvider: uploadProvider, secretKey: secretKey, PS: PS}
}

func (ctx *appCtx) GetMainDBConnection() *gorm.DB {
	return ctx.db
}
func (ctx *appCtx) UploadProvider() uploadprovider.UploadProvider {
	return ctx.uploadProvider
}
func (ctx *appCtx) SecretKey() string {
	return ctx.secretKey
}
func (ctx *appCtx) GetPubSub() pubsub.Pubsub {
	return ctx.PS
}

func (ctx *appCtx) GetRealTimeEngine() skio.RealTimeEngine {
	return ctx.rtEngine
}

// ko khởi tạo ở trên do dễ vò loop, do chưa có AppContext nên ko có realtime engine
func (ctx *appCtx) SetRealTimeEngine(rt skio.RealTimeEngine) {
	ctx.rtEngine = rt
}
