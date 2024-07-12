package restaurantlikemodel

import (
	"Food-delivery/common"
	"fmt"
	"time"
)

const EntityName = "UserLikeRestaurant"

type Like struct {
	RestaurantId int                `json:"restaurant_id" gorm:"column:restaurant_id;"`
	UserId       int                `json:"user_id" gorm:"column:user_id;"`
	CreateAt     *time.Time         `json:"create_at" gorm:"column:created_at;autoCreateTime"`
	User         *common.SimpleUser `json:"user" gorm:"preload:false"`
}

func (Like) TableName() string {
	return "restaurant_likes"
}

//func (l *Like) GetRestaurantId() int {
//	return l.RestaurantId
//}

func ErrCannotLikeRestaurant(err error) *common.AppError {
	return common.NewCustomError(
		err,
		fmt.Sprintf("cannot like restaurant"),
		fmt.Sprintf("ErrCannotLikeRestaurant"),
	)
}

func ErrCannotDisLikeRestaurant(err error) *common.AppError {
	return common.NewCustomError(
		err,
		fmt.Sprintf("cannot unlike restaurant"),
		fmt.Sprintf("ErrCannotDisLikeRestaurant"),
	)
}
