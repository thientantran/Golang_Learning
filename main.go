package main

import (
	"log"
	"os"

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

	// CREATE
	//
	//newRestaurant := Restaurant{Name: "KFC", Addr: "KFC Address"}
	//
	//if err := db.Create(&newRestaurant).Error; err != nil {
	//	//chỗ này dùng pointer, khi tạo xong thì newRestaurant sẽ có ID
	//	log.Println(err)
	//}
	//log.Println("New Restaurant ID: ", newRestaurant.Id)

	// Find
	var myRestaurant Restaurant
	if err := db.Where("id = ?", 1).First(&myRestaurant).Error; err != nil {
		log.Println(err)
	}
	log.Println("My Restaurant: ", myRestaurant)

	//UPDATE
	//phải tạo 1 struct Update, trong 1 số trường hợp giá trị là chuỗi rỗng thì cái struct cũ ko update cái chuỗi rỗng này được
	newData := ""
	updateData := RestaurantUpdate{Name: &newData}
	if err := db.Where("id = ?", 1).Updates(&updateData).Error; err != nil {
		log.Println(err)
	}
	log.Println("My Restaurant: ", myRestaurant)

	// DELETE
	if err := db.Table(Restaurant{}.TableName()).Where("id = ?", 2).Delete(nil).Error; err != nil {
		log.Println(err)
	}

}
