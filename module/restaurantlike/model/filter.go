package restaurantlikemodel

// dùng để list ra user like restaurant nào hoặc restaurant được like bởi user nào
type Filter struct {
	RestaurantId int `json:"restaurant_id" form:"restaurant_id"`
	UserId       int `json:"user_id" form:"user_id"`
}
