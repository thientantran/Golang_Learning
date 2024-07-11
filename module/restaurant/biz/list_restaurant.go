package restaurantbiz

import (
	"Food-delivery/common"
	restaurantmodel "Food-delivery/module/restaurant/model"
	"context"
	"log"
)

type ListRestaurantStore interface {
	ListDataWithCondition(
		context context.Context,
		filter *restaurantmodel.Filter,
		paging *common.Paging,
		moreKeys ...string,
	) ([]restaurantmodel.Restaurant, error)
}

type LikeRestaurantStore interface {
	GetRestaurantLikes(ctx context.Context, ids []int) (map[int]int, error)
	//GetRestaurantLikesOld(ctx context.Context, ids []int) ([]restaurantlikemodel.Like, error)
	//ko nen lam vay vì phải for qua mảng này, rồi phải for thêm 1 lần nữa để tìm số like, phức tạp thuật toán
}

type listRestaurantBiz struct {
	store     ListRestaurantStore
	likeStore LikeRestaurantStore
}

func NewListRestaurantBiz(store ListRestaurantStore, likeStore LikeRestaurantStore) *listRestaurantBiz {
	return &listRestaurantBiz{store: store, likeStore: likeStore}
}

func (biz *listRestaurantBiz) ListRestaurant(
	ctx context.Context,
	filter *restaurantmodel.Filter,
	paging *common.Paging,
) ([]restaurantmodel.Restaurant, error) {
	// slice đã là reference type (con trỏ rồi) nên không cần dùng pointer
	result, err := biz.store.ListDataWithCondition(ctx, filter, paging, "User")
	if err != nil {
		return nil, err
	}

	ids := make([]int, len(result))
	for i := range ids {
		ids[i] = result[i].Id
	}

	likeMap, err := biz.likeStore.GetRestaurantLikes(ctx, ids)

	if err != nil {
		log.Println(err)
		return result, nil
	}

	for i, item := range result {
		result[i].LikedCount = likeMap[item.Id]
	}

	return result, nil
}
