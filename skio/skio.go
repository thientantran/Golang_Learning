package skio

import (
	"Food-delivery/common"
	"net"
	"net/url"
)

// interface này lấy từ sockerio ra (copy)
type Conn interface {
	ID() string
	Close() error
	URL() url.URL
	LocalAddr() net.Addr
	RemoteAddr() net.Addr

	// Context of this connection, You can save one context for one connection, and share it between all handlers. The handlers is called in one goroutine
	// so no need to lock context if it only be accessed in one connection.

	Context() interface{}
	SetContext(v interface{})
	Namespace() string
	Emit(msg string, v ...interface{})

	//Broadcase server side apis
	Join(room string)
	Leave(room string)
	LeaveAll()
	Rooms() []string
}

type AppSocket interface {
	Conn
	common.Requester
}

type appSocket struct {
	Conn
	common.Requester
}

func NewAppSocket(conn Conn, requester common.Requester) *appSocket {
	return &appSocket{Conn: conn, Requester: requester}
}
