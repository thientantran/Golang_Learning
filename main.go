package main

import (
	"Food-delivery/component/appctx"
	"Food-delivery/middleware"
	"Food-delivery/module/restaurant/transport/ginrestaurant"
	"Food-delivery/module/upload/transport/ginupload"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

	// tạo 1 server và
	r := gin.Default() // giong nhu 1 server
	appContext := appctx.NewAppContext(db)
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
	v1.POST("/upload", ginupload.UploadImage(appContext))

	restaurants := v1.Group("/restaurants")
	restaurants.POST("", ginrestaurant.CreateRestaurant(appContext))

	// GET a restaurant
	restaurants.GET("/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		var data Restaurant
		db.Where("id = ?", id).First(&data)
		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})

	})

	// GET all restaurants
	restaurants.GET("", ginrestaurant.ListRestaurant(appContext))

	// UPDATE a restaurant
	restaurants.PATCH("/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
		var data RestaurantUpdate

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		db.Where("id = ?", id).Updates(&data)

		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	})

	restaurants.DELETE("/:id", ginrestaurant.DeleteRestaurant(appContext))
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
