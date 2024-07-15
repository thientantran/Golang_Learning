package skio

import (
	"Food-delivery/component/tokenprovider/jwt"
	userstorage "Food-delivery/module/user/storage"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"gorm.io/gorm"
	"sync"
)

type AppContext interface {
	GetMainDBConnection() *gorm.DB
	SecretKey() string
	//GetRealTimeEngine() RealTimeEngine
}

type RealTimeEngine interface {
	UserSockets(userId int) []AppSocket // 1 user có thể có nhiều socket
	EmitToRoom(room string, key string, data interface{}) error
	EmitToUser(userId int, key string, data interface{}) error
	Run(ctx AppContext, engine *gin.Engine) error
}

type rtEngine struct {
	server  *socketio.Server
	storage map[int][]AppSocket
	locker  *sync.RWMutex
}

func NewEngine() *rtEngine {
	return &rtEngine{
		storage: make(map[int][]AppSocket),
		locker:  new(sync.RWMutex),
	}
}

func (engine *rtEngine) saveAppSocket(userId int, appSck AppSocket) {
	engine.locker.Lock()
	//appSck.Jion("order-{ordID}")

	if v, ok := engine.storage[userId]; ok {
		engine.storage[userId] = append(v, appSck)
	} else {
		engine.storage[userId] = []AppSocket{appSck}
	}
	engine.locker.Unlock()
}

func (engine *rtEngine) getAppSocket(userId int) []AppSocket {
	engine.locker.RLock()
	defer engine.locker.RUnlock()
	return engine.storage[userId]
}

func (engine *rtEngine) removeAppSocket(userId int, appSck AppSocket) {
	engine.locker.Lock()
	defer engine.locker.Unlock()

	if v, ok := engine.storage[userId]; ok {
		for i := range v {
			if v[i] == appSck {
				engine.storage[userId] = append(v[:i], v[i+1:]...)
				break
			}
		}
	}
}

func (engine *rtEngine) UserSockets(userId int) []AppSocket {
	var sockets []AppSocket

	if scks, ok := engine.storage[userId]; ok {
		return scks
	}
	return sockets
}

func (engine *rtEngine) EmitToRoom(room string, key string, data interface{}) error {
	engine.server.BroadcastToRoom("/", room, key, data)
	return nil
}

func (engine *rtEngine) EmitToUser(userId int, key string, data interface{}) error {
	sockets := engine.getAppSocket(userId)
	for _, s := range sockets {
		s.Emit(key, data)
	}
	return nil
}

func (engine *rtEngine) Run(appCtx AppContext, r *gin.Engine) error {
	server, err := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{websocket.Default},
	})

	if err != nil {
		return err
	}

	engine.server = server

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("Socket connected:", s.ID(), " IP:", s.RemoteAddr(), s.ID())
		return nil
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("Meet Error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("Socket disconnected:", s.ID(), " IP:", s.RemoteAddr(), "Reason:", reason)
	})
	// setuo

	server.OnEvent("/", "authenticate", func(s socketio.Conn, token string) {
		db := appCtx.GetMainDBConnection()
		store := userstorage.NewSQLStore(db)

		tokenProvider := jwt.NewTokenJWTProvider(appCtx.SecretKey())

		payload, err := tokenProvider.Validate(token)

		if err != nil {
			s.Emit("authentication_failed", err)
			s.Close()
			return
		}

		user, err := store.FindUser(context.Background(), map[string]interface{}{"id": payload.UserID})
		if user.Status == 0 {
			s.Emit("authentication_failed", "you have been banned/deleted")
			s.Close()
			return
		}

		user.Mask(false)

		appSck := NewAppSocket(s, user)
		engine.saveAppSocket(user.Id, appSck)

		s.Emit("authenticated", user)

		//appSck.Join(user.GetRole())
		//if user.GetRole() == "admin" {
		//	appSck.Join("admin")
		//}
	})

	go server.Serve()
	r.GET("/socket.io/*any", gin.WrapH(server))
	r.POST("/socket.io/*any", gin.WrapH(server))

	return nil
}
