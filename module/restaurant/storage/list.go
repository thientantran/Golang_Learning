package restaurantstorage

import (
	"Food-delivery/common"
	restaurantmodel "Food-delivery/module/restaurant/model"
	"context"
)

// hệ thống tăng tải
// step 1: fetch 1000 ids => catching
// step 2: fetch 50 first ids => data page 1
// step 3: page2, omit 50 ids and fetch the next 50ids
// khi cache hết hạn thì làm lại step1
// điểm yếu: dữ liệu thường xuyên ko thay đổi
// điểm mạnh: tăng tải, giảm db query
// Limit OFFSET có 1 vài vấn đề do nó chỉ nhớ số dòng => nên dùng cursor
// có 1 vấn đề, khi bấm trang 2, đúng lúc db thay đổi, sẽ có dữ liệu thêm vào, gây dữ liệu qua trang, qua dòng, gây ra page có dữ liệu của page 1, load lại page 1 thì ko thấy do để qua số dòng
// ngược lại khi xoá thì 1 số dữ liệu page 2 nhảy lên, do đó, đang từ page 1 chuyển sang page 2 sẽ có dấu hiệu thiếu data
func (s *sqlStore) ListDataWithCondition(
	context context.Context,
	filter *restaurantmodel.Filter,
	paging *common.Paging,
	moreKeys ...string,
) ([]restaurantmodel.Restaurant, error) {
	var result []restaurantmodel.Restaurant

	db := s.db.Table(restaurantmodel.Restaurant{}.TableName())
	if f := filter; f != nil {
		if f.OwnerId > 0 {
			db = db.Where("owner_id = ?", f.OwnerId)
		}
		if len(f.Status) > 0 {
			db = db.Where("status in (?)", f.Status)
		}
	}

	if err := db.Count(&paging.Total).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	if v := paging.FakeCursor; v != "" {
		uid, err := common.FromBase58(v)
		if err != nil {
			return nil, common.ErrDB(err)
		}
		db = db.Where("id < ?", uid.GetLocalID())
	} else {
		offset := (paging.Page - 1) * paging.Limit
		db.Offset(offset)
	}
	if err := db.Limit(paging.Limit).Order("id desc").Find(&result).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	if len(result) > 0 {
		last := result[len(result)-1]
		last.Mask(false)
		paging.NextCursor = last.FakeId.String()
	}

	return result, nil
}
