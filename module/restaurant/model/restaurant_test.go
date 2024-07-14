package restaurantmodel

import "testing"

type testData struct {
	Input  RestaurantCreate
	Expect error
}

func TestRestaurantCreate_Validate(t *testing.T) {
	dataTable := []testData{
		{Input: RestaurantCreate{Name: ""}, Expect: ErrNameIsEmtpy},
		{Input: RestaurantCreate{Name: "TanTran"}, Expect: nil},
	}

	for _, item := range dataTable {
		err := item.Input.Validate()

		if err != item.Expect {
			t.Errorf("Validate Restaurant. Input: %v, Expect: %v, Ouput: %v", item.Input.Name, item.Expect, err)
		}
	}

}
