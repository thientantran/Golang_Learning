package main

import (
	"Food-delivery/component/appctx"
	"Food-delivery/component/uploadprovider"
	"Food-delivery/middleware"
	"Food-delivery/pubsub/localpb"
	"Food-delivery/subscriber"
	"fmt"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

type Restaurant struct {
	Id   int    `json:"id" gorm:"column:id;"`
	Name string `json:"name" gorm:"column:name;"`
	Addr string `json:"addr" gorm:"column:addr;"`
}

func (Restaurant) TableName() string {
	return "restaurants"
}

type RestaurantUpdate struct {
	Name *string `json:"name" gorm:"column:name;"`
	Addr *string `json:"addr" gorm:"column:addr;"`
}

func (RestaurantUpdate) TableName() string {
	return Restaurant{}.TableName()
}
func main() {

	//test := Restaurant{
	//	Id:   1,
	//	Name: "KFC",
	//	Addr: "KFC Address",
	//}
	//
	//jsByte, err := json.Marshal(test)
	//log.Println(string(jsByte), err) //{"id":1,"name":"KFC","addr":"KFC Address"}
	////json.Unmarshal([]byte("{\"id\":1,\"name\":\"KFC\",\"addr\":\"KFC Address\"}"), &test)
	//json.Unmarshal(jsByte, &test)
	//log.Println(test) //{1 KFC KFC Address}

	dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}
	// hiển thị log khi db hoạt động
	db = db.Debug()

	AWS_ACCESS_KEY_ID := os.Getenv("AWS_ACCESS_KEY_ID")
	AWS_SECRET_ACCESS_KEY := os.Getenv("AWS_SECRET_ACCESS_KEY")
	AWS_REGION := os.Getenv("AWS_REGION")
	BUCKET_NAME := os.Getenv("BUCKET_NAME")
	secretKet := os.Getenv("SYSTEM_SECRET")
	log.Println(AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_REGION, BUCKET_NAME)
	s3Provider := uploadprovider.NewS3Provider(BUCKET_NAME, AWS_REGION, AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, "https://tan-test-golang.s3-ap-southeast-1.amazonaws.com")
	ps := localpb.NewPubSub()
	// tạo 1 server và
	r := gin.Default() // giong nhu 1 server

	appContext := appctx.NewAppContext(db, s3Provider, secretKet, ps)

	// setup subcribers
	//subscriber.Setup(appContext, context.Background())
	_ = subscriber.NewEngine(appContext).Start()
	//co 3 cach dat middleware
	//1: toan bo
	r.Use(middleware.Recover(appContext))
	//2: theo nhom
	//v1 := r.Group("/v1", middleware.Recover(appContext)
	//3: theo tung route
	//r.GET("/ping", middleware.Recover(appContext), .....
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			// gin.H là 1 map[string]interface{}, nhu 1 dict hoac object
			"message": "pong",
		})

	})

	// POST - create a restaurant
	r.Static("/static", "./static")
	v1 := r.Group("/v1")
	setupRoute(appContext, v1)
	setupAdminRoute(appContext, v1)

	r.StaticFile("/demo/", "demo.html")
	startSocketIOServer(r, appContext)

	r.Run()

	// CREATE
	//
	//newRestaurant := Restaurant{Name: "KFC", Addr: "KFC Address"}
	//
	//if err := db.Create(&newRestaurant).Error; err != nil {
	//	//chỗ này dùng pointer, khi tạo xong thì newRestaurant sẽ có ID
	//	log.Println(err)
	//}
	//log.Println("New Restaurant ID: ", newRestaurant.Id)

	//// Find
	//var myRestaurant Restaurant
	//if err := db.Where("id = ?", 1).First(&myRestaurant).Error; err != nil {
	//	log.Println(err)
	//}
	//log.Println("My Restaurant: ", myRestaurant)
	//
	////UPDATE
	////phải tạo 1 struct Update, trong 1 số trường hợp giá trị là chuỗi rỗng thì cái struct cũ ko update cái chuỗi rỗng này được
	//newData := ""
	//updateData := RestaurantUpdate{Name: &newData}
	//// chỗ này vì restaurentUpdate sử dụng pointer nên phải truyền vào địa chỉ của biến, khi đó trỏ tới newdata là rỗng (vẫn có gía trị bộ nhớ), còn nếu cái restaurant struc cũ thì nó sẽ quét và bỏ qua false, số 0 và chuỗi rỗng

	//if err := db.Where("id = ?", 1).Updates(&updateData).Error; err != nil {
	//	log.Println(err)
	//}
	//log.Println("My Restaurant: ", myRestaurant)
	//
	//// DELETE
	//if err := db.Table(Restaurant{}.TableName()).Where("id = ?", 2).Delete(nil).Error; err != nil {
	//	log.Println(err)
	//}

}

func startSocketIOServer(engine *gin.Engine, appCtx appctx.AppContext) {
	//server, err := socketio.NewServer(nil)
	server, err := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{websocket.Default},
	})

	log.Println("Create connect socket", server)
	if err != nil {
		log.Println("socket io server error: ", err)
	}

	server.OnConnect("/", func(s socketio.Conn) error {
		//s.SetContext("")
		fmt.Println("Socket connected:", s.ID(), " IP:", s.RemoteAddr())

		//s.Join("shipper")
		//server.BroadcastToRoom("/", "shipper", "hello", "shipper")

		// sau khi kết nối xong thì gửi client
		s.Emit("test", "world")

		return nil
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
		// remove socket from socket engine (from app context)
	})

	server.OnEvent("/", "test", func(s socketio.Conn, msg string) {
		log.Println("test event: ", msg)
	})

	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	server.OnEvent("/", "notice", func(s socketio.Conn, p Person) {
		fmt.Println("server receive notice", p.Name, p.Age)

		p.Age = 33
		s.Emit("notice", p)
	})

	go server.Serve()
	//defer server.Close() // bỏ cái này vào thì ko work
	engine.GET("/socket.io/*any", gin.WrapH(server))
	engine.POST("/socket.io/*any", gin.WrapH(server))
}

//func startSocketIOServer(engine *gin.Engine, appCtx appctx.AppContext) {
//	server, err := socketio.NewServer(nil)
//	if err != nil {
//		log.Fatalf("Failed to create socket.io server: %v", err)
//	}
//
//	server.OnConnect("/", func(s socketio.Conn) error {
//		fmt.Println("Socket connected:", s.ID(), " IP:", s.RemoteAddr())
//		s.Emit("test", "world")
//		return nil
//	})
//
//	// Add more event handlers as needed...
//
//	go func() {
//		if err := server.Serve(); err != nil {
//			// Detailed logging for EOF or other errors
//			if err == io.EOF {
//				log.Println("Server stopped: client disconnected")
//			} else {
//				log.Fatalf("socketio listen error: %s\n", err)
//			}
//		}
//	}()
//	defer server.Close()
//
//	engine.GET("/socket.io/*any", gin.WrapH(server))
//	engine.POST("/socket.io/*any", gin.WrapH(server))
//}
