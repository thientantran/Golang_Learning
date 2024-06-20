package common

import (
	"database/sql/driver"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"strconv"
	"strings"
)

//UID is a method to generate an virtual unique identifier for whole system
//its struture contains 62 bits: LocalID - ObjectType - ShardID
//32 bits for localID, max( 2^32-1)
//10 bits for ObjectType
//18 bits for ShardID

type UID struct {
	localID    uint32
	objectType int
	shardID    uint32
}

func NewUID(localID uint32, objectType int, shardID uint32) UID {
	return UID{localID: localID, objectType: objectType, shardID: shardID}
}

// Shard: 1, ObjectType: 1, LocalID: 1 => 0001 0001 0001

// 1 << 8 = 0001 0000 0000
// 1 << 4 = 1 0000
// 1 << 0 =1
// ==> 0001 0001 0001
func (uid UID) String() string {
	val := uint64(uid.localID)<<28 | uint64(uid.objectType)<<18 | uint64(uid.shardID)<<0
	// dùng phép OR để gắn các bit của 3 cái lại với nhau
	return base58.Encode([]byte(fmt.Sprintf("%v", val)))
}

func (uid UID) GetLocalID() uint32 {
	return uid.localID
}

func (uid UID) GetObjectType() int {
	return uid.objectType
}

func (uid UID) GetShardID() uint32 {
	return uid.shardID
}

func DecomposeUID(s string) (UID, error) {
	uid, err := strconv.ParseUint(s, 10, 64)

	if err != nil {
		return UID{}, err
	}

	if (1 << 18) > uid {
		return UID{}, fmt.Errorf("wrong uid")
	}
	// x == 1110 1110 0101 => x >> 4 = 1110 1110 & 0000 1111 = 1110

	u := UID{
		localID:    uint32(uid >> 28),
		objectType: int(uid >> 18 & 0x3FF),     // 0x3FF=1111 1111 1111, toán tử và, bỏ phần nào thì  & với số 0, lấy phần nào thì và với số 1, chỉ giữ lại 10 bit
		shardID:    uint32(uid >> 0 & 0x3FFFF), // 0x3FFFF = 1111 1111 1111 1111 1111, phía trước của uid là 00000, & với 1111 thì sẽ ko còn gì, chỉ giữ đúng 18 bit thôi thây vì 32 bit
	}
	return u, nil
}

func FromBase58(s string) (UID, error) {
	return DecomposeUID(string(base58.Decode(s)))
	// decode s từ base58, sau đó chuyển về string, sau đó gọi hàm DecomposeUID
	// có thể tăng hiệu suất bằng cách truyền trực tiếp base58.Decode(s) vào DecomposeUID
	// có thể tăng bảo mật bằng thêm hash
}

// chưa dược giải thích
func (uid UID) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", uid.String())), nil
}

func (uid *UID) UnmarshalJSON(data []byte) error {
	decodeUID, err := FromBase58(strings.Replace(string(data), "\"", "", -1))
	if err != nil {
		return err
	}
	uid.localID = decodeUID.localID
	uid.shardID = decodeUID.shardID
	uid.objectType = decodeUID.objectType

	return nil
}

func (uid *UID) Value() (driver.Value, error) {
	if uid == nil {
		return nil, nil
	}
	return int64(uid.localID), nil
}

func (uid *UID) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	var i uint32
	switch t := value.(type) {
	case int:
		i = uint32(t)
	case int8:
		i = uint32(t)
	case int16:
		i = uint32(t)
	case int32:
		i = uint32(t)
	case int64:
		i = uint32(t)
	case uint8:
		i = uint32(t)
	case uint16:
		i = uint32(t)
	case uint32:
		i = t
	case uint64:
		i = uint32(t)
	case []byte:
		a, err := strconv.Atoi(string(t))
		if err != nil {
			return err
		}
		i = uint32(a)
	default:
		return fmt.Errorf("Invalid Scan Source")
	}

	*uid = NewUID(i, 0, 1)
	return nil

}
