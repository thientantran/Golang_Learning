package pubsub

import "context"

type Topic string

type Pubsub interface {
	Publish(ctx context.Context, channel Topic, data *Message) error
	Subscribe(ctx context.Context, channel Topic) (ch <-chan *Message, close func())
	//UnSubscribe(ctx context.Context, channel Topic) error
	// sử dụng close để unsub rồi nền ko cần unsub function
}
