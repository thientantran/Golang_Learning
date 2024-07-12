package restaurantrepo

import (
	"Food-delivery/common"
	restaurantmodel "Food-delivery/module/restaurant/model"
	"context"
)

type ListRestaurantStore interface {
	ListDataWithCondition(
		context context.Context,
		filter *restaurantmodel.Filter,
		paging *common.Paging,
		moreKeys ...string,
	) ([]restaurantmodel.Restaurant, error)
}

//type LikeRestaurantRepo interface {
//	GetRestaurantLikes(context context.Context, ids []int) (map[int]int, error)
//}

type listRestaurantRepo struct {
	store ListRestaurantStore
	//likeStore LikeRestaurantRepo
}

func NewListRestaurantRepo(store ListRestaurantStore) *listRestaurantRepo {
	return &listRestaurantRepo{store: store}
}

func (biz *listRestaurantRepo) ListRestaurant(
	context context.Context,
	filter *restaurantmodel.Filter,
	paging *common.Paging,
) ([]restaurantmodel.Restaurant, error) {
	result, err := biz.store.ListDataWithCondition(context, filter, paging, "User")
	if err != nil {
		return nil, common.ErrCannotListEntity(restaurantmodel.EntityName, err)
	}
	//ids := make([]int, len(result))
	//
	//for i := range ids {
	//	ids[i] = result[i].Id
	//}
	//
	//likeMap, err := biz.likeStore.GetRestaurantLikes(context, ids)
	//
	//if err != nil {
	//	log.Println(err)
	//	return result, nil
	//}
	//
	//for i, item := range result {
	//	result[i].LikedCount = likeMap[item.Id]
	//}

	return result, nil
}
