package uuid

import (
	"crypto/rand"
	"encoding/hex"
)

// NewString 生成新的UUID字符串
func NewString() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	// 设置版本为v4
	b[6] = (b[6] & 0x0f) | 0x40
	// 设置变体为RFC4122
	b[8] = (b[8] & 0x3f) | 0x80
	buf := make([]byte, 36)
	hex.Encode(buf[0:8], b[0:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], b[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], b[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], b[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:36], b[10:16])
	return string(buf)
}
