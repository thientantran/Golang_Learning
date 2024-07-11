package restaurantbiz

import (
	"Food-delivery/common"
	restaurantmodel "Food-delivery/module/restaurant/model"
	"context"
)

type ListRestaurantRepo interface {
	ListRestaurant(
		context context.Context,
		filter *restaurantmodel.Filter,
		paging *common.Paging,
	) ([]restaurantmodel.Restaurant, error)
}

type listRestaurantBiz struct {
	repo ListRestaurantRepo
}

func NewListRestaurantBiz(repo ListRestaurantRepo) *listRestaurantBiz {
	return &listRestaurantBiz{repo: repo}
}

func (biz *listRestaurantBiz) ListRestaurant(
	ctx context.Context,
	filter *restaurantmodel.Filter,
	paging *common.Paging,
) ([]restaurantmodel.Restaurant, error) {
	// slice đã là reference type (con trỏ rồi) nên không cần dùng pointer
	result, err := biz.repo.ListRestaurant(ctx, filter, paging)
	if err != nil {
		return nil, common.ErrCannotListEntity(restaurantmodel.EntityName, err)
	}

	// nghiệp vụ cua biz chỉ có trách nhiệm là list, chứ ko có trách nhiệm tính toán, tầng store cũng chỉ có trách nhiệm là việc với DB, do đó cần tạo 1 tầng repository để làm việc tính toán cho tầng biz
	return result, nil
}
