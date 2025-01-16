package utils

import (
	"encoding/binary"

	"github.com/jxskiss/base62"
	"github.com/sony/sonyflake"
)

var sf *sonyflake.Sonyflake

func init() {
	var st sonyflake.Settings
	sf = sonyflake.NewSonyflake(st)
	if sf == nil {
		panic("sonyflake not created")
	}
}

func GenerateUniqueId() uint64 {
	id, _ := sf.NextID()

	return id
}

func ShortenUrl(id uint64) string {
	byteSlice := make([]byte, 8)
	binary.BigEndian.PutUint64(byteSlice, id)
	encoded := base62.EncodeToString(byteSlice)

	return encoded
}
