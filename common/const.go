package common

import "log"

const (
	DbTypeRestaurant = 1
	DbTypeUser       = 2
)

const (
	CurrentUser = "user"
)

const (
	TopicUserLikeRestaurant    = "TopicUserLikeRestaurant"
	TopicUserDisLikeRestaurant = "TopicUserDisLikeRestaurant"
)

func AppRecover() {
	if err := recover(); err != nil {
		log.Println("Recovered from panic: ", err)
	}
}

type Requester interface {
	GetUserId() int
	GetEmail() string
	GetRole() string
}
