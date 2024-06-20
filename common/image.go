package common

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type Image struct {
	Id        int    `json:"id" gorm:"column:id;"`
	Url       string `json:"url" gorm:"column:url;"`
	Width     int    `json:"width" gorm:"column:width;"`
	Height    int    `json:"height" gorm:"column:height;"`
	CloudName string `json:"cloud_name,omitempty" gorm:"-"`
	Extension string `json:"extension,omitempty" gorm:"-"`
}

func (Image) TableName() string { return "images" }

// Scan đi từ database đi ra, chuyển hoá dữ liệu
func (j *Image) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	var img Image

	if err := json.Unmarshal(bytes, &img); err != nil {
		return err
	}
	// set value of pointer j to img
	*j = img
	return nil
}

// Value return json value, implement driver.Valuer interface (driver.Valuer là []byte, string, float64, int64, bool, time.Time)

func (j *Image) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan và Value là 2 method trong 2 interface được built sẵn trong golang, nếu implement 2 method này thì struct sẽ implement cả 2 interface đó
// struct đi xuống database, dùng đến 2 method này

//In this case, the Scan method is defined for the Image type. This means that when you retrieve data from the database and the data needs to be stored in an Image type, the Scan method will be automatically called to convert the database column value to an Image type.

type Images []Image

func (j *Images) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	var img []Image

	if err := json.Unmarshal(bytes, &img); err != nil {
		return err
	}
	*j = img
	return nil
}

func (j *Images) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}
