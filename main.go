package main

import (
	"Food-delivery/component/appctx"
	"Food-delivery/module/restaurant/transport/ginrestaurant"
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
	dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	log.Println(db, err)

	// tạo 1 server và
	r := gin.Default() // giong nhu 1 server
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			// gin.H là 1 map[string]interface{}, nhu 1 dict hoac object
			"message": "pong",
		})

	})
	appContext := appctx.NewAppContext(db)
	// POST - create a restaurant

	v1 := r.Group("/v1")
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
	restaurants.GET("", func(c *gin.Context) {
		var data []Restaurant

		type Paging struct {
			// bắt buộc phải có form để backend nhận dữ liệu trong body - formdata hoặc query-string
			Page  int `json:"page" form:"page"`
			Limit int `json:"limit" form:"limit"`
		}

		var pagingData Paging
		if err := c.ShouldBind(&pagingData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if pagingData.Page == 0 {
			pagingData.Page = 1
		}
		if pagingData.Limit == 0 {
			pagingData.Limit = 2
		}

		db.Offset((pagingData.Page - 1) * pagingData.Limit).Order("id desc").Limit(pagingData.Limit).Find(&data)
		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	})

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
