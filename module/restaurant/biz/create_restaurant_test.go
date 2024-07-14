package restaurantbiz

import (
	"Food-delivery/common"
	restaurantmodel "Food-delivery/module/restaurant/model"
	"context"
	"errors"
	"testing"
)

type mokeCreatStore struct{}

func (mokeCreatStore) CreateRestaurant(ctx context.Context, data *restaurantmodel.RestaurantCreate) error {
	// tự tạo để làm lỗi db,
	if data.Name == "TanTran" {
		return common.ErrDB(errors.New("Somethine went wrong in DB"))
	}

	data.Id = 200

	return nil
}

func TestNewCreateRestaurantBiz(t *testing.T) {
	biz := NewCreateRestaurantBiz(mokeCreatStore{})

	dataTest := restaurantmodel.RestaurantCreate{Name: ""}
	err := biz.CreateRestaurant(context.Background(), &dataTest)

	if err == nil || err.Error() != restaurantmodel.ErrNameIsEmtpy.Error() {
		t.Errorf("Failed")
	}

	dataTest = restaurantmodel.RestaurantCreate{Name: "TanTran"}
	err = biz.CreateRestaurant(context.Background(), &dataTest)
	if err == nil {
		t.Errorf("Failed")
	}

	dataTest = restaurantmodel.RestaurantCreate{Name: "GhostCoder"}
	err = biz.CreateRestaurant(context.Background(), &dataTest)

	if err != nil {
		t.Errorf("Failed")
	}

	//t.Log("TestNewCreateRestaurantBiz: Passed")
}
