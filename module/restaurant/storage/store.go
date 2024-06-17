package restaurantstorage

import "gorm.io/gorm"

type sqlStore struct {
	db *gorm.DB
}

// contructor để tạo 1 instance of struct sqlStore
func NewSQLStore(db *gorm.DB) *sqlStore {
	return &sqlStore{
		db: db,
	}
}
