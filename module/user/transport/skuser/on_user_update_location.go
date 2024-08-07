package skuser

import (
	"Food-delivery/common"
	socketio "github.com/googollee/go-socket.io"
	"gorm.io/gorm"
	"log"
)

type SmallAppContext interface {
	GetMainDBConnection() *gorm.DB
}

type LocationData struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// tại sao nếu import skio.AppContext ở đây thì bị loop circle
func OnUserUpdateLocation(appCtx SmallAppContext, requester common.Requester) func(s socketio.Conn, location LocationData) {
	return func(s socketio.Conn, location LocationData) {
		log.Println("User update location, user id is: ", requester.GetUserId(), "at location: ", location) // do something
	}

}
