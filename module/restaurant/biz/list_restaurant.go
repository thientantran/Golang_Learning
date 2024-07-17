package restaurantbiz

import (
	"Food-delivery/common"
	restaurantmodel "Food-delivery/module/restaurant/model"
	"context"
	"go.opencensus.io/trace"
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
	ctx1, span := trace.StartSpan(ctx, "biz.list_restaurant")
	span.AddAttributes(
		trace.Int64Attribute("page", int64(paging.Page)),
		trace.Int64Attribute("limit", int64(paging.Limit)),
		trace.StringAttribute("cursor", paging.FakeCursor),
	)
	// phải truyền lại ctx1 thay vì ctx, để jaeder phân cấp được cái span nào trong span nào, nghiệp vụ nào trong nghiệp vụ nào
	// slice đã là reference type (con trỏ rồi) nên không cần dùng pointer
	result, err := biz.repo.ListRestaurant(ctx1, filter, paging)

	span.End()

	if err != nil {
		return nil, common.ErrCannotListEntity(restaurantmodel.EntityName, err)
	}

	// nghiệp vụ cua biz chỉ có trách nhiệm là list, chứ ko có trách nhiệm tính toán, tầng store cũng chỉ có trách nhiệm là việc với DB, do đó cần tạo 1 tầng repository để làm việc tính toán cho tầng biz
	return result, nil
}
